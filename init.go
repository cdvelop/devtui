package devtui

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// Init inicializa el modelo
func (h *DevTUI) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		h.listenToMessages(),
		h.tickEverySecond(),
	)
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
