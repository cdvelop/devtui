package devtui

import "github.com/charmbracelet/lipgloss"

// calculateInputWidths calculates the width available for text input based on viewport and other elements
// Returns valueWidth (total width for the input area) and availableTextWidth (width for the text itself)
func (h *DevTUI) calculateInputWidths(fieldLabel string) (valueWidth, availableTextWidth int) {
	horizontalPadding := 1

	// Process label
	labelText := fieldLabel + ":"
	fixedWidthLabel := h.labelStyle.Render(labelText)
	paddedLabel := h.headerTitleStyle.Render(fixedWidthLabel)

	// Calculate other components
	infoWidth := lipgloss.Width(h.renderScrollInfo())
	separationSpace := horizontalPadding * 2

	// Calculate final widths
	valueWidth = h.viewport.Width - lipgloss.Width(paddedLabel) - infoWidth - separationSpace
	availableTextWidth = valueWidth - (horizontalPadding * 2)

	return valueWidth, availableTextWidth
}
