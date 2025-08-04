package devtui

import (
	"strings"
	"testing"
	"time"

	. "github.com/cdvelop/tinystring"
)

// TestWriterHandlerRegistration tests the registration of writing handlers
func TestWriterHandlerRegistration(t *testing.T) {
	h := DefaultTUIForTest() // Empty TUI

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test writing handler registration")

	// Create a test writing handler using centralized handler

	// Register the handler and get its writer
	writer := tab.NewWriter("TestWriter", false)

	if writer == nil {
		t.Fatal("RegisterHandlerWriter should return a non-nil writer")
	}

	// Verify the handler was registered
	if tab.writingHandlers == nil {
		t.Fatal("writingHandlers slice should be initialized")
	}

	if registeredHandler := tab.getWritingHandler("TestWriter"); registeredHandler == nil {
		t.Fatal("Handler should be registered in writingHandlers slice")
	}
}

// TestHandlerWriterFunctionality tests the HandlerWriter wrapper
func TestHandlerWriterFunctionality(t *testing.T) {
	h := DefaultTUIForTest() // Empty TUI

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test HandlerWriter functionality")

	// Register the handler and get its writer (basic writer without tracking)
	writer := tab.NewWriter("TestWriter", false)

	// Write a test message
	testMessage := "Test message from handler"
	n, err := writer.Write([]byte(testMessage))

	if err != nil {
		t.Fatalf("Write should not return error: %v", err)
	}

	if n != len(testMessage) {
		t.Errorf("Write should return correct byte count: expected %d, got %d", len(testMessage), n)
	}

	// Verify handler was registered (basic writer doesn't have tracking)
	if registeredHandler := tab.getWritingHandler("TestWriter"); registeredHandler == nil {
		t.Fatal("Handler should be registered in writingHandlers slice")
	}
}

// TestHandlerWriterWithTracking tests the tracking functionality
func TestHandlerWriterWithTracking(t *testing.T) {
	h := DefaultTUIForTest() // Empty TUI

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test HandlerWriter with tracking")

	// Register a writer with tracking enabled
	writer := tab.NewWriter("TrackerWriter", true)

	// Write a test message
	testMessage := "Test tracking message"
	n, err := writer.Write([]byte(testMessage))

	if err != nil {
		t.Fatalf("Write should not return error: %v", err)
	}

	if n != len(testMessage) {
		t.Errorf("Write should return correct byte count: expected %d, got %d", len(testMessage), n)
	}

	// Verify handler was registered with tracking capability
	registeredHandler := tab.getWritingHandler("TrackerWriter")
	if registeredHandler == nil {
		t.Fatal("Handler should be registered in writingHandlers slice")
	}

	// Verify the handler has tracking capability by checking if it has operation ID methods
	if registeredHandler.GetLastOperationID() == "" {
		// This is expected initially - operation ID is set when messages are sent
		t.Log("Operation ID is initially empty, which is correct")
	}

	// Simulate setting an operation ID (this would happen during message processing)
	registeredHandler.SetLastOperationID("test-op-123")

	// Verify the operation ID was set
	if registeredHandler.GetLastOperationID() != "test-op-123" {
		t.Errorf("Expected operation ID 'test-op-123', got '%s'", registeredHandler.GetLastOperationID())
	}
}

// TestHandlerNameInMessages tests that handler names appear in formatted messages
func TestHandlerNameInMessages(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test handler name in messages")

	// Create a test writing handler

	// Register the handler and get its writer
	writer := tab.NewWriter("TestWriter", false)

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

// TestExplicitWriterRegistration tests that writers must be explicitly registered using NewWriter
func TestExplicitWriterRegistration(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test explicit writer registration")

	// Create a test field handler using centralized handler
	fieldHandler := NewTestEditableHandler("TestField", "test")

	// Add field using new API (does NOT auto-register for writing anymore)
	tab.AddEditHandler(fieldHandler, 0)

	// Verify the field handler was NOT auto-registered for writing
	handlerName := fieldHandler.Name()
	if registeredHandler := tab.getWritingHandler(handlerName); registeredHandler != nil {
		t.Fatalf("Handler should NOT be auto-registered in writingHandlers slice with name '%s'", handlerName)
	}

	// Now explicitly register a writer with the same name
	writer := tab.NewWriter(handlerName, false)
	if writer == nil {
		t.Fatal("NewWriter should return a non-nil writer")
	}

	// Verify the writer was explicitly registered
	if registeredHandler := tab.getWritingHandler(handlerName); registeredHandler == nil {
		t.Fatalf("Writer should be explicitly registered in writingHandlers slice with name '%s'", handlerName)
	}
}

// TestOperationIDControl tests that handlers can control message updates vs new messages
func TestOperationIDControl(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test operation ID control")

	// Register a writer with tracking enabled for operation ID control
	writer := tab.NewWriter("TestWriter", true)

	// First write - should create new message
	writer.Write([]byte("First message"))
	time.Sleep(10 * time.Millisecond)

	// Second write - should potentially update existing message (with tracking enabled)
	writer.Write([]byte("Updated message"))
	time.Sleep(10 * time.Millisecond)

	// Verify the writer was registered with tracking capability
	registeredHandler := tab.getWritingHandler("TestWriter")
	if registeredHandler == nil {
		t.Fatal("Handler should be registered in writingHandlers slice")
	}

	// Verify messages were created
	tab.mu.RLock()
	defer tab.mu.RUnlock()

	if len(tab.tabContents) < 1 {
		t.Fatalf("Expected at least 1 message, got %d", len(tab.tabContents))
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

	// Register both handlers
	writer1 := tab.NewWriter("TestWriter1", false)
	writer2 := tab.NewWriter("TestWriter2", false)

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
		case "TestWriter1":
			writer1Messages++
		case "TestWriter2":
			writer2Messages++
		}
	}

	if writer1Messages == 0 {
		t.Error("Should have messages from TestWriter1")
	}
	if writer2Messages == 0 {
		t.Error("Should have messages from TestWriter2")
	}
}

// TestMessageTypeDetection tests that message types are still detected correctly with handler names
func TestMessageTypeDetection(t *testing.T) {
	h := DefaultTUIForTest()

	// Create a new tab for testing
	tab := h.NewTabSection("WritingTest", "Test message type detection")

	// Create a test writing handler
	writer := tab.NewWriter("TestWriter", false)

	// Test different message types
	testCases := []struct {
		message      string
		expectedType MessageType
	}{
		{"Error occurred", Msg.Error},
		{"Success! Operation completed", Msg.Success},
		{"Warning: This is a warning", Msg.Warning},
		{"Info: This is information", Msg.Info},
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
