package devtui

import "github.com/cdvelop/messagetype"

// tuiMessage represents a message sent asynchronously or synchronously in the TUI
type tuiMessage struct {
	id         string // Unique ID for the message
	Content    string
	Type       messagetype.Type
	tabSection *tabSection
}

func (t *tabSection) newTuiMessage(content string, mt messagetype.Type) tuiMessage {

	return tuiMessage{
		id:         t.tui.id.GetNewID(),
		Content:    content,
		Type:       mt,
		tabSection: t,
	}
}
