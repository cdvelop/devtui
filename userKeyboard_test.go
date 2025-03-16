package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleKeyboard(t *testing.T) {
	// Usar la función de inicialización por defecto para tests
	h := prepareForTesting()

	// Test case: Normal mode, changing tabs with tab key
	t.Run("Normal mode - Tab key", func(t *testing.T) {
		h.editModeActivated = false
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab}) // Ignoramos el comando

		if !continueParsing {
			t.Errorf("Expected continueParsing to be true, got false")
		}
	})

	// Test case: Normal mode, pressing enter to enter editing mode
	t.Run("Normal mode - Enter key", func(t *testing.T) {
		h.editModeActivated = false
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter}) // Ignoramos el comando

		if !continueParsing {
			t.Errorf("Expected continueParsing to be true, got false")
		}

		if !h.editModeActivated {
			t.Errorf("Expected editModeActivated to be true after pressing Enter")
		}
	})

	// Test case: Editing mode, pressing escape to exit
	t.Run("Editing mode - Escape key", func(t *testing.T) {
		// Setup: Enter editing mode first
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
		// Reset para esta prueba
		h := prepareForTesting()

		// Setup: Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].FieldHandlers[0]
		initialValue := field.Value // 'initial value'

		// Por defecto, el cursor estará al inicio (posición 0)
		// Simulamos escribir 'x' en esa posición
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'x'},
		})

		// El carácter debe ser insertado al inicio ya que el cursor está en la posición 0
		expectedValue := "x" + initialValue
		expectedCursor := 1 // Cursor debe moverse una posición a la derecha

		if field.tempEditValue != expectedValue {
			t.Fatalf("Expected value to be '%s', got '%s'", expectedValue, field.tempEditValue)
		}

		if field.cursor != expectedCursor {
			t.Fatalf("Expected cursor to be at position %d, got %d", expectedCursor, field.cursor)
		}

		// Ahora probemos añadiendo otro carácter 'y' en la nueva posición del cursor
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'y'},
		})

		// El carácter 'y' debe insertarse después de la 'x'
		expectedValue = "xy" + initialValue
		expectedCursor = 2

		if field.Value != expectedValue {
			t.Fatalf("Expected value to be '%s', got '%s'", expectedValue, field.Value)
		}

		if field.cursor != expectedCursor {
			t.Fatalf("Expected cursor to be at position %d, got %d", expectedCursor, field.cursor)
		}
	})

	// Test case: Editing mode, using backspace
	t.Run("Editing mode - Backspace", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Setup: Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].FieldHandlers[0]

		// Primero añadimos algunos caracteres al inicio
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		})
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'b'},
		})

		// Guardamos el valor y la posición del cursor después de añadir los caracteres
		valueBeforeBackspace := field.Value
		cursorBeforeBackspace := field.cursor

		// Verificamos que se añadieron correctamente
		expectedValueBeforeBackspace := "abinitial value"
		if valueBeforeBackspace != expectedValueBeforeBackspace {
			t.Errorf("Setup failed: Expected value to be '%s', got '%s'",
				expectedValueBeforeBackspace, valueBeforeBackspace)
		}

		// Ahora usamos backspace para eliminar el último carácter insertado ('b')
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyBackspace})

		// Verificamos que se eliminó correctamente
		expectedValueAfterBackspace := "ainitial value"
		expectedCursorAfterBackspace := cursorBeforeBackspace - 1

		if field.Value != expectedValueAfterBackspace {
			t.Errorf("After backspace: Expected value to be '%s', got '%s'",
				expectedValueAfterBackspace, field.Value)
		}

		if field.cursor != expectedCursorAfterBackspace {
			t.Errorf("After backspace: Expected cursor to be at position %d, got %d",
				expectedCursorAfterBackspace, field.cursor)
		}
	})

	// Test case: Editing mode, pressing enter to confirm edit
	t.Run("Editing mode - Enter on editable field", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Setup: Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		originalField := &h.tabSections[0].FieldHandlers[0]
		originalField.Value = "test"

		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		if continueParsing {
			t.Errorf("Expected continueParsing to be false after Enter in editing mode")
		}

		if h.editModeActivated {
			t.Errorf("Expected to exit editing mode after Enter")
		}
	})

	// Test case: Normal mode, Ctrl+C should return quit command
	t.Run("Normal mode - Ctrl+C", func(t *testing.T) {
		h := prepareForTesting() // Reset para esta prueba
		h.editModeActivated = false

		// Asegurarnos de que ExitChan está correctamente inicializado para esta prueba
		h.ExitChan = make(chan bool)

		continueParsing, cmd := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyCtrlC})

		if continueParsing {
			t.Errorf("Expected continueParsing to be false after Ctrl+C")
		}

		if cmd == nil {
			t.Errorf("Expected non-nil command (tea.Quit) after Ctrl+C")
		}
	})
}

