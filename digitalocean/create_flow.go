package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Joel-Valentine/cogo/navigation"
	"github.com/Joel-Valentine/cogo/utils"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
)

// CreateDropletFlow orchestrates the multi-step droplet creation process
// with back navigation support.
type CreateDropletFlow struct {
	client *godo.Client
	state  navigation.State
	steps  []navigation.Step
}

// NewCreateDropletFlow creates a new droplet creation flow.
func NewCreateDropletFlow(client *godo.Client) *CreateDropletFlow {
	flow := &CreateDropletFlow{
		client: client,
		state:  navigation.NewState(),
	}

	// Define all steps in order
	flow.steps = []navigation.Step{
		&DropletNameStep{},
		&ImageTypeStep{},
		&ImageSelectionStep{client: client},
		&SizeSelectionStep{client: client},
		&RegionSelectionStep{client: client},
		&SSHKeySelectionStep{client: client},
		&CreateConfirmationStep{},
	}

	return flow
}

func (f *CreateDropletFlow) Name() string {
	return "Create Droplet"
}

func (f *CreateDropletFlow) Steps() []navigation.Step {
	return f.steps
}

func (f *CreateDropletFlow) State() navigation.State {
	return f.state
}

// DropletNameStep asks for the droplet name
type DropletNameStep struct{}

func (s *DropletNameStep) Name() string {
	return "droplet_name"
}

func (s *DropletNameStep) Prompt() string {
	return "Droplet Name"
}

func (s *DropletNameStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Check if we already have a value from going back
	if existing, found := state.GetResult(s.Name()); found {
		prompt := navigation.NewInputPrompt(s.Prompt(), existing.Value.(string))
		name, err := prompt.RunWithContext(ctx)
		if err != nil {
			return navigation.Result{}, err
		}
		return navigation.NewResult(name), nil
	}

	// First time - use default
	defaultName := fmt.Sprintf("my-droplet-%d", utils.GenerateTimestamp())
	prompt := navigation.NewInputPrompt(s.Prompt(), defaultName)
	name, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	return navigation.NewResult(name), nil
}

func (s *DropletNameStep) Validate(result navigation.Result) error {
	return navigation.ValidateDropletName(result.Value)
}

func (s *DropletNameStep) Default() interface{} {
	return fmt.Sprintf("my-droplet-%d", utils.GenerateTimestamp())
}

// ImageTypeStep asks which image type (Distributions/Applications/Custom)
type ImageTypeStep struct{}

func (s *ImageTypeStep) Name() string {
	return "image_type"
}

func (s *ImageTypeStep) Prompt() string {
	return "Select Image Type"
}

func (s *ImageTypeStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	items := []string{"Distributions", "Applications", "Custom"}
	
	// Add back option if we can go back
	if state.CanGoBack() {
		items = append([]string{"← Back"}, items...)
	}
	
	prompt := navigation.NewSelectPrompt(s.Prompt(), items)
	index, selected, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Check if user selected "← Back"
	if selected == "← Back" {
		return navigation.Result{}, navigation.ErrGoBack
	}

	// Adjust index if back option was added
	if state.CanGoBack() {
		index--
	}

	// Map to codes: D, A, C
	codes := []string{"D", "A", "C"}
	code := codes[index]

	return navigation.NewResultWithMetadata(code, map[string]interface{}{
		"name":  selected,
		"index": index,
	}), nil
}

func (s *ImageTypeStep) Validate(result navigation.Result) error {
	code := result.Value.(string)
	if code != "D" && code != "A" && code != "C" {
		return navigation.NewValidationError("invalid image type")
	}
	return nil
}

func (s *ImageTypeStep) Default() interface{} {
	return "D" // Distributions
}

// ImageSelectionStep asks for specific image based on type
type ImageSelectionStep struct {
	client *godo.Client
}

func (s *ImageSelectionStep) Name() string {
	return "image_slug"
}

func (s *ImageSelectionStep) Prompt() string {
	return "Select Image"
}

func (s *ImageSelectionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Get image type from previous step
	imageTypeResult, found := state.GetResult("image_type")
	if !found {
		return navigation.Result{}, fmt.Errorf("image type not selected")
	}

	imageType := imageTypeResult.Value.(string)

	// Get appropriate image list
	var imageList []utils.SelectItem
	var err error

	switch imageType {
	case "A":
		imageList, err = imageApplicationList(ctx, s.client)
	case "D":
		imageList, err = imageDistributionList(ctx, s.client)
	case "C":
		imageList, err = imageCustomList(ctx, s.client)
	default:
		return navigation.Result{}, fmt.Errorf("invalid image type: %s", imageType)
	}

	if err != nil {
		return navigation.Result{}, fmt.Errorf("failed to get image list: %w", err)
	}

	// Check for empty state
	if len(imageList) == 0 {
		handler := &navigation.EmptyStateHandler{
			ResourceName:    "images",
			Context:         fmt.Sprintf("for %s", imageTypeResult.Metadata["name"]),
			SuggestedAction: "Try a different image type or contact DigitalOcean support.",
		}
		handler.Display()
		return navigation.Result{}, navigation.ErrEmptyState
	}

	// Convert to string slice for prompt
	items := make([]string, len(imageList))
	for i, img := range imageList {
		items[i] = img.Name
	}

	// Add back option if we can go back
	if state.CanGoBack() {
		items = append([]string{"← Back"}, items...)
	}
	
	prompt := navigation.NewSelectPrompt(s.Prompt(), items)
	index, selected, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Check if user selected "← Back"
	if selected == "← Back" {
		return navigation.Result{}, navigation.ErrGoBack
	}

	// Adjust index if back option was added
	if state.CanGoBack() {
		index--
	}

	selectedImage := imageList[index]

	return navigation.NewResultWithMetadata(selectedImage.Value, map[string]interface{}{
		"name":  selectedImage.Name,
		"index": index,
	}), nil
}

