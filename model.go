package devtui

import (
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (cf *SectionField) SetCursorAtEnd() {
	cf.cursor = len(cf.Value)
}

// Init inicializa el modelo
func (h *DevTUI) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		h.listenToMessages(),
		h.tickEverySecond(),
	)
}

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

// Add this helper function
func (h *DevTUI) addTerminalPrint(msgType MessageType, content string) {
	h.TabSections[h.activeTab].tabContents = append(
		h.TabSections[h.activeTab].tabContents,
		tabContent{
			Type:    msgType,
			Content: content,
			Time:    time.Now(),
		},
	)
}

func (h *DevTUI) cancelEditingConfig(cancel bool) {
	if cancel {
		h.tabEditingConfig = false
		h.addTerminalPrint(OkMsg, "Exited config editing mode")
	} else {
		h.tabEditingConfig = true
		h.addTerminalPrint(WarnMsg, "Entered config editing mode")
	}
}

func (h *DevTUI) updateViewport() {
	h.viewport.SetContent(h.ContentView())
	h.viewport.GotoBottom()
}

func (h *DevTUI) StartTUI(wg *sync.WaitGroup) {
	defer wg.Done()

	if _, err := h.tea.Run(); err != nil {
		fmt.Println("Error running goCompiler:", err)
		fmt.Println("\nPress any key to exit...")
		var input string
		fmt.Scanln(&input)
	}
}
