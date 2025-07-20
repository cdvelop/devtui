package devtui

import (
	"errors"
	"os"
	"strings"
	"time"
)

// GetFirstTestTabIndex returns the index of the first test tab
// This centralizes the index calculation to avoid test failures when tabs are added/removed
// Currently, NewTUI always adds SHORTCUTS tab at index 0, so test tabs start at index 1
func GetFirstTestTabIndex() int {
	return 1 // SHORTCUTS tab is always at index 0, so first test tab is at index 1
}

// GetSecondTestTabIndex returns the index of the second test tab
func GetSecondTestTabIndex() int {
	return GetFirstTestTabIndex() + 1 // Second test tab follows first test tab
}

// TestFieldHandler is a basic handler for testing purposes
type TestFieldHandler struct {
	label      string
	value      string
	editable   bool
	changeFunc func(newValue any) (string, error)
	timeout    time.Duration
}

// NewTestFieldHandler creates a new test handler with basic functionality
func NewTestFieldHandler(label, value string, editable bool, changeFunc func(newValue any) (string, error)) *TestFieldHandler {
	if changeFunc == nil {
		// Default change function like the original tests
		changeFunc = func(newValue any) (string, error) {
			return "Saved value: " + newValue.(string), nil
		}
	}

	return &TestFieldHandler{
		label:      label,
		value:      value,
		editable:   editable,
		changeFunc: changeFunc,
		timeout:    0, // No timeout by default
	}
}

// Label returns the field label
func (h *TestFieldHandler) Label() string {
	return h.label
}

// Value returns the current field value
func (h *TestFieldHandler) Value() string {
	return h.value
}

// Editable returns whether the field is editable
func (h *TestFieldHandler) Editable() bool {
	return h.editable
}

// Change processes the field value change
func (h *TestFieldHandler) Change(newValue any) (string, error) {
	if h.changeFunc != nil {
		result, err := h.changeFunc(newValue)
		if err == nil {
			// For test purposes, we need to handle different scenarios:
			// 1. If the changeFunc returns a value that looks like a transformed input (e.g., "Default Value"),
			//    and the input was empty, then the changeFunc is providing the new value
			// 2. Otherwise, use the input value as the new field value
			inputStr := newValue.(string)

			// Special case: if input is empty and result doesn't look like a status message,
			// treat result as the new value (for testing default value transformation)
			if inputStr == "" && !strings.Contains(result, ":") && !strings.HasPrefix(result, "User Input:") && !strings.HasPrefix(result, "Accepted:") && !strings.HasPrefix(result, "Saved") {
				h.value = result
			} else {
				// Normal case: use input as the new value
				h.value = inputStr
			}
		}
		return result, err
	}

	// Default behavior if no change function provided
	h.value = newValue.(string)
	return h.value, nil
}

// Timeout returns the timeout duration for async operations
func (h *TestFieldHandler) Timeout() time.Duration {
	return h.timeout
}

// SetTimeout allows setting a timeout for testing async operations
func (h *TestFieldHandler) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}

// SetValue allows updating the value for testing (simulates external changes)
func (h *TestFieldHandler) SetValue(value string) {
	h.value = value
}

// SetLabel allows updating the label for testing
func (h *TestFieldHandler) SetLabel(label string) {
	h.label = label
}

// SetEditable allows changing the editable state for testing
func (h *TestFieldHandler) SetEditable(editable bool) {
	h.editable = editable
}

// TestAsyncFieldHandler is a handler that simulates async operations
type TestAsyncFieldHandler struct {
	*TestFieldHandler
	asyncDelay time.Duration
}

// NewTestAsyncFieldHandler creates a handler that simulates async operations
func NewTestAsyncFieldHandler(label, value string, editable bool, asyncDelay time.Duration) *TestAsyncFieldHandler {
	return &TestAsyncFieldHandler{
		TestFieldHandler: NewTestFieldHandler(label, value, editable, nil),
		asyncDelay:       asyncDelay,
	}
}