func (s *ImageSelectionStep) Validate(result navigation.Result) error {
	slug := result.Value.(string)
	if slug == "" {
		return navigation.NewValidationError("image slug cannot be empty")
	}
	return nil
}

func (s *ImageSelectionStep) Default() interface{} {
	return nil
}

// SizeSelectionStep asks for droplet size
type SizeSelectionStep struct {
	client *godo.Client
}

func (s *SizeSelectionStep) Name() string {
	return "size_slug"
}

func (s *SizeSelectionStep) Prompt() string {
	return "Select Size"
}

func (s *SizeSelectionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	sizeList, err := sizeList(ctx, s.client)
	if err != nil {
		return navigation.Result{}, fmt.Errorf("failed to get size list: %w", err)
	}

	// Check for empty state
	if len(sizeList) == 0 {
		navigation.DisplayEmptyRegions() // Generic message
		return navigation.Result{}, navigation.ErrEmptyState
	}

	// Convert to string slice
	items := make([]string, len(sizeList))
	for i, size := range sizeList {
		items[i] = size.Name
	}

	// Add back option if we can go back
	if state.CanGoBack() {
		items = append([]string{"← Back"}, items...)
	}
	
	prompt := navigation.NewSelectPrompt(s.Prompt(), items)
	index, selected, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Check if user selected "← Back"
	if selected == "← Back" {
		return navigation.Result{}, navigation.ErrGoBack
	}

	// Adjust index if back option was added
	if state.CanGoBack() {
		index--
	}

	selectedSize := sizeList[index]

	return navigation.NewResultWithMetadata(selectedSize.Value, map[string]interface{}{
		"name":  selectedSize.Name,
		"index": index,
	}), nil
}

func (s *SizeSelectionStep) Validate(result navigation.Result) error {
	slug := result.Value.(string)
	if slug == "" {
		return navigation.NewValidationError("size slug cannot be empty")
	}
	return nil
}

func (s *SizeSelectionStep) Default() interface{} {
	return nil
}

// RegionSelectionStep asks for region
type RegionSelectionStep struct {
	client *godo.Client
}

func (s *RegionSelectionStep) Name() string {
	return "region_slug"
}

func (s *RegionSelectionStep) Prompt() string {
	return "Select Region"
}

func (s *RegionSelectionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	regionList, err := regionList(ctx, s.client)
	if err != nil {
		return navigation.Result{}, fmt.Errorf("failed to get region list: %w", err)
	}

	// Check for empty state
	if len(regionList) == 0 {
		navigation.DisplayEmptyRegions()
		return navigation.Result{}, navigation.ErrEmptyState
	}

	// Convert to string slice
	items := make([]string, len(regionList))
	for i, region := range regionList {
		items[i] = region.Name
	}

	// Add back option if we can go back
	if state.CanGoBack() {
		items = append([]string{"← Back"}, items...)
	}
	
	prompt := navigation.NewSelectPrompt(s.Prompt(), items)
	index, selected, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Check if user selected "← Back"
	if selected == "← Back" {
		return navigation.Result{}, navigation.ErrGoBack
	}

	// Adjust index if back option was added
	if state.CanGoBack() {
		index--
	}

	selectedRegion := regionList[index]

	return navigation.NewResultWithMetadata(selectedRegion.Value, map[string]interface{}{
		"name":  selectedRegion.Name,
		"index": index,
	}), nil
}

func (s *RegionSelectionStep) Validate(result navigation.Result) error {
	slug := result.Value.(string)
	if slug == "" {
		return navigation.NewValidationError("region slug cannot be empty")
	}
	return nil
}

func (s *RegionSelectionStep) Default() interface{} {
	return nil
}

// SSHKeySelectionStep asks for SSH key
type SSHKeySelectionStep struct {
	client *godo.Client
}

func (s *SSHKeySelectionStep) Name() string {
	return "ssh_key_id"
}

func (s *SSHKeySelectionStep) Prompt() string {
	return "Select SSH Key"
}

