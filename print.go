package devtui

import (
	"fmt"

	"github.com/cdvelop/messagetype"
	tea "github.com/charmbracelet/bubbletea"
)

// Print sends a normal Name or error to the tui in current tab
func (h *DevTUI) Print(messages ...any) {
	msgType := messagetype.DetectMessageType(messages...)
	h.sendMessage(joinMessages(messages...), msgType, &h.tabSections[h.activeTab])
}

func joinMessages(messages ...any) (Name string) {
	var space string
	for _, m := range messages {
		Name += space + fmt.Sprint(m)
		space = " "
	}
	return
}

// sendMessage envía un mensaje al tui por el canal de mensajes
func (t *DevTUI) sendMessage(content string, mt messagetype.Type, ts *tabSection) {
	t.tabContentsChan <- ts.newTuiMessage(content, mt)
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tuiMessage) string {

	timeStr := t.timeStyle.Render(t.id.UnixSecondsToTime(msg.id))

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

// ProcessFieldValueChange handles field value changes
func (h *DevTUI) ProcessFieldValueChange(field *fieldHandler, newValue string) {
	field.ExecuteValueChange(newValue, &h.tabSections[h.activeTab])
	// Start listening for async messages if not already listening
	h.tea.Send(tea.Cmd(h.listenForAsyncMessages(h.asyncMessageChan)))
}
