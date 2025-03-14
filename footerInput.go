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
	info := h.footerInfoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
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

	// Construir la representación del campo
	line := fmt.Sprintf("%s: %s", field.Label, field.Value)

	// Aplicar el estilo según el estado del campo
	var styledContent string

	// Verificar si se debe mostrar el cursor (solo si estamos en modo edición y el campo es editable)
	showCursor := h.tabEditingConfig && field.Editable

	if showCursor {
		// Asegurar que el cursor está dentro de los límites
		if field.cursor < 0 {
			field.cursor = 0
		}
		if field.cursor > len(field.Value) {
			field.cursor = len(field.Value)
		}

		// Calcular la posición del cursor en la línea completa (etiqueta + valor)
		cursorPos := field.cursor + len(field.Label) + 2 // +2 por ": "

		// Validar que la posición del cursor no exceda la longitud de la línea
		if cursorPos <= len(line) {
			line = line[:cursorPos] + "▋" + line[cursorPos:]
		}

		styledContent = h.fieldEditingStyle.Render(line)
	} else {
		// Campo seleccionado pero no en modo edición o no editable
		styledContent = h.fieldSelectedStyle.Render(line)
	}

	// Calcular el espacio restante a la derecha (asegurando que no sea negativo)
	contentWidth := lipgloss.Width(styledContent)
	remainingWidth := max(0, h.viewport.Width-contentWidth-2) // -2 por el padding izquierdo
	rightPadding := strings.Repeat(" ", remainingWidth)

	return lipgloss.JoinHorizontal(lipgloss.Center, styledContent, rightPadding)
}

// max devuelve el máximo entre dos enteros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *DevTUI) exampleRenderSectionForm() string {
	var lines []string

	for indexSection, tabSection := range t.tabSections {

		// break different index
		if indexSection != t.activeTab {
			continue
		}

		for i, field := range tabSection.FieldHandlers {
			line := fmt.Sprintf("%s: %s", field.Label, field.Value)

			if i == tabSection.indexActiveEditField {
				if t.tabEditingConfig {
					cursorPos := field.cursor + len(field.Label) + 2
					line = line[:cursorPos] + "▋" + line[cursorPos:]
					line = t.fieldEditingStyle.Render(line)
				} else {
					line = t.fieldSelectedStyle.Render(line)
				}
			} else {
				line = t.fieldLineStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
