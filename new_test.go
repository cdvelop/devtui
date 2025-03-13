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
	customSection := TabSection{
		Title: "CUSTOM1",
		FieldHandlers: []FieldHandler{
			{
				Name:     "testField",
				Label:    "Test Field",
				Value:    "test value",
				Editable: true,
				FieldValueChange: func(newValue string) (string, error) {
					return "Value updated to " + newValue, nil
				},
			},
		},
		SectionFooter: "custom footer",
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

	// Test FieldValueChange function if accessible
	result, err := customSection.FieldHandlers[0].FieldValueChange("new value")
	if err != nil {
		t.Errorf("FieldValueChange returned error: %v", err)
	}

	if result != "Value updated to new value" {
		t.Errorf("Expected 'Value updated to new value', got '%s'", result)
	}
}

func TestMultipleTabSections(t *testing.T) {
	// Test adding multiple tab sections
	section1 := TabSection{Title: "Tab1"}
	section2 := TabSection{Title: "Tab2"}

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
