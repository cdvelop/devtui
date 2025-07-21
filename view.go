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
	// return fmt.Sprintf("%s\n%s\n%s", h.headerView(), h.ContentView(), h.footerView())
}

// ContentView renderiza los mensajes para una sección de contenido
func (h *DevTUI) ContentView() string {
	if len(h.tabSections) == 0 {
		return "No tabs created yet"
	}
	if h.activeTab >= len(h.tabSections) {
		h.activeTab = 0
	}

	// Proteger el acceso a tabContents con mutex
	section := h.tabSections[h.activeTab]
	section.mu.RLock()
	tabContent := make([]tabContent, len(section.tabContents)) // Copia para evitar retener el lock
	copy(tabContent, section.tabContents)
	section.mu.RUnlock()

	var contentLines []string
	for _, content := range tabContent {
		formattedMsg := h.formatMessage(content)
		contentLines = append(contentLines, h.textContentStyle.Render(formattedMsg))
	}
	return strings.Join(contentLines, "\n")
}

func (h *DevTUI) headerView() string {
	if len(h.tabSections) == 0 {
		return h.headerTitleStyle.Render(h.AppName + "/No tabs")
	}
	if h.activeTab >= len(h.tabSections) {
		h.activeTab = 0
	}

	tab := h.tabSections[h.activeTab]

	// Truncar el título si es necesario
	headerText := h.AppName + "/" + tab.Title()
	truncatedHeader := tinystring.Convert(headerText).Truncate(h.labelWidth, 0).String()

	// Aplicar el estilo base para garantizar un ancho fijo
	fixedWidthHeader := h.labelStyle.Render(truncatedHeader)

	// Aplicar el estilo visual manteniendo el ancho fijo
	title := h.headerTitleStyle.Render(fixedWidthHeader)

	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
