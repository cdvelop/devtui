package devtui

import (
	"fmt"
	"strings"

	"github.com/cdvelop/tinystring"
	"github.com/charmbracelet/lipgloss"
)

func (h *DevTUI) View() string {
	if !h.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", h.headerView(), h.viewport.View(), h.footerView())
}

// ContentView renderiza los mensajes para una sección de contenido
func (h *DevTUI) ContentView() string {
	tabContent := h.tabSections[h.activeTab].tuiMessages
	var contentLines []string
	for _, content := range tabContent {
		formattedMsg := h.formatMessage(content)
		contentLines = append(contentLines, h.textContentStyle.Render(formattedMsg))
	}
	return strings.Join(contentLines, "\n")
}

func (h *DevTUI) headerView() string {
	tab := h.tabSections[h.activeTab]

	// Truncar el título si es necesario
	headerText := h.AppName + "/" + tab.title
	truncatedHeader := tinystring.Convert(headerText).Truncate(h.labelWidth, 0).String()

	// Aplicar el estilo base para garantizar un ancho fijo
	fixedWidthHeader := h.labelStyle.Render(truncatedHeader)

	// Aplicar el estilo visual manteniendo el ancho fijo
	title := h.headerTitleStyle.Render(fixedWidthHeader)

	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (h *DevTUI) renderTabIndicators() string {
	var indicators []string

	for i, tab := range h.tabSections {
		var indicator string
		if i == h.activeTab {
			indicator = h.activeTabStyle.Render(tab.title)
		} else {
			indicator = h.inactiveTabStyle.Render(tab.title)
		}
		indicators = append(indicators, indicator)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, indicators...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
