package devtui

import (
	"strings"
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

	// Create field name and value variables to avoid circular references
	fieldName := "Test Field"
	fieldValue := "Initial Value"

	// Declare the mockH variable before assigning to prevent circular reference
	var mockH *mockFieldHandler

	// Create a mock field handler
	mockH = &mockFieldHandler{
		name:     fieldName,
		value:    fieldValue,
		editable: true,
		changeValue: func(newValue string) <-chan MessageUpdate {
			updates := make(chan MessageUpdate)
			go func() {
				defer close(updates)
				updates <- MessageUpdate{Content: "Processing: " + fieldName + " change...", Type: messagetype.Info}
				time.Sleep(1 * time.Second)
				if newValue == "Invalid" {
					updates <- MessageUpdate{Content: "Invalid value: " + newValue, Type: messagetype.Error}
					return
				}
				fieldValue = newValue  // Update local value
				mockH.value = newValue // Update mockH value
				updates <- MessageUpdate{Content: "Value successfully changed to " + newValue, Type: messagetype.Success}

			}()
			return updates
		},
	}

	// Create a tab section with our mockH handler
	tabSection := tui.NewTabSection("Test Tab", mockH)

	// Get reference to the field handler
	fieldHandler := &tabSection.fieldHandlers[0]

	// Directly use the TUI's asyncMessageChan
	msgChan := tui.asyncMessageChan

	// Test successful value change
	t.Run("Successful Value Change", func(t *testing.T) {
		msgID := fieldHandler.ExecuteValueChange("New Test Value", tabSection)
		messages := collectMessages(t, msgChan, 2*time.Second)

		validateMessages(t, messages, msgID, "New Test Value", messagetype.Success)
	})

	// Test invalid value
	t.Run("Invalid Value", func(t *testing.T) {
		msgID := fieldHandler.ExecuteValueChange("Invalid", tabSection)
		messages := collectMessages(t, msgChan, 2*time.Second)

		validateMessages(t, messages, msgID, "Invalid", messagetype.Error)
	})
}

func validateMessages(t *testing.T, messages []tuiMessage, msgID MessageID, expectedValue string, expectedType messagetype.Type) {
	if len(messages) == 0 {
		t.Fatal("Expected at least one message, but received none")
	}

	lastMsg := messages[len(messages)-1]
	if lastMsg.Type != expectedType {
		t.Errorf("Final message should be of type %v, got %v", expectedType, lastMsg.Type)
	}

	// Verify message IDs are consistent
	for _, msg := range messages {
		if msg.id != string(msgID) {
			t.Errorf("Message ID mismatch: expected %s, got %s", msgID, msg.id)
		}
	}
}

func TestAsyncFieldIntegration(t *testing.T) {
	mockLogger := func(messageErr any) {
		t.Logf("Integration test logger: %v", messageErr)
	}

	tui := DefaultTUIForTest(mockLogger)

	// Crear variable para almacenar el valor inicial
	fieldName := "Integration Test Field"
	fieldValue := "Initial Value"

	// Definir mockHandler como una variable local para evitar referencias circulares
	var mockHandler *mockFieldHandler

	mockHandler = &mockFieldHandler{
		name:     fieldName,
		value:    fieldValue,
		editable: true,
		changeValue: func(newValue string) <-chan MessageUpdate {
			updates := make(chan MessageUpdate)
			go func() {
				defer close(updates)
				updates <- MessageUpdate{Content: "Processing: " + fieldName + " change...", Type: messagetype.Info}
				time.Sleep(1 * time.Second)
				if newValue == "Invalid" {
					updates <- MessageUpdate{Content: "Invalid value: " + newValue, Type: messagetype.Error}
					return
				}
				fieldValue = newValue        // Update the value
				mockHandler.value = newValue // Update mockHandler value
				updates <- MessageUpdate{Content: "Value successfully changed to " + newValue, Type: messagetype.Success}
			}()
			return updates
		},
	}

	// Crear un tabSection con nuestro mockHandler
	tabSection := tui.NewTabSection("Integration Test Tab", mockHandler)

	// Obtener referencia al fieldHandler
	fieldHandler := &tabSection.fieldHandlers[0]

	// Test cases
	testCases := []struct {
		name      string
		value     string
		wantType  messagetype.Type
		wantCount int
	}{
		{"Normal Change", "Integration Test Value", messagetype.Success, 2},
		{"Empty Value", "", messagetype.Success, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			messageID := fieldHandler.ExecuteValueChange(tc.value, tabSection)
			messages := collectMessages(t, tui.asyncMessageChan, 2*time.Second)

			validateMessages(t, messages, messageID, tc.value, tc.wantType)

			if len(messages) != tc.wantCount {
				t.Errorf("Expected %d messages, got %d", tc.wantCount, len(messages))
			}
		})
	}
}

// Helper function to collect messages with timeout
func collectMessages(t *testing.T, msgChan <-chan tuiMessage, timeout time.Duration) []tuiMessage {
	var messages []tuiMessage
	timeoutChan := time.After(timeout)
	collecting := true

	for collecting {
		select {
		case msg := <-msgChan:
			messages = append(messages, msg)
			if msg.Type == messagetype.Success || msg.Type == messagetype.Error {
				collecting = false
			}
		case <-timeoutChan:
			t.Logf("Timeout reached after collecting %d messages", len(messages))
			collecting = false
		}
	}
	return messages
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestAsyncFieldInRealTUI would test in a real TUI environment
// This is more of an integration test that would need to be run manually
func TestAsyncFieldInRealTUI(t *testing.T) {
	t.Skip("This test requires manual verification in a real TUI environment")

	// Este test se omite deliberadamente
}
