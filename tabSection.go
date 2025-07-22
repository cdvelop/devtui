package devtui

import (
	"io"
	"strings"
	"sync"

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

	// NEW: Async fields (always present, nil when not async)
	operationID *string // nil for sync messages, value for async operations
	isProgress  bool    // true if this is a progress update
	isComplete  bool    // true if async operation completed

	// NEW: Handler identification
	handlerName string // Handler name for message source identification
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
	mu                   sync.RWMutex // Para proteger tabContents de race conditions

	// NEW: Writing handler registry for external handlers
	writingHandlers map[string]WritingHandler // handlerName -> WritingHandler instance
	activeWriter    string                    // current active writer name for io.Writer calls
}

// Write implementa io.Writer para capturar la salida de otros procesos
func (ts *tabSection) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		// Detectar automáticamente el tipo de mensaje
		msgType := messagetype.DetectMessageType(msg)

		// NEW: Determine handler name and operation ID from active writer
		var handlerName string
		var operationID string

		if ts.activeWriter != "" && ts.writingHandlers != nil {
			if handler, exists := ts.writingHandlers[ts.activeWriter]; exists {
				handlerName = handler.Name()
				operationID = handler.GetLastOperationID()
			}
		}

		ts.tui.sendMessageWithHandler(msg, msgType, ts, handlerName, operationID)
		// Si es un error, escribirlo en el archivo de log
		if msgType == messagetype.Error {
			ts.tui.LogToFile(msg)
		}

	}
	return len(p), nil
}

// NEW: RegisterWritingHandler registers a writing handler and returns a dedicated writer
func (ts *tabSection) RegisterWritingHandler(handler WritingHandler) io.Writer {
	if ts.writingHandlers == nil {
		ts.writingHandlers = make(map[string]WritingHandler)
	}

	handlerName := handler.Name()
	ts.writingHandlers[handlerName] = handler

	// Return handler-specific writer
	return &HandlerWriter{
		tabSection:  ts,
		handlerName: handlerName,
	}
}

// NEW: SetActiveWriter sets the current active writer for general io.Writer calls
func (ts *tabSection) SetActiveWriter(handlerName string) {
	ts.activeWriter = handlerName
}

// NEW: HandlerWriter wraps tabSection with handler identification
type HandlerWriter struct {
	tabSection  *tabSection
	handlerName string
}

func (hw *HandlerWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		msgType := messagetype.DetectMessageType(msg)

		// Debug: Log the message and detected type
		println("DEBUG: Message:", msg, "Detected Type:", int(msgType), "Expected Success:", int(messagetype.Success))

		var operationID string
		if hw.tabSection.writingHandlers != nil {
			if handler, exists := hw.tabSection.writingHandlers[hw.handlerName]; exists {
				operationID = handler.GetLastOperationID()
			}
		}

		hw.tabSection.tui.sendMessageWithHandler(msg, msgType, hw.tabSection, hw.handlerName, operationID)

		if msgType == messagetype.Error {
			hw.tabSection.tui.LogToFile(msg)
		}
	}
	return len(p), nil
}

func (t *tabSection) addNewContent(msgType messagetype.Type, content string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tabContents = append(t.tabContents, t.tui.newContent(content, msgType, t))
}

// NEW: addNewContentWithHandler adds content with handler identification
func (t *tabSection) addNewContentWithHandler(msgType messagetype.Type, content string, handlerName string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tabContents = append(t.tabContents, t.tui.newContentWithHandler(content, msgType, t, handlerName))
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
// NewTabSection creates a new tab section and automatically adds it to the TUI
//
// Example:
//
//	tab := tui.NewTabSection("BUILD", "Press enter to compile")
func (t *DevTUI) NewTabSection(title, footer string) *tabSection {
	tab := &tabSection{
		title:         title,
		sectionFooter: footer,
		tui:           t,
	}

	// Automatically add to tabSections and initialize
	t.initTabSection(tab, len(t.tabSections))
	t.tabSections = append(t.tabSections, tab)

	return tab
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
