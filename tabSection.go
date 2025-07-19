package devtui

import (
	"strings"

	"github.com/cdvelop/messagetype"
)

const defaultTabName = "DEFAULT"

// Interface for handling tab field sectionFields

// tabContent imprime contenido en la tui con id único
type tabContent struct {
	Id         string
	Content    string
	Type       messagetype.Type
	tabSection *tabSection
}

// tabSection represents a tab section in the TUI with configurable fields and content
type tabSection struct {
	index         int      // index of the tab
	title         string   // eg: "BUILD", "TEST"
	fieldHandlers []*field // Field actions configured for the section
	sectionFooter string   // eg: "Press 't' to compile", "Press 'r' to run tests"
	// internal use
	tabContents          []tabContent // message contents
	indexActiveEditField int          // Índice del campo de configuración seleccionado
	tui                  *DevTUI
}

// Write implementa io.Writer para capturar la salida de otros procesos
func (ts *tabSection) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		// Detectar automáticamente el tipo de mensaje
		msgType := messagetype.DetectMessageType(msg)

		ts.tui.sendMessage(msg, msgType, ts)
		// Si es un error, escribirlo en el archivo de log
		if msgType == messagetype.Error {
			ts.tui.LogToFile(msg)
		}

	}
	return len(p), nil
}

func (t *tabSection) addNewContent(msgType messagetype.Type, content string) {
	t.tabContents = append(t.tabContents, t.tui.newContent(content, msgType, t))
}

// Title returns the tab section title
func (ts *tabSection) Title() string {
	return ts.title
}

// SetTitle sets the tab section title
func (ts *tabSection) SetTitle(title string) {
	ts.title = title
}

// Footer returns the tab section footer text
func (ts *tabSection) Footer() string {
	return ts.sectionFooter
}

// SetFooter sets the tab section footer text
func (ts *tabSection) SetFooter(footer string) {
	ts.sectionFooter = footer
}

// FieldHandlers returns the field handlers slice
func (ts *tabSection) FieldHandlers() []*field {
	return ts.fieldHandlers
}

// NewTabSection creates and initializes a new tabSection with the given title and footer
// Example:
//
//	tab := tui.NewTabSection("BUILD", "Press enter to compile")
func (t *DevTUI) NewTabSection(title, footer string) *tabSection {
	return &tabSection{
		title:         title,
		sectionFooter: footer,
		tui:           t,
	}
}

// SetIndex sets the index of the tab section
func (ts *tabSection) SetIndex(idx int) {
	ts.index = idx
}

// SetActiveEditField sets the active edit field index
func (ts *tabSection) SetActiveEditField(idx int) {
	ts.indexActiveEditField = idx
}

// AddTabSections adds one or more tabSections to the DevTUI
// If a tab with title "DEFAULT" exists, it will be replaced by the first tab section
// Deprecated: Use NewTabSection and append to tabSections directly
func (t *DevTUI) AddTabSections(sections ...*tabSection) *DevTUI {
	if len(sections) == 0 {
		return t
	}

	// Check if there's a "DEFAULT" tab to replace
	defaultTabIndex := -1
	for i, tab := range t.tabSections {
		if tab.Title() == defaultTabName {
			defaultTabIndex = i
			break
		}
	}

	// Replace DEFAULT tab if found
	if defaultTabIndex >= 0 && len(sections) > 0 {
		// Initialize first section for replacement
		t.initTabSection(sections[0], defaultTabIndex)
		t.tabSections[defaultTabIndex] = sections[0]

		// Add remaining sections
		if len(sections) > 1 {
			t.addNewTabSections(sections[1:]...)
		}
	} else {
		// Just add all sections normally
		t.addNewTabSections(sections...)
	}

	return t
}

// Helper method to initialize a single tabSection
func (t *DevTUI) initTabSection(section *tabSection, index int) {
	section.index = index
	section.tui = t

	// Initialize field handlers
	handlers := section.FieldHandlers()
	for j := range handlers {
		handlers[j].index = j
		handlers[j].cursor = 0
	}
	section.setFieldHandlers(handlers)
}

// Helper method to add multiple tab sections
func (t *DevTUI) addNewTabSections(sections ...*tabSection) {
	startIdx := len(t.tabSections)
	for i, section := range sections {
		section.index = startIdx + i
		section.tui = t
		t.tabSections = append(t.tabSections, section)
	}
}

// GetTotalTabSections returns the total number of tab sections
func (t *DevTUI) GetTotalTabSections() int {
	return len(t.tabSections)
}
