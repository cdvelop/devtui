package devtui

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Mock handlers for testing

type TestSimpleHandler struct {
	label    string
	value    string
	editable bool
	timeout  time.Duration
}

func (h *TestSimpleHandler) Label() string          { return h.label }
func (h *TestSimpleHandler) Value() string          { return h.value }
func (h *TestSimpleHandler) Editable() bool         { return h.editable }
func (h *TestSimpleHandler) Timeout() time.Duration { return h.timeout }

func (h *TestSimpleHandler) Change(newValue any) (string, error) {
	if h.editable {
		h.value = newValue.(string)
	}
	return fmt.Sprintf("Changed to: %s", h.value), nil
}

type TestSlowHandler struct {
	delay time.Duration
}

func (h *TestSlowHandler) Label() string          { return "Slow Operation" }
func (h *TestSlowHandler) Value() string          { return "Click to run" }
func (h *TestSlowHandler) Editable() bool         { return false }
func (h *TestSlowHandler) Timeout() time.Duration { return h.delay + (1 * time.Second) }

func (h *TestSlowHandler) Change(newValue any) (string, error) {
	time.Sleep(h.delay)
	return fmt.Sprintf("Operation completed after %v", h.delay), nil
}

type TestErrorHandler struct{}

func (h *TestErrorHandler) Label() string          { return "Error Operation" }
func (h *TestErrorHandler) Value() string          { return "Click to fail" }
func (h *TestErrorHandler) Editable() bool         { return false }
func (h *TestErrorHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *TestErrorHandler) Change(newValue any) (string, error) {
	return "", fmt.Errorf("simulated error occurred")
}

type TestTimeoutHandler struct {
	delay time.Duration
}

func (h *TestTimeoutHandler) Label() string          { return "Timeout Operation" }
func (h *TestTimeoutHandler) Value() string          { return "Click to timeout" }
func (h *TestTimeoutHandler) Editable() bool         { return false }
func (h *TestTimeoutHandler) Timeout() time.Duration { return 1 * time.Second } // Short timeout

func (h *TestTimeoutHandler) Change(newValue any) (string, error) {
	time.Sleep(h.delay) // Longer than timeout
	return "Should not reach here", nil
}

// Test async field operations

func TestFieldHandler_BasicOperation(t *testing.T) {
	config := &TuiConfig{
		AppName:  "Test TUI",
		ExitChan: make(chan bool, 1),
		LogToFile: func(messages ...any) {
			// Silent logger for tests
		},
	}

	tui := NewTUI(config)
	handler := &TestSimpleHandler{
		label:    "Test Field",
		value:    "initial",
		editable: true,
		timeout:  5 * time.Second,
	}

	tabSection := tui.NewTabSection("Test Tab", "Test description")
	tabSection.NewField(handler)

	// Verify field was created with handler
	if len(tabSection.fieldHandlers) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(tabSection.fieldHandlers))
	}

	field := tabSection.fieldHandlers[0]
	if field.handler != handler {
		t.Fatal("Field handler not set correctly")
	}

	// Test handler methods through field
	if field.handler.Label() != "Test Field" {
		t.Errorf("Expected label 'Test Field', got '%s'", field.handler.Label())
	}

	if field.handler.Value() != "initial" {
		t.Errorf("Expected value 'initial', got '%s'", field.handler.Value())
	}

	if !field.handler.Editable() {
		t.Error("Expected field to be editable")
	}

	if field.handler.Timeout() != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", field.handler.Timeout())
	}
}

func TestFieldHandler_AsyncExecution(t *testing.T) {
	config := &TuiConfig{
		AppName:  "Test TUI",
		ExitChan: make(chan bool, 1),
		LogToFile: func(messages ...any) {
			// Silent logger for tests
		},
	}

	tui := NewTUI(config)

	// Test slow operation
	slowHandler := &TestSlowHandler{delay: 100 * time.Millisecond}
	tabSection := tui.NewTabSection("Test Tab", "Test description")
	tabSection.NewField(slowHandler)

	field := tabSection.fieldHandlers[0]

	// Test that async state is initialized
	if field.asyncState == nil {
		t.Fatal("Async state not initialized")
	}

	if field.asyncState.isRunning {
		t.Error("Async operation should not be running initially")
	}
}

