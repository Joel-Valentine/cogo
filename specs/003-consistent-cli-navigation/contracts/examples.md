# Navigation Framework - Developer Usage Examples

This document shows how to use the navigation framework interfaces to build consistent, user-friendly CLI interactions.

## Example 1: Simple List Selection (No Back Navigation)

**Use Case**: Select a region for a droplet.

```go
package main

import (
    "context"
    "fmt"
    "github.com/Joel-Valentine/cogo/navigation"
)

// SelectRegionStep implements the Step interface
type SelectRegionStep struct {
    regions []string
}

func (s *SelectRegionStep) Name() string {
    return "select_region"
}

func (s *SelectRegionStep) Prompt() string {
    return "Select region:"
}

func (s *SelectRegionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // Check for empty state (research finding: exit 0, not error)
    if len(s.regions) == 0 {
        return navigation.Result{}, navigation.ErrEmptyState
    }

    // Use promptui wrapper with navigation support
    prompt := navigation.NewSelectPrompt(s.Prompt(), s.regions)
    index, selected, err := prompt.Run()
    
    if err != nil {
        return navigation.Result{}, err // context.Canceled, ErrCancel, or ErrGoBack
    }

    return navigation.Result{
        Value: selected,
        Metadata: map[string]interface{}{
            "index": index,
        },
    }, nil
}

func (s *SelectRegionStep) Validate(result navigation.Result) error {
    // Validation already done in Execute via prompt selection
    return nil
}

func (s *SelectRegionStep) Default() interface{} {
    if len(s.regions) > 0 {
        return s.regions[0] // First region as default
    }
    return nil
}

// Usage in command
func createDroplet(cmd *cobra.Command, args []string) error {
    ctx := context.Background()
    navigator := navigation.NewNavigator()

    // Create step
    step := &SelectRegionStep{
        regions: []string{"nyc1", "nyc3", "sfo3"},
    }

    // Run single step (simple case)
    result, err := navigator.RunStep(ctx, step)
    if err == navigation.ErrEmptyState {
        fmt.Println("No regions available.")
        return nil // Exit 0
    }
    if err != nil {
        return err
    }

    region := result.Value.(string)
    fmt.Printf("Selected region: %s\n", region)
    return nil
}
```

## Example 2: Multi-Step Flow with Back Navigation

**Use Case**: Complete droplet creation flow with region, size, image, and confirmation.

```go
package main

import (
    "context"
    "fmt"
    "github.com/Joel-Valentine/cogo/navigation"
)

// CreateDropletFlow implements the Flow interface
type CreateDropletFlow struct {
    state navigation.State
    steps []navigation.Step
}

func NewCreateDropletFlow(api *digitalocean.Client) *CreateDropletFlow {
    flow := &CreateDropletFlow{
        state: navigation.NewState(),
    }
    
    // Define steps in order
    flow.steps = []navigation.Step{
        &SelectRegionStep{api: api},
        &SelectSizeStep{api: api},
        &SelectImageStep{api: api},
        &ConfirmStep{},
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

// ConfirmStep shows summary and asks for confirmation
type ConfirmStep struct{}

func (s *ConfirmStep) Name() string {
    return "confirm"
}

func (s *ConfirmStep) Prompt() string {
    return "Confirm creation?"
}

func (s *ConfirmStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // Get previous selections from state
    region, _ := state.GetResult("select_region")
    size, _ := state.GetResult("select_size")
    image, _ := state.GetResult("select_image")

    // Show summary
    fmt.Println("\nDroplet Configuration:")
    fmt.Printf("  Region: %v\n", region.Value)
    fmt.Printf("  Size:   %v\n", size.Value)
    fmt.Printf("  Image:  %v\n", image.Value)
    fmt.Println()

    // Confirm with Y/n prompt (research pattern: show default)
    prompt := navigation.NewConfirmPrompt(s.Prompt(), true) // true = default Yes
    confirmed, err := prompt.Run()
    
    if err != nil {
        return navigation.Result{}, err
    }

    return navigation.Result{Value: confirmed}, nil
}

func (s *ConfirmStep) Validate(result navigation.Result) error {
    return nil
}

func (s *ConfirmStep) Default() interface{} {
    return true // Default to "Yes"
}

// Usage in command
func createDroplet(cmd *cobra.Command, args []string) error {
    ctx := context.Background()
    navigator := navigation.NewNavigator()
    
    api := getDigitalOceanClient()
    flow := NewCreateDropletFlow(api)

    // Run multi-step flow with automatic back navigation support
    result, err := navigator.Run(ctx, flow)
    
    // Handle special cases (research-based patterns)
    if err == navigation.ErrEmptyState {
        fmt.Println("✗ Error: No regions available")
        fmt.Println("\nCheck your DigitalOcean account status.")
        return nil // Exit 0, not error
    }
    
    if err == navigation.ErrCancel {
        fmt.Println("\nDroplet creation canceled.")
        return nil // Exit 0, user choice
    }
    
    if err == context.Canceled {
        fmt.Println() // Newline after ^C
        return nil // Exit 130 (handled by Cobra)
    }
    
    if err != nil {
        // Fatal error - show detailed message (research pattern)
        fmt.Printf("✗ Error: %v\n", err)
        return err // Exit 1
    }

    // Success - extract final confirmation
    if confirmed, ok := result.Value.(bool); ok && confirmed {
        fmt.Println("✓ Creating droplet...")
        // Actually create droplet using selections from flow.State()
    }

    return nil
}
```

