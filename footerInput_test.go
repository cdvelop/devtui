package devtui

import (
	"strings"
	"testing"
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

	// Caso 2: Tab con fields debe mostrar el campo actual como input
	t.Run("Footer with fields shows field as input", func(t *testing.T) {
		// Asegurar que hay al menos un field
		if len(h.tabSections[h.activeTab].FieldHandlers) == 0 {
			t.Skip("Se requiere al menos un campo para esta prueba")
			return
		}
		h.tabEditingConfig = true
		tabSection := &h.tabSections[h.activeTab]

		// Modificar field para la prueba
		field := &tabSection.FieldHandlers[0]
		field.Label = "TestLabel"
		field.Value = "TestValue"
		field.Editable = true
		tabSection.indexActiveEditField = 0

		// Renderizar footer
		result := h.footerView()

		// Verificar que contiene la etiqueta y valor del field
		if !strings.Contains(result, "TestLabel: TestValue") {
			t.Errorf("El footer debería mostrar 'TestLabel: TestValue', pero muestra: %s", result)
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
			t.Skip("Se requiere al menos un campo para esta prueba")
			return
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
			t.Skip("Se requiere al menos un campo para esta prueba")
			return
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

	// Caso 3: Verificar centrado del input en el footer
	t.Run("Footer input is centered", func(t *testing.T) {
		if len(h.tabSections[h.activeTab].FieldHandlers) == 0 {
			t.Skip("Se requiere al menos un campo para esta prueba")
			return
		}

		// Configurar viewport con ancho específico para prueba
		h.viewport.Width = 80
		h.tabSections[h.activeTab].indexActiveEditField = 0
		field := &h.tabSections[h.activeTab].FieldHandlers[0]
		field.Label = "X"
		field.Value = "Y" // Usar etiqueta y valor cortos para simplificar cálculos

		// Renderizar y verificar que hay espacios a ambos lados (centrado)
		result := h.renderFooterInput()

		// Comprobar que empieza y termina con espacio (está centrado)
		if !strings.HasPrefix(result, " ") || !strings.HasSuffix(result, " ") {
			t.Error("El input no parece estar centrado en el footer")
		}
	})
}
