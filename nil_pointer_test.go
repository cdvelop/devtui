package devtui

import (
	"testing"

	"github.com/cdvelop/messagetype"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNilPointerDereference(t *testing.T) {
	// Crear una configuración de TUI
	config := &TuiConfig{
		AppName:       "Test TUI",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			// No hacer nada en el test
		},
	}

	// Crear el TUI
	tui := NewTUI(config)

	// Forzar que el campo id sea nil para replicar el error
	tui.id = nil

	// Crear una sección de pestañas
	section := tui.NewTabSection("Test Section", "Test Description")

	// Intentar agregar contenido nuevo - esto YA NO debería causar panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("No se esperaba ningún panic después del fix, pero se obtuvo: %v", r)
		}
	}()

	// Esto ya no debería causar panic
	section.addNewContent(messagetype.Success, "Test content")

	// Verificar que el contenido se agregó correctamente con ID de fallback
	if len(section.tabContents) != 1 {
		t.Errorf("Se esperaba 1 elemento en tabContents, pero se obtuvo %d", len(section.tabContents))
	}

	if section.tabContents[0].Content != "Test content" {
		t.Errorf("Se esperaba 'Test content', pero se obtuvo '%s'", section.tabContents[0].Content)
	}

	if section.tabContents[0].Id != "temp-id" {
		t.Errorf("Se esperaba 'temp-id' como fallback ID, pero se obtuvo '%s'", section.tabContents[0].Id)
	}
}

func TestKeyboardHandlingWithNilID(t *testing.T) {
	// Crear una configuración de TUI
	config := &TuiConfig{
		AppName:       "Test TUI",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			t.Logf("Log: %v", messages)
		},
	}

	// Crear el TUI
	tui := NewTUI(config)

	// Forzar que el campo id sea nil para replicar el error
	tui.id = nil

	// Crear una sección de pestañas con un campo no editable
	section := tui.NewTabSection("Test Section", "Test Description")
	section.NewField("Test Field", "test value", false, func(value any) (string, error) {
		return "modified: " + value.(string), nil
	})

	// Simular presionar Enter en un campo no editable - esto YA NO debería causar panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("No se esperaba ningún panic después del fix, pero se obtuvo: %v", r)
		}
	}()

	// Simular el manejo de teclado que causaba el problema
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	tui.HandleKeyboard(keyMsg)

	// Verificar que se agregó contenido con ID de fallback
	if len(section.tabContents) > 0 && section.tabContents[0].Id != "temp-id" {
		t.Errorf("Se esperaba 'temp-id' como fallback ID, pero se obtuvo '%s'", section.tabContents[0].Id)
	}
}

func TestFixedNilPointerHandling(t *testing.T) {
	// Crear una configuración de TUI
	config := &TuiConfig{
		AppName:       "Test TUI",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			// No hacer nada en el test
		},
	}

	// Crear el TUI
	tui := NewTUI(config)

	// Verificar que el id no sea nil después de la inicialización
	if tui.id == nil {
		t.Error("El campo id no debería ser nil después de NewTUI")
	}

	// Crear una sección de pestañas
	section := tui.NewTabSection("Test Section", "Test Description")

	// Agregar contenido nuevo - esto no debería causar panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("No se esperaba ningún panic, pero se obtuvo: %v", r)
		}
	}()

	// Esto debería funcionar sin problemas
	section.addNewContent(messagetype.Success, "Test content")

	// Verificar que el contenido se agregó correctamente
	if len(section.tabContents) != 1 {
		t.Errorf("Se esperaba 1 elemento en tabContents, pero se obtuvo %d", len(section.tabContents))
	}

	if section.tabContents[0].Content != "Test content" {
		t.Errorf("Se esperaba 'Test content', pero se obtuvo '%s'", section.tabContents[0].Content)
	}
}
