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
}

type TuiConfig struct {
	AppName       string    // app name eg: "MyApp"
	TabIndexStart int       // is the index of the tab section to start default 0
	ExitChan      chan bool //  global chan to close app eg: make(chan bool)
	/* *ColorStyle style for the TUI
	 eg:
	type ColorStyle struct {
	 Foreground string // eg: #F4F4F4
	 Background string // eg: #000000
	 Highlight  string // eg: #FF6600
	 Lowlight   string // eg: #666666
	}*/
	Color *ColorStyle

	LogToFile func(messages ...any) // function to write log error
}

// NewTUI creates a new DevTUI instance.
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
//	// To start the TUI:
//	if err := tui.tea.Start(); err != nil {
//	    config.LogToFile(err)
//	}
//
//	// To close the TUI from another goroutine:
//	config.ExitChan <- true
//
// You can customize fields, tabs, and handlers after creation.
// See the documentation for more advanced usage.
func NewTUI(c *TuiConfig) *DevTUI {
	if c.AppName == "" {
		c.AppName = "DevTUI"
	}

	// Example: create a default tab section using the new API
	tmpTUI := &DevTUI{TuiConfig: c}
	defaultTab := tmpTUI.NewTabSection(defaultTabName, "build footer example")
	defaultTab.NewField(
		"Editable Field",
		"initial editable value",
		true,
		func(newValue any) (string, error) {
			strValue := newValue.(string)
			return fmt.Sprintf("Value changed to %s", strValue), nil
		},
	).
		NewField(
			"Non-Editable Field",
			"non-editable value",
			false,
			func(newValue any) (string, error) {
				return "Action executed", nil
			},
		)

	tui := &DevTUI{
		TuiConfig: c,
		focused:   true, // assume the app is focused
		tabSections: []*tabSection{
			defaultTab,
		},
		activeTab:       c.TabIndexStart,
		tabContentsChan: make(chan tabContent, 100),
		currentTime:     time.Now().Format("15:04:05"),
		tuiStyle:        newTuiStyle(c.Color),
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

	return tui
}
