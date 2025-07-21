package devtui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestEmptyFieldEnterBehavior tests the behavior when user clears a field and presses Enter
func TestEmptyFieldEnterBehavior(t *testing.T) {
	t.Run("Empty field should call changeFunc with empty string when Enter is pressed", func(t *testing.T) {
		// Setup
		h := DefaultTUIForTest()

		// Initialize viewport
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Use centralized function to get correct tab index
		testTabIndex := GetFirstTestTabIndex()
		field := h.tabSections[testTabIndex].FieldHandlers()[0]

		// The field already has "initial test value" from DefaultTUIForTest
		// No need to set it again as SetValue is deprecated

		// Switch to the test tab and enter editing mode
		h.activeTab = testTabIndex
		// Realistic: User enters edit mode by pressing Enter
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("Initial state - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// Realistic: User clears the entire field with backspace
		// Clear existing text (should be "initial test value" = 18 chars)
		for i := 0; i < 25; i++ { // More backspaces to ensure complete clearing
			h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyBackspace})
		}

		t.Logf("After clearing - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// User presses Enter to save the empty field
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		t.Logf("After pressing Enter - Value: '%s', tempEditValue: '%s'", field.Value(), field.tempEditValue)

		// The field should now have the value that the changeFunc returned for empty string
		// According to the TestField1Handler changeFunc, it should have empty string as value
		expectedValue := ""
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
		// Custom handler that allows empty values
		var receivedValue string
		customHandler := NewTestFieldHandler("Test Field", "original value", true, func(value any) (string, error) {
			receivedValue = value.(string)
			if receivedValue == "" {
				return "Field was cleared", nil
			}
			return "Field value: " + receivedValue, nil
		})

		// Create TUI with custom field
		h := DefaultTUIForTest()
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
		// Realistic: User enters edit mode by pressing Enter
		h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// Realistic: Clear the field with backspace
		// Clear existing text ("original value" = ~14 chars)
		for i := 0; i < 25; i++ { // Enough backspaces to ensure complete clearing
			h.HandleKeyboard(tea.KeyMsg{Type: tea.KeyBackspace})
		}

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
