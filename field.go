package devtui

import (
	"context"
	"fmt"
	"time"

	"github.com/cdvelop/messagetype"
)

// FieldHandler interface defines the contract for field handlers
// This replaces the individual parameters approach with a unified interface
type FieldHandler interface {
	Label() string                                                             // Field label (e.g., "Server Port")
	Value() string                                                             // Current field value (e.g., "8080")
	Editable() bool                                                            // Whether field is editable or action button
	Change(newValue any, progress ...func(string, ...float64)) (string, error) // Handler with optional progress callback
	Timeout() time.Duration                                                    // Return 0 for no timeout, or specific duration

	// NEW: WritingHandler methods (REQUIRED for all handlers)
	WritingHandler
}

// WritingHandler interface provides message source identification and operation ID management
// ALL handlers must implement this interface for message source control
type WritingHandler interface {
	Name() string                       // Handler identifier (e.g., "TinyWasm", "MainServer")
	SetLastOperationID(lastOpID string) // DevTUI calls this after processing each message
	GetLastOperationID() string         // Handler returns ID for message updates, "" for new messages
}

// Internal async state management (not exported)
type internalAsyncState struct {
	isRunning   bool
	operationID string
	cancel      context.CancelFunc
	startTime   time.Time
}

// use NewField to create a new field in the tab section
// Field represents a field in the TUI with a handler-based approach
// field represents a field in the TUI with async capabilities
type field struct {
	// NEW: Handler-based approach (replaces name, value, editable, changeFunc)
	handler   FieldHandler // Handles all field behavior
	parentTab *tabSection  // Direct reference to parent for message routing

	// NEW: Internal async state
	asyncState *internalAsyncState

	// UNCHANGED: Existing internal fields
	tempEditValue string // use for edit
	index         int
	cursor        int // cursor position in text value
}

// SetTempEditValueForTest permite modificar tempEditValue en tests
func (f *field) SetTempEditValueForTest(val string) {
	f.tempEditValue = val
}

// SetCursorForTest permite modificar el cursor en tests
func (f *field) SetCursorForTest(cursor int) {
	f.cursor = cursor
}

// NewField creates a new field with handler-based approach, adds it to the tabSection, and returns the tabSection for chaining.
// Example usage:
//
//	tab.NewField(&MyHandler{})
func (ts *tabSection) NewField(handler FieldHandler) *tabSection {
	f := &field{
		handler:    handler,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}

	// AUTO-REGISTER: FieldHandlers are automatically registered for writing
	// Initialize writingHandlers map if needed
	if ts.writingHandlers == nil {
		ts.writingHandlers = make(map[string]WritingHandler)
	}
	// Register the handler for writing capability
	ts.writingHandlers[handler.Name()] = handler

	ts.addFields(f)
	return ts
}

// setFieldHandlers sets the field handlers slice (mainly for testing)
// Only for internal/test use
func (ts *tabSection) setFieldHandlers(handlers []*field) {
	ts.fieldHandlers = handlers
}

// addFields adds one or more field handlers to the section (private)
func (ts *tabSection) addFields(fields ...*field) {
	ts.fieldHandlers = append(ts.fieldHandlers, fields...)
}

func (f *field) Name() string {
	if f.handler != nil {
		return f.handler.Label()
	}
	return ""
}

func (f *field) SetName(name string) {
	// This method is deprecated with handler-based approach
	// Handler manages its own label state
}

func (f *field) Value() string {
	if f.handler != nil {
		return f.handler.Value()
	}
	return ""
}

func (f *field) SetValue(value string) {
	// This method is deprecated with handler-based approach
	// Handler manages its own value state
}

// GetHandlerForTest returns the handler for testing purposes
func (f *field) GetHandlerForTest() FieldHandler {
	return f.handler
}

func (f *field) Editable() bool {
	if f.handler != nil {
		return f.handler.Editable()
	}
	return false
}

func (f *field) SetEditable(editable bool) {
	// This method is deprecated with handler-based approach
	// Handler manages its own editable state
}

func (f *field) SetCursorAtEnd() {
	// Calculate cursor position based on rune count, not byte count
	if f.handler != nil {
		f.cursor = len([]rune(f.handler.Value()))
	}
}

// getCurrentValue returns the appropriate value for Change() method
func (f *field) getCurrentValue() any {
	if f.handler == nil {
		return ""
	}

	if f.handler.Editable() {
		// For editable fields, return the edited text (tempEditValue or current value)
		// This matches current field behavior with tempEditValue
		// Check if we're in editing mode by looking at parent tab's edit state
		if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.editModeActivated {
			// In edit mode, always use tempEditValue (even if empty string)
			return f.tempEditValue
		}
		return f.handler.Value()
	} else {
		// For non-editable fields (action buttons), return the original value
		return f.handler.Value()
	}
}

// sendProgressMessage sends a progress message through parent tab
func (f *field) sendProgressMessage(content string) {
	if f.parentTab != nil && f.parentTab.tui != nil && f.asyncState != nil {
		handlerName := ""
		if f.handler != nil {
			handlerName = f.handler.Name()
		}

		f.parentTab.tui.sendMessageWithHandler(content, messagetype.Info, f.parentTab, handlerName, f.asyncState.operationID)
	}
}

