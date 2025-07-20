package devtui

import (
	"errors"
	"os"
	"time"
)

// DefaultTestHandlers for testing purposes

// TestField1Handler - Default editable field
type TestField1Handler struct {
	currentValue string
}

func (h *TestField1Handler) Label() string { return "Field 1  (Editable)" }
func (h *TestField1Handler) Value() string { return h.currentValue }
func (h *TestField1Handler) Editable() bool { return true }
func (h *TestField1Handler) Timeout() time.Duration { return 0 }

func (h *TestField1Handler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "Saved value: " + strValue, nil
}

// TestField2Handler - Default action field
type TestField2Handler struct{}

func (h *TestField2Handler) Label() string { return "Field 2 (Non-Editable)" }
func (h *TestField2Handler) Value() string { return "special action" }
func (h *TestField2Handler) Editable() bool { return false }
func (h *TestField2Handler) Timeout() time.Duration { return 0 }

func (h *TestField2Handler) Change(newValue any) (string, error) {
	return "Action executed", nil
}

// TestTab2Field1Handler - Tab 2 editable field
type TestTab2Field1Handler struct {
	currentValue string
}

func (h *TestTab2Field1Handler) Label() string { return "Field 1" }
func (h *TestTab2Field1Handler) Value() string { return h.currentValue }
func (h *TestTab2Field1Handler) Editable() bool { return true }
func (h *TestTab2Field1Handler) Timeout() time.Duration { return 0 }

func (h *TestTab2Field1Handler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "Tab 2 saved: " + strValue, nil
}

// TestTab2Field2Handler - Tab 2 error field
type TestTab2Field2Handler struct {
	currentValue string
}

func (h *TestTab2Field2Handler) Label() string { return "Field 2" }
func (h *TestTab2Field2Handler) Value() string { return h.currentValue }
func (h *TestTab2Field2Handler) Editable() bool { return true }
func (h *TestTab2Field2Handler) Timeout() time.Duration { return 0 }

func (h *TestTab2Field2Handler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "", errors.New("Error message test field 2 " + strValue)
}

// NewDefaultTUI creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messages ...any)) *DevTUI {
	// Create basic tabSections for testing/demo
	tmpTUI := &DevTUI{TuiConfig: &TuiConfig{}}
	tab1 := tmpTUI.NewTabSection("Tab 1", "")

	// Create handlers for tab 1
	field1Handler := &TestField1Handler{currentValue: "initial test value"}
	field2Handler := &TestField2Handler{}

	tab1.NewField(field1Handler).
		NewField(field2Handler)
		
	tab1.SetIndex(0)
	tab1.SetActiveEditField(0)

	tab2 := tmpTUI.NewTabSection("Tab 2", "")

	// Create handlers for tab 2
	tab2Field1Handler := &TestTab2Field1Handler{currentValue: "tab 2 value 1"}
	tab2Field2Handler := &TestTab2Field2Handler{currentValue: "error value"}

	tab2.NewField(tab2Field1Handler).
		NewField(tab2Field2Handler)
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

	// Set initial value for the field (deprecated - handlers manage their own state)
	// h.tabSections[0].FieldHandlers()[0].SetValue("initial value")
	// Values are now managed by the handlers themselves

	return h
}
