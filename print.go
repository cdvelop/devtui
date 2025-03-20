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
func (t *DevTUI) sendMessage(content string, mt messagetype.Type, tabSection *TabSection) {
	t.tabContentsChan <- t.newTuiMessage(content, mt, tabSection)
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg TuiMessage) string {

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
	case messagetype.OK:
		msg.Content = t.okStyle.Render(msg.Content)
		// default:
		// 	msg.Content= msg.Content
	}

	return fmt.Sprintf("%s %s", timeStr, msg.Content)
}

// ProcessFieldValueChange handles both synchronous and asynchronous field value changes
func (h *DevTUI) ProcessFieldValueChange(field *FieldHandler, newValue string) {
	if field.IsAsync && field.AsyncFieldValueChange != nil {
		// Start a goroutine to handle async processing
		go field.AsyncFieldValueChange(newValue, h.asyncMessageChan)

		// Start listening for async messages if not already listening
		h.tea.Send(tea.Cmd(h.listenForAsyncMessages(h.asyncMessageChan)))
	} else if field.ChangeValue != nil {
		// Handle synchronous field value change
		execMessage, err := field.ChangeValue(newValue)
		tabSection := &h.tabSections[h.activeTab]

		if err != nil {
			tabSection.addNewContent(messagetype.Error, err.Error())
		} else if execMessage != "" {
			tabSection.addNewContent(messagetype.Info, execMessage)
		}
	}
}
