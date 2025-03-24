package devtui

import (
	"fmt"
	"strings"

	"github.com/cdvelop/tinystring"
	"github.com/charmbracelet/lipgloss"
)

// calculateInputWidths returns the widths for label and value parts of an input field
func (h *DevTUI) calculateInputWidths(fieldName string) (labelWidth, valueWidth int) {
	// Usar el ancho estándar de etiquetas definido en el estilo
	labelWidth = h.labelWidth

	// Calcular el ancho disponible para el valor (con padding)
	horizontalPadding := 1 // Este valor debe ser consistente con los estilos definidos

	// Restar del viewport width:
	// - ancho de la etiqueta
	// - ancho del indicador de scroll (típicamente 4 caracteres para "100%")
	// - padding horizontal para los elementos
	scrollInfoWidth := 4 // "100%" típicamente

	// Ajustar para que quepa dentro del viewport menos los otros elementos
	valueWidth = h.viewport.Width - labelWidth - scrollInfoWidth - 2*horizontalPadding

	// Asegurar un mínimo razonable para el valor
	if valueWidth < 10 {
		valueWidth = 10
	}

	return labelWidth, valueWidth
}

// footerView renderiza la vista del footer
// Si hay campos activos, muestra el campo actual como input
// Si no hay campos, muestra una barra de desplazamiento estándar

func (h *DevTUI) footerView() string {
	// Si hay campos disponibles, mostrar el input (independiente de si estamos en modo edición)
	if len(h.tabSections[h.activeTab].fieldHandlers) > 0 {
		return h.renderFooterInput()
	}

	// Si no hay campos, mostrar scrollbar estándar
	info := h.renderScrollInfo()
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// renderScrollInfo returns the formatted scroll percentage
func (h *DevTUI) renderScrollInfo() string {
	percentValue := h.viewport.ScrollPercent() * 100
	percentText := fmt.Sprintf("%.0f%%", percentValue)
	return h.footerInfoStyle.Render(percentText)
}

// renderFooterInput renderiza un campo de entrada en el footer
// Si el campo es editable y estamos en modo edición, muestra un cursor en la posición actual
func (h *DevTUI) renderFooterInput() string {
	// Obtener el campo activo
	tabSection := &h.tabSections[h.activeTab]

	// Verificar que el índice activo esté en rango
	if tabSection.indexActiveEditField >= len(tabSection.fieldHandlers) {
		tabSection.indexActiveEditField = 0 // Reiniciar a 0 si está fuera de rango
	}

	field := &tabSection.fieldHandlers[tabSection.indexActiveEditField]

	// Usar el ancho estándar de etiquetas definido en el estilo
	labelWidth := h.labelWidth

	// Obtener el padding utilizado en el header/footer para mantener consistencia
	horizontalPadding := 1 // Este valor viene del Padding(0, 1) en headerTitleStyle

	// Truncar la etiqueta si es necesario y añadir ":" al final
	labelText := tinystring.Convert(field.Name()).Truncate(labelWidth-1, 0).String() + ":"

	// Aplicar el estilo base para garantizar un ancho fijo
	fixedWidthLabel := h.labelStyle.Render(labelText)

	// Aplicar el estilo visual (colores) manteniendo el ancho fijo
	paddedLabel := h.headerTitleStyle.Render(fixedWidthLabel)

	// Obtener el indicador de porcentaje con el estilo actual
	info := h.renderScrollInfo()

	// OR if you need truncation:
	labelText = tinystring.Convert(field.Name()).Truncate(labelWidth-1, 0).String()
	valueWidth, _ := h.calculateInputWidths(labelText)

	var showCursor bool
	// Preparar el valor del campo
	valueText := field.Value()
	// solo si estamos en modo edición y el campo es editable
	if h.editModeActivated && field.Editable() {
		showCursor = true
		valueText = field.tempEditValue
	}

	// Definir el estilo para el valor del campo
	inputValueStyle := lipgloss.NewStyle().
		Width(valueWidth).
		Padding(0, horizontalPadding). // Añadir padding consistente
		Background(lipgloss.Color(h.Lowlight)).
		Foreground(lipgloss.Color(h.Background))

	// si el campo no es editable cambiar el color del fondo y del texto
	if !field.Editable() {
		inputValueStyle = inputValueStyle.Background(lipgloss.Color(h.ForeGround)).
			Foreground(lipgloss.Color(h.Background))
	}

	// Añadir cursor si corresponde
	if showCursor {
		// Si está en modo edición, cambiar el color del texto a ForeGround
		inputValueStyle = inputValueStyle.Foreground(lipgloss.Color(h.ForeGround))
		// Asegurar que el cursor está dentro de los límites
		runes := []rune(field.tempEditValue)
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
			valueText = field.tempEditValue + "▋"
		}
	}

	// Renderizar el valor con el estilo adecuado
	styledValue := inputValueStyle.Render(valueText)

	// Crear un estilo para el espacio entre elementos
	spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")

	// Unir todos los componentes horizontalmente con espacios consistentes
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		paddedLabel,
		spacerStyle, // Espacio entre label y value
		styledValue,
		spacerStyle, // Espacio entre value e info
		info,
	)
}

// renderInputField renders an input field with current value and cursor if editable
func (h *DevTUI) renderInputField(field *fieldHandler) string {
	// Necesitamos el ancho de la etiqueta y del valor
	labelWidth, valueWidth := h.calculateInputWidths(field.Name())

	// Preparar la etiqueta con el ancho fijo
	labelText := tinystring.Convert(field.Name()).Truncate(labelWidth, 0).String()
	label := h.inputLabelStyle.Render(
		lipgloss.PlaceHorizontal(
			labelWidth,
			lipgloss.Left,
			labelText,
		),
	)

	// Decidir qué valor mostrar (original o temporal)
	var valueText string
	var cursor int

	if field.tempEditValue != "" {
		valueText = field.tempEditValue
		cursor = field.cursor
	} else {
		valueText = field.Value()
		cursor = field.cursor
	}

	// Truncar si es necesario
	valueStr := tinystring.Convert(valueText).Truncate(valueWidth, 0).String()

	// Añadir cursor si estamos en modo edición y el campo es editable
	var renderedValue string
	if field.Editable() {
		// Limitar la posición del cursor al rango válido
		runeValue := []rune(valueStr)
		if cursor > len(runeValue) {
			cursor = len(runeValue)
		}

		// Dividir el texto en la posición del cursor
		var beforeCursor, afterCursor string
		if cursor <= len(runeValue) && cursor >= 0 {
			beforeCursor = string(runeValue[:cursor])
			if cursor < len(runeValue) {
				afterCursor = string(runeValue[cursor:])
			}
		} else {
			beforeCursor = valueStr
		}

		// Renderizar con cursor
		renderedValue = h.inputValueStyle.Render(beforeCursor) +
			h.cursorStyle.Render("▋") +
			h.inputValueStyle.Render(afterCursor)
	} else {
		renderedValue = h.inputValueStyle.Render(valueStr)
	}

	// Combinar ambas partes
	return lipgloss.JoinHorizontal(lipgloss.Top, label, renderedValue)
}
