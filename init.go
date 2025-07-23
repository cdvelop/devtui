package devtui

import (
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/messagetype"
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
}

type TuiConfig struct {
	AppName       string    // app name eg: "MyApp"
	TabIndexStart int       // is the index of the tab section to start default 0
	ExitChan      chan bool //  global chan to close app eg: make(chan bool)
	/* *ColorStyle style for the TUI
	 if nil it will use default style:
	type ColorStyle struct {
	 Foreground string // eg: #F4F4F4
	 Background string // eg: #000000
	 Highlight  string // eg: #FF6600
	 Lowlight   string // eg: #666666
	}*/
	Color *ColorStyle

	LogToFile func(messages ...any) // function to write log error
	TestMode  bool                  // only used in tests to enable synchronous behavior
}

// NewTUI creates a new DevTUI instance and initializes it.
//
// Usage Example:
//
//	config := &TuiConfig{
//	    AppName: "MyApp",
//	    TabIndexStart: 0,
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
//	go tui.Run(&wg)
//	wg.Wait()
func NewTUI(c *TuiConfig) *DevTUI {
	if c.AppName == "" {
		c.AppName = "DevTUI"
	}

	// Initialize the unique ID generator first
	id, err := unixid.NewUnixID()
	if err != nil {
		if c.LogToFile != nil {
			c.LogToFile("Error initializing unixid:", err)
		}
		// id will remain nil, but newContent method will handle this
	} else {
		if c.LogToFile != nil {
			c.LogToFile("UnixID initialized successfully")
		}
	}

	tui := &DevTUI{
		TuiConfig:       c,
		focused:         true, // assume the app is focused
		tabSections:     []*tabSection{},
		activeTab:       c.TabIndexStart,
		tabContentsChan: make(chan tabContent, 100),
		currentTime:     time.Now().Format("15:04:05"),
		tuiStyle:        newTuiStyle(c.Color),
		id:              id, // Set the ID here
	}

	// Always add SHORTCUTS tab first
	shortcutsTab := tui.NewTabSection("SHORTCUTS", "Keyboard navigation instructions")
	shortcutsHandler := NewShortcutsHandler()
	shortcutsTab.NewField(shortcutsHandler)

	// Automatically display shortcuts content when tab is created (unless in test mode)
	// Use sendMessageWithHandler to respect readonly handler formatting
	if !c.TestMode {
		tui.sendMessageWithHandler(shortcutsHandler.shortcuts, messagetype.Info, shortcutsTab, shortcutsHandler.Name(), "")
	}

	tui.tea = tea.NewProgram(tui,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
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

	// If user didn't specify a custom TabIndexStart and we have more than 1 tab,
	// default to tab 1 (skip SHORTCUTS which is at index 0)
	if h.TuiConfig.TabIndexStart == 0 && len(h.tabSections) > 1 {
		h.activeTab = 1
	}

	if _, err := h.tea.Run(); err != nil {
		fmt.Println("Error running DevTUI:", err)
		fmt.Println("\nPress any key to exit...")
		var input string
		fmt.Scanln(&input)
	}
}
