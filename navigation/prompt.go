package navigation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// SelectPrompt wraps promptui.Select with navigation support.
// Handles:
// - Back navigation via "← Back" option selection
// - Ctrl+C (context cancellation)
// - Arrow key navigation
// - Invalid key silently ignored (no spam)
type SelectPrompt struct {
	Label       string
	Items       []string
	HideHelp    bool
	Searcher    func(input string, index int) bool
	cursorPos   int
	defaultHelp string
}

// NewSelectPrompt creates a new selection prompt with navigation support.
func NewSelectPrompt(label string, items []string) *SelectPrompt {
	return &SelectPrompt{
		Label:       label,
		Items:       items,
		HideHelp:    false,
		defaultHelp: "Use arrow keys to navigate, Enter to select, Ctrl+C to cancel",
	}
}

// WithSearcher adds a custom search function to filter items.
func (p *SelectPrompt) WithSearcher(searcher func(input string, index int) bool) *SelectPrompt {
	p.Searcher = searcher
	return p
}

// WithoutHelp disables the help text.
func (p *SelectPrompt) WithoutHelp() *SelectPrompt {
	p.HideHelp = true
	return p
}

// Run executes the selection prompt.
// Returns (index, selected_item, error).
//
// Possible errors:
// - ErrGoBack: User pressed 'b' or ←
// - ErrCancel: User pressed 'q' or Esc
// - context.Canceled: User pressed Ctrl+C
func (p *SelectPrompt) Run() (int, string, error) {
	return p.RunWithContext(context.Background())
}

// RunWithContext executes the selection prompt with a context for cancellation.
func (p *SelectPrompt) RunWithContext(ctx context.Context) (int, string, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return -1, "", ctx.Err()
	default:
	}

	// Create custom templates to show our help text
	// Note: We use templates instead of adding help to the label to avoid duplication
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✔ {{ . | green }}",
	}
	
	// Only add help if not hidden
	if !p.HideHelp {
		templates.Help = "↑↓: Navigate | Enter: Select | Ctrl+C: Cancel"
	}
	
	// Create promptui select
	prompt := promptui.Select{
		Label:     p.Label,
		Items:     p.Items,
		Templates: templates,
		CursorPos: p.cursorPos,
		Size:      10, // Limit visible items to reduce flickering
	}

	if p.Searcher != nil {
		prompt.Searcher = p.Searcher
	}

	// Run prompt
	index, result, err := prompt.Run()

	// Handle errors
	if err != nil {
		// Ctrl+C or interrupt
		if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(err.Error(), "interrupt") {
			return -1, "", context.Canceled
		}

		// Any other error
		return -1, "", err
	}

	// Check if user selected a special navigation option
	// (if we added "← Back" or "Quit" to the items list)
	resultLower := strings.ToLower(result)
	if strings.Contains(resultLower, "back") && (strings.HasPrefix(resultLower, "←") || strings.HasPrefix(resultLower, "<-")) {
		return -1, "", ErrGoBack
	}
	if resultLower == "quit" || resultLower == "q" {
		return -1, "", ErrCancel
	}

	return index, result, nil
}

// AddBackOption adds a "← Back" option at the beginning of the items list.
// Returns the modified prompt for chaining.
func (p *SelectPrompt) AddBackOption() *SelectPrompt {
	p.Items = append([]string{"← Back"}, p.Items...)
	return p
}

// AddQuitOption adds a "Quit" option at the end of the items list.
// Returns the modified prompt for chaining.
func (p *SelectPrompt) AddQuitOption() *SelectPrompt {
	p.Items = append(p.Items, "Quit")
	return p
}

// InputPrompt wraps promptui.Prompt for text input with validation.
type InputPrompt struct {
	Label     string
	Default   string
	Validator Validator
	Mask      rune // For password input (e.g., '*')
	HideHelp  bool
}

// NewInputPrompt creates a new input prompt.
func NewInputPrompt(label, defaultValue string) *InputPrompt {
	return &InputPrompt{
		Label:   label,
		Default: defaultValue,
	}
}

