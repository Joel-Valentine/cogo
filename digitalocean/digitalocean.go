package digitalocean

import (
	"context"
	"fmt"
	"github.com/Midnight-Conqueror/cogo/config"
	"github.com/Midnight-Conqueror/cogo/utils"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"strconv"
)

// CreateDroplet will create a droplet and return the created droplet
func CreateDroplet() (*godo.Droplet, error) {

	digitalOceanToken, tokenError := getToken()

	if tokenError != nil {
		return nil, tokenError
	}

	client := godo.NewFromToken(digitalOceanToken)

	ctx := context.TODO()

	promptDropletName := promptui.Prompt{
		Label:    "Droplet Name",
		Validate: utils.ValidateDropletName,
	}

	dropletName, promptDropletError := promptDropletName.Run()

	if promptDropletError != nil {
		fmt.Printf("Droplet name prompt failed %v\n", promptDropletError)
		return nil, promptDropletError
	}

	imageList, imageListError := ImageList(ctx, client)

	if imageListError != nil {
		fmt.Printf("Something bad happened getting image list: %s\n\n", imageListError)
		return nil, imageListError
	}

	imagePrompt := utils.CreateCustomSelectPrompt("Image Select", imageList)

	imageIndex, _, imagePromptError := imagePrompt.Run()

	if imagePromptError != nil {
		fmt.Printf("Prompt failed %v\n", imagePromptError)
		return nil, imagePromptError
	}

	selectedImage := imageList[imageIndex]

	sizeList, sizeListError := SizeList(ctx, client)

	if sizeListError != nil {
		fmt.Printf("Something bad happened getting size list: %s\n\n", sizeListError)
		return nil, sizeListError
	}

	sizePrompt := utils.CreateCustomSelectPrompt("Size Select", sizeList)

	sizeIndex, _, sizePromptError := sizePrompt.Run()

	if sizePromptError != nil {
		fmt.Printf("Prompt failed %v\n", sizePromptError)
		return nil, sizePromptError
	}

	selectedSize := sizeList[sizeIndex]

	regionList, regionListError := regionList(ctx, client)

	if regionListError != nil {
		fmt.Printf("Something bad happened getting region list: %s\n\n", regionListError)
		return nil, regionListError
	}

	regionPrompt := utils.CreateCustomSelectPrompt("Region Select", regionList)

	regionIndex, _, regionPromptError := regionPrompt.Run()

	if regionPromptError != nil {
		fmt.Printf("Prompt failed %v\n", regionPromptError)
		return nil, regionPromptError
	}

	selectedRegion := regionList[regionIndex]

	keyList, sshKeyListError := SSHKeyList(ctx, client)

	if sshKeyListError != nil {
		fmt.Printf("Something bad happened getting region list: %s\n\n", sshKeyListError)
		return nil, sshKeyListError
	}

	sshKeyPrompt := utils.CreateCustomSelectPrompt("SSH Key Select", keyList)

	keyIndex, _, sshKeyPromptError := sshKeyPrompt.Run()

	if sshKeyPromptError != nil {
		fmt.Printf("Prompt failed %v\n", sshKeyPromptError)
		return nil, sshKeyPromptError
	}

	selectedKey := keyList[keyIndex]

	sshKeyID, strconvError := strconv.Atoi(selectedKey.Value)

	if strconvError != nil {
		fmt.Println("ssh key id was not an int")
		return nil, strconvError
	}

	promptAreYouSure := promptui.Prompt{
		Label:    "Are you sure? (y/n)",
		Validate: utils.ValidateAreYouSure,
	}

	areYouSure, areYouSureError := promptAreYouSure.Run()

	if areYouSureError != nil {
		return nil, areYouSureError
	}

	if areYouSure != "y" {
		return nil, areYouSureError
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: selectedRegion.Value,
		Size:   selectedSize.Value,
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{ID: sshKeyID},
		},
		Image: godo.DropletCreateImage{
			Slug: selectedImage.Value,
		},
	}

	newDroplet, _, createDropletError := client.Droplets.Create(ctx, createRequest)

	return newDroplet, createDropletError
}

func getToken() (string, error) {
	configSetup, err := config.Config()

	var digitalOceanToken string

	if err != nil && err.Code == 01 {
		promptDigitalOceanToken := promptui.Prompt{
			Label: "Enter your Digital Ocean API Token",
		}

		token, err := promptDigitalOceanToken.Run()

		digitalOceanToken = token

		if err != nil {
			fmt.Printf("Token prompt failed %v\n", err)
			return "", err
		}

		saveTokenPrompt := promptui.Prompt{
			Label:    "Do you want to save this to ./.cogo_config.json? (y/n)",
			Validate: utils.ValidateAreYouSure,
		}

		saveToken, err := saveTokenPrompt.Run()

		if saveToken == "y" {
			config.SaveConfigFile(digitalOceanToken)
		}

		if err != nil {
			fmt.Printf("Save Token prompt failed %v\n", err)
			return "", err
		}
	} else {
		digitalOceanToken = configSetup.GetString("digitalOceanToken")
	}
	return digitalOceanToken, nil
}

// DisplayDropletList gets all the droplets and formats it nicely to be printed
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
		for _, d := range droplets {
			list = append(list, d)
		}

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

// RegionList will list all the regions from DO
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
		for _, d := range regions {
			list = append(list, d)
		}

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

// ImageList will list all the regions from DO
func ImageList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Image{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Images.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's regions to our list
		for _, d := range images {
			list = append(list, d)
		}

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

// SizeList will list all the regions from DO
func SizeList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Size{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Sizes.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's regions to our list
		for _, d := range images {
			list = append(list, d)
		}

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

// SSHKeyList will list all the regions from DO
func SSHKeyList(ctx context.Context, client *godo.Client) ([]utils.SelectItem, error) {
	// create a list to hold our droplets
	list := []godo.Key{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Keys.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's regions to our list
		for _, d := range images {
			list = append(list, d)
		}

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
