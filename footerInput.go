package devtui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// footerView renderiza la vista del footer
// Si hay campos activos, muestra el campo actual como input
// Si no hay campos, muestra una barra de desplazamiento estándar

func (h *DevTUI) footerView() string {
	// Si hay campos disponibles, mostrar el input (independiente de si estamos en modo edición)
	if len(h.tabSections[h.activeTab].FieldHandlers) > 0 {
		return h.renderFooterInput()
	}

	// Si no hay campos, mostrar scrollbar estándar
	info := h.renderScrollInfo()
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// renderScrollInfo returns the formatted scroll percentage
func (h *DevTUI) renderScrollInfo() string {
	return h.footerInfoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
}

// renderFooterInput renderiza un campo de entrada en el footer
// Si el campo es editable y estamos en modo edición, muestra un cursor en la posición actual
func (h *DevTUI) renderFooterInput() string {
	// Obtener el campo activo
	tabSection := &h.tabSections[h.activeTab]
	if len(tabSection.FieldHandlers) == 0 {
		return "" // No hay campos disponibles
	}

	// Verificar que el índice activo esté en rango
	if tabSection.indexActiveEditField >= len(tabSection.FieldHandlers) {
		tabSection.indexActiveEditField = 0 // Reiniciar a 0 si está fuera de rango
	}

	field := &tabSection.FieldHandlers[tabSection.indexActiveEditField]

	// Calcular el ancho para la etiqueta (15% del viewport)
	labelWidth := int(float64(h.viewport.Width) * 0.15)

	// Truncar la etiqueta si es necesario (aseguramos que se trunca en una sola línea)
	labelText := field.Label
	if len(labelText) > labelWidth-1 { // -1 para dejar espacio para el ":"
		if labelWidth > 3 {
			labelText = labelText[:labelWidth-3-1] + "..." // -1 para el ":"
		} else {
			labelText = labelText[:max(1, labelWidth-1)]
		}
	}

	// Crear un estilo para la etiqueta igual al del header
	// Usar el headerTitleStyle como base para mantener consistencia visual

	// Formatear la etiqueta usando el estilo del header
	paddedLabel := h.headerTitleStyle.Render(labelText + ":")

	// Obtener el indicador de porcentaje con el estilo actual
	info := h.renderScrollInfo()
	// Calcular el espacio disponible para el valor del campo
	infoWidth := lipgloss.Width(info)
	valueWidth := h.viewport.Width - lipgloss.Width(paddedLabel) - infoWidth

	// Verificar si se debe mostrar el cursor (solo si estamos en modo edición y el campo es editable)
	showCursor := h.tabEditingConfig && field.Editable

	// Preparar el valor del campo
	valueText := field.Value

	// Añadir cursor si corresponde
	if showCursor {
		// Asegurar que el cursor está dentro de los límites
		runes := []rune(field.Value)
		if field.cursor < 0 {
			field.cursor = 0
		}
		if field.cursor > len(runes) {
			field.cursor = len(runes)
		}

		// Insertar el cursor en la posición correcta usando slices de runes para manejar
		// correctamente caracteres multibyte
		if field.cursor <= len(runes) {
			beforeCursor := string(runes[:field.cursor])
			afterCursor := string(runes[field.cursor:])
			valueText = beforeCursor + "▋" + afterCursor
		} else {
			valueText = field.Value + "▋"
		}
	}

	// Definir el estilo para el valor del campo
	valueStyle := lipgloss.NewStyle().
		Width(valueWidth).
		Background(lipgloss.Color(h.Lowlight)).
		Foreground(lipgloss.Color(h.ForeGround))

	// Si está en modo edición, cambiar el color del texto a Highlight
	if showCursor {
		valueStyle = valueStyle.Foreground(lipgloss.Color(h.Highlight))
	}

	// Renderizar el valor con el estilo adecuado
	styledValue := valueStyle.Render(valueText)

	// Unir todos los componentes horizontalmente (sin saltos de línea)
	return lipgloss.JoinHorizontal(lipgloss.Left, paddedLabel, styledValue, info)
}

// max devuelve el máximo entre dos enteros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
