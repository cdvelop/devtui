package devtui

import (
	"strings"
	"testing"

	"github.com/cdvelop/messagetype"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Para realizar pruebas en footerInput
type testField struct {
	name       string
	value      string
	editable   bool
	changeFunc func(string) <-chan MessageUpdate
}

func (f *testField) Name() string {
	return f.name
}

func (f *testField) Value() string {
	return f.value
}

func (f *testField) Editable() bool {
	return f.editable
}

func (f *testField) ChangeValue(newValue string) <-chan MessageUpdate {
	if f.changeFunc != nil {
		return f.changeFunc(newValue)
	}

	// Default implementation
	updates := make(chan MessageUpdate)
	go func() {
		defer close(updates)
		// Update the field value
		f.value = newValue
		// Send a success message
		updates <- MessageUpdate{
			Content: "Changed " + f.name + " to " + newValue,
			Type:    messagetype.Success,
		}
	}()
	return updates
}

// TestFooterView verifica el comportamiento del renderizado del footer
func TestFooterView(t *testing.T) {
	h := prepareForTesting()

	// Caso 1: Tab sin fields debe mostrar el scrollbar estándar
	t.Run("Footer with no fields shows scrollbar", func(t *testing.T) {
		// Guardar estado actual para restaurar después de la prueba
		originalFields := h.tabSections[h.activeTab].fieldHandlers

		// Configurar pestaña sin fields
		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{}

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene indicador de porcentaje típico del scrollbar
		if !strings.Contains(result, "%") {
			t.Error("El footer sin campos debería mostrar indicador de porcentaje")
		}

		// Restaurar estado
		h.tabSections[h.activeTab].fieldHandlers = originalFields
	})

	// Caso 2: Tab con fields debe mostrar el campo actual como input (ahora siempre, no solo en modo edición)
	t.Run("Footer with fields shows field as input even when not editing", func(t *testing.T) {

		// Asegurar que hay al menos un field
		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "TestLabel", value: "TestValue Rendered", editable: true}},
		}
		field := &h.tabSections[h.activeTab].fieldHandlers[0]

		// Desactivar modo edición para verificar que aún así se muestra el campo
		h.editModeActivated = false
		tabSection := &h.tabSections[h.activeTab]
		tabSection.indexActiveEditField = 0

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene la etiqueta y valor del field
		if !strings.Contains(result, field.fieldHandlerAdapter.Value()) {
			t.Errorf("El footer debería mostrar: %v incluso sin estar en modo edición, pero muestra: %s",
				field.fieldHandlerAdapter.Value(), result)
		}
	})
}

// TestRenderFooterInput verifica el comportamiento específico del renderizado del input
func TestRenderFooterInput(t *testing.T) {
	h := prepareForTesting()

	// Asegurar que hay al menos un field en la pestaña activa
	h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
		{fieldHandlerAdapter: &testField{name: "TestLabel", value: "TestValue", editable: true}},
	}

	// Caso 1: Campo editable en modo edición debe mostrar cursor
	t.Run("Editable field in edit mode shows cursor", func(t *testing.T) {
		h.editModeActivated = true
		field := &h.tabSections[h.activeTab].fieldHandlers[0]
		field.cursor = 2 // Cursor en posición "Va|lue"
		field.tempEditValue = "test value"

		// Renderizar input
		result := h.renderFooterInput()

		// Verificar que te aparece antes del cursor y st value después del cursor
		cursor := "▋"
		if !strings.Contains(result, "te"+cursor+"st value") {
			t.Errorf("El cursor no se renderiza correctamente en la posición esperada (te▋st value), resultado: %s", result)
		}
	})

	// Caso 2: Campo no editable no debe mostrar cursor
	t.Run("Non-editable field doesn't show cursor", func(t *testing.T) {
		// Configurar campo no editable
		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value", editable: false}},
		}

		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		// field := &h.tabSections[h.activeTab].fieldHandlers[0]

		// Renderizar input
		result := h.renderFooterInput()

		// No debe contener cursor
		if strings.Contains(result, "▋") {
			t.Error("Campo no editable no debería mostrar cursor")
		}
	})

	// Nuevo test - Caso 4: Verificar que se maneja correctamente el índice fuera de rango
	t.Run("Index out of range is handled correctly", func(t *testing.T) {
		// Configurar un índice activo fuera de rango
		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value index Success", editable: true}},
		}
		h.tabSections[h.activeTab].indexActiveEditField = 5 // Índice fuera de rango

		// Renderizar - no debería producir pánico
		result := h.renderFooterInput()

		// Verificar que se resetea el índice y se muestra el primer campo
		if !strings.Contains(result, "Value index Success") {
			t.Fatalf("No se manejó correctamente el índice fuera de rango result: %s", result)
		}
	})

	// Nuevo test - Caso 5: Verificar el estilo correcto cuando está seleccionado pero no en modo edición
	t.Run("Field has correct style when selected but not in edit mode", func(t *testing.T) {
		h := prepareForTesting()

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value", editable: true}},
		}
		h.tabSections[h.activeTab].indexActiveEditField = 0
		h.editModeActivated = false // No en modo edición

		originalFieldSelectedStyle := h.fieldSelectedStyle
		originalFieldEditingStyle := h.fieldEditingStyle

		h.fieldSelectedStyle = h.fieldSelectedStyle.Background(lipgloss.Color("blue"))
		h.fieldEditingStyle = h.fieldEditingStyle.Background(lipgloss.Color("red"))

		result := h.renderFooterInput()

		h.fieldSelectedStyle = originalFieldSelectedStyle
		h.fieldEditingStyle = originalFieldEditingStyle

		// En modo no edición no debería mostrar cursor
		if strings.Contains(result, "▋") {
			t.Error("Campo seleccionado pero no en modo edición no debería mostrar cursor")
		}
	})
}

