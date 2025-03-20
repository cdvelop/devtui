package devtui

import (
	"testing"
	"time"

	"github.com/cdvelop/messagetype"
)

func TestAsyncFieldProcessing(t *testing.T) {
	// Create a mock logger function
	mockLogger := func(messageErr any) {
		t.Logf("Logger called with: %v", messageErr)
	}

	// Initialize the TUI with our mock logger
	tui := DefaultTUIForTest(mockLogger)

	// Get the async field from Tab 2, Field 2
	asyncField := &tui.GetTabSections()[1].FieldHandlers[1]

	// Verify the field is configured as async
	if !asyncField.IsAsync {
		t.Fatalf("Field should be configured as async")
	}

	if asyncField.AsyncFieldValueChange == nil {
		t.Fatalf("AsyncFieldValueChange should be configured")
	}

	// Create a channel to receive messages
	msgChan := make(chan tuiMessage, 10) // Buffered to avoid blocking

	// Use a shorter timeout to prevent tests from hanging
	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	// Start the async processing in a goroutine
	go func() {
		// Manually invoke the async function with our message channel
		asyncField.AsyncFieldValueChange("TestValue", msgChan)
		done <- true
	}()

	// Collect messages until either done signal or timeout
	var messages []tuiMessage
	collecting := true

	for collecting {
		select {
		case msg := <-msgChan:
			messages = append(messages, msg)
		case <-done:
			// Allow a little more time for any final messages
			time.Sleep(100 * time.Millisecond)
			collecting = false
		case <-timeout:
			t.Logf("Timeout reached after collecting %d messages", len(messages))
			collecting = false
		}
	}

	// Verify we got at least one message
	if len(messages) == 0 {
		t.Fatalf("Expected at least one message, but received none")
	}

	t.Logf("Collected %d messages", len(messages))

	// Check if the last message is a completion message
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		if lastMsg.Type != messagetype.Success {
			t.Errorf("Final message should be of type Success, got %v", lastMsg.Type)
		}
	}

	// Verify all non-final messages are progress messages
	for i := 0; i < len(messages)-1; i++ {
		if messages[i].Type != messagetype.Info {
			t.Errorf("tuiMessage %d should be of type Info, got %v", i+1, messages[i].Type)
		}
	}
}

func TestAsyncFieldIntegration(t *testing.T) {
	// Create a mock logger function
	mockLogger := func(messageErr any) {
		t.Logf("Integration test logger: %v", messageErr)
	}

	// Initialize the TUI with our mock logger
	tui := DefaultTUIForTest(mockLogger)

	// Use the helper functions from test_helpers.go
	tabIndex := 1
	fieldIndex := 1
	testValue := "TestIntegration"

	// We need to use RunAsyncFieldTest instead of CollectAsyncMessages
	// because RunAsyncFieldTest handles the channel setup correctly
	messages := RunAsyncFieldTest(t, tui, tabIndex, fieldIndex, testValue)

	// Verify we got at least one message
	if len(messages) == 0 {
		t.Fatalf("Expected at least one async message, but received none")
	}

	t.Logf("Integration test collected %d messages", len(messages))

	// Verify the last message is a completion message
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		if lastMsg.Type != messagetype.Success {
			t.Errorf("Final message should be of type Success, got %v", lastMsg.Type)
		}

		// Check if the message contains our test value
		if !contains(lastMsg.Content, testValue) {
			t.Errorf("Last message should contain '%s', got: '%s'", testValue, lastMsg.Content)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestAsyncFieldInRealTUI would test in a real TUI environment
// This is more of an integration test that would need to be run manually
func TestAsyncFieldInRealTUI(t *testing.T) {
	t.Skip("This test requires manual verification in a real TUI environment")

	// Este test se omite deliberadamente
}
