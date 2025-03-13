package devtui

import (
	"fmt"
	"time"

	. "github.com/cdvelop/messagetype"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// channelMsg es un tipo especial para mensajes del canal
type channelMsg tabContent

// Print representa un mensaje de actualización
type tickMsg time.Time

// tabContent imprime un mensaje en la tui
type tabContent struct {
	Content    string
	Type       MessageType
	Time       time.Time
	tabSection *TabSection
}

// DevTUI mantiene el estado de la aplicación
type DevTUI struct {
	*TuiConfig
	*tuiStyle

	ready    bool
	viewport viewport.Model

	activeTab        int  // current tab index
	tabEditingConfig bool // global flag to edit config

	SectionFooter   string
	currentTime     string
	tabContentsChan chan tabContent
	tea             *tea.Program
}

// represent the tab section in the tui
type TabSection struct {
	index         int            // index of the tab
	Title         string         // eg: "BUILD", "TEST"
	FieldHanlders []FieldHanlder // Field actions configured for the section
	SectionFooter string         // eg: "Press 't' to compile", "Press 'r' to run tests"
	// internal use
	tabContents          []tabContent // message contents
	indexActiveEditField int          // Índice del campo de configuración seleccionado
	tui                  *DevTUI
}

// Interface for handling tab field sectionFields
type FieldHanlder struct {
	Name             string                                               // eg: "port", "Server Port", "8080"
	Label            string                                               // eg: "Server Port"
	Value            string                                               // eg: "8080"
	Editable         bool                                                 // if no editable eject the action FieldValueChange directly
	FieldValueChange func(newValue string) (runMessage string, err error) //eg: "8080" -> "9090" runMessage: "Port changed from 8080 to 9090"
	//internal use
	index  int
	cursor int // cursor position in text value
}

type TuiConfig struct {
	TabIndexStart int          // is the index of the tab to start
	ExitChan      chan bool    //  global chan to close app
	TabSections   []TabSection // represent sections in the tui
	Color         *ColorStyle

	LogToFile func(messageErr string) // function to write log error
}

func NewTUI(c *TuiConfig) *DevTUI {

	// Create default tab if no tabs provided
	if len(c.TabSections) == 0 {
		defaultTab := TabSection{
			Title: "BUILD",
			FieldHanlders: []FieldHanlder{
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
		}
		c.TabSections = append(c.TabSections, defaultTab)

		testTab := TabSection{
			Title:         "DEPLOY",
			FieldHanlders: []FieldHanlder{},
			SectionFooter: "deploy footer example",
			tabContents:   []tabContent{},
		}
		c.TabSections = append(c.TabSections, testTab)

		c.TabIndexStart = 0 // Set the default tab index to 0
	}

	tui := &DevTUI{
		TuiConfig:       c,
		activeTab:       c.TabIndexStart,
		tabContentsChan: make(chan tabContent, 100),
		currentTime:     time.Now().Format("15:04:05"),
		tuiStyle:        newTuiStyle(c.Color),
	}

	// Recorremos c.TabSections y actualizamos el índice de cada campo.
	for i := range c.TabSections {
		section := &c.TabSections[i]
		section.index = i
		section.tui = tui
		for j := range section.FieldHanlders {
			section.FieldHanlders[j].index = j
			section.FieldHanlders[j].cursor = 0
		}
		// Si es necesario asignar otros valores, se hace aquí.
	}

	tui.tea = tea.NewProgram(tui,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	return tui
}
