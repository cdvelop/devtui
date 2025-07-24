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
  • Backspace   - Delete character

Viewport Navigation:
  • Up Arrow    - Scroll viewport up line by line
  • Down Arrow  - Scroll viewport down line by line
  • Page Up     - Scroll viewport up page by page
  • Page Down   - Scroll viewport down page by page
  • Mouse Wheel - Scroll viewport (when available)

Application:
  • Ctrl+C      - Exit application

Note: Text selection enabled for copying error messages and logs.
Mouse scroll may work depending on bubbletea version and terminal capabilities.
`
	return &ShortcutsHandler{shortcuts: shortcuts}
}

func (h *ShortcutsHandler) Name() string    { return "DevTUI Help & Navigation Guide" }
func (h *ShortcutsHandler) Content() string { return h.shortcuts }
