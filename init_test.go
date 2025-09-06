package devtui

import (
	"testing"
)

func TestNewTUI(t *testing.T) {
	// Test configuration with default tabs
	config := &TuiConfig{
		ExitChan: make(chan bool),
		Color:    &ColorPalette{}, // Usando un ColorPalette vacío
		Logger: func(messages ...any) {
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

func TestMultipleTabSections(t *testing.T) {
	// Test that NewTUI correctly adds multiple tab sections
	config := &TuiConfig{
		Color: &ColorPalette{},
	}

	tui := NewTUI(config)

	// Enable test mode for synchronous execution
	tui.SetTestMode(true)

	// Create two more sections using NewTabSection
	tui.NewTabSection("Tab1", "Description 1")
	tui.NewTabSection("Tab2", "Description 2")

	totalSections := len(tui.tabSections)

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
		Color: &ColorPalette{},
	}

	tui := NewTUI(config)

	// We can only test that the TUI was created successfully
	if tui == nil {
		t.Error("Failed to create TUI with channel functionality")
	}
}
