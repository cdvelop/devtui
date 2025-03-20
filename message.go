package devtui

import "github.com/cdvelop/messagetype"

// Message represents a message sent asynchronously or synchronously in the TUI
type Message struct {
	id         string // Unique ID for the message
	Content    string
	Type       messagetype.Type
	TabSection *TabSection
}