## Example 3: Empty State Handling

**Use Case**: Destroy droplet when none exist (the user's reported bug).

```go
func destroyDroplet(cmd *cobra.Command, args []string) error {
    ctx := context.Background()
    api := getDigitalOceanClient()

    // Fetch droplets
    droplets, err := api.ListDroplets(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch droplets: %w", err)
    }

    // Research pattern: Check for empty BEFORE prompting (prevent crash)
    if len(droplets) == 0 {
        // Follow research recommendations: friendly message + suggestion
        fmt.Println("No droplets found in your DigitalOcean account.")
        fmt.Println()
        fmt.Println("Run 'cogo create' to create a droplet.")
        return nil // Exit 0 (not an error)
    }

    // Normal flow: Create and run step
    navigator := navigation.NewNavigator()
    step := &SelectDropletStep{droplets: droplets, action: "destroy"}
    
    result, err := navigator.RunStep(ctx, step)
    if err != nil {
        return err
    }

    // ... continue with destroy confirmation
}
```

## Example 4: Input Validation (On-Enter, Not Per-Keystroke)

**Use Case**: Enter droplet name with validation.

```go
type NameInputStep struct{}

func (s *NameInputStep) Name() string {
    return "enter_name"
}

func (s *NameInputStep) Prompt() string {
    return "Droplet name:"
}

func (s *NameInputStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // Create input prompt with validation
    prompt := navigation.NewInputPrompt(s.Prompt(), s.Default().(string))
    
    // Validation happens on Enter, not per-keystroke (research pattern)
    for {
        input, err := prompt.Run()
        if err != nil {
            return navigation.Result{}, err // Cancel or Ctrl+C
        }

        // Validate AFTER user presses Enter
        result := navigation.Result{Value: input}
        if valErr := s.Validate(result); valErr != nil {
            // Show error and re-prompt (research pattern: don't crash)
            fmt.Printf("✗ %v\n", valErr)
            fmt.Println()
            continue // Re-prompt
        }

        return result, nil
    }
}

func (s *NameInputStep) Validate(result navigation.Result) error {
    name := result.Value.(string)
    
    // Validation rules (research: clear, actionable messages)
    if len(name) == 0 {
        return fmt.Errorf("name cannot be empty")
    }
    if len(name) > 63 {
        return fmt.Errorf("name must be 63 characters or less")
    }
    if !isValidDropletName(name) {
        return fmt.Errorf("name can only contain letters, numbers, and hyphens")
    }
    
    return nil
}

func (s *NameInputStep) Default() interface{} {
    return fmt.Sprintf("my-droplet-%d", time.Now().Unix())
}
```

## Example 5: Error Handling (Research-Based Patterns)

**Use Case**: Handle different error types consistently.

```go
func runCommand(cmd *cobra.Command, args []string) error {
    result, err := runSomeOperation()

    // Pattern 1: Empty State (exit 0)
    if err == navigation.ErrEmptyState {
        fmt.Println("No resources found.")
        fmt.Println()
        fmt.Println("Run 'cogo create' to get started.")
        return nil
    }

    // Pattern 2: User Cancellation (exit 0)
    if err == navigation.ErrCancel {
        fmt.Println("\nOperation canceled.")
        return nil
    }

    // Pattern 3: Ctrl+C (exit 130 via Cobra)
    if err == context.Canceled {
        fmt.Println() // Clean newline
        return nil
    }

    // Pattern 4: Validation Error (shouldn't reach here, but handle gracefully)
    if err == navigation.ErrValidationFailed {
        fmt.Printf("✗ Invalid input: %v\n", err)
        return nil
    }

    // Pattern 5: API Error (detailed message from research)
    if apiErr, ok := err.(*digitalocean.APIError); ok {
        fmt.Printf("✗ Error: %s\n", apiErr.Message)
        fmt.Println()
        fmt.Printf("API returned: %d %s\n", apiErr.StatusCode, apiErr.Status)
        fmt.Println()
        fmt.Println("Check your API token: cogo config status")
        return err // Exit 1
    }

    // Pattern 6: Generic Error (still helpful)
    if err != nil {
        fmt.Printf("✗ Error: %v\n", err)
        fmt.Println()
        fmt.Println("Run 'cogo <command> --help' for more information.")
        return err // Exit 1
    }

    // Success
    fmt.Println("✓ Operation completed successfully")
    return nil
}
```

## Navigation Keys (Research Findings Applied)

All prompts in the framework support these keys:

- **↑/↓** or **j/k**: Navigate list items
- **Enter**: Select/confirm
- **Ctrl+C**: Immediate cancel (exit 130)
- **Esc** or **q**: Quit flow (exit 0)
- **b** or **←**: Go back to previous step (gcloud pattern)
- **Invalid keys**: Silently ignored (no error spam)

These are shown in help text on every interactive prompt:

```
? Select region: (Use arrow keys, 'b' for back, 'q' to quit)
  › nyc1 - New York 1
    nyc3 - New York 3
    sfo3 - San Francisco 3
```

## Testing Navigation Flows

```go
func TestCreateDropletFlow(t *testing.T) {
    // Create mock API
    api := &mockDigitalOceanAPI{
        regions: []string{"nyc1", "nyc3"},
        sizes:   []string{"s-1vcpu-1gb", "s-2vcpu-2gb"},
        images:  []string{"ubuntu-20-04-x64"},
    }

    // Create flow
    flow := NewCreateDropletFlow(api)

    // Simulate user selections
    state := navigation.NewState()
    state.AddResult(flow.Steps()[0], navigation.Result{Value: "nyc3"})
    state.AddResult(flow.Steps()[1], navigation.Result{Value: "s-1vcpu-1gb"})
    state.AddResult(flow.Steps()[2], navigation.Result{Value: "ubuntu-20-04-x64"})

    // Test back navigation
    assert.True(t, state.CanGoBack())
    index, err := state.Back()
    assert.NoError(t, err)
    assert.Equal(t, 2, index) // Back to image selection

    // Test history
    history := state.History()
    assert.Len(t, history, 3)
    assert.Equal(t, "nyc3", history[0].Result.Value)
}
```

---

**Key Takeaways**:

1. ✅ **Check empty states before prompting** (prevents crashes)
2. ✅ **Validate on Enter, not per-keystroke** (prevents spam)
3. ✅ **Support back navigation** (only modern CLI pattern)
4. ✅ **Handle errors gracefully** (never crash, always helpful)
5. ✅ **Show defaults and help text** (universal pattern)
6. ✅ **Exit 0 for user actions** (cancel, empty, etc.)
7. ✅ **Exit 1 only for real errors** (API failures, validation, etc.)

