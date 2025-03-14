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
	// Si estamos en modo edición y hay campos disponibles, mostrar el input
	if h.tabEditingConfig && len(h.tabSections[h.activeTab].FieldHandlers) > 0 {
		return h.renderFooterInput()
	}

	// Si no hay campos o no estamos editando, mostrar scrollbar estándar
	info := h.footerInfoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// renderFooterInput renderiza un campo de entrada en el footer
// Si el campo es editable, muestra un cursor en la posición actual
func (h *DevTUI) renderFooterInput() string {
	// Obtener el campo activo
	tabSection := h.tabSections[h.activeTab]
	if tabSection.indexActiveEditField >= len(tabSection.FieldHandlers) {
		return "" // Protección contra índices fuera de rango
	}

	field := tabSection.FieldHandlers[tabSection.indexActiveEditField]

	// Construir la representación del campo
	var displayContent string

	// Si el campo es editable y estamos en modo edición, mostrar cursor
	if field.Editable {
		// Insertar el cursor en la posición correcta dentro del valor
		if field.cursor > len(field.Value) {
			field.cursor = len(field.Value)
		}

		beforeCursor := field.Value[:field.cursor]
		afterCursor := ""
		if field.cursor < len(field.Value) {
			afterCursor = field.Value[field.cursor:]
		}

		// Construir el texto con el cursor
		displayContent = fmt.Sprintf("%s: %s▋%s", field.Label, beforeCursor, afterCursor)
	} else {
		// Campo no editable, mostrar sin cursor
		displayContent = fmt.Sprintf("%s: %s", field.Label, field.Value)
	}

	// Aplicar estilo
	styledContent := h.footerInputStyle.Render(displayContent)

	// Calcular ancho visual del contenido estilizado
	contentWidth := lipgloss.Width(styledContent)

	// Calcular padding para centrado exacto
	padding := (h.viewport.Width - contentWidth) / 2
	leftPadding := strings.Repeat(" ", max(0, padding))
	rightPadding := strings.Repeat(" ", max(0, h.viewport.Width-contentWidth-padding))

	return leftPadding + styledContent + rightPadding

	return lipgloss.JoinHorizontal(lipgloss.Center, line, styledContent)
}

// max devuelve el máximo entre dos enteros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *DevTUI) renderLeftSectionForm() string {
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
