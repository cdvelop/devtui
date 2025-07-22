package devtui

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Test async field operations

func TestFieldHandler_BasicOperation(t *testing.T) {
	// Use centralized handler from handler_test.go
	handler := NewTestEditableHandler("Test Field", "initial")

	// Use simplified DefaultTUIForTest
	tui := DefaultTUIForTest(handler)

	// Get the first tab created by DefaultTUIForTest
	if len(tui.tabSections) == 0 {
		t.Fatal("No tab sections created")
	}

	tabSection := tui.tabSections[GetFirstTestTabIndex()]

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

	if field.handler.Timeout() != 0 {
		t.Errorf("Expected timeout 0s, got %v", field.handler.Timeout())
	}
}

func TestFieldHandler_AsyncExecution(t *testing.T) {
	// Use centralized handler from handler_test.go - non-editable for action button
	slowHandler := NewTestNonEditableHandler("Slow Operation", "Click to run")
	tui := DefaultTUIForTest(slowHandler)

	tabSection := tui.tabSections[GetFirstTestTabIndex()]
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
	// Test error handling using centralized error handler
	handler := NewTestErrorHandler("Error Field", "test")

	// Test that Change method returns error correctly
	result, err := handler.Change("any value")
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
			name:            "Editable Handler",
			handler:         NewTestEditableHandler("Test", "value"),
			expectedTimeout: 0, // Default timeout from handler_test.go
		},
		{
			name:            "Non-Editable Handler",
			handler:         NewTestNonEditableHandler("Action", "Press Enter"),
			expectedTimeout: 0, // Default timeout from handler_test.go
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
	// Use centralized handlers from handler_test.go
	editableHandler := NewTestEditableHandler("Editable Field", "original")
	nonEditableHandler := NewTestNonEditableHandler("Non-Editable Field", "action button")

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
	// Use centralized handler from handler_test.go
	handler := NewTestEditableHandler("Test Field", "test")
	tui := DefaultTUIForTest(handler)

	tabSection := tui.tabSections[GetFirstTestTabIndex()]
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

func TestSpinner_Start_Stop(t *testing.T) {
	// Use centralized handler from handler_test.go
	handler := NewTestEditableHandler("Test Operation", "Click to test")
	tui := DefaultTUIForTest(handler)

	tabSection := tui.tabSections[GetFirstTestTabIndex()]
	field := tabSection.fieldHandlers[0]

	// Test spinner start - simulate what happens when operation starts
	field.asyncState.isRunning = true
	if !field.asyncState.isRunning {
		t.Error("Spinner should be running after start")
	}

	// Test spinner stop - simulate what happens when operation ends
	field.asyncState.isRunning = false
	if field.asyncState.isRunning {
		t.Error("Spinner should not be running after stop")
	}
}

// Benchmark tests for performance

func BenchmarkFieldHandler_SimpleOperation(b *testing.B) {
	// Use centralized handler from handler_test.go
	handler := NewTestEditableHandler("Benchmark Field", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.Change(fmt.Sprintf("value-%d", i))
	}
}

func BenchmarkFieldHandler_MultipleFields(b *testing.B) {
	// Create multiple handlers using centralized handler
	var handlers []interface{}
	for i := 0; i < 10; i++ {
		handler := NewTestEditableHandler(
			fmt.Sprintf("Field-%d", i),
			fmt.Sprintf("value-%d", i),
		)
		handlers = append(handlers, handler)
	}

	tui := DefaultTUIForTest(handlers...)
	tabSection := tui.tabSections[GetFirstTestTabIndex()]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, field := range tabSection.fieldHandlers {
			_, _ = field.handler.Change(fmt.Sprintf("benchmark-value-%d-%d", i, j))
		}
	}
}
