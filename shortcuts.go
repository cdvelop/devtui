package devtui

import "time"

// ShortcutsHandler - Shows keyboard navigation instructions
type ShortcutsHandler struct {
	shortcuts string
	lastOpID  string
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
  • Backspace   - Delete character left
  • Space       - Insert space
  • Characters  - Insert text at cursor

Viewport Navigation:
  • Up Arrow    - Scroll viewport up
  • Down Arrow  - Scroll viewport down

Application:
  • Ctrl+C      - Exit application
  • Ctrl+L      - Clear current tab content

Field Types:
  • Editable    - Press Enter to edit, Esc to cancel
  • Action      - Press Enter to execute, shows spinner during async operations`

	return &ShortcutsHandler{shortcuts: shortcuts}
}

func (h *ShortcutsHandler) Label() string          { return "Keyboard Shortcuts" }
func (h *ShortcutsHandler) Value() string          { return "Press Enter to view" }
func (h *ShortcutsHandler) Editable() bool         { return false }
func (h *ShortcutsHandler) Timeout() time.Duration { return 0 }

func (h *ShortcutsHandler) Change(newValue any) (string, error) {
	return h.shortcuts, nil
}

// WritingHandler methods
func (h *ShortcutsHandler) Name() string                       { return "Shortcuts" }
func (h *ShortcutsHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *ShortcutsHandler) GetLastOperationID() string         { return "" } // Always create new messages
