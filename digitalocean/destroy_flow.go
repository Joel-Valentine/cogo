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

type DestroyDropletFlow struct {
	client   *godo.Client
	droplets []godo.Droplet
	state    navigation.State
	steps    []navigation.Step
}

func NewDestroyDropletFlow(client *godo.Client, droplets []godo.Droplet) *DestroyDropletFlow {
	flow := &DestroyDropletFlow{
		client:   client,
		droplets: droplets,
		state:    navigation.NewState(),
	}

	// Define all steps in order
	flow.steps = []navigation.Step{
		&SelectDropletToDestroyStep{droplets: droplets},
		&ConfirmDestroyStep{},
		&ReEnterDropletNameStep{droplets: droplets}, // Final confirmation - typing name is enough
	}

	return flow
}

func (f *DestroyDropletFlow) Name() string {
	return "Destroy Droplet"
}

func (f *DestroyDropletFlow) Steps() []navigation.Step {
	return f.steps
}

func (f *DestroyDropletFlow) State() navigation.State {
	return f.state
}

type SelectDropletToDestroyStep struct {
	droplets []godo.Droplet
}

func (s *SelectDropletToDestroyStep) Name() string {
	return "select_droplet"
}

func (s *SelectDropletToDestroyStep) Prompt() string {
	return "Select droplet to delete"
}

func (s *SelectDropletToDestroyStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Convert droplets to SelectItems
	selectItems := utils.ParseDropletListResults(s.droplets)

	// Convert to string slice
	items := make([]string, len(selectItems))
	for i, item := range selectItems {
		items[i] = item.Name
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

	selectedDroplet := s.droplets[index]

	return navigation.NewResultWithMetadata(selectedDroplet.ID, map[string]interface{}{
		"name":  selectedDroplet.Name,
		"index": index,
	}), nil
}

func (s *SelectDropletToDestroyStep) Validate(result navigation.Result) error {
	dropletID := result.Value.(int)
	if dropletID <= 0 {
		return navigation.NewValidationError("invalid droplet ID")
	}
	return nil
}

func (s *SelectDropletToDestroyStep) Default() interface{} {
	return nil
}

type ConfirmDestroyStep struct{}

func (s *ConfirmDestroyStep) Name() string {
	return "confirm_destroy"
}

func (s *ConfirmDestroyStep) Prompt() string {
	return "Are you sure?"
}

func (s *ConfirmDestroyStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Get selected droplet
	dropletResult, _ := state.GetResult("select_droplet")
	dropletName := dropletResult.Metadata["name"].(string)

	fmt.Println()
	color.Yellow("⚠️  WARNING: You are about to delete droplet: %s", dropletName)
	fmt.Println()

	prompt := navigation.NewConfirmPrompt(s.Prompt(), false) // Default to No for safety
	confirmed, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	if !confirmed {
		return navigation.Result{}, navigation.ErrCancel
	}

	return navigation.NewResult(confirmed), nil
}

func (s *ConfirmDestroyStep) Validate(result navigation.Result) error {
	return nil
}

func (s *ConfirmDestroyStep) Default() interface{} {
	return false
}

type ReEnterDropletNameStep struct {
	droplets []godo.Droplet
}

func (s *ReEnterDropletNameStep) Name() string {
	return "reenter_name"
}

func (s *ReEnterDropletNameStep) Prompt() string {
	return "Re-enter droplet name to confirm delete (WARNING: DROPLET WILL BE DELETED FOREVER)"
}

func (s *ReEnterDropletNameStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
	// Get selected droplet
	dropletResult, _ := state.GetResult("select_droplet")
	expectedName := dropletResult.Metadata["name"].(string)

	// Show the droplet name clearly so user can see what to type
	fmt.Println()
	color.Red("⚠️  FINAL CONFIRMATION: Type the droplet name and press Enter to DELETE FOREVER")
	color.Red("⚠️  This action is IRREVERSIBLE - all data will be lost!")
	fmt.Println()
	color.Cyan("Droplet name: %s", expectedName)
	fmt.Println()

	prompt := navigation.NewInputPrompt("Type droplet name to DELETE (pressing Enter will destroy it)", "")
	enteredName, err := prompt.RunWithContext(ctx)
	if err != nil {
		return navigation.Result{}, err
	}

	// Validate name matches
	if enteredName != expectedName {
		color.Red("✗ Name doesn't match! Expected: %s, Got: %s", expectedName, enteredName)
		fmt.Println()
		// Re-prompt by returning validation error
		return navigation.Result{}, navigation.NewValidationError(fmt.Sprintf("name must match exactly: %s", expectedName))
	}

	return navigation.NewResult(enteredName), nil
}

func (s *ReEnterDropletNameStep) Validate(result navigation.Result) error {
	name := result.Value.(string)
	if name == "" {
		return navigation.NewValidationError("must enter the name of the droplet you want to delete")
	}
	return nil
}

func (s *ReEnterDropletNameStep) Default() interface{} {
	return ""
}

func ExecuteDestroyFlow(client *godo.Client) (*utils.SelectItem, error) {
	ctx := context.Background()

	// Get droplet list first
	droplets, err := dropletList(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get droplet list: %w", err)
	}

	// Check for empty state
	if len(droplets) == 0 {
		fmt.Println("No droplets found in your DigitalOcean account.")
		fmt.Println()
		fmt.Println("Run 'cogo create' to create a droplet.")
		return nil, nil
	}

	// Create and run flow
	navigator := navigation.NewNavigator()
	flow := NewDestroyDropletFlow(client, droplets)

	_, err = navigator.Run(ctx, flow)

	// Handle special cases
	if err == navigation.ErrEmptyState {
		// Already displayed message
		return nil, nil
	}

	if err == navigation.ErrCancel {
		color.Cyan("\nDroplet destruction canceled.")
		return nil, nil
	}

	if err == context.Canceled {
		fmt.Println() // Clean newline after ^C
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("destroy flow failed: %w", err)
	}

	// User has confirmed by typing the droplet name - proceed with deletion
	// Extract droplet ID and perform deletion
	state := flow.State()
	dropletResult, _ := state.GetResult("select_droplet")
	dropletID := dropletResult.Value.(int)
	dropletName := dropletResult.Metadata["name"].(string)

	// Delete the droplet
	_, err = client.Droplets.Delete(ctx, dropletID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete droplet: %w", err)
	}

	// Return deleted droplet info
	return &utils.SelectItem{
		Name:  dropletName,
		Value: strconv.Itoa(dropletID),
	}, nil
}

