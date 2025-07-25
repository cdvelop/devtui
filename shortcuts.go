package devtui

// createShortcutsTab creates and registers the shortcuts tab with its handler
func createShortcutsTab(tui *DevTUI) {
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")

	shortcuts := tui.AppName + ` Keyboard Commands:

Tabs:
  • Tab/Shift+Tab  - Switch tabs

Fields:
  • Left/Right     - Navigate fields
  • Enter          - Edit/Execute
  • Esc            - Cancel

Text Edit:
  • Left/Right     - Move cursor
  • Backspace      - Create space
  • Space/Letters  - Insert char

Viewport:
  • Up/Down        - Scroll line
  • PgUp/PgDown    - Scroll page
  • Mouse Wheel    - Scroll (optional)

Exit:
  • Ctrl+C         - Quit

Text selection enabled for copy/paste.
`

	handler := &shortcutsHandler{shortcuts: shortcuts}
	shortcutsTab.RegisterHandlerDisplay(handler)
}

// shortcutsHandler - Shows keyboard navigation instructions
type shortcutsHandler struct {
	shortcuts string
}

func (h *shortcutsHandler) Name() string    { return "DevTUI Help & Navigation Guide" }
func (h *shortcutsHandler) Content() string { return h.shortcuts }
