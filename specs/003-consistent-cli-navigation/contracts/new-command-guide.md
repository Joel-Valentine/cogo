# Guide: Adding New Commands to cogo

**Last Updated**: 2026-01-18

## Overview

This guide walks you through adding a new command to cogo that follows our navigation standards and provides a consistent user experience.

## Prerequisites

Before starting, familiarize yourself with:
1. `specs/003-consistent-cli-navigation/navigation-patterns.md` - Standard patterns
2. `specs/003-consistent-cli-navigation/contracts/examples.md` - Code examples
3. `navigation/` package - Framework documentation
4. Existing flows: `digitalocean/create_flow.go`, `digitalocean/destroy_flow.go`

## Step-by-Step Guide

### 1. Define Your Command

**File**: `cmd/root.go` (or create new file `cmd/mycommand.go`)

```go
var myCommand = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description (< 80 chars)",
    Long: `Detailed description.
    
    Explain what the command does, when to use it, and any notes.
    
    Examples:
      cogo mycommand
      cogo mycommand --flag value`,
    Run: runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCommand)
}
```

### 2. Create Flow File

**File**: `digitalocean/mycommand_flow.go` (or provider-specific directory)

```go
package digitalocean

import (
    "context"
    "fmt"
    
    "github.com/Joel-Valentine/cogo/navigation"
    "github.com/digitalocean/godo"
    "github.com/fatih/color"
)

// MyCommandFlow orchestrates the multi-step process
type MyCommandFlow struct {
    client *godo.Client
    state  navigation.State
    steps  []navigation.Step
}

// NewMyCommandFlow creates a new flow
func NewMyCommandFlow(client *godo.Client) *MyCommandFlow {
    flow := &MyCommandFlow{
        client: client,
        state:  navigation.NewState(),
    }
    
    // Define steps in order
    flow.steps = []navigation.Step{
        &Step1{client: client},
        &Step2{client: client},
        &ConfirmationStep{},
    }
    
    return flow
}

func (f *MyCommandFlow) Name() string {
    return "My Command"
}

func (f *MyCommandFlow) Steps() []navigation.Step {
    return f.steps
}

func (f *MyCommandFlow) State() navigation.State {
    return f.state
}
```

### 3. Implement Each Step

For each step in your flow:

```go
// Step1 asks for some input
type Step1 struct {
    client *godo.Client
}

func (s *Step1) Name() string {
    return "step1_name"  // Unique identifier, snake_case
}

func (s *Step1) Prompt() string {
    return "What is your selection?"
}

func (s *Step1) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // 1. Fetch data (if needed)
    items, err := fetchItems(ctx, s.client)
    if err != nil {
        return navigation.Result{}, fmt.Errorf("failed to fetch items: %w", err)
    }
    
    // 2. Check for empty state (REQUIRED)
    if len(items) == 0 {
        handler := &navigation.EmptyStateHandler{
            ResourceName:     "items",
            Context:          "in your account",
            SuggestedCommand: "cogo create-item",
        }
        handler.Display()
        return navigation.Result{}, navigation.ErrEmptyState
    }
    
    // 3. Convert to string slice for prompt
    itemNames := make([]string, len(items))
    for i, item := range items {
        itemNames[i] = item.Name
    }
    
    // 4. Create prompt
    prompt := navigation.NewSelectPrompt(s.Prompt(), itemNames)
    
    // 5. Add back option if available
    if state.CanGoBack() {
        prompt = prompt.AddBackOption()
    }
    
    // 6. Run prompt
    index, _, err := prompt.RunWithContext(ctx)
    if err != nil {
        return navigation.Result{}, err  // Handles Ctrl+C, Esc, etc.
    }
    
    // 7. Get selected item
    selected := items[index]
    
    // 8. Return result with metadata
    return navigation.NewResultWithMetadata(selected.ID, map[string]interface{}{
        "name":  selected.Name,
        "index": index,
    }), nil
}

func (s *Step1) Validate(result navigation.Result) error {
    // Validation logic (if needed)
    // This is called by Navigator after Execute returns
    return nil
}

func (s *Step1) Default() interface{} {
    // Return default value if any, or nil
    return nil
}
```

