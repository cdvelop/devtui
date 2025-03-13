package devtui

import (
	"fmt"
	"time"

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
	tabSections []TabSection // represent sections in the tui
	*tuiStyle

	ready    bool
	viewport viewport.Model

	activeTab        int  // current tab index
	tabEditingConfig bool // global flag to edit config

	currentTime     string
	tabContentsChan chan tabContent
	tea             *tea.Program
}

type TuiConfig struct {
	TabIndexStart int       // is the index of the tab to start
	ExitChan      chan bool //  global chan to close app
	Color         *ColorStyle

	LogToFile func(messageErr any) // function to write log error
}

func NewTUI(c *TuiConfig) *DevTUI {

	tui := &DevTUI{
		TuiConfig: c,
		tabSections: []TabSection{ // default tab section
			{
				Title: "DEFAULT",
				FieldHandlers: []FieldHandler{
					{
						Name:     "editableField",
						Label:    "Editable Field",
						Value:    "initial editable value",
						Editable: true,
						FieldValueChange: func(newValue string) (string, error) {
							// Agregar la lógica de cambio de valor deseada
							return fmt.Sprintf("Value changed to %s", newValue), nil
						},
					},
					{
						Name:     "nonEditableField",
						Label:    "Non-Editable Field",
						Value:    "non-editable value",
						Editable: false,
						FieldValueChange: func(newValue string) (string, error) {
							// Agregar la acción deseada para el campo no editable
							return "Action executed", nil
						},
					},
				},
				SectionFooter: "build footer example",
				tabContents:   []tabContent{},
			},
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

	return tui
}
