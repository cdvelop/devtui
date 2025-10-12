package devtui

import "time"

// AddHandler is the ONLY method to register handlers of any type.
// It accepts any handler interface and internally detects the type.
// Does NOT return anything - enforces complete decoupling.
//
// Supported handler interfaces (from interfaces.go):
//   - HandlerDisplay: Static/dynamic content display
//   - HandlerEdit: Interactive text input fields
//   - HandlerExecution: Action buttons
//   - HandlerInteractive: Combined display + interaction
//   - HandlerLogger: Basic line-by-line logging (via MessageTracker detection)
//
// Optional interfaces (detected automatically):
//   - MessageTracker: Enables message update tracking
//   - ShortcutProvider: Registers global keyboard shortcuts
//
// Parameters:
//   - handler: ANY handler implementing one of the supported interfaces
//   - timeout: Operation timeout (used for Edit/Execution/Interactive handlers, ignored for Display)
//   - color: Hex color for handler messages (e.g., "#1e40af", empty string for default)
//
// Example:
//   tab.AddHandler(myEditHandler, 2*time.Second, "#3b82f6")
//   tab.AddHandler(myDisplayHandler, 0, "") // timeout ignored for display
//   tab.AddHandler(myExecutionHandler, 5*time.Second, "#10b981")
func (ts *tabSection) AddHandler(handler any, timeout time.Duration, color string) {
	// Type detection and routing
	switch h := handler.(type) {

	case HandlerDisplay:
		ts.registerDisplayHandler(h, color)

	case HandlerInteractive:
		ts.registerInteractiveHandler(h, timeout, color)

	case HandlerExecution:
		ts.registerExecutionHandler(h, timeout, color)

	case HandlerEdit:
		ts.registerEditHandler(h, timeout, color)

	case HandlerLogger:
		// Logger detection: check for MessageTracker to determine tracking capability
		_, hasTracking := handler.(MessageTracker)
		ts.registerLoggerHandler(h, color, hasTracking)

	default:
		// Invalid handler type - log error or panic
		if ts.tui != nil && ts.tui.Logger != nil {
			ts.tui.Logger("ERROR: Unknown handler type provided to AddHandler:", handler)
		}
	}
}

// Internal registration methods (private)

func (ts *tabSection) registerDisplayHandler(handler HandlerDisplay, color string) {
	anyH := NewDisplayHandler(handler, color)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
}

func (ts *tabSection) registerEditHandler(handler HandlerEdit, timeout time.Duration, color string) {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := NewEditHandler(handler, timeout, tracker, color)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)

	// Check for shortcut support
	ts.registerShortcutsIfSupported(handler, len(ts.fieldHandlers)-1)
}

func (ts *tabSection) registerExecutionHandler(handler HandlerExecution, timeout time.Duration, color string) {
	anyH := NewExecutionHandler(handler, timeout, color)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
}

func (ts *tabSection) registerInteractiveHandler(handler HandlerInteractive, timeout time.Duration, color string) {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := NewInteractiveHandler(handler, timeout, tracker, color)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
}

func (ts *tabSection) registerLoggerHandler(handler HandlerLogger, color string, hasTracking bool) {
	var anyH *anyHandler

	if hasTracking {
		// Handler implements MessageTracker
		if tracker, ok := handler.(interface {
			Name() string
			GetLastOperationID() string
			SetLastOperationID(string)
		}); ok {
			anyH = NewWriterTrackerHandler(tracker, color)
		} else {
			// This should not happen if hasTracking is true, but as a fallback:
			anyH = NewWriterHandler(handler, color)
		}
	} else {
		// Basic logger without tracking
		anyH = NewWriterHandler(handler, color)
	}

	// Register in writing handlers list
	ts.mu.Lock()
	ts.writingHandlers = append(ts.writingHandlers, anyH)
	ts.mu.Unlock()
}


// AddLogger creates a logger function with the given name and tracking capability
// enableTracking: true = can update existing lines, false = always creates new lines
//
// Example:
//
//	log := tab.AddLogger("BuildProcess", true, "#1e40af")
//	log("Starting build...")
//	log("Compiling", 42, "files")
//	log("Build completed successfully")
func (ts *tabSection) AddLogger(name string, enableTracking bool, color string) func(message ...any) {
	if enableTracking {
		handler := &simpleWriterTrackerHandler{name: name}
		return ts.registerLoggerFunc(handler, color)
	} else {
		handler := &simpleWriterHandler{name: name}
		return ts.registerLoggerFunc(handler, color)
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
