package devtui

import (
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/unixid"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// channelMsg es un tipo especial para mensajes del canal
type channelMsg tabContent

// Print representa un mensaje de actualización
type tickMsg time.Time

// DevTUI mantiene el estado de la aplicación
type DevTUI struct {
	*TuiConfig
	*tuiStyle

	id *unixid.UnixID

	ready    bool
	viewport viewport.Model

	focused bool // is the app focused

	tabSections       []*tabSection // represent sections in the tui
	activeTab         int           // current tab index
	editModeActivated bool          // global flag to edit config

	currentTime     string
	tabContentsChan chan tabContent
	tea             *tea.Program
	testMode        bool // private: only used in tests to enable synchronous behavior
}

type TuiConfig struct {
	AppName  string    // app name eg: "MyApp"
	ExitChan chan bool //  global chan to close app eg: make(chan bool)
	/*// *ColorStyle style for the TUI
	  // if nil it will use default style:
	type ColorStyle struct {
	 Foreground string // eg: #F4F4F4
	 Background string // eg: #000000
	 Highlight  string // eg: #FF6600
	 Lowlight   string // eg: #666666
	}*/
	Color *ColorStyle

	LogToFile func(messages ...any) // function to write log error
}

// NewTUI creates a new DevTUI instance and initializes it.
//
// Usage Example:
//
//	config := &TuiConfig{
//	    AppName: "MyApp",
//	    ExitChan: make(chan bool),
//	    Color: nil, // or your *ColorStyle
//	    LogToFile: func(err any) { fmt.Println(err) },
//	}
//	tui := NewTUI(config)
//
//	// Configure your sections and fields:
//	tui.NewTabSection("My Section", "Description").
//		NewField("Field1", "value", true, nil)
//
//	// Start the TUI:
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go tui.Start(&wg)
//	wg.Wait()
func NewTUI(c *TuiConfig) *DevTUI {
	if c.AppName == "" {
		c.AppName = "DevTUI"
	}

	// Initialize the unique ID generator first
	id, err := unixid.NewUnixID()
	if err != nil {
		if c.LogToFile != nil {
			c.LogToFile("Critical: Error initializing unixid:", err, "- timestamp generation will use fallback")
		}
		// id will remain nil, but createTabContent method will handle this gracefully now
	} else {
		if c.LogToFile != nil {
			c.LogToFile("Success: UnixID initialized correctly")
		}
	}

	tui := &DevTUI{
		TuiConfig:       c,
		focused:         true, // assume the app is focused
		tabSections:     []*tabSection{},
		activeTab:       0, // Will be adjusted in Start() method
		tabContentsChan: make(chan tabContent, 100),
		currentTime:     time.Now().Format("15:04:05"),
		tuiStyle:        newTuiStyle(c.Color),
		id:              id, // Set the ID here
	}

	// Always add SHORTCUTS tab first
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")
	shortcutsHandler := NewShortcutsHandler()
	shortcutsTab.NewDisplayHandler(shortcutsHandler).Register()

	// FIXED: Removed manual content sending to prevent duplication
	// HandlerDisplay automatically shows Content() when field is selected
	// No need for manual sendMessageWithHandler() call

	tui.tea = tea.NewProgram(tui,
		tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
		// Mouse support disabled to enable terminal text selection
	)

	return tui
}

// Init initializes the terminal UI application.
func (h *DevTUI) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		h.listenToMessages(),
		h.tickEverySecond(),
	)
}

// Start initializes and runs the terminal UI application.
//
// It accepts optional variadic arguments of any type. If a *sync.WaitGroup
// is provided among these arguments, Start will call its Done() method
// before returning.
//
// The method runs the UI using the internal tea engine, and handles any
// errors that may occur during execution. If an error occurs, it will be
// displayed on the console and the application will wait for user input
// before exiting.
//
// Parameters:
//   - args ...any: Optional arguments. Can include a *sync.WaitGroup for synchronization.
func (h *DevTUI) Start(args ...any) {
	// Check if a WaitGroup was passed
	for _, arg := range args {
		if wg, ok := arg.(*sync.WaitGroup); ok {
			defer wg.Done()
			break
		}
	}

	// Start with tab 1 (skip SHORTCUTS which is at index 0) if there are multiple tabs
	if len(h.tabSections) > 1 {
		h.activeTab = 1
	}

	if _, err := h.tea.Run(); err != nil {
		fmt.Println("Error running DevTUI:", err)
		fmt.Println("\nPress any key to exit...")
		var input string
		fmt.Scanln(&input)
	}
}

// SetTestMode enables or disables test mode for synchronous behavior in tests.
// This should only be used in test files to make tests deterministic.
func (h *DevTUI) SetTestMode(enabled bool) {
	h.testMode = enabled
}

// isTestMode returns true if the TUI is running in test mode (synchronous execution).
// This is an internal method used by field handlers to determine execution mode.
func (h *DevTUI) isTestMode() bool {
	return h.testMode
}
