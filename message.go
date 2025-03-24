package devtui

import "github.com/cdvelop/messagetype"

// IDGenerator represents a unique ID generator
type IDGenerator interface {
	GetNewID() string
}

// Use the existing UnixID as an IDGenerator implementation
func NewIDGenerator() IDGenerator {
	return nil // Will be replaced with tui.id in NewTUI
}

func (t *tabSection) newTuiMessage(content string, mt messagetype.Type) tuiMessage {
	return tuiMessage{
		id:         t.tui.id.GetNewID(),
		Content:    content,
		Type:       mt,
		tabSection: t,
	}
}

// messageTracker keeps track of messages and their IDs
type messageTracker struct {
	messages map[MessageID]*tuiMessage
	idGen    IDGenerator
}

// NewMessageTracker creates a new message tracker
func NewMessageTracker() *messageTracker {
	return &messageTracker{
		messages: make(map[MessageID]*tuiMessage),
		// idGen will be set in NewTUI
	}
}

// SetIDGenerator sets the ID generator for the message tracker
func (mt *messageTracker) SetIDGenerator(idGen IDGenerator) {
	mt.idGen = idGen
}

// TrackMessage adds a message to be tracked
func (mt *messageTracker) TrackMessage(msg *tuiMessage) MessageID {
	id := MessageID(mt.idGen.GetNewID())
	msg.id = string(id)
	mt.messages[id] = msg
	return id
}

// UpdateMessage updates a tracked message with new content
func (mt *messageTracker) UpdateMessage(id MessageID, update MessageUpdate) bool {
	msg, exists := mt.messages[id]
	if !exists {
		return false
	}

	msg.Content = update.Content
	msg.Type = update.Type
	return true
}

// Add this helper function
func (t *tabSection) addNewContent(msgType messagetype.Type, content string) {
	t.tuiMessages = append(t.tuiMessages, t.newTuiMessage(content, msgType))
}
