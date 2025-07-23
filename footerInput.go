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

	// Check if this is a HandlerDisplay for special layout
	if field.isDisplayOnly() {
		// Special layout for HandlerDisplay: full width label, no separate value section
		horizontalPadding := 1
		info := h.renderScrollInfo()

		// Use full width for label content (field content from HandlerDisplay.Content())
		fullWidth := h.viewport.Width - lipgloss.Width(info) - horizontalPadding*2
		labelText := tinystring.Convert(field.handler.Value()).Truncate(fullWidth-1, 0).String()

		// Layout: [Full Width Label Content] [ScrollInfo]
		styledLabel := h.headerTitleStyle.Render(labelText)
		spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")

		return lipgloss.JoinHorizontal(lipgloss.Left, styledLabel, spacerStyle, info)
	}

	// Normal layout for Edit/Run handlers: [Label] [Value] [ScrollInfo]
	// Usar el ancho estándar de etiquetas definido en el estilo
	labelWidth := h.labelWidth

	// Obtener el padding utilizado en el header/footer para mantener consistencia
	horizontalPadding := 1 // Este valor viene del Padding(0, 1) en headerTitleStyle

	// Truncar la etiqueta si es necesario (sin agregar ":")
	labelText := tinystring.Convert(field.handler.Label()).Truncate(labelWidth-1, 0).String()

	// Aplicar el estilo base para garantizar un ancho fijo
	fixedWidthLabel := h.labelStyle.Render(labelText)

	// Aplicar el estilo visual (colores) manteniendo el ancho fijo
	paddedLabel := h.headerTitleStyle.Render(fixedWidthLabel)

	// Obtener el indicador de porcentaje con el estilo actual
	info := h.renderScrollInfo()

	// OR if you need truncation:
	labelText = tinystring.Convert(field.handler.Label()).Truncate(labelWidth-1, 0).String()
	valueWidth, _ := h.calculateInputWidths(labelText)

	var showCursor bool
	// Preparar el valor del campo
	valueText := field.Value()
	// Usar tempEditValue si existe (modo edición)
	if field.tempEditValue != "" {
		valueText = field.tempEditValue
	}

	// Truncar el valor para que no afecte el diseño del footer
	valueText = tinystring.Convert(valueText).Truncate(valueWidth, 0).String()

	// Mostrar cursor solo si estamos en modo edición y el campo es editable y NO es readonly
	if h.editModeActivated && field.Editable() && !field.isDisplayOnly() {
		showCursor = true
	}

	// Definir el estilo para el valor del campo
	inputValueStyle := lipgloss.NewStyle().
		Width(valueWidth).
		Padding(0, horizontalPadding) // Añadir padding consistente

	// Aplicar estilos según el estado
	if field.isDisplayOnly() { // NEW: Empty label detection (exactly "")
		// Use fieldReadOnlyStyle - highlight background with clear text
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Highlight)).
			Foreground(lipgloss.Color(h.Foreground)) // Clear text on highlight
		// No cursor allowed, no interaction
	} else if h.editModeActivated && field.Editable() {
		// Estilo para edición activa
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Lowlight)).
			Foreground(lipgloss.Color(h.Foreground))
	} else if !field.Editable() {
		// Estilo para campos no editables (action buttons)
		inputValueStyle = inputValueStyle.
			Background(lipgloss.Color(h.Foreground)).
			Foreground(lipgloss.Color(h.Background))
	} else {
		// Estilo para campos editables pero no en modo edición
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
