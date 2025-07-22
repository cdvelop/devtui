package devtui

import (
	"strings"
	"testing"
	"time"

	"github.com/cdvelop/messagetype"
)

// TestWriterHandlerRegistration tests the registration of writing handlers
func TestWriterHandlerRegistration(t *testing.T) {
	h := DefaultTUIForTest() // Empty TUI

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test writing handler registration")

	// Create a test writing handler using centralized handler
	handler := NewTestWriterHandler("TestWriter")

	// Register the handler and get its writer
	writer := tab.RegisterWritingHandler(handler)

	if writer == nil {
		t.Fatal("RegisterWritingHandler should return a non-nil writer")
	}

	// Verify the handler was registered
	if tab.writingHandlers == nil {
		t.Fatal("writingHandlers map should be initialized")
	}

	if registeredHandler, exists := tab.writingHandlers["TestWriter"]; !exists {
		t.Fatal("Handler should be registered in writingHandlers map")
	} else if registeredHandler != handler {
		t.Error("Registered handler should be the same instance")
	}
}

// TestHandlerWriterFunctionality tests the HandlerWriter wrapper
func TestHandlerWriterFunctionality(t *testing.T) {
	h := DefaultTUIForTest() // Empty TUI

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test HandlerWriter functionality")

	// Create a test writing handler using centralized handler
	handler := NewTestWriterHandler("TestWriter")

	// Register the handler and get its writer
	writer := tab.RegisterWritingHandler(handler)

	// Write a test message
	testMessage := "Test message from handler"
	n, err := writer.Write([]byte(testMessage))

	if err != nil {
		t.Fatalf("Write should not return error: %v", err)
	}

	if n != len(testMessage) {
		t.Errorf("Write should return correct byte count: expected %d, got %d", len(testMessage), n)
	}

	// Verify handler's SetLastOperationID was called
	if handler.lastOpID == "" {
		t.Error("Handler's SetLastOperationID should have been called")
	}
}

// TestHandlerNameInMessages tests that handler names appear in formatted messages
func TestHandlerNameInMessages(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test handler name in messages")

	// Create a test writing handler
	handler := &TestWriterHandler{
		name: "TestWriter",
	}

	// Register the handler and get its writer
	writer := tab.RegisterWritingHandler(handler)

	// Write a test message
	testMessage := "Test message with handler name"
	writer.Write([]byte(testMessage))

	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)

	// Check if the message contains the handler name
	// Note: We need to check the formatted message in the tab contents
	tab.mu.RLock()
	defer tab.mu.RUnlock()

	if len(tab.tabContents) == 0 {
		t.Fatal("No messages found in tab contents")
	}

	lastContent := tab.tabContents[len(tab.tabContents)-1]
	if lastContent.handlerName != "TestWriter" {
		t.Errorf("Message should have handler name 'TestWriter', got '%s'", lastContent.handlerName)
	}

	if !strings.Contains(lastContent.Content, testMessage) {
		t.Errorf("Message content should contain test message: %s", lastContent.Content)
	}
}

// TestFieldHandlerAutoRegistration tests that FieldHandlers are automatically registered for writing
func TestFieldHandlerAutoRegistration(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test FieldHandler auto-registration")

	// Create a test field handler using centralized handler
	fieldHandler := NewTestEditableHandler("TestField", "test")

	// Add field (should auto-register for writing)
	tab.NewField(fieldHandler)

	// Verify the field handler was auto-registered for writing
	if tab.writingHandlers == nil {
		t.Fatal("writingHandlers map should be initialized")
	}

	handlerName := fieldHandler.Name()
	if registeredHandler, exists := tab.writingHandlers[handlerName]; !exists {
		t.Fatalf("FieldHandler should be auto-registered in writingHandlers map with name '%s'", handlerName)
	} else if registeredHandler != fieldHandler {
		t.Error("Auto-registered handler should be the same instance")
	}
}

// TestOperationIDControl tests that handlers can control message updates vs new messages
func TestOperationIDControl(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test operation ID control")

	// Create a test writing handler
	handler := &TestWriterHandler{
		name: "TestWriter",
	}

	// Register the handler and get its writer
	writer := tab.RegisterWritingHandler(handler)

	// First write - should create new message
	writer.Write([]byte("First message"))
	time.Sleep(10 * time.Millisecond)

	// Enable update mode
	handler.updateMode = true

	// Second write - should update existing message (same operation ID)
	writer.Write([]byte("Updated message"))
	time.Sleep(10 * time.Millisecond)

	// Verify handler received operation IDs
	if handler.lastOpID == "" {
		t.Error("Handler should have received operation ID")
	}

	// Verify messages were created with correct operation ID behavior
	tab.mu.RLock()
	defer tab.mu.RUnlock()

	if len(tab.tabContents) < 2 {
		t.Fatalf("Expected at least 2 messages, got %d", len(tab.tabContents))
	}

	// Check that the handler name is preserved in messages
	for _, content := range tab.tabContents {
		if content.handlerName != "TestWriter" {
			t.Errorf("All messages should have handler name 'TestWriter', got '%s'", content.handlerName)
		}
	}
}

// TestMultipleHandlersInSameTab tests multiple handlers writing to the same tab
func TestMultipleHandlersInSameTab(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test multiple handlers")

	// Create multiple test writing handlers
	handler1 := &TestWriterHandler{name: "Writer1"}
	handler2 := &TestWriterHandler{name: "Writer2"}

	// Register both handlers
	writer1 := tab.RegisterWritingHandler(handler1)
	writer2 := tab.RegisterWritingHandler(handler2)

	// Write messages from both handlers
	writer1.Write([]byte("Message from Writer1"))
	writer2.Write([]byte("Message from Writer2"))

	time.Sleep(10 * time.Millisecond)

	// Verify both handlers are registered
	if len(tab.writingHandlers) != 2 {
		t.Errorf("Expected 2 registered handlers, got %d", len(tab.writingHandlers))
	}

	// Verify messages from both handlers are present
	tab.mu.RLock()
	defer tab.mu.RUnlock()

	var writer1Messages, writer2Messages int
	for _, content := range tab.tabContents {
		switch content.handlerName {
		case "Writer1":
			writer1Messages++
		case "Writer2":
			writer2Messages++
		}
	}

	if writer1Messages == 0 {
		t.Error("Should have messages from Writer1")
	}
	if writer2Messages == 0 {
		t.Error("Should have messages from Writer2")
	}
}

// TestMessageTypeDetection tests that message types are still detected correctly with handler names
func TestMessageTypeDetection(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test message type detection")

	// Create a test writing handler
	handler := &TestWriterHandler{name: "TestWriter"}
	writer := tab.RegisterWritingHandler(handler)

	// Test different message types
	testCases := []struct {
		message      string
		expectedType messagetype.Type
	}{
		{"Error occurred", messagetype.Error},
		{"Success! Operation completed", messagetype.Success},
		{"Warning: This is a warning", messagetype.Warning},
		{"Info: This is information", messagetype.Info},
	}

	for _, tc := range testCases {
		writer.Write([]byte(tc.message))
		time.Sleep(5 * time.Millisecond)

		// Check the last message
		tab.mu.RLock()
		if len(tab.tabContents) > 0 {
			lastMessage := tab.tabContents[len(tab.tabContents)-1]
			if lastMessage.Type != tc.expectedType {
				t.Errorf("Message '%s' should have type %v, got %v", tc.message, tc.expectedType, lastMessage.Type)
			}
		}
		tab.mu.RUnlock()
	}
}
