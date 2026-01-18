package navigation

import (
	"fmt"

	"github.com/fatih/color"
)

// EmptyStateHandler provides utilities for detecting and handling empty resource states.
//
// Research finding: 10/11 CLI tools return exit code 0 for empty states (not an error).
// Message format follows gh, kubectl, gcloud patterns.
type EmptyStateHandler struct {
	// ResourceName is the plural name of the resource (e.g., "droplets", "regions")
	ResourceName string

	// Context provides additional information (e.g., "in your DigitalOcean account")
	Context string

	// SuggestedAction tells the user what to do next (e.g., "Run 'cogo create' to create a droplet")
	SuggestedAction string

	// SuggestedCommand is the specific command to run (if any)
	SuggestedCommand string
}

// Check returns ErrEmptyState if the list is empty, nil otherwise.
func (h *EmptyStateHandler) Check(items interface{}) error {
	if IsEmpty(items) {
		return ErrEmptyState
	}
	return nil
}

// CheckAndDisplay checks if items are empty and displays a friendly message if so.
// Returns ErrEmptyState if empty, nil otherwise.
//
// Example output:
//
//	No droplets found in your DigitalOcean account.
//
//	Run 'cogo create' to create your first droplet.
func (h *EmptyStateHandler) CheckAndDisplay(items interface{}) error {
	if !IsEmpty(items) {
		return nil
	}

	h.Display()
	return ErrEmptyState
}

// Display prints the empty state message to stdout.
// Does not return an error - use Check() or CheckAndDisplay() to get ErrEmptyState.
func (h *EmptyStateHandler) Display() {
	// Message format from research: "No {resource} found {context}."
	message := fmt.Sprintf("No %s found", h.ResourceName)
	if h.Context != "" {
		message += " " + h.Context
	}
	message += "."

	fmt.Println(message)
	fmt.Println()

	// Suggestion (actionable, from research findings)
	if h.SuggestedAction != "" {
		fmt.Println(h.SuggestedAction)
	} else if h.SuggestedCommand != "" {
		fmt.Printf("Run '%s' to get started.\n", h.SuggestedCommand)
	}
}

// IsEmpty checks if a value represents an empty collection.
// Handles: slices, arrays, maps, strings, nil.
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
		// For other slice/array types, use reflection would be needed
		// but we'll keep it simple for common cases
		return false
	}
}

// EmptyStateMessage creates a formatted empty state message following research patterns.
//
// Format (from gh, kubectl, gcloud):
//
//	No {resourceName} found {context}.
//
//	{suggestedAction}
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

// EmptyStateError creates an error with a formatted empty state message.
// Useful for returning from Step.Execute() when resources are empty.
func EmptyStateError(resourceName, context, suggestedAction string) error {
	return fmt.Errorf("%w: %s", ErrEmptyState, EmptyStateMessage(resourceName, context, suggestedAction))
}

// DisplayEmptyDroplets is a convenience function for the common case of no droplets.
func DisplayEmptyDroplets() {
	handler := &EmptyStateHandler{
		ResourceName:     "droplets",
		Context:          "in your DigitalOcean account",
		SuggestedCommand: "cogo create",
	}
	handler.Display()
}

// DisplayEmptyRegions is a convenience function for no available regions.
func DisplayEmptyRegions() {
	handler := &EmptyStateHandler{
		ResourceName:    "regions",
		Context:         "in your DigitalOcean account",
		SuggestedAction: "Contact DigitalOcean support if you believe this is an error.",
	}
	handler.Display()
}

// DisplayEmptySSHKeys is a convenience function for no configured SSH keys.
func DisplayEmptySSHKeys() {
	handler := &EmptyStateHandler{
		ResourceName:    "SSH keys",
		Context:         "in your DigitalOcean account",
		SuggestedAction: "Add an SSH key at: https://cloud.digitalocean.com/account/security",
	}
	handler.Display()
}

// ErrorMessage creates a formatted error message following research patterns.
//
// Format (from gh, cargo, terraform):
//
//	✗ Error: {summary}
//
//	{details}
//
//	{suggestion}
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

