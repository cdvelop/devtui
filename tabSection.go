package devtui

import (
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/messagetype"
)

const defaultTabName = "DEFAULT"

// Interface for handling tab field sectionFields

// tabContent imprime contenido en la tui con id único
type tabContent struct {
	Id         string // unix number id eg: "1234567890" - INMUTABLE
	Timestamp  string // unix nano timestamp - MUTABLE (se actualiza en cada cambio)
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
	mu                   sync.RWMutex // Para proteger tabContents y writingHandlers de race conditions

	// Writing handler registry for external handlers using new interfaces
	writingHandlers []*anyHandler // CAMBIO: slice en lugar de map para thread-safety
	activeWriter    string        // current active writer name for io.Writer calls
}

// getWritingHandler busca un handler por nombre en el slice thread-safe
func (ts *tabSection) getWritingHandler(name string) *anyHandler {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, h := range ts.writingHandlers {
		if h.Name() == name {
			return h
		}
	}
	return nil
}

// registerAnyHandler registers an anyHandler in the thread-safe slice
func (ts *tabSection) registerAnyHandler(handler *anyHandler) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.writingHandlers = append(ts.writingHandlers, handler)
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

		if ts.activeWriter != "" {
			if handler := ts.getWritingHandler(ts.activeWriter); handler != nil {
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

// RegisterHandlerWriter registers a basic writer handler and returns a dedicated writer
func (ts *tabSection) RegisterHandlerWriter(handler HandlerWriter) io.Writer {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var anyH *anyHandler

	// Check if handler also implements HandlerTrackerWriter interface
	if trackerHandler, ok := handler.(HandlerTrackerWriter); ok {
		anyH = newTrackerWriterHandler(trackerHandler)
	} else {
		anyH = newWriterHandler(handler)
	}

	ts.writingHandlers = append(ts.writingHandlers, anyH)
	return &handlerWriter{tabSection: ts, handlerName: anyH.Name()}
}

// RegisterHandlerTrackerWriter registers an advanced writer handler with message tracking and returns a dedicated writer
func (ts *tabSection) RegisterHandlerTrackerWriter(handler HandlerTrackerWriter) io.Writer {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	anyH := newTrackerWriterHandler(handler)
	ts.writingHandlers = append(ts.writingHandlers, anyH)
	return &handlerWriter{tabSection: ts, handlerName: anyH.Name()}
}

// NEW: SetActiveWriter sets the current active writer for general io.Writer calls
func (ts *tabSection) SetActiveWriter(handlerName string) {
	ts.activeWriter = handlerName
}

// HandlerWriter wraps tabSection with handler identification
type handlerWriter struct {
	tabSection  *tabSection
	handlerName string
}

func (hw *handlerWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		msgType := messagetype.DetectMessageType(msg)

		var operationID string
		if handler := hw.tabSection.getWritingHandler(hw.handlerName); handler != nil {
			operationID = handler.GetLastOperationID()
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
	t.tabContents = append(t.tabContents, t.tui.createTabContent(content, msgType, t, "", ""))
}

// NEW: updateOrAddContentWithHandler updates existing content by operationID or adds new if not found
// Returns true if content was updated, false if new content was added
func (t *tabSection) updateOrAddContentWithHandler(msgType messagetype.Type, content string, handlerName string, operationID string) (updated bool, newContent tabContent) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// If operationID is provided, try to find and update existing content
	if operationID != "" {
		for i := range t.tabContents {
			// Match by both operationID and handlerName to ensure each handler updates its own message
			if t.tabContents[i].operationID != nil &&
				*t.tabContents[i].operationID == operationID &&
				t.tabContents[i].handlerName == handlerName {
				// Update existing content
				t.tabContents[i].Content = content
				t.tabContents[i].Type = msgType
				// Actualizar timestamp usando GetNewID directamente
				if t.tui.id != nil {
					t.tabContents[i].Timestamp = t.tui.id.GetNewID()
				} else {
					// Log the issue before using fallback
					if t.tui.LogToFile != nil {
						t.tui.LogToFile("Warning: unixid not initialized, using fallback timestamp for content update:", content)
					}
					// Graceful fallback when unixid initialization failed
					t.tabContents[i].Timestamp = time.Now().Format("15:04:05")
				}
				return true, t.tabContents[i]
			}
		}
	}

	// If not found or no operationID, add new content
	newContent = t.tui.createTabContent(content, msgType, t, handlerName, operationID)
	t.tabContents = append(t.tabContents, newContent)
	return false, newContent
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