// TestAdditionalKeyboardFeatures prueba características adicionales del teclado
func TestAdditionalKeyboardFeatures(t *testing.T) {
	h := prepareForTesting()

	// Test: Cancelación de edición con ESC debe restaurar el valor original
	t.Run("Editing mode - Cancel with ESC discards changes", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Setup: Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].FieldHandlers[0]
		originalValue := "Original value"
		field.Value = originalValue

		// Modificamos el valor añadiendo caracteres
		h.HandleKeyboard(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'m', 'o', 'd', 'i', 'f', 'i', 'e', 'd'},
		})

		expectedTemValue := "modified"

		// Verificamos que el campo tempEditValue fue modificado
		if field.tempEditValue != expectedTemValue {
			t.Fatalf("Setup failed: Expected tempEditValue to be '%s', got '%s'", expectedTemValue, field.tempEditValue)
		}

		// Ahora presionamos ESC para cancelar
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEsc})

		if continueParsing {
			t.Errorf("Expected continueParsing to be false after ESC in editing mode")
		}

		if h.editModeActivated {
			t.Errorf("Expected to exit editing mode after ESC")
		}

		// Verificamos que el valor volvió al original
		if field.Value != originalValue {
			t.Errorf("After ESC: Expected value to be restored to '%s', got '%s'",
				originalValue, field.Value)
		}

		// Verificamos que el campo tempEditValue fue limpiado
		if field.tempEditValue != "" {
			t.Errorf("After ESC: Expected tempEditValue to be empty, got '%s'", field.tempEditValue)
		}
	})

	// Test: Navegación entre campos con flechas up y down no afecta a los inputs
	t.Run("Arrow keys in normal mode", func(t *testing.T) {
		// Configuración inicial - normal mode
		h.editModeActivated = false
		h.tabSections[0].indexActiveEditField = 0
		initialIndex := h.tabSections[0].indexActiveEditField

		// Intentar navegar con flechas up o down - no debería cambiar inputs
		continueParsing, _ := h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyDown})
		if !continueParsing {
			t.Errorf("Expected continueParsing to be true after Down key")
		}
		if h.tabSections[0].indexActiveEditField != initialIndex {
			t.Errorf("Expected indexActiveEditField to remain %d, but got %d",
				initialIndex, h.tabSections[0].indexActiveEditField)
		}

		continueParsing, _ = h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyUp})
		if !continueParsing {
			t.Errorf("Expected continueParsing to be true after Up key")
		}
		if h.tabSections[0].indexActiveEditField != initialIndex {
			t.Errorf("Expected indexActiveEditField to remain %d, but got %d",
				initialIndex, h.tabSections[0].indexActiveEditField)
		}
	})

	// Test: Movimiento del cursor en modo edición
	t.Run("Cursor movement in edit mode", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Configuración inicial - modo edición
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0
		field := &h.tabSections[0].FieldHandlers[0]
		field.Value = "hello"
		field.cursor = 2 // Cursor en medio (he|llo)

		// Mover cursor a la izquierda
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		if field.cursor != 1 {
			t.Errorf("Expected cursor to move left to position 1, got %d", field.cursor)
		}

		// Mover cursor a la derecha
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRight})
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		if field.cursor != 3 {
			t.Errorf("Expected cursor to move right to position 3, got %d", field.cursor)
		}
	})
}
