package devtui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TestFooterView verifica el comportamiento del renderizado del footer
func TestFooterView(t *testing.T) {
	h := prepareForTesting()

	// Caso 1: Tab sin fields debe mostrar el scrollbar estándar
	t.Run("Footer with no fields shows scrollbar", func(t *testing.T) {
		// Guardar estado actual para restaurar después de la prueba
		originalFields := h.tabSections[h.activeTab].FieldHandlers

		// Configurar pestaña sin fields
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{}

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene indicador de porcentaje típico del scrollbar
		if !strings.Contains(result, "%") {
			t.Error("El footer sin campos debería mostrar indicador de porcentaje")
		}

		// Restaurar estado
		h.tabSections[h.activeTab].FieldHandlers = originalFields
	})

	// Caso 2: Tab con fields debe mostrar el campo actual como input (ahora siempre, no solo en modo edición)
	t.Run("Footer with fields shows field as input even when not editing", func(t *testing.T) {
		// Asegurar que hay al menos un field
		if len(h.tabSections[h.activeTab].FieldHandlers) == 0 {
			// Crear un campo de prueba
			h.tabSections[h.activeTab].FieldHandlers = append(h.tabSections[h.activeTab].FieldHandlers, FieldHandler{
				Label: "TestLabel",
				Value: "TestValue",
			})
		} else {
			// Modificar field existente para la prueba
			field := &h.tabSections[h.activeTab].FieldHandlers[0]
			field.Label = "TestLabel"
			field.Value = "TestValue"
		}

		// Desactivar modo edición para verificar que aún así se muestra el campo
		h.tabEditingConfig = false
		tabSection := &h.tabSections[h.activeTab]
		tabSection.indexActiveEditField = 0

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene la etiqueta y valor del field
		if !strings.Contains(result, "TestLabel: TestValue") {
			t.Errorf("El footer debería mostrar 'TestLabel: TestValue' incluso sin estar en modo edición, pero muestra: %s", result)
		}
	})
}

// TestRenderFooterInput verifica el comportamiento específico del renderizado del input
func TestRenderFooterInput(t *testing.T) {
	h := prepareForTesting()

	// Caso 1: Campo editable en modo edición debe mostrar cursor
	t.Run("Editable field in edit mode shows cursor", func(t *testing.T) {
		// Configurar campo editable en modo edición
		if len(h.tabSections[h.activeTab].FieldHandlers) == 0 {
			// Crear un campo si no hay ninguno
			h.tabSections[h.activeTab].FieldHandlers = append(h.tabSections[h.activeTab].FieldHandlers, FieldHandler{
				Label:    "Test",
				Value:    "Value",
				Editable: true,
			})
		}

		h.tabEditingConfig = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := &h.tabSections[h.activeTab].FieldHandlers[0]
		field.Label = "Test"
		field.Value = "Value"
		field.Editable = true
		field.cursor = 2 // Cursor en posición "Va|lue"

		// Renderizar input
		result := h.renderFooterInput()

		// Comprobar que se renderiza con cursor
		// El cursor '▋' debe estar después de las primeras dos letras del valor
		if !strings.Contains(result, "Test: Va▋lue") {
			t.Errorf("El cursor no se renderiza correctamente, resultado: %s", result)
		}
	})

	// Caso 2: Campo no editable no debe mostrar cursor
	t.Run("Non-editable field doesn't show cursor", func(t *testing.T) {
		// Configurar campo no editable
		if len(h.tabSections[h.activeTab].FieldHandlers) == 0 {
			// Crear un campo si no hay ninguno
			h.tabSections[h.activeTab].FieldHandlers = append(h.tabSections[h.activeTab].FieldHandlers, FieldHandler{
				Label:    "Test",
				Value:    "Value",
				Editable: false,
			})
		}

		h.tabEditingConfig = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := &h.tabSections[h.activeTab].FieldHandlers[0]
		field.Label = "Test"
		field.Value = "Value"
		field.Editable = false

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
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value"},
		}
		h.tabSections[h.activeTab].indexActiveEditField = 5 // Índice fuera de rango

		// Renderizar - no debería producir pánico
		result := h.renderFooterInput()

		// Verificar que se resetea el índice y se muestra el primer campo
		if !strings.Contains(result, "Test: Value") {
			t.Error("No se manejó correctamente el índice fuera de rango")
		}
	})

	// Nuevo test - Caso 5: Verificar el estilo correcto cuando está seleccionado pero no en modo edición
	t.Run("Field has correct style when selected but not in edit mode", func(t *testing.T) {
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value", Editable: true},
		}
		h.tabSections[h.activeTab].indexActiveEditField = 0
		h.tabEditingConfig = false // No en modo edición

		// El estilo debe ser fieldSelectedStyle en vez de fieldEditingStyle
		originalFieldSelectedStyle := h.fieldSelectedStyle
		originalFieldEditingStyle := h.fieldEditingStyle

		// Modificar temporalmente los estilos para distinguirlos claramente
		h.fieldSelectedStyle = h.fieldSelectedStyle.Background(lipgloss.Color("blue"))
		h.fieldEditingStyle = h.fieldEditingStyle.Background(lipgloss.Color("red"))

		result := h.renderFooterInput()

		// Restaurar estilos originales
		h.fieldSelectedStyle = originalFieldSelectedStyle
		h.fieldEditingStyle = originalFieldEditingStyle

		// Verificar que no contiene el cursor de edición
		if strings.Contains(result, "▋") {
			t.Error("Campo seleccionado pero no en modo edición no debería mostrar cursor")
		}
	})
}

