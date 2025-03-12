package devtui

import (
	"fmt"
	"strings"
	"time"

	. "github.com/cdvelop/messagetype"
)

// Print sends a normal Label or error to the tui
func (h *DevTUI) Print(messages ...any) {
	msgType := Normal
	newMessages := make([]any, 0, len(messages))

	for _, msg := range messages {
		if str, isString := msg.(string); isString {

			switch strings.ToLower(str) {
			case "error":
				msgType = Error
				continue
			case "warning", "debug":
				msgType = Warning
				continue
			case "info":
				msgType = Info
				continue
			case "ok":
				msgType = OK
				continue
			}
		}
		if _, isError := msg.(error); isError {
			msgType = Error
		}

		newMessages = append(newMessages, msg)
	}

	h.sendMessage(joinMessages(newMessages...), msgType)
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
func (t *DevTUI) sendMessage(content string, msgType messageType, tabSection *TabSection) {

	t.tabContentsChan <- tabContent{
		Content:    content,
		Type:       msgType,
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
