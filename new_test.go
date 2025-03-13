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
		LogToFile: func(messageErr string) {
			// Mock function for logging
		},
	}

	tui := NewTUI(config)

	// Check if TUI was created correctly
	if tui == nil {
		t.Fatal("TUI was not created correctly")
	}

	// Check if default tabs were created
	if len(tui.TabSections) != 2 {
		t.Fatalf("Expected 2 default tabs, got %d", len(tui.TabSections))
	}

	// Check if tab titles are correct
	if tui.TabSections[0].Title != "BUILD" {
		t.Errorf("Expected first tab title 'BUILD', got '%s'", tui.TabSections[0].Title)
	}

	if tui.TabSections[1].Title != "DEPLOY" {
		t.Errorf("Expected second tab title 'DEPLOY', got '%s'", tui.TabSections[1].Title)
	}

	// Check if the active tab is set correctly
	if tui.activeTab != config.TabIndexStart {
		t.Errorf("Expected active tab to be %d, got %d", config.TabIndexStart, tui.activeTab)
	}
}

func TestCustomTabs(t *testing.T) {
	// Create a custom configuration with custom tabs
	tabs := []TabSection{
		{
			Title: "CUSTOM1",
			FieldHanlders: []FieldHanlder{
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
		},
	}

	config := &TuiConfig{
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		TabSections:   tabs,
		Color:         &ColorStyle{},
	}

	tui := NewTUI(config)

	// Check if custom tab was set correctly
	if len(tui.TabSections) != 1 {
		t.Fatalf("Expected 1 custom tab, got %d", len(tui.TabSections))
	}

	if tui.TabSections[0].Title != "CUSTOM1" {
		t.Errorf("Expected custom tab title 'CUSTOM1', got '%s'", tui.TabSections[0].Title)
	}

	// Check if section fields were set correctly
	if len(tui.TabSections[0].FieldHanlders) != 1 {
		t.Fatalf("Expected 1 section field, got %d", len(tui.TabSections[0].FieldHanlders))
	}

	field := tui.TabSections[0].FieldHanlders[0]
	if field.Name != "testField" || field.Value != "test value" {
		t.Errorf("Field not set correctly. Expected name 'testField' and value 'test value', got '%s' and '%s'",
			field.Name, field.Value)
	}

	// Test FieldValueChange function
	result, err := field.FieldValueChange("new value")
	if err != nil {
		t.Errorf("FieldValueChange returned error: %v", err)
	}

	if result != "Value updated to new value" {
		t.Errorf("Expected 'Value updated to new value', got '%s'", result)
	}
}

func TestSectionFieldIndexing(t *testing.T) {
	// Create tab with multiple fields
	fields := []FieldHanlder{
		{Name: "field1", Value: "value1"},
		{Name: "field2", Value: "value2"},
		{Name: "field3", Value: "value3"},
	}

	tabs := []TabSection{
		{
			Title:         "TestTab",
			FieldHanlders: fields,
		},
	}

	config := &TuiConfig{
		TabSections: tabs,
		Color:       &ColorStyle{},
	}

	tui := NewTUI(config)

	// Check if field indexes were set correctly
	for i, field := range tui.TabSections[0].FieldHanlders {
		if field.index != i {
			t.Errorf("Field %s: Expected index %d, got %d", field.Name, i, field.index)
		}

		// Check that cursor starts at 0
		if field.cursor != 0 {
			t.Errorf("Field %s: Expected cursor to be 0, got %d", field.Name, field.cursor)
		}
	}
}

func TestTabSectionReferences(t *testing.T) {
	// Verifica que las referencias entre TabSection y DevTUI se establezcan correctamente
	config := &TuiConfig{
		TabSections: []TabSection{
			{Title: "Tab1"},
			{Title: "Tab2"},
		},
		Color: &ColorStyle{},
	}

	tui := NewTUI(config)

	for i, section := range tui.TabSections {
		if section.tui != tui {
			t.Errorf("Tab %d: tui reference not set correctly", i)
		}
	}
}

func TestChannelFunctionality(t *testing.T) {
	// Prueba básica de la creación de canal
	config := &TuiConfig{
		Color: &ColorStyle{},
	}

	tui := NewTUI(config)

	if tui.tabContentsChan == nil {
		t.Error("Tab content channel was not created")
	}

	// Verificamos que el canal tenga la capacidad correcta
	// Nota: Esta prueba es limitada ya que no podemos acceder a la capacidad del canal directamente
	select {
	case tui.tabContentsChan <- tabContent{Content: "Test message"}:
		// El canal aceptó el mensaje
	default:
		t.Error("Channel could not accept message")
	}
}
