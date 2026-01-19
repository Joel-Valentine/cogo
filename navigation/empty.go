package navigation

import (
	"fmt"

	"github.com/fatih/color"
)

type EmptyStateHandler struct {
	ResourceName     string
	Context          string
	SuggestedAction  string
	SuggestedCommand string
}

func (h *EmptyStateHandler) Check(items interface{}) error {
	if IsEmpty(items) {
		return ErrEmptyState
	}
	return nil
}

func (h *EmptyStateHandler) CheckAndDisplay(items interface{}) error {
	if !IsEmpty(items) {
		return nil
	}

	h.Display()
	return ErrEmptyState
}

func (h *EmptyStateHandler) Display() {
	message := fmt.Sprintf("No %s found", h.ResourceName)
	if h.Context != "" {
		message += " " + h.Context
	}
	message += "."

	fmt.Println(message)
	fmt.Println()

	if h.SuggestedAction != "" {
		fmt.Println(h.SuggestedAction)
	} else if h.SuggestedCommand != "" {
		fmt.Printf("Run '%s' to get started.\n", h.SuggestedCommand)
	}
}

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case []string:
		return len(val) == 0
	case []int:
		return len(val) == 0
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	case string:
		return len(val) == 0
	default:
		return false
	}
}

func EmptyStateMessage(resourceName, context, suggestedAction string) string {
	message := fmt.Sprintf("No %s found", resourceName)
	if context != "" {
		message += " " + context
	}
	message += ".\n\n"

	if suggestedAction != "" {
		message += suggestedAction
	}

	return message
}

func EmptyStateError(resourceName, context, suggestedAction string) error {
	return fmt.Errorf("%w: %s", ErrEmptyState, EmptyStateMessage(resourceName, context, suggestedAction))
}

func DisplayEmptyDroplets() {
	handler := &EmptyStateHandler{
		ResourceName:     "droplets",
		Context:          "in your DigitalOcean account",
		SuggestedCommand: "cogo create",
	}
	handler.Display()
}

func DisplayEmptyRegions() {
	handler := &EmptyStateHandler{
		ResourceName:    "regions",
		Context:         "in your DigitalOcean account",
		SuggestedAction: "Contact DigitalOcean support if you believe this is an error.",
	}
	handler.Display()
}

func DisplayEmptySSHKeys() {
	handler := &EmptyStateHandler{
		ResourceName:    "SSH keys",
		Context:         "in your DigitalOcean account",
		SuggestedAction: "Add an SSH key at: https://cloud.digitalocean.com/account/security",
	}
	handler.Display()
}

func ErrorMessage(summary, details, suggestion string) string {
	red := color.New(color.FgRed).SprintFunc()
	message := red("✗ Error: ") + summary + "\n"

	if details != "" {
		message += "\n" + details + "\n"
	}

	if suggestion != "" {
		message += "\n" + suggestion
	}

	return message
}

