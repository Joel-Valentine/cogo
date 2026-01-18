// Package digitalocean is used for interacting with digitalocean.
// currently allows for creation and listing of droplets
package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"github.com/Joel-Valentine/cogo/credentials"
	"github.com/Joel-Valentine/cogo/utils"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"strconv"
)

var imageFork = []utils.SelectItem{{Name: "Distributions", Value: "D"}, {Name: "Applications", Value: "A"}, {Name: "Custom", Value: "C"}}

// CreateDroplet will ask the user a series of questions to determine what kind of
// droplet they would like to be create
// 1. Asks for a digital ocean api token
// 2. Asks what name you would like for the droplet
// 3. Asks what Image you would like to use on the droplet (ubuntu, centos...)
// 4. Asks what size you would like the droplet to be (1gb RAM 1 CPU..)
// 5. Asks what region you want the droplet to be hosted in (London, Amsterdam...)
// 6. Asks what SSH Key you would like to use to access the droplet
// 7. Asks if you are sure with a y/n answer. It will not create a droplet if you chose n
// Finally the droplet is created and returned
func CreateDroplet() (*godo.Droplet, error) {
	digitalOceanToken, tokenError := getToken()

	if tokenError != nil {
		return nil, tokenError
	}

	client := godo.NewFromToken(digitalOceanToken)

	ctx := context.TODO()

	promptDropletName := promptui.Prompt{
		Label: "Droplet Name",
	}

	dropletName, promptDropletError := promptDropletName.Run()

	if promptDropletError != nil {
		fmt.Printf("Droplet name prompt failed %v\n", promptDropletError)
		return nil, promptDropletError
	}

	// Validate after Enter (research pattern: prevents keystroke spam)
	if err := utils.ValidateDropletName(dropletName); err != nil {
		color.Red("✗ Invalid droplet name: %v", err)
		return nil, err
	}

	distAppCustom, distAppCustomErr := utils.AskAndAnswerCustomSelect("Select Image Type", imageFork)

	if distAppCustomErr != nil {
		fmt.Printf("Could not get Image or Distribtion %v\n", distAppCustomErr)
		return nil, distAppCustomErr
	}

	var selectedImage string

	if distAppCustom == "A" {
		selected, err := getSelectedImageApplicationSlug(ctx, client)

		if err != nil {
			fmt.Printf("Failed to get application image slug: %s", err)
			return nil, err
		}

		selectedImage = selected
	}

	if distAppCustom == "D" {
		selected, err := getSelectedImageDistributionSlug(ctx, client)

		if err != nil {
			fmt.Printf("Failed to get distribution image slug: %s", err)
			return nil, err
		}

		selectedImage = selected
	}

	if distAppCustom == "C" {
		selected, err := getSelectedCustomImageSlug(ctx, client)

		if err != nil {
			fmt.Printf("Failed to get custom image slug: %s", err)
			return nil, err
		}

		selectedImage = selected
	}

	selectedSize, err := getSelectedSizeSlug(ctx, client)

	if err != nil {
		fmt.Printf("Failed to get size slug: %s", err)
		return nil, err
	}

	selectedRegion, err := getSelectedRegionSlug(ctx, client)

	if err != nil {
		fmt.Printf("Failed to get region slug: %s", err)
		return nil, err
	}

	sshKeyID, err := getSelectedSSHKeyID(ctx, client)

	if err != nil {
		fmt.Printf("Failed to get SSH key ID: %s", err)
		return nil, err
	}

	shouldCreate, err := confirmCreate("Are you sure? (y/n)")

	if err != nil {
		return nil, err
	}

	if !shouldCreate {
		return nil, err
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: selectedRegion,
		Size:   selectedSize,
		SSHKeys: []godo.DropletCreateSSHKey{
			{ID: sshKeyID},
		},
		Image: godo.DropletCreateImage{
			Slug: selectedImage,
		},
	}

	newDroplet, _, createDropletError := client.Droplets.Create(ctx, createRequest)

	return newDroplet, createDropletError
}

