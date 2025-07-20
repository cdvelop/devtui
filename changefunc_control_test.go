package devtui

import (
	"errors"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestChangeFuncControlsEmptyFieldBehavior demonstrates that changeFunc has full control
// over what happens when a field is cleared, not DevTUI
func TestChangeFuncControlsEmptyFieldBehavior(t *testing.T) {
	t.Run("changeFunc can reject empty values", func(t *testing.T) {
		// Custom handler that rejects empty values
		customHandler := NewTestFieldHandler("Required Field", "initial value", true, func(value any) (string, error) {
			strValue := value.(string)
			if strValue == "" {
				return "", errors.New("Field cannot be empty")
			}
			return "Accepted: " + strValue, nil
		})

		// Create TUI with custom field
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Replace the default field with our custom one
		testTabIndex := GetFirstTestTabIndex()
		tab := h.tabSections[testTabIndex]
		tab.setFieldHandlers([]*field{})
		tab.NewField(customHandler)

		field := tab.FieldHandlers()[0]

		// Switch to test tab and enter editing mode
		h.activeTab = testTabIndex
		h.editModeActivated = true
		h.tabSections[testTabIndex].indexActiveEditField = 0

		// Initialize editing
		field.tempEditValue = field.Value()
		field.cursor = len([]rune(field.Value()))

		// Clear the field
		field.tempEditValue = ""
		field.cursor = 0

		// Press Enter - changeFunc should reject the empty value
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// The field should still have the original value because changeFunc rejected the empty value
		expectedValue := "initial value"
		if field.Value() != expectedValue {
			t.Errorf("Expected field to keep original value '%s' after changeFunc rejects empty, got '%s'", expectedValue, field.Value())
		}

		// Edit mode should be deactivated even if changeFunc fails
		if h.editModeActivated {
			t.Error("Expected edit mode to be deactivated after Enter, even when changeFunc fails")
		}
	})

	t.Run("changeFunc can accept and transform empty values", func(t *testing.T) {
		// Ensure TEST_MODE is set for synchronous execution
		os.Setenv("TEST_MODE", "true")

		// Custom handler that accepts empty values and transforms them
		customHandler := NewTestFieldHandler("Optional Field", "original value", true, func(value any) (string, error) {
			strValue := value.(string)
			if strValue == "" {
				return "Default Value", nil
			}
			return "User Input: " + strValue, nil
		})

		// Create TUI with custom field
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Replace the default field
		testTabIndex := GetFirstTestTabIndex()
		tab := h.tabSections[testTabIndex]
		tab.setFieldHandlers([]*field{})
		tab.NewField(customHandler)

		field := tab.FieldHandlers()[0]

		// Switch to test tab and enter editing mode
		h.activeTab = testTabIndex
		h.editModeActivated = true
		h.tabSections[testTabIndex].indexActiveEditField = 0

		// Initialize editing
		field.tempEditValue = field.Value()
		field.cursor = len([]rune(field.Value()))

		// Clear the field
		field.tempEditValue = ""
		field.cursor = 0

		// Press Enter - changeFunc should accept and transform the empty value
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// The field should have the transformed value from changeFunc
		expectedValue := "Default Value"
		if field.Value() != expectedValue {
			t.Errorf("Expected field value to be '%s' after changeFunc transforms empty value, got '%s'", expectedValue, field.Value())
		}
	})

	t.Run("changeFunc can preserve empty values", func(t *testing.T) {
		// Custom handler that allows and preserves empty values
		customHandler := NewTestFieldHandler("Clearable Field", "some value", true, func(value any) (string, error) {
			strValue := value.(string)
			return strValue, nil // Return exactly what was input, including empty string
		})

		// Create TUI with custom field
		h := DefaultTUIForTest(func(messages ...any) {})
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Replace the default field
		testTabIndex := GetFirstTestTabIndex()
		tab := h.tabSections[testTabIndex]
		tab.setFieldHandlers([]*field{})
		tab.NewField(customHandler)

		field := tab.FieldHandlers()[0]

		// Switch to test tab and enter editing mode
		h.activeTab = testTabIndex
		h.editModeActivated = true
		h.tabSections[testTabIndex].indexActiveEditField = 0

		// Initialize editing
		field.tempEditValue = field.Value()
		field.cursor = len([]rune(field.Value()))

		// Clear the field
		field.tempEditValue = ""
		field.cursor = 0

		// Press Enter - changeFunc should preserve the empty value
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// The field should be empty
		expectedValue := ""
		if field.Value() != expectedValue {
			t.Errorf("Expected field value to be empty '%s', got '%s'", expectedValue, field.Value())
		}
	})
}