### 4. Add Confirmation Step

**Best Practice**: Always show summary before destructive operations

```go
type ConfirmationStep struct{}

func (s *ConfirmationStep) Name() string {
    return "confirm"
}

func (s *ConfirmationStep) Prompt() string {
    return "Proceed with operation?"
}

func (s *ConfirmationStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // Get previous selections
    step1Result, _ := state.GetResult("step1_name")
    step2Result, _ := state.GetResult("step2_name")
    
    // Display summary
    fmt.Println()
    color.Cyan("=== Operation Summary ===")
    fmt.Printf("Selection 1: %s\n", step1Result.Metadata["name"])
    fmt.Printf("Selection 2: %s\n", step2Result.Metadata["name"])
    color.Cyan("========================")
    fmt.Println()
    
    // For destructive operations, default to No (false)
    // For creation operations, default to Yes (true)
    prompt := navigation.NewConfirmPrompt(s.Prompt(), false)
    confirmed, err := prompt.RunWithContext(ctx)
    if err != nil {
        return navigation.Result{}, err
    }
    
    return navigation.NewResult(confirmed), nil
}

func (s *ConfirmationStep) Validate(result navigation.Result) error {
    return nil
}

func (s *ConfirmationStep) Default() interface{} {
    return false  // Safe default
}
```

### 5. Create Execute Function

**File**: `digitalocean/mycommand_flow.go`

```go
// ExecuteMyCommandFlow runs the entire flow
func ExecuteMyCommandFlow(client *godo.Client) (*Result, error) {
    ctx := context.Background()
    navigator := navigation.NewNavigator()
    flow := NewMyCommandFlow(client)
    
    // Run the flow
    result, err := navigator.Run(ctx, flow)
    
    // Handle special cases (REQUIRED PATTERN)
    if err == navigation.ErrEmptyState {
        // Message already displayed in step
        return nil, nil
    }
    
    if err == navigation.ErrCancel {
        color.Cyan("\nOperation canceled.")
        return nil, nil
    }
    
    if err == context.Canceled {
        fmt.Println()  // Clean newline after ^C
        return nil, nil
    }
    
    if err != nil {
        return nil, fmt.Errorf("command flow failed: %w", err)
    }
    
    // Check if user confirmed
    confirmed := result.Value.(bool)
    if !confirmed {
        color.Cyan("Operation canceled.")
        return nil, nil
    }
    
    // Extract selections from state
    state := flow.State()
    selection1, _ := state.GetResult("step1_name")
    selection2, _ := state.GetResult("step2_name")
    
    // Perform actual operation
    result, err := performOperation(ctx, client, selection1, selection2)
    if err != nil {
        return nil, fmt.Errorf("operation failed: %w", err)
    }
    
    return result, nil
}
```

### 6. Wire Up in Command

**File**: `cmd/root.go`

```go
func runMyCommand(cmd *cobra.Command, args []string) {
    // 1. Get credentials
    ctx := context.Background()
    credManager := credentials.NewManager()
    token, _, err := credManager.GetToken(ctx)
    if err != nil {
        color.Red("✗ Error: Unable to get API token\n")
        fmt.Println()
        fmt.Println("Run 'cogo config set-token' to configure your token.")
        return
    }
    
    // 2. Create client
    client := godo.NewFromToken(token)
    
    // 3. Execute flow
    result, err := do.ExecuteMyCommandFlow(client)
    
    // 4. Handle errors
    if err != nil {
        color.Red("✗ Error: %v\n", err)
        return
    }
    
    if result == nil {
        // Canceled or empty state (message already shown)
        return
    }
    
    // 5. Show success
    color.Green("✓ Operation completed successfully!")
    fmt.Printf("Result: %v\n", result)
}
```

