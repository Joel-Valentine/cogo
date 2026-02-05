# Standard Navigation Patterns for cogo

**Created**: 2026-01-18  
**Status**: Active Standard

## Overview

This document defines the standard navigation patterns that ALL cogo commands must follow, ensuring a consistent user experience across all cloud providers and operations.

## Universal Keyboard Shortcuts

**ALL commands MUST support**:

| Key | Action | Context |
|-----|--------|---------|
| **Ctrl+C** | Immediate cancel | Any prompt, any time |
| **Esc** | Quit current flow | Multi-step operations |
| **q** | Quit current flow | Selection lists (alternative to Esc) |
| **b** or **←** | Go back | Multi-step operations (when "← Back" shown) |
| **↑** / **↓** | Navigate items | Selection lists |
| **Enter** | Confirm/Continue | All prompts |

## Standard Message Formats

### Empty State Messages

**Format**:
```
No {resources} found {context}.

{Actionable suggestion}
```

**Examples**:
```bash
No droplets found in your DigitalOcean account.

Run 'cogo create' to create a droplet.
```

```bash
No SSH keys found in your DigitalOcean account.

Add an SSH key at: https://cloud.digitalocean.com/account/security
```

**Rules**:
- Always exit with code 0 (not an error condition)
- Always suggest next action
- Keep message concise and helpful
- Use consistent formatting

### Error Messages

**Format**:
```
✗ Error: {Summary}

{Details}

{Actionable suggestion}
```

**Examples**:
```bash
✗ Error: Failed to create droplet

The DigitalOcean API returned: 422 Unprocessable Entity
You cannot create droplets in the nyc1 region with this account.

Try a different region: cogo create --region nyc3
```

```bash
✗ Error: Unable to get DigitalOcean API token

No token found in environment, keychain, or config file.

Run 'cogo config set-token' to configure your token.
```

**Rules**:
- Always use ✗ symbol at start
- Red color for error text
- Exit with code 1
- Include actionable fix if possible
- Keep summary on first line

### Success Messages

**Format**:
```
✓ {Action completed successfully}
```

**Examples**:
```bash
✓ Droplet [my-server] was created!
```

```bash
✓ Droplet [old-server] has been destroyed
```

**Rules**:
- Always use ✓ symbol
- Green color for success text
- Keep concise
- Include resource name in brackets

### Warning Messages

**Format**:
```
⚠️  WARNING: {Warning message}
```

**Examples**:
```bash
⚠️  WARNING: You are about to delete droplet: my-production-server
```

**Rules**:
- Always use ⚠️ symbol
- Yellow color for warning text
- Use for destructive operations
- Be clear about consequences

## Multi-Step Flow Pattern

### Flow Structure

All multi-step operations MUST:
1. Use the navigation framework (`navigation.Navigator`)
2. Implement `Flow` and `Step` interfaces
3. Support back navigation with "← Back" option
4. Show summary before final confirmation
5. Handle empty states gracefully

### Step Pattern

Each step MUST:
```go
type MyStep struct {
    // Dependencies
}

func (s *MyStep) Name() string {
    return "step_identifier"  // Snake case, unique
}

func (s *MyStep) Prompt() string {
    return "User-facing prompt"  // Clear, concise
}

func (s *MyStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // 1. Check for empty state
    if len(items) == 0 {
        DisplayEmptyMessage()
        return navigation.Result{}, navigation.ErrEmptyState
    }
    
    // 2. Create prompt with back option
    prompt := navigation.NewSelectPrompt(s.Prompt(), items)
    if state.CanGoBack() {
        prompt = prompt.AddBackOption()
    }
    
    // 3. Run prompt
    index, selected, err := prompt.RunWithContext(ctx)
    if err != nil {
        return navigation.Result{}, err
    }
    
    // 4. Return result with metadata
    return navigation.NewResultWithMetadata(value, map[string]interface{}{
        "name": selected,
        "index": index,
    }), nil
}

func (s *MyStep) Validate(result navigation.Result) error {
    // Validation logic (called after Execute)
    return nil
}

func (s *MyStep) Default() interface{} {
    return nil  // Or default value
}
```

### Confirmation Pattern

Final confirmation steps MUST:
1. Show complete summary of selections
2. Default to "No" for destructive operations
3. Default to "Yes" for creation operations
4. Use colored output for summary

**Example**:
```go
// Display summary
fmt.Println()
color.Cyan("=== Droplet Configuration ===")
fmt.Printf("Name:     %s\n", name.Value)
fmt.Printf("Image:    %s\n", image.Metadata["name"])
fmt.Printf("Size:     %s\n", size.Metadata["name"])
fmt.Printf("Region:   %s\n", region.Metadata["name"])
color.Cyan("=============================")
fmt.Println()

// Confirm
prompt := navigation.NewConfirmPrompt("Create this droplet?", true)
confirmed, err := prompt.RunWithContext(ctx)
```

## Help Text Standards

### Prompt Help Text

**Format**: `(Use arrow keys, 'b' for back, 'q' to quit)`

