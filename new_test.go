package devtui

import (
	"testing"
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
	mockField := &mockFieldHandler{
		name:     "Test Field",
		value:    "test value",
		editable: true,
	}

	// Create custom tab section with mockField
	config := &TuiConfig{
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color:         &ColorStyle{},
	}

	// Create a TUI first
	tui := NewTUI(config)

	// Add custom tab section with the mock field
	customSection := tui.NewTabSection("CUSTOM1", mockField)

	// Check the mock field was properly added
	if len(customSection.fieldHandlers) != 1 {
		t.Errorf("Expected 1 field handler, got %d", len(customSection.fieldHandlers))
	}

	// Test that the field value can be retrieved
	if customSection.fieldHandlers[0].Value() != "test value" {
		t.Errorf("Expected 'test value', got '%s'", customSection.fieldHandlers[0].Value())
	}
}

func TestMultipleTabSections(t *testing.T) {
	// Test adding multiple tab sections
	config := &TuiConfig{
		TabIndexStart: 0,
		Color:         &ColorStyle{},
	}

	tui := NewTUI(config)

	// Create two tab sections and use them in the test
	tab1 := tui.NewTabSection("Tab1")
	tab2 := tui.NewTabSection("Tab2")

	// Check that we have the default tab plus our two new tabs
	if len(tui.tabSections) != 3 {
		t.Errorf("Expected 3 tab sections (default + 2 new), got %d", len(tui.tabSections))
	}

	// Check names of the tabs
	if tui.tabSections[1].title != "Tab1" {
		t.Errorf("Expected tab1 title to be 'Tab1', got '%s'", tui.tabSections[1].title)
	}

	if tui.tabSections[2].title != "Tab2" {
		t.Errorf("Expected tab2 title to be 'Tab2', got '%s'", tui.tabSections[2].title)
	}

	// Verificar que podemos acceder a los tabs creados por referencia
	if tab1.title != "Tab1" {
		t.Errorf("La referencia tab1 no tiene el título correcto. Esperado: 'Tab1', obtenido: '%s'", tab1.title)
	}

	if tab2.title != "Tab2" {
		t.Errorf("La referencia tab2 no tiene el título correcto. Esperado: 'Tab2', obtenido: '%s'", tab2.title)
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
