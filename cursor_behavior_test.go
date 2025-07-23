package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestCursorBehaviorInEditMode verifica el comportamiento del cursor durante la edición
func TestCursorBehaviorInEditMode(t *testing.T) {
	t.Run("Cursor position affects character insertion correctly", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest(func(messages ...any) {})
		portHandler := &PortTestHandler{currentPort: "8080"}
		h.NewTabSection("Server", "Config").NewEditHandler(portHandler).Register()

		// Set viewport size properly for calculation
		h.viewport.Width = 80
		h.viewport.Height = 24

		serverTabIndex := len(h.tabSections) - 1
		h.activeTab = serverTabIndex
		field := h.tabSections[serverTabIndex].FieldHandlers()[0]

		// Enter edit mode
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// Test insertion at beginning
		field.cursor = 0
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'9'},
		})

		if field.tempEditValue != "98080" {
			t.Errorf("Expected '98080' when typing at beginning, got '%s'", field.tempEditValue)
		}

		// Test insertion at end
		field.tempEditValue = "8080"
		field.cursor = len([]rune(field.tempEditValue))

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'9'},
		})

		if field.tempEditValue != "80809" {
			t.Errorf("Expected '80809' when typing at end, got '%s'", field.tempEditValue)
		}
	})
}
