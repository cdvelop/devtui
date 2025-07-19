package devtui

import (
	"errors"
	"os"
)

// NewDefaultTUI creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messages ...any)) *DevTUI {
	// Create basic tabSections for testing/demo
	tmpTUI := &DevTUI{TuiConfig: &TuiConfig{}}
	tab1 := tmpTUI.NewTabSection("Tab 1", "")

	tab1.NewField(
		"Field 1  (Editable)",
		"initial test value",
		true,
		func(value any) (string, error) {
			strValue := value.(string)
			return "Saved value: " + strValue, nil
		},
	).
		NewField(
			"Field 2 (Non-Editable)",
			"special action",
			false,
			func(value any) (string, error) {
				return "Action executed", nil
			},
		)
	tab1.SetIndex(0)
	tab1.SetActiveEditField(0)

	tab2 := tmpTUI.NewTabSection("Tab 2", "")

	tab2.NewField(
		"Field 1",
		"tab 2 value 1",
		true,
		func(value any) (string, error) {
			strValue := value.(string)
			return "Tab 2 saved: " + strValue, nil
		},
	).
		NewField(
			"Field 2",
			"error value",
			true,
			func(value any) (string, error) {
				strValue := value.(string)
				return "", errors.New("Error message test field 2 " + strValue)
			},
		)
	tab2.SetIndex(1)
	tab2.SetActiveEditField(0)

	tabSections := []*tabSection{tab1, tab2}

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
	testLogger := func(messages ...any) {
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
