package devtui

// createShortcutsTab creates and registers the shortcuts tab with its handler
import (
	. "github.com/cdvelop/tinystring"
)

func createShortcutsTab(tui *DevTUI) {
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")

	handler := &shortcutsInteractiveHandler{
		appName:            tui.AppName,
		lang:               OutLang(), // Get current language automatically
		needsLanguageInput: false,     // Initially show help content
	}
	// Use AddInteractiveHandler instead of AddEditHandler
	shortcutsTab.AddInteractiveHandler(handler, 0)
}

// shortcutsInteractiveHandler - Interactive handler for language selection and help display
type shortcutsInteractiveHandler struct {
	appName            string
	lang               string // e.g. "EN", "ES", etc.
	needsLanguageInput bool   // Controls when to activate edit mode
	lastOpID           string // Operation ID for tracking
}

func (h *shortcutsInteractiveHandler) Name() string {
	return "shortcutsGuide"
}

func (h *shortcutsInteractiveHandler) Label() string {
	return T(D.Language, ":").String()
}

// MessageTracker implementation for operation tracking
func (h *shortcutsInteractiveHandler) GetLastOperationID() string   { return h.lastOpID }
func (h *shortcutsInteractiveHandler) SetLastOperationID(id string) { h.lastOpID = id }

func (h *shortcutsInteractiveHandler) Value() string { return h.lang }

// Change handles both content display and user input via progress()
func (h *shortcutsInteractiveHandler) Change(newValue string, progress func(msgs ...any)) {
	if newValue == "" && !h.needsLanguageInput {
		// Display help content when field is selected (not in edit mode)
		progress(h.generateHelpContent())
		return
	}

	// Handle language change
	lang := OutLang(newValue)
	h.lang = lang
	h.needsLanguageInput = false

	// Show updated help content
	progress(h.generateHelpContent())
}

func (h *shortcutsInteractiveHandler) WaitingForUser() bool {
	return h.needsLanguageInput
}

// generateHelpContent creates the help content string
func (h *shortcutsInteractiveHandler) generateHelpContent() string {
	return T(h.appName, D.Shortcuts, D.Keyboard, `("`+h.lang+`"):

Tabs:
  • Tab/Shift+Tab  -`, D.Switch, ` tabs

`, D.Fields, `:
  • Left/Right     - Navigate fields
  • Enter          - Edit/Execute
  • Esc            -`, D.Cancel, `

Text Edit:
  • Left/Right     -`, D.Move, `cursor
  • Backspace      -`, D.Create, D.Space, `
  • Space/Letters  -`, D.Insert, D.Character, `

Viewport:
  • Up/Down        - Scroll line
  • PgUp/PgDown    - Scroll`, D.Page, `
  • Mouse Wheel    - Scroll`, D.Page, `

Scroll Status Icons:
  •  ■  - All content visible
  •  ▼  - Can scroll down
  •  ▲  - Can scroll up
  • ▼ ▲ - Can scroll both ways

Exit:
  • Ctrl+C         - Quit

Text selection enabled for copy/paste.

`, D.Language, D.Supported, `: EN, ES, ZH, HI, AR, PT, FR, DE, RU`).String()
}
