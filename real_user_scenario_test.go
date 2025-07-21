package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestRealUserScenario simula exactamente lo que describe el usuario
func TestRealUserScenario(t *testing.T) {
	t.Run("Change port from 8080 to 80 like user described", func(t *testing.T) {
		// Setup exactly like main.go
		config := &TuiConfig{
			AppName:       "DevTUI - New Async API Demo",
			TabIndexStart: 0,
			ExitChan:      make(chan bool),
			Color: &ColorStyle{
				Foreground: "#F4F4F4",
				Background: "#000000",
				Highlight:  "#FF6600",
				Lowlight:   "#666666",
			},
			LogToFile: func(messages ...any) {
				t.Logf("DevTUI Log: %v", messages)
			},
		}

		tui := NewTUI(config)

		// Create port handler exactly like main.go
		portHandler := &PortTestHandler{currentPort: "8080"}

		// Configure tab exactly like main.go
		tui.NewTabSection("Server", "Server configuration").
			NewField(portHandler)

		// Initialize viewport
		tui.viewport.Width = 80
		tui.viewport.Height = 24

		// Navigate to Server tab (should be index 1 after SHORTCUTS)
		tui.activeTab = 1
		portField := tui.tabSections[1].FieldHandlers()[0]

		t.Logf("=== SIMULATING USER SCENARIO ===")
		t.Logf("Step 1: Initial state - field.Value(): '%s'", portField.Value())

		// User sees "8080" and wants to change it to "80"
		// User presses Enter to edit
		t.Logf("Step 2: User presses Enter to edit...")
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("Step 3: Now in edit mode - tempEditValue: '%s', cursor: %d",
			portField.tempEditValue, portField.cursor)

		// User clears the field (simulating selecting all and deleting)
		t.Logf("Step 4: User clears field...")
		portField.tempEditValue = ""
		portField.cursor = 0

		// User types "8"
		t.Logf("Step 5: User types '8'...")
		tui.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'8'},
		})

		// User types "0"
		t.Logf("Step 6: User types '0'...")
		tui.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'0'},
		})

		t.Logf("Step 7: After typing '80' - tempEditValue: '%s', field.Value(): '%s'",
			portField.tempEditValue, portField.Value())

		// At this point user sees "80" in the field but field.Value() still returns "8080"
		// This is expected during editing

		// User presses Enter to confirm
		t.Logf("Step 8: User presses Enter to confirm...")
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("Step 9: After pressing Enter - tempEditValue: '%s', field.Value(): '%s'",
			portField.tempEditValue, portField.Value())

		// NOW THE CRITICAL TEST: field.Value() should return "80"
		finalValue := portField.Value()
		if finalValue != "80" {
			t.Errorf("CRITICAL BUG: After editing, field.Value() should return '80', got '%s'", finalValue)
		}

		// Handler should also be updated
		if portHandler.currentPort != "80" {
			t.Errorf("Handler not updated: expected currentPort '80', got '%s'", portHandler.currentPort)
		}

		// If user enters edit mode again, they should see "80"
		t.Logf("Step 10: Test re-entering edit mode...")
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("Step 11: Re-entered edit mode - tempEditValue: '%s', field.Value(): '%s'",
			portField.tempEditValue, portField.Value())

		// This is where the user's issue manifests: when re-entering edit mode,
		// tempEditValue should be "80", not "8080"
		if portField.tempEditValue != "80" {
			t.Errorf("BUG FOUND: When re-entering edit mode, tempEditValue should be '80', got '%s'",
				portField.tempEditValue)
		}

		if portField.Value() != "80" {
			t.Errorf("BUG FOUND: When re-entering edit mode, field.Value() should be '80', got '%s'",
				portField.Value())
		}
	})

	t.Run("Test what UI actually displays during editing", func(t *testing.T) {
		// Setup
		config := &TuiConfig{
			AppName:  "DevTUI - UI Test",
			ExitChan: make(chan bool),
			Color: &ColorStyle{
				Foreground: "#F4F4F4",
				Background: "#000000",
				Highlight:  "#FF6600",
				Lowlight:   "#666666",
			},
		}

		tui := NewTUI(config)
		portHandler := &PortTestHandler{currentPort: "8080"}
		tui.NewTabSection("Server", "Server configuration").NewField(portHandler)
		tui.viewport.Width = 80
		tui.viewport.Height = 24
		tui.activeTab = 1

		portField := tui.tabSections[1].FieldHandlers()[0]

		// Test the UI rendering during different phases
		t.Logf("=== TESTING UI RENDERING ===")

		// Phase 1: Before editing
		content1 := tui.ContentView()
		t.Logf("Phase 1 - Before editing:\n%s", content1)

		// Phase 2: During editing
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})
		portField.tempEditValue = "80" // User has typed "80"

		content2 := tui.ContentView()
		t.Logf("Phase 2 - During editing (user typed '80'):\n%s", content2)

		// Phase 3: After saving
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		content3 := tui.ContentView()
		t.Logf("Phase 3 - After saving:\n%s", content3)

		// The UI should now show the updated value
		finalValue := portField.Value()
		if finalValue != "80" {
			t.Errorf("UI should show updated value '80', but field.Value() is '%s'", finalValue)
		}
	})
}
