package devtui

import (
	"strings"
	"time"

	"github.com/cdvelop/messagetype"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AsyncMessageMsg is a message from an async field handler
type AsyncMessageMsg tuiMessage

func (cf *fieldHandler) SetCursorAtEnd() {
	// Calculate cursor position based on rune count, not byte count
	cf.cursor = len([]rune(cf.Value()))
}

// listenToMessages crea un comando para escuchar mensajes del canal
func (h *DevTUI) listenToMessages() tea.Cmd {
	return func() tea.Msg {
		msg := <-h.tabContentsChan
		return channelMsg(msg)
	}
}

// listenForAsyncMessages creates a command to wait for async field messages
func (h *DevTUI) listenForAsyncMessages(asyncMsgChan chan tuiMessage) tea.Cmd {
	return func() tea.Msg {
		return AsyncMessageMsg(<-asyncMsgChan)
	}
}

// tickEverySecond crea un comando para actualizar el tiempo
func (h *DevTUI) tickEverySecond() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update maneja las actualizaciones del estado
func (h *DevTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg: // Al presionar una tecla
		continueProcessing, keyCmd := h.HandleKeyboard(msg)
		if !continueProcessing {
			if keyCmd != nil {
				return h, keyCmd
			}
			return h, nil
		}

		if keyCmd != nil {
			cmds = append(cmds, keyCmd)
		}

	case AsyncMessageMsg:
		// Process async message from field handler
		asyncMsg := tuiMessage(msg)

		// Check if this is an update to an existing message
		if asyncMsg.id != "" {
			// Find and update the existing message
			for i, existingMsg := range h.tabSections[h.activeTab].tuiMessages {
				if existingMsg.id == asyncMsg.id {
					h.tabSections[h.activeTab].tuiMessages[i].Content = asyncMsg.Content
					h.tabSections[h.activeTab].tuiMessages[i].Type = asyncMsg.Type
					break
				}
			}
		} else {
			// If no ID, treat as a new message
			h.sendMessage(asyncMsg.Content, asyncMsg.Type, asyncMsg.tabSection)
		}

		// Update the viewport to show the changes
		h.updateViewport()

		// Continue listening for more async messages
		cmds = append(cmds, h.listenForAsyncMessages(h.asyncMessageChan))

	case channelMsg: // Handle messages from the channel
		// Start listening for new messages again after processing the current one
		cmds = append(cmds, h.listenToMessages())

		// Convert the channel message to a tabContent type
		tc := tuiMessage(msg)

		// Only update the viewport if the message belongs to the currently active tab
		if tc.tabSection.index == h.activeTab {
			h.updateViewport()
		}

	case tea.WindowSizeMsg: // update the viewport size

		headerHeight := lipgloss.Height(h.headerView())
		footerHeight := lipgloss.Height(h.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !h.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			h.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			h.viewport.YPosition = headerHeight
			h.viewport.SetContent(h.ContentView())
			h.ready = true
		} else {
			h.viewport.Width = msg.Width
			h.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tickMsg: // update the time every second
		h.currentTime = time.Now().Format("15:04:05")
		cmds = append(cmds, h.tickEverySecond())

	case tea.FocusMsg:
		h.focused = true
	case tea.BlurMsg:
		h.focused = false

	}
	// Handle keyboard and mouse events in the viewport
	h.viewport, cmd = h.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return h, tea.Batch(cmds...)
}

// Write implementa io.Writer para capturar la salida de otros procesos
func (ts *tabSection) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		// Detectar automáticamente el tipo de mensaje
		msgType := messagetype.DetectMessageType(msg)

		ts.tui.sendMessage(msg, msgType, ts)
		// Si es un error, escribirlo en el archivo de log
		if msgType == messagetype.Error {
			ts.tui.LogToFile(msg)
		}

	}
	return len(p), nil
}

func (h *DevTUI) updateViewport() {
	h.viewport.SetContent(h.ContentView())
	h.viewport.GotoBottom()
}

func (h *DevTUI) editingConfigOpen(open bool, currentField *fieldHandler, msg string) {

	if open {
		h.editModeActivated = true
	} else {
		h.editModeActivated = false
	}

	if currentField != nil {
		currentField.SetCursorAtEnd()
	}

	if msg != "" {
		tabSection := &h.tabSections[h.activeTab]
		tabSection.addNewContent(messagetype.Warning, msg)
	}
}
