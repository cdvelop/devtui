package devtui

import (
	"fmt"
	"strings"

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
	tabContent := h.tabSections[h.activeTab].tabContents
	var contentLines []string
	for _, content := range tabContent {
		formattedMsg := h.formatMessage(content)
		contentLines = append(contentLines, h.textContentStyle.Render(formattedMsg))
	}
	return strings.Join(contentLines, "\n")
}

func (h *DevTUI) headerView() string {
	tab := h.tabSections[h.activeTab]
	title := h.headerTitleStyle.Render(h.AppName + "/" + tab.Title)
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (h *DevTUI) footerView() string {
	info := h.footerInfoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (t *DevTUI) renderLeftSectionForm() string {
	var lines []string

	style := lipgloss.NewStyle().
		Padding(0, 2)

	selectedStyle := style
	selectedStyle = selectedStyle.
		Bold(true).
		Background(lipgloss.Color(t.Highlight)).
		Foreground(lipgloss.Color(t.ForeGround))

	editingStyle := selectedStyle
	editingStyle = editingStyle.
		Foreground(lipgloss.Color(t.Background))

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
					line = editingStyle.Render(line)
				} else {
					line = selectedStyle.Render(line)
				}
			} else {
				line = style.Render(line)
			}

			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
