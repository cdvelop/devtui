package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Helper para debuguear el estado de los campos durante los tests
func debugFieldState(t *testing.T, prefix string, field *fieldHandler) {
	t.Logf("%s - Value: '%s', tempEditValue: '%s', cursor: %d",
		prefix, field.Value(), field.tempEditValue, field.cursor)
}

// Helper para inicializar un campo para testing
func prepareFieldForEditing(t *testing.T, h *DevTUI) *fieldHandler {
	h.editModeActivated = true
	h.tabSections[0].indexActiveEditField = 0
	field := &h.tabSections[0].fieldHandlers[0]
	field.tempEditValue = field.Value() // Inicializar tempEditValue con el valor actual
	field.cursor = 0                    // Inicializar cursor
	return field
}

func TestHandleKeyboard(t *testing.T) {
	// Usar la función de inicialización por defecto para tests
	h := prepareForTesting()

	// Test case: Normal mode, changing tabs with tab key
	t.Run("Normal mode - Tab key", func(t *testing.T) {
		h.editModeActivated = false
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})

		if !continueParsing {
			t.Errorf("Expected continueParsing to be true, got false")
		}
	})

	// Test case: Normal mode, pressing enter to enter editing mode
	t.Run("Normal mode - Enter key", func(t *testing.T) {
		h.editModeActivated = false
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		if !continueParsing {
			t.Errorf("Expected continueParsing to be true, got false")
		}

		if !h.editModeActivated {
			t.Errorf("Expected editModeActivated to be true after pressing Enter")
		}
	})

	// Test case: Editing mode, pressing escape to exit
	t.Run("Editing mode - Escape key", func(t *testing.T) {
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEsc})

		if continueParsing {
			t.Errorf("Expected continueParsing to be false after Esc in editing mode")
		}

		if h.editModeActivated {
			t.Errorf("Expected to exit editing mode after Esc")
		}
	})

	// Test case: Editing mode, modifying text
	t.Run("Editing mode - Text input", func(t *testing.T) {
		h := prepareForTesting()
		field := prepareFieldForEditing(t, h)
		initialValue := "initial value"
		field.tempEditValue = initialValue
		field.cursor = 0

		// Escribir 'x'
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'x'},
		})

		expectedValue := "x" + initialValue
		expectedCursor := 1

		if field.tempEditValue != expectedValue {
			t.Errorf("Expected tempEditValue '%s', got '%s'", expectedValue, field.tempEditValue)
		}
		if field.cursor != expectedCursor {
			t.Errorf("Expected cursor %d, got %d", expectedCursor, field.cursor)
		}

		// Escribir 'y'
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'y'},
		})

		expectedValueAfterY := "xy" + initialValue
		expectedCursorAfterY := 2

		if field.tempEditValue != expectedValueAfterY {
			t.Errorf("Expected tempEditValue '%s', got '%s'", expectedValueAfterY, field.tempEditValue)
		}
		if field.cursor != expectedCursorAfterY {
			t.Errorf("Expected cursor %d, got %d", expectedCursorAfterY, field.cursor)
		}

		// Confirmar con Enter
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		if field.Value() != expectedValueAfterY {
			t.Errorf("After Enter: Expected Value '%s', got '%s'", expectedValueAfterY, field.Value())
		}
	})

	// Test case: Editing mode, using backspace
	t.Run("Editing mode - Backspace", func(t *testing.T) {
		h := prepareForTesting()
		field := prepareFieldForEditing(t, h)
		initialValue := field.Value()
		field.tempEditValue = initialValue
		field.cursor = 0

		// Escribir 'a'
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		})

		expectedValueAfterA := "a" + initialValue

		if field.tempEditValue != expectedValueAfterA {
			t.Errorf("After 'a': Expected tempEditValue '%s', got '%s'", expectedValueAfterA, field.tempEditValue)
		}

		// Escribir 'b'
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'b'},
		})

		expectedValueAfterB := "ab" + initialValue

		if field.tempEditValue != expectedValueAfterB {
			t.Errorf("After 'b': Expected tempEditValue '%s', got '%s'", expectedValueAfterB, field.tempEditValue)
		}

		// Backspace
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyBackspace})

		expectedValueAfterBackspace := "a" + initialValue

		if field.tempEditValue != expectedValueAfterBackspace {
			t.Errorf("After backspace: Expected tempEditValue '%s', got '%s'", expectedValueAfterBackspace, field.tempEditValue)
		}
	})

	// Test case: Editing mode, pressing enter to confirm edit
	t.Run("Editing mode - Enter on editable field", func(t *testing.T) {
		h := prepareForTesting()
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].fieldHandlers[0]
		originalValue := "test"
		field.tempEditValue = originalValue + "modified"

		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		expectedFinalValue := "testmodified"
		if field.Value() != expectedFinalValue {
			t.Errorf("Expected Value '%s', got '%s'", expectedFinalValue, field.Value())
		}
	})

	// Test case: Normal mode, Ctrl+C should return quit command
	t.Run("Normal mode - Ctrl+C", func(t *testing.T) {
		h := prepareForTesting()
		h.editModeActivated = false
		h.ExitChan = make(chan bool)

		continueParsing, cmd := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyCtrlC})

		if continueParsing {
			t.Errorf("Expected continueParsing to be false after Ctrl+C")
		}
		if cmd == nil {
			t.Errorf("Expected non-nil command (tea.Quit)")
		}
	})
}

// TestAdditionalKeyboardFeatures prueba características adicionales del teclado
func TestAdditionalKeyboardFeatures(t *testing.T) {

	t.Run("Editing mode - Cancel with ESC discards changes", func(t *testing.T) {
		h := prepareForTesting()
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].fieldHandlers[0]
		originalValue := "Original value"
		field.tempEditValue = "modified"
		field.cursor = len(field.tempEditValue)

		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEsc})

		if field.Value() != originalValue {
			t.Errorf("Expected Value '%s', got '%s'", originalValue, field.Value())
		}
		if field.tempEditValue != originalValue {
			t.Errorf("Expected tempEditValue '%s', got '%s'", originalValue, field.tempEditValue)
		}
	})

	t.Run("Cursor movement in edit mode", func(t *testing.T) {
		h := prepareForTesting()
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].fieldHandlers[0]
		field.tempEditValue = "hello"
		field.cursor = 2

		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyLeft})
		if field.cursor != 1 {
			t.Errorf("Expected cursor 1, got %d", field.cursor)
		}

		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRight})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRight})
		if field.cursor != 3 {
			t.Errorf("Expected cursor 3, got %d", field.cursor)
		}
	})
}
