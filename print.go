package devtui

import (
	"fmt"
	"time"

	. "github.com/cdvelop/messagetype"
)

// Print sends a normal Label or error to the tui in current tab
func (h *DevTUI) Print(messages ...any) {
	msgType := DetectMessageType(messages...)
	h.sendMessage(joinMessages(messages...), msgType, &h.TabSections[h.activeTab])
}

func joinMessages(messages ...any) (Label string) {
	var space string
	for _, m := range messages {
		Label += space + fmt.Sprint(m)
		space = " "
	}
	return
}

// sendMessage envía un mensaje al tui
func (t *DevTUI) sendMessage(content string, mt MessageType, tabSection *TabSection) {

	t.tabContentsChan <- tabContent{
		Content:    content,
		Type:       mt,
		Time:       time.Now(),
		tabSection: tabSection,
	}
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tabContent) string {
	timeStr := t.timeStyle.Render(msg.Time.Format("15:04:05"))
	// content := fmt.Sprintf("[%s] %s", timeStr, msg.Content)

	switch msg.Type {
	case Error:
		msg.Content = t.errStyle.Render(msg.Content)
	case Warning:
		msg.Content = t.warnStyle.Render(msg.Content)
	case Info:
		msg.Content = t.infoStyle.Render(msg.Content)
	case OK:
		msg.Content = t.okStyle.Render(msg.Content)
		// default:
		// 	msg.Content= msg.Content
	}

	return fmt.Sprintf("%s %s", timeStr, msg.Content)
}
