package devtui

import (
	"errors"
	"os"
)

// NewDefaultTUI creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messageErr any)) *DevTUI {

	// Create basic tabSections for testing/demo
	// Create temporary minimal DevTUI instance just for NewTabSection
	tmpTUI := &DevTUI{TuiConfig: &TuiConfig{}}

	tab1 := tmpTUI.NewTabSection("Tab 1", "")
	tab1.index = 0
	tab1.SetFieldHandlers([]Field{
		*NewField(
			"Field 1  (Editable)",
			"initial test value",
			true,
			func(value string) (string, error) {
				return "Saved value: " + value, nil
			},
		),
		*NewField(
			"Field 2 (Non-Editable)",
			"special action",
			false,
			func(value string) (string, error) {
				return "Action executed", nil
			},
		),
	})
	tab1.indexActiveEditField = 0

	tab2 := tmpTUI.NewTabSection("Tab 2", "")
	tab2.index = 1
	tab2.SetFieldHandlers([]Field{
		*NewField(
			"Field 1",
			"tab 2 value 1",
			true,
			func(value string) (string, error) {
				return "Tab 2 saved: " + value, nil
			},
		),
		*NewField(
			"Field 2",
			"error value",
			true,
			func(value string) (string, error) {
				return "", errors.New("Error message test field 2 " + value)
			},
		),
	})
	tab2.indexActiveEditField = 0

	tabSections := []TabSection{tab1, tab2}

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
	h.tabSections[0].FieldHandlers()[0].SetValue("initial value")

	return h
}
