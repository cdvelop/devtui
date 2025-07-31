package example

import (
	"strings"
	"time"
)

type SimpleChatHandler struct {
	Messages           []ChatMessage
	CurrentInput       string
	WaitingForUserFlag bool
	IsProcessing       bool
}

type ChatMessage struct {
	IsUser bool
	Text   string
	Time   time.Time
}

func (h *SimpleChatHandler) Name() string { return "SimpleChat" }

func (h *SimpleChatHandler) Label() string {
	if h.IsProcessing {
		return "Processing..."
	}
	if h.WaitingForUserFlag {
		return "Type message"
	}
	return "Chat (Press Enter)"
}

func (h *SimpleChatHandler) Value() string        { return h.CurrentInput }
func (h *SimpleChatHandler) WaitingForUser() bool { return h.WaitingForUserFlag && !h.IsProcessing }

func (h *SimpleChatHandler) Change(newValue string, progress func(msgs ...any)) {
	// Display content when field selected
	if newValue == "" && !h.WaitingForUserFlag && !h.IsProcessing {
		if len(h.Messages) == 0 {
			progress("Welcome")
		} else {
			for _, msg := range h.Messages {
				if msg.IsUser {
					progress("U: " + msg.Text)
				} else {
					progress("A: " + msg.Text)
				}
			}
		}
		return
	}

	// Handle user input
	if newValue != "" && strings.TrimSpace(newValue) != "" {
		userMsg := strings.TrimSpace(newValue)

		h.Messages = append(h.Messages, ChatMessage{
			IsUser: true,
			Text:   userMsg,
			Time:   time.Now(),
		})

		h.WaitingForUserFlag = false
		h.IsProcessing = true
		h.CurrentInput = ""

		progress("U: " + userMsg)
		progress("Processing...")

		go h.generateAIResponse(userMsg, progress)
		return
	}

	// Empty input while waiting
	if newValue == "" && h.WaitingForUserFlag && !h.IsProcessing {
		progress("Type message")
		return
	}
}

func (h *SimpleChatHandler) generateAIResponse(userMessage string, progress func(msgs ...any)) {
	time.Sleep(500 * time.Millisecond) // Short delay for testing

	var response string
	switch strings.ToLower(userMessage) {
	case "hello", "hi":
		response = "Hello"
	case "help":
		response = "Help available"
	case "test":
		response = "Test OK"
	default:
		response = "Response: " + userMessage
	}

	h.Messages = append(h.Messages, ChatMessage{
		IsUser: false,
		Text:   response,
		Time:   time.Now(),
	})

	h.IsProcessing = false
	h.WaitingForUserFlag = true

	progress("A: " + response)
}
