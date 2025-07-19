package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestSpaceKeyInEditMode tests that the space key works correctly in edit mode
func TestSpaceKeyInEditMode(t *testing.T) {
	t.Run("Space key should insert space in edit mode", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		// Initialize viewport with a reasonable size for testing
		h.viewport.Width = 80
		h.viewport.Height = 24

		field := h.tabSections[0].FieldHandlers()[0]

		// Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		// Clear field
		field.tempEditValue = ""
		field.cursor = 0

		// Type "hello"
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'h'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'e'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'l'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'l'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'o'},
		})

		// Should have "hello"
		if field.tempEditValue != "hello" {
			t.Errorf("Expected 'hello', got '%s'", field.tempEditValue)
		}

		// Now try to add a space using KeySpace
		h.HandleKeyboard(tea.KeyMsg{
			Type: tea.KeySpace,
		})

		t.Logf("After pressing space - tempEditValue: '%s', cursor: %d", field.tempEditValue, field.cursor)

		// Should have "hello " (hello with space)
		expectedAfterSpace := "hello "
		if field.tempEditValue != expectedAfterSpace {
			t.Errorf("Expected '%s', got '%s'", expectedAfterSpace, field.tempEditValue)
		}

		// Continue typing "world"
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'w'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'o'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'r'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'l'},
		})

		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'d'},
		})

		// Should have "hello world"
		expectedFinal := "hello world"
		if field.tempEditValue != expectedFinal {
			t.Errorf("Expected '%s', got '%s'", expectedFinal, field.tempEditValue)
		}
	})

	t.Run("Space using Runes should work", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		field := h.tabSections[0].FieldHandlers()[0]

		// Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		// Clear field
		field.tempEditValue = ""
		field.cursor = 0

		// Type "hello" using Runes
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'h', 'e', 'l', 'l', 'o'},
		})

		// Add space using Runes
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{' '},
		})

		t.Logf("After adding space via Runes - tempEditValue: '%s', cursor: %d", field.tempEditValue, field.cursor)

		// Add "world" using Runes
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'w', 'o', 'r', 'l', 'd'},
		})

		// Should have "hello world"
		expectedFinal := "hello world"
		if field.tempEditValue != expectedFinal {
			t.Errorf("Expected '%s', got '%s'", expectedFinal, field.tempEditValue)
		}
	})

	t.Run("Complete text editing with spaces should work", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		field := h.tabSections[0].FieldHandlers()[0]

		// Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		// Start with initial value and edit it
		field.tempEditValue = field.Value()       // Start with current value
		field.cursor = len([]rune(field.Value())) // Cursor at end

		// Clear the field first
		field.tempEditValue = ""
		field.cursor = 0

		// Type a complete sentence with spaces
		// "Hello world from Go"
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeySpace}) // First space
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeySpace}) // Second space
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeySpace}) // Third space
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})

		expectedText := "Hello world from Go"
		if field.tempEditValue != expectedText {
			t.Errorf("Expected '%s', got '%s'", expectedText, field.tempEditValue)
		}

		if field.cursor != len([]rune(expectedText)) {
			t.Errorf("Expected cursor at position %d, got %d", len([]rune(expectedText)), field.cursor)
		}

		t.Logf("Successfully typed: '%s'", field.tempEditValue)
	})
}
