package devtui

import (
	"errors"
	"os"
)

// NewDefaultTUI creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messageErr any)) *DevTUI {

	// Create basic tabSections for testing/demo
	tabSections := []TabSection{
		{
			Title: "Tab 1",
			index: 0,
			FieldHandlers: []FieldHandler{
				{
					Name:     "Field 1  (Editable)",
					Value:    "initial test value",
					Editable: true,
					cursor:   0,
					FieldValueChange: func(value string) (string, error) {
						return "Saved value: " + value, nil
					},
				},
				{
					Name:     "Field 2 (Non-Editable)",
					Value:    "special action",
					Editable: false,
					FieldValueChange: func(value string) (string, error) {
						return "Action executed", nil
					},
				},
			},
			indexActiveEditField: 0,
		},
		{
			Title: "Tab 2",
			index: 1,
			FieldHandlers: []FieldHandler{
				{
					Name:     "Field 1",
					Value:    "tab 2 value 1",
					Editable: true,
					cursor:   0,
					FieldValueChange: func(value string) (string, error) {
						return "Tab 2 saved: " + value, nil
					},
				},
				{
					Name:     "Field 2",
					Value:    "error value",
					Editable: true,
					cursor:   0,
					FieldValueChange: func(value string) (string, error) {
						return "", errors.New("Error message test field 2 " + value)
					},
				},
			},
			indexActiveEditField: 0,
		},
	}

	// Initialize the UI
	h := NewTUI(&TuiConfig{
		TabIndexStart: 0,               // Start with the first tab
		ExitChan:      make(chan bool), // Channel to signal exit
		Color:         nil,             // Use default colors
		LogToFile:     LogToFile,
	}).AddTabSections(tabSections...)

	return h
}

// prepareForTesting configures a DevTUI instance for use in unit tests
func prepareForTesting() *DevTUI {
	// Create a logger that doesn't do anything during tests
	testLogger := func(messageErr any) {
		// In test mode, we don't need to log
		if os.Getenv("TEST_MODE") != "true" {
			// This is a no-op logger for tests
		}
	}

	// Get default TUI instance
	h := DefaultTUIForTest(testLogger)

	// Set up test environment
	os.Setenv("TEST_MODE", "true")

	// Set initial value for the field
	h.tabSections[0].FieldHandlers[0].Value = "initial value"

	return h
}
