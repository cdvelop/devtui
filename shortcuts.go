package devtui

// ShortcutsHandler - Shows keyboard navigation instructions
type ShortcutsHandler struct {
	shortcuts string
}

func NewShortcutsHandler() *ShortcutsHandler {
	shortcuts := `Keyboard Navigation Commands:

Navigation Between Tabs:
  • Tab         - Next tab
  • Shift+Tab   - Previous tab

Navigation Between Fields:
  • Left Arrow  - Previous field (cycle)
  • Right Arrow - Next field (cycle)

Field Editing:
  • Enter       - Edit field / Execute action
  • Esc         - Cancel editing / Exit field

Text Editing (when in edit mode):
  • Left Arrow  - Move cursor left
  • Right Arrow - Move cursor right
  • Backspace   - Create space

Viewport Navigation:
  • Up Arrow/Mouse Wheel - Scroll viewport up
  • Down Arrow/Mouse Wheel - Scroll viewport down

Application:
  • Ctrl+C      - Exit application
`
	return &ShortcutsHandler{shortcuts: shortcuts}
}

func (h *ShortcutsHandler) Name() string    { return "Shortcuts" }
func (h *ShortcutsHandler) Label() string   { return "Help" }
func (h *ShortcutsHandler) Content() string { return h.shortcuts }
