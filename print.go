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
func (d *DevTUI) sendMessage(content string, mt messagetype.Type, tabSection *tabSection) {

	tabSection.addNewContent(mt, content)

	newContent := d.newContent(content, mt, tabSection)

	d.tabContentsChan <- newContent
}

func (h *DevTUI) newContent(content string, mt messagetype.Type, tabSection *tabSection) tabContent {
	var id string
	if h.id != nil {
		id = h.id.GetNewID()
	} else {
		// Fallback ID if unixid is not available
		id = "temp-id"
		h.LogToFile("Warning: unixid not initialized, using fallback ID")
	}

	return tabContent{
		Id:         id,
		Content:    content,
		Type:       mt,
		tabSection: tabSection,
	}
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tabContent) string {

	timeStr := t.timeStyle.Render(t.id.UnixNanoToTime(msg.Id))

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

	return fmt.Sprintf("%s %s", timeStr, msg.Content)
}