// WithValidator adds a validator to the prompt.
func (p *InputPrompt) WithValidator(validator Validator) *InputPrompt {
	p.Validator = validator
	return p
}

// WithMask sets a mask character for password input.
func (p *InputPrompt) WithMask(mask rune) *InputPrompt {
	p.Mask = mask
	return p
}

// Run executes the input prompt with validation.
// Validation happens on Enter (not per-keystroke) per research findings.
// Re-prompts on validation error until valid input or cancellation.
func (p *InputPrompt) Run() (string, error) {
	return p.RunWithContext(context.Background())
}

// RunWithContext executes the input prompt with a context for cancellation.
func (p *InputPrompt) RunWithContext(ctx context.Context) (string, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Show help text once (not in loop to avoid spam)
	fmt.Println()
	if p.Default != "" {
		color.Cyan("(Press Enter to use default, or type to override | Ctrl+C: quit)")
	} else {
		color.Cyan("(Type your text, press Enter to confirm | Ctrl+C: quit)")
	}
	fmt.Println()
	
	for {
		// Create promptui prompt
		prompt := promptui.Prompt{
			Label:   p.Label,
			Default: p.Default,
		}

		if p.Mask != 0 {
			prompt.Mask = p.Mask
		}

		// No per-keystroke validation (research finding: causes spam)
		// We'll validate after user presses Enter

		// Run prompt
		result, err := prompt.Run()

		// Handle errors
		if err != nil {
			errStr := err.Error()

			// Ctrl+C or interrupt
			if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(errStr, "interrupt") {
				return "", context.Canceled
			}

			// EOF (Ctrl+D)
			if errors.Is(err, promptui.ErrEOF) || strings.Contains(errStr, "EOF") {
				return "", ErrCancel
			}

			return "", err
		}

		// Validate on Enter (research pattern: 9/10 tools)
		if p.Validator != nil {
			if err := p.Validator(result); err != nil {
				// Show validation error and re-prompt
				fmt.Printf("✗ %v\n\n", err)
				continue
			}
		}

		return result, nil
	}
}

// ConfirmPrompt wraps promptui.Prompt for yes/no confirmation.
type ConfirmPrompt struct {
	Label      string
	Default    bool // true = Yes, false = No
	HideHelp   bool
}

// NewConfirmPrompt creates a new confirmation prompt.
func NewConfirmPrompt(label string, defaultYes bool) *ConfirmPrompt {
	return &ConfirmPrompt{
		Label:   label,
		Default: defaultYes,
	}
}

// Run executes the confirmation prompt.
// Accepts: y, yes, n, no (case-insensitive)
// Returns (confirmed bool, error)
func (p *ConfirmPrompt) Run() (bool, error) {
	return p.RunWithContext(context.Background())
}

// RunWithContext executes the confirmation prompt with a context.
func (p *ConfirmPrompt) RunWithContext(ctx context.Context) (bool, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Use Select instead of Prompt to avoid checkmark spam
	// This is cleaner and matches common CLI patterns (kubectl, gh, etc.)
	items := []string{"No", "Yes"}
	defaultIndex := 0
	if p.Default {
		defaultIndex = 1
	}

	prompt := promptui.Select{
		Label:     p.Label,
		Items:     items,
		CursorPos: defaultIndex,
		Size:      2,
		HideHelp:  p.HideHelp,
		Templates: &promptui.SelectTemplates{
			Help: "{{ \"↑↓: Navigate | Enter: Confirm | Ctrl+C: Cancel\" | faint }}",
		},
	}

	_, result, err := prompt.Run()

	// Handle errors
	if err != nil {
		errStr := err.Error()

		// Ctrl+C
		if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(errStr, "interrupt") {
			return false, context.Canceled
		}

		// EOF or Esc
		if errors.Is(err, promptui.ErrEOF) || strings.Contains(errStr, "EOF") || strings.Contains(errStr, "esc") {
			return false, ErrCancel
		}

		return false, err
	}

	return result == "Yes", nil
}