func (s *SSHKeySelectionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	keyList, err := sshKeyList(ctx, s.client)
	if err != nil {
		return navigation.Result{}, fmt.Errorf("failed to get SSH key list: %w", err)
	}

	// Check for empty state
	if len(keyList) == 0 {
		navigation.DisplayEmptySSHKeys()
		return navigation.Result{}, navigation.ErrEmptyState
	}

	// Convert to string slice
	items := make([]string, len(keyList))
	for i, key := range keyList {
		items[i] = key.Name
	}

	// Add back option if we can go back
	if state.CanGoBack() {
		items = append([]string{"← Back"}, items...)
	}
	
	prompt := navigation.NewSelectPrompt(s.Prompt(), items)
	index, selected, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Check if user selected "← Back"
	if selected == "← Back" {
		return navigation.Result{}, navigation.ErrGoBack
	}

	// Adjust index if back option was added
	if state.CanGoBack() {
		index--
	}

	selectedKey := keyList[index]

	// Convert to int
	keyID, convErr := strconv.Atoi(selectedKey.Value)
	if convErr != nil {
		return navigation.Result{}, fmt.Errorf("invalid SSH key ID: %w", convErr)
	}

	return navigation.NewResultWithMetadata(keyID, map[string]interface{}{
		"name":  selectedKey.Name,
		"index": index,
	}), nil
}

func (s *SSHKeySelectionStep) Validate(result navigation.Result) error {
	keyID := result.Value.(int)
	if keyID <= 0 {
		return navigation.NewValidationError("SSH key ID must be positive")
	}
	return nil
}

func (s *SSHKeySelectionStep) Default() interface{} {
	return nil
}

// CreateConfirmationStep shows summary and asks for confirmation
type CreateConfirmationStep struct{}

func (s *CreateConfirmationStep) Name() string {
	return "confirm"
}

func (s *CreateConfirmationStep) Prompt() string {
	return "Create this droplet?"
}

func (s *CreateConfirmationStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Get all previous selections
	name, _ := state.GetResult("droplet_name")
	imageSlug, _ := state.GetResult("image_slug")
	sizeSlug, _ := state.GetResult("size_slug")
	regionSlug, _ := state.GetResult("region_slug")
	sshKeyID, _ := state.GetResult("ssh_key_id")

	// Display summary
	fmt.Println()
	color.Cyan("=== Droplet Configuration ===")
	fmt.Printf("Name:     %s\n", name.Value)
	fmt.Printf("Image:    %s\n", imageSlug.Metadata["name"])
	fmt.Printf("Size:     %s\n", sizeSlug.Metadata["name"])
	fmt.Printf("Region:   %s\n", regionSlug.Metadata["name"])
	fmt.Printf("SSH Key:  %s\n", sshKeyID.Metadata["name"])
	color.Cyan("=============================")
	fmt.Println()

	// Confirm (default to No for safety)
	prompt := navigation.NewConfirmPrompt(s.Prompt(), false)
	confirmed, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	return navigation.NewResult(confirmed), nil
}

func (s *CreateConfirmationStep) Validate(result navigation.Result) error {
	return nil // Boolean result, always valid
}

func (s *CreateConfirmationStep) Default() interface{} {
	return true
}

// ExecuteCreateFlow runs the entire create droplet flow with back navigation
func ExecuteCreateFlow(client *godo.Client) (*godo.Droplet, error) {
	ctx := context.Background()
	navigator := navigation.NewNavigator()
	flow := NewCreateDropletFlow(client)

	// Run the flow
	result, err := navigator.Run(ctx, flow)

	// Handle special cases
	if err == navigation.ErrEmptyState {
		// Already displayed message in step
		return nil, nil
	}

	if err == navigation.ErrCancel {
		color.Cyan("\nDroplet creation canceled.")
		return nil, nil
	}

	if err == context.Canceled {
		fmt.Println() // Clean newline after ^C
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("create flow failed: %w", err)
	}

	// Check if user confirmed
	confirmed := result.Value.(bool)
	if !confirmed {
		color.Cyan("Droplet creation canceled.")
		return nil, nil
	}

	// Extract all selections from state
	state := flow.State()
	nameResult, _ := state.GetResult("droplet_name")
	imageSlugResult, _ := state.GetResult("image_slug")
	sizeSlugResult, _ := state.GetResult("size_slug")
	regionSlugResult, _ := state.GetResult("region_slug")
	sshKeyIDResult, _ := state.GetResult("ssh_key_id")

	// Create the droplet using godo
	dropletName := nameResult.Value.(string)
	imageSlug := imageSlugResult.Value.(string)
	sizeSlug := sizeSlugResult.Value.(string)
	regionSlug := regionSlugResult.Value.(string)
	sshKeyID := sshKeyIDResult.Value.(int)

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: regionSlug,
		Size:   sizeSlug,
		Image: godo.DropletCreateImage{
			Slug: imageSlug,
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{ID: sshKeyID},
		},
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create droplet: %w", err)
	}

	return droplet, nil
}

