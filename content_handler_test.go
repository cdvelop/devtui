package devtui

import (
	"strings"
	"testing"
)

// TestHandlerInteractiveInterface verifies the HandlerInteractive interface behavior
func TestHandlerInteractiveInterface(t *testing.T) {
	t.Run("HandlerInteractive should process input and send progress messages", func(t *testing.T) {
		// Use DefaultTUIForTest to get a TUI instance
		h := NewTUI(&TuiConfig{
			AppName:  "TestApp",
			ExitChan: make(chan bool),
			Logger: func(messages ...any) {
				// Silent logging for tests
			},
		})

		// Initialize viewport
		h.viewport.Width = 80
		h.viewport.Height = 24

		// The shortcuts tab and its handler are created automatically by DefaultTUIForTest
		shortcutsTab := h.TabSections[0]
		if len(shortcutsTab.fieldHandlers) == 0 {
			t.Fatal("Shortcuts handler not found")
		}
		field := shortcutsTab.fieldHandlers[0]

		// 1. Verify that the handler is interactive
		if !field.isInteractiveHandler() {
			t.Fatal("Shortcuts handler should be detected as HandlerInteractive")
		}

		// 2. Simulate a change and verify the handler's value is updated
		initialValue := field.handler.Value()
		if initialValue != "en" {
			t.Fatalf("Expected initial language to be 'en', got '%s'", initialValue)
		}

		// Simulate changing the language to "es"
		field.executeChangeSyncWithTracking("es")

		finalValue := field.handler.Value()
		if finalValue != "es" {
			t.Errorf("Expected handler value to be 'es' after change, got '%s'", finalValue)
		}

		// 3. Verify that a progress message was sent
		// After the change, there should be at least one message in the tab contents
		if len(shortcutsTab.tabContents) == 0 {
			t.Fatal("Expected progress message to be sent, but tab contents are empty")
		}

		// Check the content of the last message
		lastMessage := shortcutsTab.tabContents[len(shortcutsTab.tabContents)-1]
		if !strings.Contains(lastMessage.Content, "es") {
			t.Errorf("Expected progress message to contain 'es', but got '%s'", lastMessage.Content)
		}

		t.Logf("SUCCESS: HandlerInteractive correctly processed input, changed value to '%s', and sent progress message.", finalValue)
	})
}