// Nuevo test para verificar el modo automático de edición
func TestAutoEditMode(t *testing.T) {
	h := prepareForTesting()

	t.Run("Auto edit mode activates with single editable field", func(t *testing.T) {

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value", editable: true}},
		}
		h.editModeActivated = false

		h.checkAutoEditMode()

		if !h.editModeActivated {
			t.Error("El modo edición debería activarse automáticamente con un solo campo editable")
		}
	})

	t.Run("Auto edit mode does not activate with multiple fields", func(t *testing.T) {

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{
				fieldHandlerAdapter: &testField{
					name:     "Test1",
					value:    "Value1",
					editable: true,
				},
			},
			{
				fieldHandlerAdapter: &testField{
					name:     "Test2",
					value:    "Value2",
					editable: true,
				},
			},
		}
		h.editModeActivated = false

		h.checkAutoEditMode()

		if h.editModeActivated {
			t.Error("El modo edición NO debería activarse automáticamente con múltiples campos")
		}
	})

	t.Run("Auto edit mode does not activate with non-editable field", func(t *testing.T) {

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{
				fieldHandlerAdapter: &testField{
					name:     "Test",
					value:    "Value",
					editable: false,
				},
			},
		}
		h.editModeActivated = false

		h.checkAutoEditMode()

		if h.editModeActivated {
			t.Error("El modo edición NO debería activarse automáticamente con un campo no editable")
		}
	})
}

// Nuevos tests para la navegación y comportamiento de teclas
func TestInputNavigation(t *testing.T) {
	h := prepareForTesting()

	h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
		{
			fieldHandlerAdapter: &testField{
				name:     "Field1",
				value:    "Value1",
				editable: true,
			},
		},
		{
			fieldHandlerAdapter: &testField{
				name:     "Field2",
				value:    "Value2",
				editable: true,
			},
		},
		{
			fieldHandlerAdapter: &testField{
				name:     "Field3",
				value:    "Value3",
				editable: true,
			},
		},
	}
	h.tabSections[h.activeTab].indexActiveEditField = 0
	h.editModeActivated = false

	t.Run("Right key navigates to next field", func(t *testing.T) {

		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		if h.tabSections[h.activeTab].indexActiveEditField != 1 {
			t.Errorf("La tecla derecha debería navegar al siguiente campo, pero se quedó en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Left key navigates to previous field", func(t *testing.T) {

		h.tabSections[h.activeTab].indexActiveEditField = 1

		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		if h.tabSections[h.activeTab].indexActiveEditField != 0 {
			t.Errorf("La tecla izquierda debería navegar al campo anterior, pero se quedó en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Cyclical navigation wraps around at boundaries", func(t *testing.T) {

		h.tabSections[h.activeTab].indexActiveEditField = 0

		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		if h.tabSections[h.activeTab].indexActiveEditField != 2 {
			t.Errorf("La navegación cíclica debería ir al último campo, pero está en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}

		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		if h.tabSections[h.activeTab].indexActiveEditField != 0 {
			t.Errorf("La navegación cíclica debería volver al primer campo, pero está en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Enter enters edit mode", func(t *testing.T) {

		h := prepareForTesting()

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value", editable: true}},
		}

		h.editModeActivated = false
		h.tabSections[h.activeTab].indexActiveEditField = 0

		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		if !h.editModeActivated {
			t.Error("Enter debería activar el modo edición")
		}
	})

	t.Run("Esc exits edit mode", func(t *testing.T) {

		h := prepareForTesting()

		h.tabSections[h.activeTab].fieldHandlers = []fieldHandler{
			{fieldHandlerAdapter: &testField{name: "Test", value: "Value", editable: true}},
		}

		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0

		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyEscape})

		if h.editModeActivated {
			t.Error("Esc debería salir del modo edición")
		}
	})

	t.Run("Left/right moves cursor in edit mode", func(t *testing.T) {

		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := &h.tabSections[h.activeTab].fieldHandlers[0]
		field.cursor = 3
		field.fieldHandlerAdapter = &testField{
			name:     "Field1",
			value:    "Value1",
			editable: true,
		}
		field.tempEditValue = "Value1"

		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		if field.cursor != 2 {
			t.Errorf("La tecla izquierda en modo edición debería mover el cursor a la izquierda, pero quedó en: %d",
				field.cursor)
		}

		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		if field.cursor != 3 {
			t.Errorf("La tecla derecha en modo edición debería mover el cursor a la derecha, pero quedó en: %d",
				field.cursor)
		}
	})
}
