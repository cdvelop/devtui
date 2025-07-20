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
		LogToFile: func(messages ...any) {
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

func TestCustomTabs(t *testing.T) { // Create a custom configuration with custom tabs
	customSection := NewTUI(&TuiConfig{}).NewTabSection("CUSTOM1", "custom footer")

	// Create handler for testing
	testHandler := NewTestFieldHandler("Test Field", "test value", true, func(newValue any) (string, error) {
		strValue := newValue.(string)
		return "Value updated to " + strValue, nil
	})

	customSection.NewField(testHandler)

	config := &TuiConfig{
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color:         &ColorStyle{},
	}

	// Add custom tab section
	NewTUI(config).AddTabSections(customSection)

	// Since internal fields are not accessible in real usage,
	// we can only test that the TUI was modified successfully
}

func TestMultipleTabSections(t *testing.T) {
	// Test that NewTUI correctly adds multiple tab sections
	config := &TuiConfig{
		TabIndexStart: 0,
		Color:         &ColorStyle{},
		TestMode:      true, // Evitar mensaje automático de shortcuts
	}

	tui := NewTUI(config)

	// Create two more sections using NewTabSection
	tui.NewTabSection("Tab1", "Description 1")
	tui.NewTabSection("Tab2", "Description 2")

	totalSections := tui.GetTotalTabSections()

	// Expected: 1 (SHORTCUTS) + 2 (Tab1, Tab2) = 3
	expected := 3
	if totalSections != expected {
		t.Errorf("Expected %d tab sections, got %d", expected, totalSections)

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
