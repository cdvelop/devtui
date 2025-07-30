package devtui

// createShortcutsTab creates and registers the shortcuts tab with its handler
import (
	. "github.com/cdvelop/tinystring"
)

func createShortcutsTab(tui *DevTUI) {
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")

	handler := &shortcutsEditHandler{
		appName: tui.AppName,
		lang:    "EN",
	}
	// Provide a suitable time.Duration, e.g., 0 for no timeout
	shortcutsTab.AddEditHandler(handler, 0)
}

// shortcutsEditHandler - Editable handler for language selection and help
type shortcutsEditHandler struct {
	appName string
	lang    string // e.g. "EN", "ES", etc.
}

func (h *shortcutsEditHandler) Name() string  { return "DevTUI Help & Navigation Guide" }
func (h *shortcutsEditHandler) Label() string { return "Language (idioma)" }
func (h *shortcutsEditHandler) Value() string { return h.lang }

// Change actualiza el idioma global usando OutLang
func (h *shortcutsEditHandler) Change(newValue string, progress func(msgs ...any)) {
	OutLang(newValue)
	h.lang = newValue
	progress(D.Language, D.Changed, D.To, newValue)
}

func (h *shortcutsEditHandler) Content() string {
	return h.appName + ` Keyboard Commands ("` + h.lang + `"):

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

Scroll Status Icons:
  •  ■  - All content visible
  •  ▼  - Can scroll down
  •  ▲  - Can scroll up
  • ▼ ▲ - Can scroll both ways

Exit:
  • Ctrl+C         - Quit

Text selection enabled for copy/paste.`
}