// sendErrorMessage sends an error message through parent tab
func (f *field) sendErrorMessage(content string) {
	if f.parentTab != nil && f.parentTab.tui != nil {
		var operationID string
		if f.asyncState != nil {
			operationID = f.asyncState.operationID
		}

		handlerName := ""
		if f.handler != nil {
			handlerName = f.handler.Name()
		}

		f.parentTab.tui.sendMessageWithHandler(content, messagetype.Error, f.parentTab, handlerName, operationID)
	}
}

// sendSuccessMessage sends a success message through parent tab
func (f *field) sendSuccessMessage(content string) {
	if f.parentTab != nil && f.parentTab.tui != nil {
		var operationID string
		if f.asyncState != nil {
			operationID = f.asyncState.operationID
		}

		handlerName := ""
		if f.handler != nil {
			handlerName = f.handler.Name()
		}

		f.parentTab.tui.sendMessageWithHandler(content, messagetype.Success, f.parentTab, handlerName, operationID)
	}
}

// executeAsyncChange executes the handler's Change method asynchronously
func (f *field) executeAsyncChange(valueToSave any) {
	if f.handler == nil || f.asyncState == nil {
		return
	}

	// In test mode, execute synchronously for predictable test behavior
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.TestMode {
		f.executeChangeSyncWithValue(valueToSave)
		return
	}

	// Create internal context with timeout from handler
	timeout := f.handler.Timeout()
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	f.asyncState.cancel = cancel
	f.asyncState.isRunning = true

	// Generate ONE operation ID for the entire async operation OR reuse existing one
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.id != nil {
		// Check if handler has existing operationID to reuse (for updates)
		if existingID := f.handler.GetLastOperationID(); existingID != "" {
			f.asyncState.operationID = existingID
		} else {
			// Generate new ID for new operations
			f.asyncState.operationID = f.parentTab.tui.id.GetNewID()
		}
	}
	f.asyncState.startTime = time.Now()

	// Create progress callback for handler
	progressCallback := func(message string, percent ...float64) {
		f.sendProgressMessage(message)
	}

	// Use the pre-captured value instead of getCurrentValue()
	currentValue := valueToSave

	// Execute user's Change method with context monitoring
	resultChan := make(chan struct {
		result string
		err    error
	}, 1)

	go func() {
		result, err := f.handler.Change(currentValue, progressCallback)
		resultChan <- struct {
			result string
			err    error
		}{result, err}
	}()

	// Wait for completion or timeout
	select {
	case res := <-resultChan:
		// Operation completed normally
		f.asyncState.isRunning = false

		if res.err != nil {
			// Handler decides error message content
			f.sendErrorMessage(res.err.Error())
		} else {
			// Handler decides success message content
			f.sendSuccessMessage(res.result)
		}

	case <-ctx.Done():
		// Operation timed out
		f.asyncState.isRunning = false

		if ctx.Err() == context.DeadlineExceeded {
			f.sendErrorMessage(fmt.Sprintf("Operation timed out after %v", timeout))
		} else {
			f.sendErrorMessage("Operation was cancelled")
		}
	}

	cancel() // Clean up context
}

// executeChangeSync executes the handler's Change method synchronously (for tests)
func (f *field) executeChangeSync() {
	if f.handler == nil {
		return
	}

	// Generate operation ID for message routing (same as async version)
	if f.asyncState != nil && f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.id != nil {
		// Check if handler has existing operationID to reuse (for updates)
		if existingID := f.handler.GetLastOperationID(); existingID != "" {
			f.asyncState.operationID = existingID
		} else {
			// Generate new ID for new operations
			f.asyncState.operationID = f.parentTab.tui.id.GetNewID()
		}
	}

	// Get current value based on field type
	currentValue := f.getCurrentValue()

	// Create empty progress callback for sync execution
	progressCallback := func(message string, percent ...float64) {
		// In sync mode, we could log progress but not send to UI to avoid race conditions
	}

	// Execute user's Change method synchronously
	result, err := f.handler.Change(currentValue, progressCallback)

	if err != nil {
		// Handler decides error message content
		f.sendErrorMessage(err.Error())
	} else {
		// Handler decides success message content
		f.sendSuccessMessage(result)
	}
}

// executeChangeSyncWithValue executes the handler's Change method synchronously with pre-captured value
func (f *field) executeChangeSyncWithValue(valueToSave any) {
	if f.handler == nil {
		return
	}

	// In sync test mode, we don't generate operation IDs or send messages to avoid race conditions
	// Use the pre-captured value directly

	// Create empty progress callback for sync test execution
	progressCallback := func(message string, percent ...float64) {
		// In sync test mode, we don't send messages to avoid race conditions
	}

	_, err := f.handler.Change(valueToSave, progressCallback)

	// In test mode, we don't send messages to UI to avoid race conditions
	// The test can verify the handler's internal state directly
	_ = err // We still execute the handler but don't send UI messages
}

// handleEnter triggers async operation when user presses Enter
func (f *field) handleEnter() {
	if f.handler == nil {
		return
	}

	// Capture the current value BEFORE any state changes
	valueToSave := f.getCurrentValue()

	// In test mode, execute synchronously without goroutine
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.TestMode {
		f.executeChangeSyncWithValue(valueToSave)
		return
	}

	// DevTUI handles async internally - user doesn't see this complexity
	go f.executeAsyncChange(valueToSave)
}