### 7. Add Tests

**File**: `digitalocean/mycommand_flow_test.go`

```go
package digitalocean

import (
    "context"
    "testing"
    
    "github.com/Joel-Valentine/cogo/navigation"
    "github.com/stretchr/testify/assert"
)

func TestMyCommandFlow_EmptyState(t *testing.T) {
    // Test empty state handling
    client := &mockClient{items: []Item{}}
    flow := NewMyCommandFlow(client)
    
    navigator := navigation.NewNavigator()
    _, err := navigator.Run(context.Background(), flow)
    
    assert.ErrorIs(t, err, navigation.ErrEmptyState)
}

func TestMyCommandFlow_BackNavigation(t *testing.T) {
    // Test back navigation
    // ...
}

func TestMyCommandFlow_Cancellation(t *testing.T) {
    // Test Ctrl+C and Esc
    // ...
}

func TestMyCommandFlow_FullFlow(t *testing.T) {
    // Test complete flow
    // ...
}
```

## Common Patterns

### Input Prompts (Text Entry)

```go
prompt := navigation.NewInputPrompt("Enter name", "default-name")
name, err := prompt.RunWithContext(ctx)
if err != nil {
    return navigation.Result{}, err
}

// Validate after Enter
if err := ValidateName(name); err != nil {
    color.Red("✗ %v", err)
    return navigation.Result{}, err
}
```

### Selection Prompts (List)

```go
items := []string{"Option 1", "Option 2", "Option 3"}
prompt := navigation.NewSelectPrompt("Choose option", items)

if state.CanGoBack() {
    prompt = prompt.AddBackOption()
}

index, selected, err := prompt.RunWithContext(ctx)
if err != nil {
    return navigation.Result{}, err
}
```

### Confirmation Prompts (Yes/No)

```go
prompt := navigation.NewConfirmPrompt("Are you sure?", false)
confirmed, err := prompt.RunWithContext(ctx)
if err != nil {
    return navigation.Result{}, err
}

if !confirmed {
    return navigation.Result{}, navigation.ErrCancel
}
```

### Accessing Previous Results

```go
// In a later step, access previous selections
previousResult, found := state.GetResult("previous_step_name")
if !found {
    return navigation.Result{}, fmt.Errorf("previous step result not found")
}

// Access the value
value := previousResult.Value.(string)

// Access metadata
name := previousResult.Metadata["name"].(string)
index := previousResult.Metadata["index"].(int)
```

## Checklist

Before submitting your PR:

- [ ] Command defined in `cmd/root.go` with proper help text
- [ ] Flow file created in appropriate directory
- [ ] All steps implement `Step` interface
- [ ] Empty state checking in ALL steps
- [ ] Back navigation support ("← Back" option)
- [ ] Confirmation step with summary
- [ ] Execute function with proper error handling
- [ ] Success/error messages follow standards
- [ ] Tests written for all scenarios
- [ ] Manual testing performed:
  - [ ] Empty state displays message
  - [ ] Can go back at each step
  - [ ] Ctrl+C cancels cleanly
  - [ ] Esc/q quits flow
  - [ ] Invalid input handled gracefully
  - [ ] Success message on completion
- [ ] Code follows `navigation-patterns.md` standards
- [ ] Documentation updated if needed

## Getting Help

- Check existing flows: `digitalocean/create_flow.go`, `digitalocean/destroy_flow.go`
- See examples: `specs/003-consistent-cli-navigation/contracts/examples.md`
- Read standards: `specs/003-consistent-cli-navigation/navigation-patterns.md`
- Ask in: Project discussions/Slack/Discord

## Examples

For complete working examples, see:
- `digitalocean/create_flow.go` - 7-step creation flow
- `digitalocean/destroy_flow.go` - 4-step destruction flow with safety checks
- `cmd/root.go` - Command integration examples