// DestroyDroplet will show the user a list of servers
// upon selecting the server you will have to confirm with y/n
// Once confirmed the user will then have to type in the name of the droplet to make sure they are aware of what they're deleting
// Once entered they will have to do another y/n to confirm that they definitely want it gone
// The deleted droplets name is returned
func DestroyDroplet() (*utils.SelectItem, error) {
	digitalOceanToken, tokenError := getToken()

	if tokenError != nil {
		return nil, tokenError
	}

	client := godo.NewFromToken(digitalOceanToken)

	ctx := context.TODO()

	droplets, err := dropletList(ctx, client)

	if err != nil {
		return nil, err
	}

	// Check for empty state before prompting (prevents crash)
	if len(droplets) == 0 {
		fmt.Println("No droplets found in your DigitalOcean account.")
		fmt.Println()
		fmt.Println("Run 'cogo create' to create a droplet.")
		return nil, nil
	}

	selectItemDroplets := utils.ParseDropletListResults(droplets)

	selectDropletPrompt := utils.CreateCustomSelectPrompt("Select droplet to delete", selectItemDroplets)

	selectedDropletIndex, _, err := selectDropletPrompt.Run()

	if err != nil {
		return nil, err
	}

	selectedDroplet := selectItemDroplets[selectedDropletIndex]

	selectedDropletID, err := strconv.Atoi(selectedDroplet.Value)

	if err != nil {
		return nil, err
	}

	areYouSure, err := confirmCreate("Are you sure? (y/n)")

	if err != nil {
		fmt.Printf("Something went wrong asking you to confirm: %s", err)
		return nil, err
	}

	if !areYouSure {
		fmt.Println("You decided not to delete this droplet")
		return nil, err
	}

	promptReEnterDropletName := promptui.Prompt{
		Label: "Re enter droplet name to confirm delete (WARNING DROPLET WILL BE DELETED FOREVER)",
	}

	enteredDropletName, err := promptReEnterDropletName.Run()

	if err != nil {
		fmt.Printf("Droplet name prompt failed %v\n", err)
		return nil, err
	}

	// Validate after submission instead of on every keystroke
	if len(enteredDropletName) == 0 {
		return nil, errors.New("Must enter the name of the droplet you want to delete")
	}
	if enteredDropletName != selectedDroplet.Name {
		color.Red("✗ Name doesn't match! Expected: %s, Got: %s", selectedDroplet.Name, enteredDropletName)
		return nil, errors.New("Must enter the exact same name to delete")
	}

	fullDropletInfo := droplets[selectedDropletIndex]

	selectedDropletIP, err := fullDropletInfo.PublicIPv4()

	if err != nil {
		fmt.Printf("Something went wrong getting ip of droplet to delete: %s", err)
		return nil, err
	}

	color.Cyan("Name: %s\nSize: %s\nRegion: %s\nImage: %s\nIP: %s", fullDropletInfo.Name, fullDropletInfo.Size.Slug, fullDropletInfo.Region.Name, fullDropletInfo.Image.Name, selectedDropletIP)

	areYouReallyReallySure, err := confirmCreate("Are you really really sure you want to delete this droplet? (y/n)")

	if err != nil {
		fmt.Printf("Something went wrong asking you to confirm deletion: %s", err)
		return nil, err
	}

	if !areYouReallyReallySure {
		fmt.Println("You decided not to delete this droplet")
		return nil, err
	}

	if _, err := client.Droplets.Delete(ctx, selectedDropletID); err != nil {
		fmt.Printf("Something went wrong deleting droplet: %s", err)
		return nil, err
	}

	if enteredDropletName != selectedDroplet.Name {
		fmt.Printf("You entered the droplet name incorrectly")
		return nil, errors.New("Incorrect droplet name")
	}

	return &selectedDroplet, nil
}

// getToken retrieves the DigitalOcean API token using the modern credential manager
// Priority order: CLI flag → Env var → Keychain → Config file → Interactive prompt
func getToken() (string, error) {
	ctx := context.TODO()

	// Create credential manager with all providers including prompt
	manager := credentials.NewManager(
		credentials.NewFlagProvider(""), // Flag support for future use
		credentials.NewEnvProvider(),
		credentials.NewKeychainProvider(),
		credentials.NewFileProvider(),
		credentials.NewPromptProvider(),
	)

	token, source, err := manager.GetToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	// If token came from prompt, offer to save it
	if source.Provider == "prompt" {
		offerToSaveToken(ctx, token)
	}

	return token, nil
}

// offerToSaveToken asks the user if they want to save the token they just entered
func offerToSaveToken(ctx context.Context, token string) {
	prompt := promptui.Prompt{
		Label:     "Save token securely in keychain for future use?",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		// User declined
		return
	}

	// Try keychain first
	keychainProvider := credentials.NewKeychainProvider()
	if keychainProvider.Available() {
		if err := keychainProvider.SetToken(ctx, token); err == nil {
			color.Green("✓ Token saved securely in keychain")
			return
		}
	}

	// Fallback to file if keychain not available
	color.Yellow("⚠  Keychain not available, using file storage")
	fileProvider := credentials.NewFileProvider()
	if err := fileProvider.SetToken(ctx, token); err != nil {
		color.Red("✗ Failed to save token: %v", err)
	}
}

