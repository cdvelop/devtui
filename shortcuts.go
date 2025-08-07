package devtui

// createShortcutsTab creates and registers the shortcuts tab with its handler
import (
	"fmt"

	. "github.com/cdvelop/tinystring"
)

func createShortcutsTab(tui *DevTUI) {
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")

	handler := &shortcutsInteractiveHandler{
		appName:            tui.AppName,
		lang:               OutLang(), // Get current language automatically
		needsLanguageInput: false,     // Initially show help content
		tui:                tui,       // NEW: Reference to TUI for shortcut registry access
	}
	// Use AddInteractiveHandler instead of AddEditHandler
	shortcutsTab.AddInteractiveHandler(handler, 0)
}

// shortcutsInteractiveHandler - Interactive handler for language selection and help display
type shortcutsInteractiveHandler struct {
	appName            string
	lang               string  // e.g. "EN", "ES", etc.
	needsLanguageInput bool    // Controls when to activate edit mode
	lastOpID           string  // Operation ID for tracking
	tui                *DevTUI // NEW: Reference to TUI for shortcut registry access
}

func (h *shortcutsInteractiveHandler) Name() string {
	return "shortcutsGuide"
}

func (h *shortcutsInteractiveHandler) Label() string {
	return Translate(D.Language, ":").String()
}

// MessageTracker implementation for operation tracking
func (h *shortcutsInteractiveHandler) GetLastOperationID() string   { return h.lastOpID }
func (h *shortcutsInteractiveHandler) SetLastOperationID(id string) { h.lastOpID = id }

func (h *shortcutsInteractiveHandler) Value() string { return Convert(h.lang).Low().String() }

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
	content := Translate(h.appName, D.Shortcuts, D.Keyboard, `:

`, D.Content, D.Tab, `:
  • Tab/Shift+Tab  -`, D.Switch, D.Content, `

`, D.Fields, `:
  • `, D.Arrow, D.Left, `/`, D.Right, `     -`, D.Switch, D.Field, `
  • Enter          				-`, D.Edit, `/`, D.Execute, `
  • Esc            				-`, D.Cancel, `

`, D.Edit, D.Text, `:
  • `, D.Arrow, D.Left, `/`, D.Right, `   -`, D.Move, `cursor
  • Backspace      			-`, D.Create, D.Space, `

Viewport:
  • `, D.Arrow, D.Up, "/", D.Down, `    - Scroll`, D.Line, D.Text, `
  • PgUp/PgDown    		- Scroll`, D.Page, `
  • Mouse Wheel    		- Scroll`, D.Page, `

Scroll `, D.Status, D.Icons, `:
  •  ■  - `, D.All, D.Content, D.Visible, `
  •  ▼  - `, D.Can, `scroll`, D.Down, `
  •  ▲  - `, D.Can, `scroll`, D.Up, `
  • ▼ ▲ - `, D.Can, `scroll`, D.Down, `/`, D.Up, `

`, D.Quit, `:
  • Ctrl+C         - `, D.Quit, `
`).String()

	// NEW: Add registered shortcuts section
	if h.tui != nil && h.tui.shortcutRegistry != nil {
		shortcuts := h.getRegisteredShortcuts()
		if len(shortcuts) > 0 {
			content += "\n\nRegistered Shortcuts:\n"
			for key, description := range shortcuts {
				content += fmt.Sprintf("  • %s - %s\n", key, description)
			}
		}
	}

	content += "\n" + Translate(D.Language, D.Supported, `: en, es, zh, hi, ar, pt, fr, de, ru`).String()
	return content
}

// getRegisteredShortcuts returns all registered shortcuts with descriptions
func (h *shortcutsInteractiveHandler) getRegisteredShortcuts() map[string]string {
	shortcuts := make(map[string]string)
	if h.tui != nil && h.tui.shortcutRegistry != nil {
		allEntries := h.tui.shortcutRegistry.GetAll()
		for key, entry := range allEntries {
			shortcuts[key] = entry.Description
		}
	}
	return shortcuts
}
