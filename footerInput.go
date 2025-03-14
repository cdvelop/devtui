package devtui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (h *DevTUI) footerView() string {
	info := h.footerInfoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
	line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
