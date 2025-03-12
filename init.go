package devtui

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the terminal UI application.
func (h *DevTUI) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		h.listenToMessages(),
		h.tickEverySecond(),
	)
}

// InitTUI initializes and runs the terminal UI application.
//
// It accepts optional variadic arguments of any type. If a *sync.WaitGroup
// is provided among these arguments, InitTUI will call its Done() method
// before returning.
//
// The method runs the UI using the internal tea engine, and handles any
// errors that may occur during execution. If an error occurs, it will be
// displayed on the console and the application will wait for user input
// before exiting.
//
// Parameters:
//   - args ...any: Optional arguments. Can include a *sync.WaitGroup for synchronization.
func (h *DevTUI) InitTUI(args ...any) {
	// Check if a WaitGroup was passed
	for _, arg := range args {
		if wg, ok := arg.(*sync.WaitGroup); ok {
			defer wg.Done()
			break
		}
	}

	if _, err := h.tea.Run(); err != nil {
		fmt.Println("Error running goCompiler:", err)
		fmt.Println("\nPress any key to exit...")
		var input string
		fmt.Scanln(&input)
	}
}
