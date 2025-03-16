package devtui

import (
	"time"

	"github.com/cdvelop/messagetype"
)

const defaultTabName = "DEFAULT"

// Interface for handling tab field sectionFields
type FieldHandler struct {
	Name             string                                                // eg: "port", "Server Port", "8080"
	Label            string                                                // eg: "Server Port"
	Value            string                                                // eg: "8080"
	tempEditValue    string                                                // use for edit
	Editable         bool                                                  // if no editable eject the action FieldValueChange directly
	FieldValueChange func(newValue string) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port changed from 8080 to 9090"
	//internal use
	index  int
	cursor int // cursor position in text value
}

// tabContent imprime un mensaje en la tui
type tabContent struct {
	Content    string
	Type       messagetype.MessageType
	Time       time.Time
	tabSection *TabSection
}

// represent the tab section in the tui
type TabSection struct {
	index         int            // index of the tab
	Title         string         // eg: "BUILD", "TEST"
	FieldHandlers []FieldHandler // Field actions configured for the section
	SectionFooter string         // eg: "Press 't' to compile", "Press 'r' to run tests"
	// internal use
	tabContents          []tabContent // message contents
	indexActiveEditField int          // Índice del campo de configuración seleccionado
	tui                  *DevTUI
}

// AddTabSections adds one or more TabSections to the DevTUI
// If a tab with title "DEFAULT" exists, it will be replaced by the first tab section
func (t *DevTUI) AddTabSections(sections ...TabSection) *DevTUI {
	if len(sections) == 0 {
		return t
	}

	// Check if there's a "DEFAULT" tab to replace
	defaultTabIndex := -1
	for i, tab := range t.tabSections {
		if tab.Title == defaultTabName {
			defaultTabIndex = i
			break
		}
	}

	// Replace DEFAULT tab if found
	if defaultTabIndex >= 0 && len(sections) > 0 {
		// Initialize first section for replacement
		newSection := t.initTabSection(sections[0], defaultTabIndex)
		t.tabSections[defaultTabIndex] = newSection

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

// Helper method to initialize a single TabSection
func (t *DevTUI) initTabSection(section TabSection, index int) TabSection {
	newSection := section
	newSection.index = index
	newSection.tui = t

	// Initialize field handlers
	for j := range newSection.FieldHandlers {
		newSection.FieldHandlers[j].index = j
		newSection.FieldHandlers[j].cursor = 0
	}

	return newSection
}

// Helper method to add multiple tab sections
func (t *DevTUI) addNewTabSections(sections ...TabSection) {
	startIdx := len(t.tabSections)
	for i, section := range sections {
		newSection := t.initTabSection(section, startIdx+i)
		t.tabSections = append(t.tabSections, newSection)
	}
}

// GetTotalTabSections returns the total number of tab sections
func (t *DevTUI) GetTotalTabSections() int {
	return len(t.tabSections)
}
