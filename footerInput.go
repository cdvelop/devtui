package devtui

import (
	"fmt"
	"strings"

	"github.com/cdvelop/tinystring"
	"github.com/charmbracelet/lipgloss"
)

// footerView renderiza la vista del footer
// Si hay campos activos, muestra el campo actual como input
// Si no hay campos, muestra una barra de desplazamiento estándar

func (h *DevTUI) footerView() string {
	// Verificar que haya tabs disponibles
	if len(h.tabSections) == 0 {
		return h.footerInfoStyle.Render("No tabs available")
	}
	if h.activeTab >= len(h.tabSections) {
		h.activeTab = 0
	}

	// Si hay campos disponibles, mostrar el input (independiente de si estamos en modo edición)
	if len(h.tabSections[h.activeTab].FieldHandlers()) > 0 {
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
	tabSection := h.tabSections[h.activeTab]

	// Verificar que el índice activo esté en rango
	fieldHandlers := tabSection.FieldHandlers()
	if tabSection.indexActiveEditField >= len(fieldHandlers) {
		tabSection.indexActiveEditField = 0 // Reiniciar a 0 si está fuera de rango
	}

	field := fieldHandlers[tabSection.indexActiveEditField]
	info := h.renderScrollInfo()
	horizontalPadding := 1

	// Check if this handler uses expanded footer (Display only)
	if field.isDisplayOnly() {
		// Layout for Display: [Label expandido usando resto del espacio] [Scroll%]
		remainingWidth := h.viewport.Width - lipgloss.Width(info) - horizontalPadding
		labelText := tinystring.Convert(field.getExpandedFooterLabel()).Truncate(remainingWidth-1, 0).String()

		// Display: [Label expandido] [Scroll%]
		displayStyle := lipgloss.NewStyle().
			Width(remainingWidth).
			Padding(0, horizontalPadding)
		styledLabel := displayStyle.Render(labelText)

		spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")
		return lipgloss.JoinHorizontal(lipgloss.Left, styledLabel, spacerStyle, info)
	}

	// Normal layout for Edit and Execution handlers: [Scroll%] [Label] [Value]
	labelWidth := h.labelWidth

	// Truncar la etiqueta si es necesario
	labelText := tinystring.Convert(field.handler.Label()).Truncate(labelWidth-1, 0).String()

	// Aplicar el estilo base para garantizar un ancho fijo
	fixedWidthLabel := h.labelStyle.Render(labelText)
	paddedLabel := h.headerTitleStyle.Render(fixedWidthLabel)

	// Calcular ancho para el valor usando el espacio restante
	usedWidth := lipgloss.Width(info) + lipgloss.Width(paddedLabel) + horizontalPadding*2
	valueWidth := h.viewport.Width - usedWidth
	if valueWidth < 10 {
		valueWidth = 10 // Mínimo
	}

	var showCursor bool
	// Preparar el valor del campo
	valueText := field.Value()
	// Usar tempEditValue si existe (modo edición)
	if field.tempEditValue != "" {
		valueText = field.tempEditValue
	}

	// Truncar el valor para que no afecte el diseño del footer
	// Descontar el padding que se aplicará al estilo
	textWidth := valueWidth - (horizontalPadding * 2)
	if textWidth < 1 {
		textWidth = 1
	}
	valueText = tinystring.Convert(valueText).Truncate(textWidth, 0).String()

	// Mostrar cursor solo si estamos en modo edición y el campo es editable
	if h.editModeActivated && field.Editable() {
		showCursor = true
	}

	// Definir el estilo para el valor del campo
	inputValueStyle := lipgloss.NewStyle().
		Width(valueWidth).
		Padding(0, horizontalPadding)

	// Aplicar estilos según el estado y tipo de handler
	if field.isExecutionHandler() {
		// Execution: Fondo blanco con letras grises (botón ejecutable)
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Foreground)).
			Foreground(lipgloss.Color(h.Background))
	} else if h.editModeActivated && field.Editable() {
		// Edit en modo edición activa
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Lowlight)).
			Foreground(lipgloss.Color(h.Foreground))
	} else {
		// Edit en modo no edición
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Lowlight)).
			Foreground(lipgloss.Color(h.Background))
	}

	// Añadir cursor si corresponde
	if showCursor {
		// Asegurar que el cursor está dentro de los límites
		runes := []rune(field.tempEditValue)
		if field.cursor < 0 {
			field.cursor = 0
		}
		if field.cursor > len(runes) {
			field.cursor = len(runes)
		}

		// Insertar el cursor en la posición correcta
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

	// Layout: [Label] [Value] [Scroll%] - scroll siempre a la derecha
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		paddedLabel,
		spacerStyle, // Espacio entre label y value
		styledValue,
		spacerStyle, // Espacio entre value y scroll
		info,        // Scroll % siempre a la derecha
	)
}
