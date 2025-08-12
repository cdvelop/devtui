package devtui

import (
	"io"
	"time"
)

/*
AddDisplayHandler registers a HandlerDisplay directly

	type HandlerDisplay interface {
	    Name() string    // Full text to display in footer (handler responsible for content) eg. "System Status Information Display"
	    Content() string // Display content (e.g., "help\n1-..\n2-...", "executing deploy wait...")
	}
*/
func (ts *tabSection) AddDisplayHandler(handler HandlerDisplay) *tabSection {
	anyH := newDisplayHandler(handler)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

/*
AddEditHandler registers a HandlerEdit with mandatory timeout.

	type HandlerEdit interface {
	    Name() string                                       // Identificador: "ServerPort", "DatabaseURL"
	    Label() string                                      // Field label (e.g., "Server Port", "Host Configuration")
	    Value() string                                      // Current/initial value (e.g., "8080", "localhost")
	    Change(newValue string, progress func(msgs ...any)) // value to change
	}

ts.AddEditHandler(myEditHandler, 2*time.Second)
*/
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration) *tabSection {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := newEditHandler(handler, timeout, tracker)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)

	// NEW: Check for shortcut support and register shortcuts
	ts.registerShortcutsIfSupported(handler, len(ts.fieldHandlers)-1)

	// REMOVED: Auto-register handler for writing if it implements HandlerWriterTracker (obsolete)

	return ts
}

/*
AddEditHandlerTracking registers a HandlerEditTracker with mandatory timeout

	 type HandlerExecution interface {
	    Name() string                       // Identificador : "DeployProd", "BuildProject"
	    Label() string                      // Button label (e.g., "Deploy to Production", "Build Project")
	    Execute(progress func(msgs ...any)) // Nueva firma: sin error, sin variádico
	}

eg: ts.AddEditHandlerTracking(myEditHandler, 2*time.Second)
*/
func (ts *tabSection) AddExecutionHandler(handler HandlerExecution, timeout time.Duration) *tabSection {
	anyH := newExecutionHandler(handler, timeout)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

// AddExecutionHandlerTracking registers a HandlerExecutionTracker with mandatory timeout
// NewWriter creates a writer with the given name and tracking capability
// enableTracking: true = can update existing lines, false = always creates new lines
func (ts *tabSection) NewWriter(name string, enableTracking bool) io.Writer {
	if enableTracking {
		handler := &simpleWriterTrackerHandler{name: name}
		return ts.registerWriter(handler)
	} else {
		handler := &simpleWriterHandler{name: name}
		return ts.registerWriter(handler)
	}
}

// Internal simple handler implementations
type simpleWriterHandler struct {
	name string
}

func (w *simpleWriterHandler) Name() string {
	return w.name
}

type simpleWriterTrackerHandler struct {
	name            string
	lastOperationID string
}

func (w *simpleWriterTrackerHandler) Name() string {
	return w.name
}

func (w *simpleWriterTrackerHandler) GetLastOperationID() string {
	return w.lastOperationID
}

func (w *simpleWriterTrackerHandler) SetLastOperationID(id string) {
	w.lastOperationID = id
}

/*
AddInteractiveHandler registers a HandlerInteractive with mandatory timeout

	type HandlerInteractive interface {
	    Name() string                                       // Identifier for logging: "ChatBot", "ConfigWizard"
	    Label() string                                      // Field label (updates dynamically)
	    Value() string                                      // Current input value
	    Change(newValue string, progress func(msgs ...any)) // Handle user input + content display via progress
	    WaitingForUser() bool                               // Should edit mode be auto-activated?
	}
*/
func (ts *tabSection) AddInteractiveHandler(handler HandlerInteractive, timeout time.Duration) *tabSection {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := newInteractiveHandler(handler, timeout, tracker)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

// AddInteractiveHandlerTracking registers a HandlerInteractiveTracker with mandatory timeout
// registerShortcutsIfSupported checks if handler implements shortcut interface and registers shortcuts
func (ts *tabSection) registerShortcutsIfSupported(handler HandlerEdit, fieldIndex int) {
	// Check if handler implements shortcut interface
	if shortcutProvider, hasShortcuts := handler.(ShortcutProvider); hasShortcuts {
		shortcuts := shortcutProvider.Shortcuts()
		for key, description := range shortcuts {
			entry := &ShortcutEntry{
				Key:         key,
				Description: description,
				TabIndex:    ts.index,
				FieldIndex:  fieldIndex,
				HandlerName: handler.Name(),
				Value:       key, // Use the key as the value by default
			}
			ts.tui.shortcutRegistry.Register(key, entry)
		}
	}
}