// Change simulates an async operation with delay
func (h *TestAsyncFieldHandler) Change(newValue any) (string, error) {
	if h.asyncDelay > 0 {
		time.Sleep(h.asyncDelay)
	}

	result := "Async result: " + newValue.(string)
	h.value = result
	return result, nil
}

// Timeout returns the async delay as timeout
func (h *TestAsyncFieldHandler) Timeout() time.Duration {
	return h.asyncDelay + (100 * time.Millisecond) // Add buffer for timeout
}

// TestErrorFieldHandler is a handler that always returns errors for testing error handling
type TestErrorFieldHandler struct {
	*TestFieldHandler
	errorMessage string
}

// NewTestErrorFieldHandler creates a handler that always returns errors
func NewTestErrorFieldHandler(label, value string, errorMessage string) *TestErrorFieldHandler {
	return &TestErrorFieldHandler{
		TestFieldHandler: NewTestFieldHandler(label, value, true, nil),
		errorMessage:     errorMessage,
	}
}

// Change always returns an error for testing error handling
func (h *TestErrorFieldHandler) Change(newValue any) (string, error) {
	return "", errors.New(h.errorMessage)
}

// Specific handlers for DefaultTUIForTest

// TestField1Handler - Default editable field
type TestField1Handler struct {
	currentValue string
}

func (h *TestField1Handler) Label() string          { return "Field 1  (Editable)" }
func (h *TestField1Handler) Value() string          { return h.currentValue }
func (h *TestField1Handler) Editable() bool         { return true }
func (h *TestField1Handler) Timeout() time.Duration { return 0 }

func (h *TestField1Handler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "Saved value: " + strValue, nil
}

// TestField2Handler - Default action field
type TestField2Handler struct{}

func (h *TestField2Handler) Label() string          { return "Field 2 (Non-Editable)" }
func (h *TestField2Handler) Value() string          { return "special action" }
func (h *TestField2Handler) Editable() bool         { return false }
func (h *TestField2Handler) Timeout() time.Duration { return 0 }

func (h *TestField2Handler) Change(newValue any) (string, error) {
	return "Action executed", nil
}

// TestTab2Field1Handler - Tab 2 editable field
type TestTab2Field1Handler struct {
	currentValue string
}

func (h *TestTab2Field1Handler) Label() string          { return "Field 1" }
func (h *TestTab2Field1Handler) Value() string          { return h.currentValue }
func (h *TestTab2Field1Handler) Editable() bool         { return true }
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

func (h *TestTab2Field2Handler) Label() string          { return "Field 2" }
func (h *TestTab2Field2Handler) Value() string          { return h.currentValue }
func (h *TestTab2Field2Handler) Editable() bool         { return true }
func (h *TestTab2Field2Handler) Timeout() time.Duration { return 0 }

func (h *TestTab2Field2Handler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "", errors.New("Error message test field 2 " + strValue)
}

// DefaultTUIForTest creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messages ...any)) *DevTUI {
	// Set up test environment first
	os.Setenv("TEST_MODE", "true")

	// Create basic tabSections for testing/demo
	tmpTUI := &DevTUI{TuiConfig: &TuiConfig{}}
	tab1 := tmpTUI.NewTabSection("Tab 1", "")

	// Create handlers for tab 1
	field1Handler := &TestField1Handler{currentValue: "initial test value"}
	field2Handler := &TestField2Handler{}

	tab1.NewField(field1Handler).
		NewField(field2Handler)

	tab1.SetIndex(GetFirstTestTabIndex()) // Use centralized index
	tab1.SetActiveEditField(0)

	tab2 := tmpTUI.NewTabSection("Tab 2", "")

	// Create handlers for tab 2
	tab2Field1Handler := &TestTab2Field1Handler{currentValue: "tab 2 value 1"}
	tab2Field2Handler := &TestTab2Field2Handler{currentValue: "error value"}

	tab2.NewField(tab2Field1Handler).
		NewField(tab2Field2Handler)
	tab2.SetIndex(GetSecondTestTabIndex()) // Use centralized index
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