**When to show**:
- All selection prompts in multi-step flows
- When back navigation is available

**When NOT to show**:
- Simple yes/no confirmations
- Single-step operations
- Use `HideHelp()` method to suppress

### Command Help Text

All commands MUST have:
```go
var myCommand = &cobra.Command{
    Use:   "command-name",
    Short: "One-line description",  // < 80 chars
    Long:  `Detailed description.
    
    Explains what the command does, when to use it, and any important notes.
    
    Examples:
      cogo command-name
      cogo command-name --flag value`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

## Validation Standards

### Input Validation

**Rules**:
1. ✅ **Validate AFTER Enter** (not per-keystroke)
2. ✅ Show clear error message with ✗ symbol
3. ✅ Re-prompt automatically (don't exit)
4. ✅ Use `navigation.NewValidationError()` for consistency

**Example**:
```go
for {
    input, err := prompt.Run()
    if err != nil {
        return navigation.Result{}, err
    }
    
    // Validate AFTER Enter
    if err := ValidateName(input); err != nil {
        color.Red("✗ %v", err)
        fmt.Println()
        continue  // Re-prompt
    }
    
    return navigation.NewResult(input), nil
}
```

### Selection Validation

**Rules**:
1. ✅ Check bounds before accessing arrays
2. ✅ Handle empty lists before prompting
3. ✅ Validate selected value exists

## Error Handling Standards

### Error Return Pattern

```go
// Empty state (not fatal)
if len(items) == 0 {
    DisplayEmptyMessage()
    return navigation.Result{}, navigation.ErrEmptyState
}

// User canceled
if err == navigation.ErrCancel {
    color.Cyan("\nOperation canceled.")
    return nil
}

// Ctrl+C
if err == context.Canceled {
    fmt.Println()  // Clean newline
    return nil
}

// Fatal error
if err != nil {
    color.Red("✗ Error: %v", err)
    return err
}
```

### Exit Codes

| Code | Meaning | When to Use |
|------|---------|-------------|
| **0** | Success | Normal completion, empty states, user cancellation |
| **1** | Error | API failures, validation errors, configuration issues |
| **130** | Ctrl+C | User interrupt (handled by Cobra automatically) |

## Color Standards

| Color | Use For | Package |
|-------|---------|---------|
| **Green** | Success messages (✓) | `color.Green()` |
| **Red** | Error messages (✗) | `color.Red()` |
| **Yellow** | Warnings (⚠️) | `color.Yellow()` |
| **Cyan** | Informational, summaries | `color.Cyan()` |

## State Management

### State Preservation

When going back:
1. ✅ Previous selections MUST be preserved
2. ✅ Show previous value as default
3. ✅ Allow user to change selection
4. ✅ Truncate future history (git rebase style)

### State Access

```go
// Get previous result
if existing, found := state.GetResult("step_name"); found {
    // Use existing.Value
}

// Add new result
state.AddResult("step_name", result)

// Check if can go back
if state.CanGoBack() {
    // Show "← Back" option
}
```

## Command Structure Standards

### File Organization

```
cmd/
  root.go          # Command definitions
digitalocean/
  create_flow.go   # Create flow implementation
  destroy_flow.go  # Destroy flow implementation
  digitalocean.go  # Legacy code (to be refactored)
```

### Command Implementation

```go
var myCommand = &cobra.Command{
    Use:   "command",
    Short: "Short description",
    Long:  `Long description`,
    Run: func(cmd *cobra.Command, args []string) {
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
        result, err := do.ExecuteMyFlow(client)
        
        // 4. Handle errors (following error handling standards)
        if err != nil {
            color.Red("✗ Error: %v\n", err)
            return
        }
        
        if result == nil {
            return  // Canceled or empty (message already shown)
        }
        
        // 5. Show success
        color.Green("✓ Operation completed!")
    },
}
```

## Testing Standards

### Flow Testing

Each flow MUST have tests for:
1. Empty state handling
2. Back navigation
3. Cancellation (Ctrl+C, Esc)
4. Validation
5. Full flow completion

### Manual Testing Checklist

For each command:
- [ ] Empty state displays message and exits 0
- [ ] Can go back at each step
- [ ] Ctrl+C cancels cleanly
- [ ] Esc/q quits flow
- [ ] Invalid input re-prompts (doesn't crash)
- [ ] Success message shows on completion
- [ ] Error messages are clear and actionable

## Compliance Checklist

Before merging, ensure:
- [ ] Uses navigation framework
- [ ] Implements Step interface properly
- [ ] Handles empty states
- [ ] Supports back navigation
- [ ] Uses standard message formats
- [ ] Follows color standards
- [ ] Has proper help text
- [ ] Validates after Enter only
- [ ] Handles Ctrl+C gracefully
- [ ] Returns correct exit codes
- [ ] Has tests

## Examples

See `specs/003-consistent-cli-navigation/contracts/examples.md` for complete code examples.

---

**Questions?** See the navigation framework documentation or check existing flows (`create_flow.go`, `destroy_flow.go`) for reference implementations.

