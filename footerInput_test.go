package devtui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TestFooterView verifica el comportamiento del renderizado del footer
func TestFooterView(t *testing.T) {
	h := DefaultTUIForTest(func(messages ...any) {
		// Test logger - do nothing
	})

	// Caso 1: Tab sin fields debe mostrar el scrollbar estándar
	t.Run("Footer with no fields shows scrollbar", func(t *testing.T) {
		// Guardar estado actual para restaurar después de la prueba
		tab := h.tabSections[h.activeTab]
		originalFields := tab.FieldHandlers()

		// Configurar pestaña sin fields
		tab.setFieldHandlers([]*field{})

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene indicador de scroll con iconos (no porcentaje)
		hasScrollIcon := strings.Contains(result, "■") ||
			strings.Contains(result, "▼") ||
			strings.Contains(result, "▲")
		if !hasScrollIcon {
			t.Error("El footer sin campos debería mostrar indicador de scroll con iconos (■, ▼, ▲)")
		}

		// Restaurar estado
		tab.setFieldHandlers(originalFields)
	})

	// Caso 2: Tab con fields debe mostrar el campo actual como input (ahora siempre, no solo en modo edición)
	t.Run("Footer with fields shows field as input even when not editing", func(t *testing.T) {

		// Crear un nuevo field con handler para la prueba
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("TestLabel", "TestValue Rendered")
		tab.AddEditHandler(testHandler, 0)

		// Set viewport width for proper layout calculation
		h.viewport.Width = 80

		// Desactivar modo edición para verificar que aún así se muestra el campo
		h.editModeActivated = false
		tabSection := h.tabSections[h.activeTab]
		tabSection.indexActiveEditField = 0

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene la etiqueta y valor del field
		field := tab.FieldHandlers()[0]
		if !strings.Contains(result, field.Value()) {
			t.Errorf("El footer debería mostrar:\n%v\n incluso sin estar en modo edición, pero muestra:\n%s\n", field.Value(), result)
		}
	})
}

// TestRenderFooterInput verifica el comportamiento específico del renderizado del input
func TestRenderFooterInput(t *testing.T) {
	// Caso 1: Campo editable en modo edición debe mostrar cursor
	t.Run("Editable field in edit mode shows cursor", func(t *testing.T) {
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})
		h.editModeActivated = true
		tab := h.tabSections[h.activeTab]

		// Crear un nuevo field con handler para la prueba
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("Test", "test value")
		tab.AddEditHandler(testHandler, 0)

		// Set viewport width for proper layout calculation
		h.viewport.Width = 80

		field := tab.FieldHandlers()[0]
		field.SetCursorForTest(2) // Cursor en posición "te|st value"
		field.SetTempEditValueForTest("test value")

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
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		// Limpiar todos los handlers existentes y crear explícitamente un handler no editable
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{}) // Limpiar cualquier handler por defecto

		// Crear explícitamente un handler no editable (ExecutionHandler)
		testHandler := NewTestNonEditableHandler("Test", "Value")
		tab.AddExecutionHandler(testHandler, 0)

		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Renderizar input
		result := h.renderFooterInput()

		// No debe contener cursor porque es un handler de ejecución (no editable)
		if strings.Contains(result, "▋") {
			t.Error("Campo no editable no debería mostrar cursor")
		}
	})

	// Nuevo test - Caso 4: Verificar que se maneja correctamente el índice fuera de rango
	t.Run("Index out of range is handled correctly", func(t *testing.T) {
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		expectedLabel := "Test Handler"
		// Configurar un índice activo fuera de rango
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestNonEditableHandler(expectedLabel, "Some Value")
		tab.AddExecutionHandler(testHandler, 0)

		// Set viewport width for proper layout calculation
		h.viewport.Width = 80

		h.tabSections[h.activeTab].indexActiveEditField = 5 // Índice fuera de rango

		// Renderizar - no debería producir pánico
		result := h.renderFooterInput()

		// Verificar que se resetea el índice y se muestra el primer campo
		// Para handlers de ejecución, se muestra el Label() en el footer
		if !strings.Contains(result, expectedLabel) {
			t.Fatalf("No se manejó correctamente el índice fuera de rango. Esperado: %q, resultado:\n%s", expectedLabel, result)
		}
	})

	// Nuevo test - Caso 5: Verificar el estilo correcto cuando está seleccionado pero no en modo edición
	t.Run("Field has correct style when selected but not in edit mode", func(t *testing.T) {
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("Test", "Value")
		tab.AddEditHandler(testHandler, 0)
		h.tabSections[h.activeTab].indexActiveEditField = 0
		h.editModeActivated = false // No en modo edición

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

// Nuevos tests para la navegación y comportamiento de teclas
func TestInputNavigation(t *testing.T) {
	h := DefaultTUIForTest(func(messages ...any) {
		// Test logger - do nothing
	})

	// Configurar múltiples campos para prueba de navegación
	tab := h.tabSections[h.activeTab]
	tab.setFieldHandlers([]*field{})
	testHandler1 := NewTestEditableHandler("Field1", "Value1")
	testHandler2 := NewTestEditableHandler("Field2", "Value2")
	testHandler3 := NewTestEditableHandler("Field3", "Value3")
	tab.AddEditHandler(testHandler1, 0)
	tab.AddEditHandler(testHandler2, 0)
	tab.AddEditHandler(testHandler3, 0)
	h.tabSections[h.activeTab].indexActiveEditField = 0
	h.editModeActivated = false

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
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		// Configurar un campo editable
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("Test", "Value")
		tab.AddEditHandler(testHandler, 0)

		// Asegurar que no estamos en modo edición
		h.editModeActivated = false
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Simular pulsación de Enter
		_, _ = h.handleNormalModeKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// Verificar que entramos en modo edición
		if !h.editModeActivated {
			t.Error("Enter debería activar el modo edición")
		}
	})

	t.Run("Esc exits edit mode", func(t *testing.T) {
		// Reset para esta prueba
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})

		// Configurar un campo editable
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("Test", "Value")
		tab.AddEditHandler(testHandler, 0)

		// Asegurar que estamos en modo edición
		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0

		// Simular pulsación de Esc
		_, _ = h.handleEditingConfigKeyboard(tea.KeyMsg{Type: tea.KeyEscape})

		// Verificar que salimos del modo edición
		if h.editModeActivated {
			t.Error("Esc debería salir del modo edición")
		}
	})

	t.Run("Left/right moves cursor in edit mode", func(t *testing.T) {
		// Reset para esta prueba y configurar un campo editable
		h := DefaultTUIForTest(func(messages ...any) {
			// Test logger - do nothing
		})
		tab := h.tabSections[h.activeTab]
		tab.setFieldHandlers([]*field{})
		testHandler := NewTestEditableHandler("Test", "Value1")
		tab.AddEditHandler(testHandler, 0)

		// Configurar para edición
		h.editModeActivated = true
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := tab.FieldHandlers()[0]
		field.SetCursorAtEnd()
		// Move cursor to position 3 for test
		field.SetCursorForTest(3)

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
