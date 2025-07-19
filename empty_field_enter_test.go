package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestEmptyFieldEnterBehavior tests the behavior when user clears a field and presses Enter
func TestEmptyFieldEnterBehavior(t *testing.T) {
	t.Run("Empty field should call changeFunc with empty string when Enter is pressed", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		// Initialize viewport
		h.viewport.Width = 80
		h.viewport.Height = 24

		field := h.tabSections[0].FieldHandlers()[0]
		originalValue := "initial test value"

		// Ensure field has the original value
		field.SetValue(originalValue)

		// Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		// Initialize editing with current value
		field.tempEditValue = field.Value()
		field.cursor = len([]rune(field.Value()))

		t.Logf("Initial state - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// User clears the entire field
		field.tempEditValue = ""
		field.cursor = 0

		t.Logf("After clearing - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// User presses Enter to save the empty field
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("After pressing Enter - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// The field should now have the value that the changeFunc returned for empty string
		// According to the default changeFunc, it should return "Saved value: " (empty string)
		expectedValue := "Saved value: "
		if field.Value() != expectedValue {
			t.Errorf("Expected field value to be '%s', got '%s'", expectedValue, field.Value())
		}

		// tempEditValue should be cleared after pressing Enter
		if field.tempEditValue != "" {
			t.Errorf("Expected tempEditValue to be empty after Enter, got '%s'", field.tempEditValue)
		}

		// Edit mode should be deactivated
		if h.editModeActivated {
			t.Error("Expected edit mode to be deactivated after Enter")
		}
	})

	t.Run("Field should NOT revert to original value when cleared and Enter is pressed", func(t *testing.T) {
		// Custom changeFunc that allows empty values
		var receivedValue string
		customChangeFunc := func(value any) (string, error) {
			receivedValue = value.(string)
			if receivedValue == "" {
				return "Field was cleared", nil
			}
			return "Field value: " + receivedValue, nil
		}

		// Create TUI with custom field
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Replace the default field with our custom one
		tab := h.tabSections[0]
		tab.setFieldHandlers([]*field{})
		tab.NewField("Test Field", "original value", true, customChangeFunc)

		field := tab.FieldHandlers()[0]

		// Enter editing mode
		h.editModeActivated = true
		h.tabSections[0].indexActiveEditField = 0

		// Initialize editing
		field.tempEditValue = field.Value()
		field.cursor = len([]rune(field.Value()))

		// Clear the field
		field.tempEditValue = ""
		field.cursor = 0

		// Press Enter
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// The changeFunc should have received an empty string
		if receivedValue != "" {
			t.Errorf("Expected changeFunc to receive empty string, got '%s'", receivedValue)
		}

		// The field should have the value returned by changeFunc for empty string
		expectedValue := "Field was cleared"
		if field.Value() != expectedValue {
			t.Errorf("Expected field value to be '%s', got '%s'", expectedValue, field.Value())
		}

		// The field should NOT have reverted to the original value
		if field.Value() == "original value" {
			t.Error("BUG: Field reverted to original value instead of calling changeFunc with empty string")
		}
	})
}