// DisplayDropletList gets all the droplets and formats it with some colours.
// Finally priting it to the terminal
func DisplayDropletList() {
	digitalOceanToken, tokenError := getToken()

	if tokenError != nil {
		fmt.Println("Unable to get DigitalOcean API Token")
	}

	client := godo.NewFromToken(digitalOceanToken)

	ctx := context.TODO()

	dropletList, dropletListError := dropletList(ctx, client)

	if dropletListError != nil {
		fmt.Println("Unable to get a list of droplets")
	}

	// Check for empty state
	if len(dropletList) == 0 {
		fmt.Println("No droplets found in your DigitalOcean account.")
		fmt.Println()
		fmt.Println("Run 'cogo create' to create a droplet.")
		return
	}

	color.Green("\nYour droplets:\n\n")
	for index, element := range dropletList {
		ip, err := element.PublicIPv4()
		if err != nil {
			fmt.Printf("%v Name: %s\n\n", index, element.Name)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			color.Cyan("%v  Name: %s\n   IP: %s\n\n", cyan(index), red(element.Name), red(ip))
		}
	}
}

// dropletList will return a list of droplets for an account using the godo client
func dropletList(ctx context.Context, client *godo.Client) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		list = append(list, droplets...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

// regionList will return a list of regions using the godo client
func regionList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Region{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		regions, resp, err := client.Regions.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's regions to our list
		list = append(list, regions...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseRegionListresults(list)

	return selectList, nil
}

// sizeList will return a list of sizes using the godo client
func sizeList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Size{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Sizes.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's images to our list
		list = append(list, images...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseSizeListResults(list)

	return selectList, nil
}

// imageDistributionList will return a list of distribution images using the godo client
func imageDistributionList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Image{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Images.ListDistribution(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's images to our list
		list = append(list, images...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseImageListResults(list)

	return selectList, nil
}

// imageApplicationList will return a list of application images using the godo client
func imageApplicationList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Image{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Images.ListApplication(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's images to our list
		list = append(list, images...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseImageListResults(list)

	return selectList, nil
}

// imageCustomList will return a list of custom user images using the godo client
func imageCustomList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Image{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Images.ListUser(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's images to our list
		list = append(list, images...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseImageListResults(list)

	return selectList, nil
}

// sshKeyList will return a list of available SSH keys on your account
func sshKeyList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Key{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Keys.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's images to our list
		list = append(list, images...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	selectList := utils.ParseSSHKeyListResults(list)

	return selectList, nil
}

// getSelectedSSHKeyID will get all ssh keys on the account
// asks the user to select one
// once one is selected, convert it into an int
// return the ID (11111111)
func getSelectedSSHKeyID(ctx context.Context, client *godo.Client) (int, error) {
	keyList, err := sshKeyList(ctx, client)

	if err != nil {
		fmt.Printf("Something bad happened getting region list: %s\n\n", err)
		return -1, err
	}

	// Check for empty state
	if len(keyList) == 0 {
		fmt.Println("No SSH keys found in your DigitalOcean account.")
		fmt.Println()
		fmt.Println("Add an SSH key at: https://cloud.digitalocean.com/account/security")
		return -1, errors.New("no SSH keys available")
	}

	selectedKey, err := utils.AskAndAnswerCustomSelect("SSH Key Select", keyList)

	if err != nil {
		fmt.Printf("Failed to ask SSH key question: %s", err)
		return -1, err
	}

	sshKeyID, strconvError := strconv.Atoi(selectedKey)

	if strconvError != nil {
		fmt.Println("ssh key id was not an int")
		return -1, strconvError
	}

	return sshKeyID, nil
}

// getSelectedRegionSlug will get all the regions
// ask the user to chose one
// returns the slug of the region (nyc1)
func getSelectedRegionSlug(ctx context.Context, client *godo.Client) (string, error) {
	regionList, regionListError := regionList(ctx, client)

	if regionListError != nil {
		fmt.Printf("Something bad happened getting region list: %s\n\n", regionListError)
		return "", regionListError
	}

	// Check for empty state
	if len(regionList) == 0 {
		fmt.Println("No regions found in your DigitalOcean account.")
		fmt.Println()
		fmt.Println("Contact DigitalOcean support if you believe this is an error.")
		return "", errors.New("no regions available")
	}

	selectedRegion, err := utils.AskAndAnswerCustomSelect("Region Select", regionList)

	if err != nil {
		fmt.Printf("Failed to ask region question: %s", err)
		return "", err
	}

	return selectedRegion, nil
}

// getSelectedSizeSlug will get all sizes of droplets
// ask the user to chose one
// returns the slug of the chose size (s-1vcpu-1gb)
func getSelectedSizeSlug(ctx context.Context, client *godo.Client) (string, error) {
	sizeList, sizeListError := sizeList(ctx, client)

	if sizeListError != nil {
		fmt.Printf("Something bad happened getting size list: %s\n\n", sizeListError)
		return "", sizeListError
	}

	// Check for empty state
	if len(sizeList) == 0 {
		fmt.Println("No droplet sizes found.")
		fmt.Println()
		fmt.Println("Contact DigitalOcean support if you believe this is an error.")
		return "", errors.New("no droplet sizes available")
	}

	selectedSize, err := utils.AskAndAnswerCustomSelect("Size Select", sizeList)

	if err != nil {
		fmt.Printf("Failed to ask size question, %s", err)
		return "", err
	}

	return selectedSize, nil
}

// getSelectedImageApplicationSlug will get the application images
// asks the use to chose one
// returns the chosen image slug (cassandra, centos)
func getSelectedImageApplicationSlug(ctx context.Context, client *godo.Client) (string, error) {
	imageList, imageListError := imageApplicationList(ctx, client)

	if imageListError != nil {
		fmt.Printf("Something bad happened getting image list: %s\n\n", imageListError)
		return "", imageListError
	}

	// Check for empty state
	if len(imageList) == 0 {
		fmt.Println("No application images found.")
		fmt.Println()
		fmt.Println("Try selecting a different image type or contact DigitalOcean support.")
		return "", errors.New("no application images available")
	}

	selectedImage, err := utils.AskAndAnswerCustomSelect("Image Select", imageList)

	if err != nil {
		fmt.Printf("Failed to ask image question, %s", err)
		return "", err
	}

	return selectedImage, nil
}

// getSelectedImageDistributionSlug will get all the distribution  images
// asks the use to chose one
// returns the chosen image slug (ubuntu-19-10-x64)
func getSelectedImageDistributionSlug(ctx context.Context, client *godo.Client) (string, error) {
	imageList, imageListError := imageDistributionList(ctx, client)

	if imageListError != nil {
		fmt.Printf("Something bad happened getting image list: %s\n\n", imageListError)
		return "", imageListError
	}

	// Check for empty state
	if len(imageList) == 0 {
		fmt.Println("No distribution images found.")
		fmt.Println()
		fmt.Println("Try selecting a different image type or contact DigitalOcean support.")
		return "", errors.New("no distribution images available")
	}

	selectedImage, err := utils.AskAndAnswerCustomSelect("Image Select", imageList)

	if err != nil {
		fmt.Printf("Failed to ask image question, %s", err)
		return "", err
	}

	return selectedImage, nil
}

// getSelectedCustomImageSlug will get all the available images
// asks the use to chose one
// returns the chosen image slug (ubuntu-19-10-x64)
func getSelectedCustomImageSlug(ctx context.Context, client *godo.Client) (string, error) {
	imageList, imageListError := imageCustomList(ctx, client)

	if imageListError != nil {
		fmt.Printf("Something bad happened getting image list: %s\n\n", imageListError)
		return "", imageListError
	}

	// Check for empty state
	if len(imageList) == 0 {
		fmt.Println("No custom images found.")
		fmt.Println()
		fmt.Println("Upload a custom image at: https://cloud.digitalocean.com/images")
		return "", errors.New("no custom images available")
	}

	selectedImage, err := utils.AskAndAnswerCustomSelect("Image Select", imageList)

	if err != nil {
		fmt.Printf("Failed to ask image question, %s", err)
		return "", err
	}

	return selectedImage, nil
}

// confirmCreate asks the user if they are sure they want to create the droplet
// answering with a "y" will return true
func confirmCreate(label string) (bool, error) {
	promptAreYouSure := promptui.Prompt{
		Label: label,
	}

	areYouSure, err := promptAreYouSure.Run()

	if err != nil {
		return false, err
	}

	// Validate after Enter (research pattern: prevents keystroke spam)
	if err := utils.ValidateAreYouSure(areYouSure); err != nil {
		color.Red("✗ %v", err)
		return false, err
	}

	return areYouSure == "y", nil
}
