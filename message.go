package devtui

import "github.com/cdvelop/messagetype"

// TuiMessage represents a message sent asynchronously or synchronously in the TUI
type TuiMessage struct {
	id         string // Unique ID for the message
	Content    string
	Type       messagetype.Type
	TabSection *TabSection
}

func (h *DevTUI) newTuiMessage(content string, mt messagetype.Type, tabSection *TabSection) TuiMessage {

	return TuiMessage{
		id:         h.id.GetNewID(),
		Content:    content,
		Type:       mt,
		TabSection: tabSection,
	}
}
