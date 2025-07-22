package devtui

import (
	"fmt"

	"github.com/cdvelop/messagetype"
)

// Print sends a normal Label or error to the tui in current tab
func (h *DevTUI) Print(messages ...any) {
	msgType := messagetype.DetectMessageType(messages...)
	h.sendMessage(joinMessages(messages...), msgType, h.tabSections[h.activeTab])
}

func joinMessages(messages ...any) (Label string) {
	var space string
	for _, m := range messages {
		Label += space + fmt.Sprint(m)
		space = " "
	}
	return
}

// sendMessage envía un mensaje al tui por el canal de mensajes
func (d *DevTUI) sendMessage(content string, mt messagetype.Type, tabSection *tabSection, operationID ...string) {

	tabSection.addNewContent(mt, content)

	newContent := d.newContent(content, mt, tabSection, operationID...)

	d.tabContentsChan <- newContent
}

// NEW: sendMessageWithHandler sends a message with handler identification
func (d *DevTUI) sendMessageWithHandler(content string, mt messagetype.Type, tabSection *tabSection, handlerName string, operationID string) {
	// Use update or add function that handles operationID reuse
	_, newContent := tabSection.updateOrAddContentWithHandler(mt, content, handlerName, operationID)

	// Always send to channel to trigger UI update, regardless of whether content was updated or added new
	d.tabContentsChan <- newContent

	// Call SetLastOperationID on the handler after processing
	if tabSection.writingHandlers != nil {
		if handler, exists := tabSection.writingHandlers[handlerName]; exists {
			handler.SetLastOperationID(newContent.Id)
		}
	}
}

func (h *DevTUI) newContent(content string, mt messagetype.Type, tabSection *tabSection, operationID ...string) tabContent {
	var id string
	var opID *string

	if len(operationID) > 0 && operationID[0] != "" {
		// Use provided operation ID for async operations
		id = operationID[0]
		opID = &operationID[0]
	} else {
		// Generate new ID for regular operations (current behavior)
		if h.id != nil {
			id = h.id.GetNewID()
		} else {
			id = "temp-id"
			h.LogToFile("Warning: unixid not initialized, using fallback ID")
		}
		opID = nil // Not an async operation
	}

	return tabContent{
		Id:          id,
		Content:     content,
		Type:        mt,
		tabSection:  tabSection,
		operationID: opID,
		isProgress:  false, // Will be set by specific async methods
		isComplete:  false, // Will be set by specific async methods
	}
}

// NEW: newContentWithHandler creates tabContent with handler identification
func (h *DevTUI) newContentWithHandler(content string, mt messagetype.Type, tabSection *tabSection, handlerName string, operationID ...string) tabContent {
	var id string
	var opID *string

	if len(operationID) > 0 && operationID[0] != "" {
		// Use provided operation ID for async operations
		id = operationID[0]
		opID = &operationID[0]
	} else {
		// Generate new ID for regular operations
		if h.id != nil {
			id = h.id.GetNewID()
		} else {
			id = "temp-id"
			h.LogToFile("Warning: unixid not initialized, using fallback ID")
		}
		opID = nil // Not an async operation
	}

	return tabContent{
		Id:          id,
		Content:     content,
		Type:        mt,
		tabSection:  tabSection,
		operationID: opID,
		isProgress:  false,       // Will be set by specific async methods
		isComplete:  false,       // Will be set by specific async methods
		handlerName: handlerName, // NEW: Include handler name
	}
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tabContent) string {

	var timeStr string
	if t.id != nil {
		timeStr = t.timeStyle.Render(t.id.UnixNanoToTime(msg.Id))
	} else {
		// When unixid is not initialized, use a simple timestamp format
		timeStr = t.timeStyle.Render("--:--:--")
	}

	// NEW: Include handler name in message format
	var handlerName string
	if msg.handlerName != "" {
		handlerName = fmt.Sprintf("[%s] ", msg.handlerName)
	}

	// timeStr := t.timeStyle.Render(msg.Time.Format("15:04:05"))
	// content := fmt.Sprintf("[%s] %s", timeStr, msg.Content)

	switch msg.Type {
	case messagetype.Error:
		msg.Content = t.errStyle.Render(msg.Content)
	case messagetype.Warning:
		msg.Content = t.warnStyle.Render(msg.Content)
	case messagetype.Info:
		msg.Content = t.infoStyle.Render(msg.Content)
	case messagetype.Success:
		msg.Content = t.okStyle.Render(msg.Content)
		// default:
		// 	msg.Content= msg.Content
	}

	return fmt.Sprintf("%s %s%s", timeStr, handlerName, msg.Content)
}
