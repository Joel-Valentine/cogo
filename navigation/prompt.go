package navigation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// bellSkipper filters terminal bell characters and cursor hide/show codes
// to prevent macOS error sounds and window flashing during navigation.
type bellSkipper struct{}

func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7
	
	filtered := make([]byte, 0, len(b))
	for _, ch := range b {
		if ch != charBell {
			filtered = append(filtered, ch)
		}
	}
	
	s := string(filtered)
	
	if strings.Contains(s, "\033[?25l") || strings.Contains(s, "\033[?25h") {
		return len(b), nil
	}
	
	if len(filtered) == 0 {
		return len(b), nil
	}
	
	return os.Stderr.Write(filtered)
}

func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}

type SelectPrompt struct {
	Label       string
	Items       []string
	HideHelp    bool
	Searcher    func(input string, index int) bool
	cursorPos   int
	defaultHelp string
}

func NewSelectPrompt(label string, items []string) *SelectPrompt {
	return &SelectPrompt{
		Label:       label,
		Items:       items,
		HideHelp:    false,
		defaultHelp: "Use arrow keys to navigate, Enter to select, Ctrl+C to cancel",
	}
}

func (p *SelectPrompt) WithSearcher(searcher func(input string, index int) bool) *SelectPrompt {
	p.Searcher = searcher
	return p
}

func (p *SelectPrompt) WithoutHelp() *SelectPrompt {
	p.HideHelp = true
	return p
}

func (p *SelectPrompt) Run() (int, string, error) {
	return p.RunWithContext(context.Background())
}

func (p *SelectPrompt) RunWithContext(ctx context.Context) (int, string, error) {
	select {
	case <-ctx.Done():
		return -1, "", ctx.Err()
	default:
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✔ {{ . | green }}",
	}
	
	if !p.HideHelp {
		templates.Help = "↑↓: Navigate | Enter: Select | Ctrl+C: Cancel"
	}
	
	prompt := promptui.Select{
		Label:     p.Label,
		Items:     p.Items,
		Templates: templates,
		CursorPos: p.cursorPos,
		Size:      10,
		Stdout:    &bellSkipper{},
		Stdin:     os.Stdin,
	}

	if p.Searcher != nil {
		prompt.Searcher = p.Searcher
	}

	index, result, err := prompt.Run()

	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(err.Error(), "interrupt") {
			return -1, "", context.Canceled
		}
		return -1, "", err
	}

	resultLower := strings.ToLower(result)
	if strings.Contains(resultLower, "back") && (strings.HasPrefix(resultLower, "←") || strings.HasPrefix(resultLower, "<-")) {
		return -1, "", ErrGoBack
	}
	if resultLower == "quit" || resultLower == "q" {
		return -1, "", ErrCancel
	}

	return index, result, nil
}

func (p *SelectPrompt) AddBackOption() *SelectPrompt {
	p.Items = append([]string{"← Back"}, p.Items...)
	return p
}

func (p *SelectPrompt) AddQuitOption() *SelectPrompt {
	p.Items = append(p.Items, "Quit")
	return p
}

type InputPrompt struct {
	Label     string
	Default   string
	Validator Validator
	Mask      rune
	HideHelp  bool
}

func NewInputPrompt(label, defaultValue string) *InputPrompt {
	return &InputPrompt{
		Label:   label,
		Default: defaultValue,
	}
}

func (p *InputPrompt) WithValidator(validator Validator) *InputPrompt {
	p.Validator = validator
	return p
}

func (p *InputPrompt) WithMask(mask rune) *InputPrompt {
	p.Mask = mask
	return p
}

func (p *InputPrompt) Run() (string, error) {
	return p.RunWithContext(context.Background())
}

func (p *InputPrompt) RunWithContext(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	fmt.Println()
	if p.Default != "" {
		color.Cyan("(Press Enter to use default, or type to override | Ctrl+C: quit)")
	} else {
		color.Cyan("(Type your text, press Enter to confirm | Ctrl+C: quit)")
	}
	fmt.Println()
	
	for {
		prompt := promptui.Prompt{
			Label:   p.Label,
			Default: p.Default,
			Stdout:  &bellSkipper{},
		}

		if p.Mask != 0 {
			prompt.Mask = p.Mask
		}

		result, err := prompt.Run()

		if err != nil {
			errStr := err.Error()

			if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(errStr, "interrupt") {
				return "", context.Canceled
			}

			if errors.Is(err, promptui.ErrEOF) || strings.Contains(errStr, "EOF") {
				return "", ErrCancel
			}

			return "", err
		}

		if p.Validator != nil {
			if err := p.Validator(result); err != nil {
				fmt.Printf("✗ %v\n\n", err)
				continue
			}
		}

		return result, nil
	}
}

type ConfirmPrompt struct {
	Label    string
	Default  bool
	HideHelp bool
}

func NewConfirmPrompt(label string, defaultYes bool) *ConfirmPrompt {
	return &ConfirmPrompt{
		Label:   label,
		Default: defaultYes,
	}
}

func (p *ConfirmPrompt) Run() (bool, error) {
	return p.RunWithContext(context.Background())
}

func (p *ConfirmPrompt) RunWithContext(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	items := []string{"No", "Yes"}
	defaultIndex := 0
	if p.Default {
		defaultIndex = 1
	}

	color.Yellow(p.Label)
	fmt.Println()
	
	if !p.HideHelp {
		color.Cyan("(↑↓ to navigate, Enter to select, Ctrl+C to cancel)")
		fmt.Println()
	}

	prompt := promptui.Select{
		Label:     "Select",
		Items:     items,
		CursorPos: defaultIndex,
		Size:      2,
		HideHelp:  true,
		Stdout:    &bellSkipper{},
		Stdin:     os.Stdin,
	}

	_, result, err := prompt.Run()

	if err != nil {
		errStr := err.Error()

		if errors.Is(err, promptui.ErrInterrupt) || strings.Contains(errStr, "interrupt") {
			return false, context.Canceled
		}

		if errors.Is(err, promptui.ErrEOF) || strings.Contains(errStr, "EOF") || strings.Contains(errStr, "esc") {
			return false, ErrCancel
		}

		return false, err
	}

	return result == "Yes", nil
}

