package devtui

import (
	"testing"
	"time"
)

// TestOperationIDReuse verifies that when a handler returns an existing operationID,
// the message is updated in place instead of creating a new line
func TestOperationIDReuse(t *testing.T) {
	t.Run("handler with existing operationID should update same line", func(t *testing.T) {
		// Create handler that returns the same operationID
		handler := &TestOperationIDHandler{
			label:       "Test Field",
			value:       "initial",
			operationID: "test-op-123", // Fixed operationID
		}

		h := DefaultTUIForTest()
		h.SetTestMode(true) // Enable test mode

		// Create a test tab since DefaultTUIForTest only creates SHORTCUTS tab
		testTab := h.NewTabSection("Test Tab", "Test description")
		testTab.AddEditHandler(handler).Register()

		field := testTab.FieldHandlers()[0]

		// First execution - should create new message with the fixed operationID
		field.executeChangeSyncWithTracking(field.Value())

		// Verify the handler was called with SetLastOperationID
		if handler.lastSetOperationID == "" {
			t.Error("Handler should have received SetLastOperationID call")
		}

		// Get initial message count
		initialMessageCount := len(testTab.tabContents)

		// Second execution - should reuse the same operationID and update existing message
		field.executeChangeSyncWithTracking(field.Value())

		// Verify message count didn't increase (updated existing instead of creating new)
		finalMessageCount := len(testTab.tabContents)
		if finalMessageCount != initialMessageCount {
			t.Errorf("Expected message count to remain %d, got %d. Messages should be updated in place, not create new lines",
				initialMessageCount, finalMessageCount)
		}

		// Verify the operationID was reused
		if handler.getLastOperationIDCallCount < 2 {
			t.Error("GetLastOperationID should have been called multiple times to check for reusable ID")
		}
	})

	t.Run("multiple handlers with same operationID should maintain separate messages", func(t *testing.T) {
		// Create two handlers with the same operationID but different names
		handler1 := &TestOperationIDHandler{
			label:       "Handler 1",
			value:       "value1",
			operationID: "shared-op-456", // Same operationID
		}

		handler2 := &TestOperationIDHandler{
			label:       "Handler 2",
			value:       "value2",
			operationID: "shared-op-456", // Same operationID
		}
		// Override Name() to make them different
		handler1.handlerName = "TestHandler1"
		handler2.handlerName = "TestHandler2"

		h := DefaultTUIForTest()
		h.SetTestMode(true)

		// Register handlers using the new API
		testTab := h.NewTabSection("Test Tab 2", "Test description for multiple handlers")
		testTab.AddEditHandler(handler1).Register()
		testTab.AddEditHandler(handler2).Register()

		field1 := testTab.FieldHandlers()[0]
		field2 := testTab.FieldHandlers()[1]

		// Execute both handlers
		field1.executeChangeSyncWithTracking(field1.Value())
		field2.executeChangeSyncWithTracking(field2.Value())

		// Should have 2 messages, one per handler
		messageCount := len(testTab.tabContents)
		if messageCount != 2 {
			t.Errorf("Expected 2 messages (one per handler), got %d", messageCount)
		}

		// Execute first handler again - should update its message, not create new
		field1.executeChangeSyncWithTracking(field1.Value())

		finalMessageCount := len(testTab.tabContents)
		if finalMessageCount != 2 {
			t.Errorf("Expected message count to remain 2, got %d. Each handler should maintain its own message", finalMessageCount)
		}
	})

	t.Run("handler without existing operationID should create new lines", func(t *testing.T) {
		// Create handler that always returns empty operationID (new operations each time)
		handler := &TestNewOperationHandler{
			label: "New Operation Field",
			value: "initial",
		}

		h := DefaultTUIForTest()
		h.SetTestMode(true)

		// Register handler using the new API
		testTab := h.NewTabSection("Test Tab 3", "Test description for new operation ID")
		testTab.AddEditHandler(handler).Register()

		field := testTab.FieldHandlers()[0]

		// First execution
		field.executeChangeSyncWithTracking(field.Value())
		initialMessageCount := len(testTab.tabContents)

		// Second execution - should create new message since no existing operationID
		field.executeChangeSyncWithTracking(field.Value())
		finalMessageCount := len(testTab.tabContents)

		// Verify new message was created
		if finalMessageCount <= initialMessageCount {
			t.Errorf("Expected message count to increase from %d to create new line, got %d",
				initialMessageCount, finalMessageCount)
		}
	})
}

// TestOperationIDHandler returns a fixed operationID to simulate updating same operation
type TestOperationIDHandler struct {
	label                       string
	value                       string
	operationID                 string // Fixed operationID to simulate updates
	lastSetOperationID          string
	getLastOperationIDCallCount int
	handlerName                 string // Customizable handler name
}

func (h *TestOperationIDHandler) Name() string {
	if h.handlerName != "" {
		return h.handlerName
	}
	return "TestOperationIDHandler"
}
func (h *TestOperationIDHandler) Label() string          { return h.label }
func (h *TestOperationIDHandler) Value() string          { return h.value }
func (h *TestOperationIDHandler) Editable() bool         { return true }
func (h *TestOperationIDHandler) Timeout() time.Duration { return 0 }

func (h *TestOperationIDHandler) Change(newValue string, progress func(msgs ...any)) {
	if progress != nil {
		progress("Operation completed")
	}
}

func (h *TestOperationIDHandler) SetLastOperationID(id string) {
	h.lastSetOperationID = id
}

func (h *TestOperationIDHandler) GetLastOperationID() string {
	h.getLastOperationIDCallCount++
	return h.operationID // Always return the same ID to simulate updates
}

// TestNewOperationHandler always returns empty operationID to simulate new operations each time
type TestNewOperationHandler struct {
	label    string
	value    string
	lastOpID string
}

func (h *TestNewOperationHandler) Name() string           { return "TestNewOperationHandler" }
func (h *TestNewOperationHandler) Label() string          { return h.label }
func (h *TestNewOperationHandler) Value() string          { return h.value }
func (h *TestNewOperationHandler) Editable() bool         { return true }
func (h *TestNewOperationHandler) Timeout() time.Duration { return 0 }

func (h *TestNewOperationHandler) Change(newValue string, progress func(msgs ...any)) {
	if progress != nil {
		progress("New operation completed")
	}
}

func (h *TestNewOperationHandler) SetLastOperationID(id string) {
	h.lastOpID = id
}

func (h *TestNewOperationHandler) GetLastOperationID() string {
	return "" // Always return empty to simulate new operations each time
}