func TestFieldHandler_ErrorHandling(t *testing.T) {
	handler := &TestErrorHandler{}

	// Test that Change method returns error correctly
	result, err := handler.Change(nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got '%s'", result)
	}

	expectedError := "simulated error occurred"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestFieldHandler_TimeoutConfiguration(t *testing.T) {
	testCases := []struct {
		name            string
		handler         FieldHandler
		expectedTimeout time.Duration
	}{
		{
			name: "Simple Handler",
			handler: &TestSimpleHandler{
				timeout: 10 * time.Second,
			},
			expectedTimeout: 10 * time.Second,
		},
		{
			name: "Slow Handler",
			handler: &TestSlowHandler{
				delay: 2 * time.Second,
			},
			expectedTimeout: 3 * time.Second, // delay + 1s
		},
		{
			name:            "Timeout Handler",
			handler:         &TestTimeoutHandler{},
			expectedTimeout: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeout := tc.handler.Timeout()
			if timeout != tc.expectedTimeout {
				t.Errorf("Expected timeout %v, got %v", tc.expectedTimeout, timeout)
			}
		})
	}
}

func TestFieldHandler_EditableFields(t *testing.T) {
	editableHandler := &TestSimpleHandler{
		label:    "Editable Field",
		value:    "original",
		editable: true,
	}

	nonEditableHandler := &TestSimpleHandler{
		label:    "Non-Editable Field",
		value:    "button",
		editable: false,
	}

	// Test editable field
	if !editableHandler.Editable() {
		t.Error("Handler should be editable")
	}

	result, err := editableHandler.Change("new value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if editableHandler.Value() != "new value" {
		t.Errorf("Expected value 'new value', got '%s'", editableHandler.Value())
	}

	if !strings.Contains(result, "new value") {
		t.Errorf("Expected result to contain 'new value', got '%s'", result)
	}

	// Test non-editable field (button)
	if nonEditableHandler.Editable() {
		t.Error("Handler should not be editable")
	}

	originalValue := nonEditableHandler.Value()
	result, err = nonEditableHandler.Change("attempted change")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if nonEditableHandler.Value() != originalValue {
		t.Error("Non-editable field value should not change")
	}
}

func TestAsyncState_Management(t *testing.T) {
	config := &TuiConfig{
		AppName:  "Test TUI",
		ExitChan: make(chan bool, 1),
		LogToFile: func(messages ...any) {
			// Silent logger for tests
		},
	}

	tui := NewTUI(config)
	handler := &TestSimpleHandler{
		label:   "Test Field",
		value:   "test",
		timeout: 5 * time.Second,
	}

	tabSection := tui.NewTabSection("Test Tab", "Test description")
	tabSection.NewField(handler)

	field := tabSection.fieldHandlers[0]

	// Test initial async state
	if field.asyncState == nil {
		t.Fatal("Async state should be initialized")
	}

	if field.asyncState.isRunning {
		t.Error("Field should not be running initially")
	}

	if field.asyncState.operationID != "" {
		t.Error("Operation ID should be empty initially")
	}

	if field.asyncState.cancel != nil {
		t.Error("Cancel function should be nil initially")
	}
}

func TestSpinner_Integration(t *testing.T) {
	config := &TuiConfig{
		AppName:  "Test TUI",
		ExitChan: make(chan bool, 1),
		LogToFile: func(messages ...any) {
			// Silent logger for tests
		},
	}

	tui := NewTUI(config)
	handler := &TestSlowHandler{delay: 100 * time.Millisecond}

	tabSection := tui.NewTabSection("Test Tab", "Test description")
	tabSection.NewField(handler)

	field := tabSection.fieldHandlers[0]

	// Test spinner initialization
	if field.spinner.Spinner.Frames == nil {
		t.Error("Spinner should be initialized with frames")
	}

	// Spinner should initially not be active
	if field.asyncState.isRunning {
		t.Error("Spinner should not be active initially")
	}
}

// Benchmark tests for performance

func BenchmarkFieldHandler_SimpleOperation(b *testing.B) {
	handler := &TestSimpleHandler{
		label:    "Benchmark Field",
		value:    "test",
		editable: true,
		timeout:  5 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.Change(fmt.Sprintf("value-%d", i))
	}
}

func BenchmarkFieldHandler_MultipleFields(b *testing.B) {
	config := &TuiConfig{
		AppName:   "Benchmark TUI",
		ExitChan:  make(chan bool, 1),
		LogToFile: func(messages ...any) {},
	}

	tui := NewTUI(config)
	tabSection := tui.NewTabSection("Benchmark Tab", "Benchmark description")

	// Create multiple fields
	for i := 0; i < 10; i++ {
		handler := &TestSimpleHandler{
			label:    fmt.Sprintf("Field-%d", i),
			value:    fmt.Sprintf("value-%d", i),
			editable: true,
			timeout:  5 * time.Second,
		}
		tabSection.NewField(handler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, field := range tabSection.fieldHandlers {
			_, _ = field.handler.Change(fmt.Sprintf("benchmark-value-%d-%d", i, j))
		}
	}
}
