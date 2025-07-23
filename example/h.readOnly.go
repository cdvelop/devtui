package example

import "time"

// WelcomeHandler - Readonly information display (empty label)
type WelcomeHandler struct {
}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{}
}

// WritingHandler implementation
func (h *WelcomeHandler) Name() string                 { return "WelcomeHandler" }
func (h *WelcomeHandler) SetLastOperationID(id string) {}
func (h *WelcomeHandler) GetLastOperationID() string   { return "" }

// FieldHandler implementation
func (h *WelcomeHandler) Label() string { return "" } // EMPTY = readonly display
func (h *WelcomeHandler) Value() string {
	return "DevTUI Features"
}
func (h *WelcomeHandler) Editable() bool         { return false }
func (h *WelcomeHandler) Timeout() time.Duration { return 0 }

func (h *WelcomeHandler) Change(newValue any, progress ...func(string)) (string, error) {
	// For readonly fields, Change() shows clean content without timestamp
	return "DevTUI Features:\n• Async operations with dynamic progress messages\n• Configurable timeouts\n• Error handling\n• Real-time progress feedback\n• Handler-based architecture", nil
}
