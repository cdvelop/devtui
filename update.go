package devtui

import (
	"time"

	. "github.com/cdvelop/tinystring"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// listenToMessages crea un comando para escuchar mensajes del canal
func (h *DevTUI) listenToMessages() tea.Cmd {
	return func() tea.Msg {
		msg := <-h.tabContentsChan
		return channelMsg(msg)
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
		continueProcessing, keyCmd := h.handleKeyboard(msg)
		if !continueProcessing {
			if keyCmd != nil {
				return h, keyCmd
			}
			return h, nil
		}

		if keyCmd != nil {
			cmds = append(cmds, keyCmd)
		}

	case channelMsg: // Handle messages from the channel
		// Start listening for new messages again after processing the current one
		cmds = append(cmds, h.listenToMessages())

		// Convert the channel message to a tabContent type
		tc := tabContent(msg)

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
			// Disable mouse wheel to enable terminal text selection
			h.viewport.MouseWheelEnabled = false
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

	// Update viewport with all messages since mouse is disabled
	h.viewport, cmd = h.viewport.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return h, tea.Batch(cmds...)
}

func (h *DevTUI) updateViewport() {
	h.viewport.SetContent(h.ContentView())
	h.viewport.GotoBottom()
}

func (h *DevTUI) editingConfigOpen(open bool, currentField *field, msg string) {

	if open {
		h.editModeActivated = true
	} else {
		h.editModeActivated = false
	}

	if currentField != nil {
		currentField.setCursorAtEnd()
	}

	if msg != "" {
		tabSection := h.tabSections[h.activeTab]
		tabSection.addNewContent(Msg.Warning, msg)
	}

}
