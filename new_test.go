package devtui_test

import (
	"testing"

	. "github.com/cdvelop/devtui"
)

func TestNewTUI(t *testing.T) {
	// Test configuration with default tabs
	config := &TuiConfig{
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color:         &ColorStyle{}, // Usando un ColorStyle vacío
		LogToFile: func(messageErr any) {
			// Mock function for logging
		},
	}

	tui := NewTUI(config)

	// Check if TUI was created correctly
	if tui == nil {
		t.Fatal("TUI was not created correctly")
	}

	// Since internal fields are not accessible in real usage, we can only test
	// that the TUI was created successfully
	// The default tab should be titled "DEFAULT" according to new.go
}

func TestCustomTabs(t *testing.T) {
	// Create a custom configuration with custom tabs
	customSection := NewTUI(&TuiConfig{}).NewTabSection("CUSTOM1", "custom footer")
	customFields := []Field{
		*NewField(
			"Test Field",
			"test value",
			true,
			func(newValue string) (string, error) {
				return "Value updated to " + newValue, nil
			},
		),
	}

	config := &TuiConfig{
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color:         &ColorStyle{},
	}

	// Add custom tab section
	NewTUI(config).AddTabSections(customSection)

	// Since internal fields are not accessible in real usage,
	// we can only test that the TUI was modified successfully

	// Verify field was created with expected initial value
	customSection.SetFieldHandlers(customFields)
	field := &customFields[0]
	if field.Value() != "test value" {
		t.Errorf("Expected initial value 'test value', got '%s'", field.Value())
	}
}

func TestMultipleTabSections(t *testing.T) {
	// Test adding multiple tab sections
	section1 := NewTUI(&TuiConfig{}).NewTabSection("Tab1", "")
	section2 := NewTUI(&TuiConfig{}).NewTabSection("Tab2", "")

	config := &TuiConfig{
		TabIndexStart: 0,
		Color:         &ColorStyle{},
	}

	totalSections := NewTUI(config).AddTabSections(section1, section2).GetTotalTabSections()

	if totalSections != 2 {
		t.Errorf("Expected 2 tab sections, got %d", totalSections)

	}

}

func TestChannelFunctionality(t *testing.T) {
	// Since the channel is internal to the TUI, we can't directly test it
	// This test should be modified to test observable behavior or removed

	config := &TuiConfig{
		Color: &ColorStyle{},
	}

	tui := NewTUI(config)

	// We can only test that the TUI was created successfully
	if tui == nil {
		t.Error("Failed to create TUI with channel functionality")
	}
}
