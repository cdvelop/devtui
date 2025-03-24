package devtui

import (
	"sync"
	"time"

	"github.com/cdvelop/unixid"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// channelMsg es un tipo especial para mensajes del canal
type channelMsg tuiMessage

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

	tabSections       []tabSection // represent sections in the tui
	activeTab         int          // current tab index
	editModeActivated bool         // global flag to edit config

	currentTime     string
	tabContentsChan chan tuiMessage
	// Channel for async field value change messages
	asyncMessageChan chan tuiMessage
	// Message tracker for handling message updates
	messageTracker *messageTracker
	tea            *tea.Program
}

type TuiConfig struct {
	AppName       string    // app name eg: "MyApp"
	TabIndexStart int       // is the index of the tab section to start
	ExitChan      chan bool //  global chan to close app
	Color         *ColorStyle

	LogToFile func(messageErr any) // function to write log error
}

func NewTUI(c *TuiConfig) *DevTUI {
	if c.AppName == "" {
		c.AppName = "DevTUI"
	}

	// Create message tracker
	msgTracker := NewMessageTracker()

	tui := &DevTUI{
		TuiConfig: c,
		focused:   true, // assume the app is focused
		tabSections: []tabSection{ // default tab section
			{
				title:         defaultTabName,
				fieldHandlers: []fieldHandler{},
				tuiMessages:   []tuiMessage{},
			},
		},
		activeTab:        c.TabIndexStart,
		tabContentsChan:  make(chan tuiMessage, 100),
		asyncMessageChan: make(chan tuiMessage),
		messageTracker:   msgTracker,
		currentTime:      time.Now().Format("15:04:05"),
		tuiStyle:         newTuiStyle(c.Color),
	}

	tui.tea = tea.NewProgram(tui,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	// Initialize the unique ID generator
	id, err := unixid.NewUnixID(sync.Mutex{})
	if err != nil {
		c.LogToFile(err)
	}
	tui.id = id

	// Set the ID generator for the message tracker
	tui.messageTracker.SetIDGenerator(tui.id)

	// Initialize the default tab section
	tui.tabSections[0].tui = tui

	return tui
}