// Nuevo test para verificar el modo automático de edición
func TestAutoEditMode(t *testing.T) {
	h := prepareForTesting()

	t.Run("Auto edit mode activates with single editable field", func(t *testing.T) {
		// Configurar un solo campo editable
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value", Editable: true},
		}
		h.tabEditingConfig = false // Iniciar no en modo edición

		// Llamar al método que verifica si debe activar modo edición automático
		h.checkAutoEditMode()

		// Verificar que se activó el modo edición
		if !h.tabEditingConfig {
			t.Error("El modo edición debería activarse automáticamente con un solo campo editable")
		}
	})

	t.Run("Auto edit mode does not activate with multiple fields", func(t *testing.T) {
		// Configurar múltiples campos
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test1", Value: "Value1", Editable: true},
			{Label: "Test2", Value: "Value2", Editable: true},
		}
		h.tabEditingConfig = false // Iniciar no en modo edición

		// Llamar al método que verifica si debe activar modo edición automático
		h.checkAutoEditMode()

		// Verificar que NO se activó el modo edición
		if h.tabEditingConfig {
			t.Error("El modo edición NO debería activarse automáticamente con múltiples campos")
		}
	})

	t.Run("Auto edit mode does not activate with non-editable field", func(t *testing.T) {
		// Configurar un solo campo NO editable
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value", Editable: false},
		}
		h.tabEditingConfig = false // Iniciar no en modo edición

		// Llamar al método que verifica si debe activar modo edición automático
		h.checkAutoEditMode()

		// Verificar que NO se activó el modo edición
		if h.tabEditingConfig {
			t.Error("El modo edición NO debería activarse automáticamente con un campo no editable")
		}
	})
}

// Nuevos tests para la navegación y comportamiento de teclas
func TestInputNavigation(t *testing.T) {
	h := prepareForTesting()

	// Configurar múltiples campos para prueba de navegación
	h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
		{Label: "Field1", Value: "Value1", Editable: true},
		{Label: "Field2", Value: "Value2", Editable: true},
		{Label: "Field3", Value: "Value3", Editable: true},
	}
	h.tabSections[h.activeTab].indexActiveEditField = 0
	h.tabEditingConfig = false

	t.Run("Right key navigates to next field", func(t *testing.T) {
		// Simular pulsación de tecla derecha
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		// Verificar que nos movimos al siguiente campo
		if h.tabSections[h.activeTab].indexActiveEditField != 1 {
			t.Errorf("La tecla derecha debería navegar al siguiente campo, pero se quedó en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Left key navigates to previous field", func(t *testing.T) {
		// Nos aseguramos de estar en el campo del medio
		h.tabSections[h.activeTab].indexActiveEditField = 1

		// Simular pulsación de tecla izquierda
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		// Verificar que nos movimos al campo anterior
		if h.tabSections[h.activeTab].indexActiveEditField != 0 {
			t.Errorf("La tecla izquierda debería navegar al campo anterior, pero se quedó en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Cyclical navigation wraps around at boundaries", func(t *testing.T) {
		// Ir al primer campo
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Simular pulsación de tecla izquierda (debe ir al último campo)
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		// Verificar que se movió al último campo
		if h.tabSections[h.activeTab].indexActiveEditField != 2 {
			t.Errorf("La navegación cíclica debería ir al último campo, pero está en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}

		// Simular pulsación de tecla derecha (debe volver al primer campo)
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		// Verificar que volvió al primer campo
		if h.tabSections[h.activeTab].indexActiveEditField != 0 {
			t.Errorf("La navegación cíclica debería volver al primer campo, pero está en: %d",
				h.tabSections[h.activeTab].indexActiveEditField)
		}
	})

	t.Run("Enter enters edit mode", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Configurar un campo editable
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value", Editable: true},
		}

		// Asegurar que no estamos en modo edición
		h.tabEditingConfig = false
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Simular pulsación de Enter
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// Verificar que entramos en modo edición
		if !h.tabEditingConfig {
			t.Error("Enter debería activar el modo edición")
		}
	})

	t.Run("Esc exits edit mode", func(t *testing.T) {
		// Reset para esta prueba
		h := prepareForTesting()

		// Configurar un campo editable
		h.tabSections[h.activeTab].FieldHandlers = []FieldHandler{
			{Label: "Test", Value: "Value", Editable: true},
		}

		// Asegurar que estamos en modo edición
		h.tabEditingConfig = true
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Simular pulsación de Esc
		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyEscape})

		// Verificar que salimos del modo edición
		if h.tabEditingConfig {
			t.Error("Esc debería salir del modo edición")
		}
	})

	t.Run("Left/right moves cursor in edit mode", func(t *testing.T) {
		// Configurar para edición
		h.tabEditingConfig = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := &h.tabSections[h.activeTab].FieldHandlers[0]
		field.cursor = 3
		field.Value = "Value1"

		// Simular pulsación de tecla izquierda
		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyLeft})

		// Verificar que el cursor se movió a la izquierda
		if field.cursor != 2 {
			t.Errorf("La tecla izquierda en modo edición debería mover el cursor a la izquierda, pero quedó en: %d",
				field.cursor)
		}

		// Simular pulsación de tecla derecha
		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyRight})

		// Verificar que el cursor volvió a la posición original
		if field.cursor != 3 {
			t.Errorf("La tecla derecha en modo edición debería mover el cursor a la derecha, pero quedó en: %d",
				field.cursor)
		}
	})
}
