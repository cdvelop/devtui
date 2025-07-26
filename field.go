package devtui

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/messagetype"
)

// ============================================================================
// PRIVATE IMPLEMENTATION - anyHandler Structure
// ============================================================================

type handlerType int

const (
	handlerTypeDisplay handlerType = iota
	handlerTypeEdit
	handlerTypeExecution
	handlerTypeWriter
	handlerTypeTrackerWriter
)

// anyHandler - Estructura privada que unifica todos los handlers
type anyHandler struct {
	handlerType handlerType
	timeout     time.Duration // Solo edit/execution
	lastOpID    string        // Tracking interno
	mu          sync.RWMutex  // Protección para lastOpID

	origHandler interface{} // Store original handler for type assertions

	// Function pointers - solo los necesarios poblados
	nameFunc     func() string              // Todos
	labelFunc    func() string              // Display/Edit/Execution
	valueFunc    func() string              // Edit/Display
	contentFunc  func() string              // Display únicamente
	editableFunc func() bool                // Por tipo
	changeFunc   func(string, func(string)) // Edit/Execution (nueva firma)
	executeFunc  func(func(string))         // Execution únicamente (nueva firma)
	timeoutFunc  func() time.Duration       // Edit/Execution
	getOpIDFunc  func() string              // Tracking
	setOpIDFunc  func(string)               // Tracking
}

// ============================================================================
// anyHandler Methods - Replaces fieldHandler interface
// ============================================================================

func (a *anyHandler) Name() string {
	if a.nameFunc != nil {
		return a.nameFunc()
	}
	return ""
}

func (a *anyHandler) Label() string {
	if a.labelFunc != nil {
		return a.labelFunc()
	}
	return ""
}

func (a *anyHandler) Value() string {
	if a.valueFunc != nil {
		return a.valueFunc()
	}
	return ""
}

func (a *anyHandler) Editable() bool {
	if a.editableFunc != nil {
		return a.editableFunc()
	}
	return false
}

func (a *anyHandler) Change(newValue string, progress func(string)) {
	if a.changeFunc != nil {
		a.changeFunc(newValue, progress)
	}
}

func (a *anyHandler) Timeout() time.Duration {
	if a.timeoutFunc != nil {
		return a.timeoutFunc()
	}
	return a.timeout
}

func (a *anyHandler) SetLastOperationID(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.lastOpID = id
	if a.setOpIDFunc != nil {
		a.setOpIDFunc(id)
	}
}

func (a *anyHandler) GetLastOperationID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.getOpIDFunc != nil {
		return a.getOpIDFunc()
	}
	return a.lastOpID
}

// ============================================================================
// Factory Methods
// ============================================================================

func newEditHandler(h HandlerEdit, timeout time.Duration, tracker MessageTracker) *anyHandler {
	anyH := &anyHandler{
		handlerType:  handlerTypeEdit,
		timeout:      timeout,
		nameFunc:     h.Name,
		labelFunc:    h.Label,
		valueFunc:    h.Value,
		editableFunc: func() bool { return true },
		changeFunc:   h.Change,
		timeoutFunc:  func() time.Duration { return timeout },
		origHandler:  h,
	}

	// Configurar tracking opcional
	if tracker != nil {
		anyH.getOpIDFunc = tracker.GetLastOperationID
		anyH.setOpIDFunc = tracker.SetLastOperationID
	} else {
		anyH.getOpIDFunc = func() string { return "" }
		anyH.setOpIDFunc = func(string) {}
	}

	return anyH
}

func newDisplayHandler(h HandlerDisplay) *anyHandler {
	return &anyHandler{
		handlerType:  handlerTypeDisplay,
		timeout:      0,         // Display no requiere timeout
		nameFunc:     h.Name,    // Solo Name()
		valueFunc:    h.Content, // Content como Value para compatibilidad interna
		contentFunc:  h.Content, // Solo Content()
		editableFunc: func() bool { return false },
		getOpIDFunc:  func() string { return "" },
		setOpIDFunc:  func(string) {},
	}
}

func newExecutionHandler(h HandlerExecution, timeout time.Duration) *anyHandler {
	anyH := &anyHandler{
		handlerType:  handlerTypeExecution,
		timeout:      timeout,
		nameFunc:     h.Name,
		labelFunc:    h.Label,
		editableFunc: func() bool { return false },
		executeFunc:  h.Execute,
		changeFunc: func(_ string, progress func(string)) {
			h.Execute(progress)
		},
		timeoutFunc: func() time.Duration { return timeout },
		origHandler: h,
	}

	// Check if handler implements MessageTracker interface for operation tracking
	if tracker, ok := h.(MessageTracker); ok {
		anyH.getOpIDFunc = tracker.GetLastOperationID
		anyH.setOpIDFunc = tracker.SetLastOperationID
	} else {
		anyH.getOpIDFunc = func() string { return "" }
		anyH.setOpIDFunc = func(string) {}
	}

	// Check if handler also implements Value() method (like TestNonEditableHandler)
	if valuer, ok := h.(interface{ Value() string }); ok {
		anyH.valueFunc = valuer.Value
	} else {
		anyH.valueFunc = h.Label // Fallback to Label
	}

	return anyH
}

func newWriterHandler(h HandlerWriter) *anyHandler {
	return &anyHandler{
		handlerType: handlerTypeWriter,
		nameFunc:    h.Name,
		getOpIDFunc: func() string { return "" }, // Siempre nuevas líneas
		setOpIDFunc: func(string) {},
	}
}

func newTrackerWriterHandler(h HandlerWriterTracker) *anyHandler {
	return &anyHandler{
		handlerType: handlerTypeTrackerWriter,
		nameFunc:    h.Name,
		getOpIDFunc: h.GetLastOperationID,
		setOpIDFunc: h.SetLastOperationID,
	}
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
	// NEW: Handler-based approach with anyHandler (replaces fieldHandler)
	handler   *anyHandler // Handles all field behavior
	parentTab *tabSection // Direct reference to parent for message routing

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
// DEPRECATED: Use NewEditHandler, NewExecutionHandler, or NewDisplayHandler instead
// This method will be removed in future versions. Use the type-safe builder methods instead.
//
// Example migration:
//
//	OLD: tab.NewField(handler)
//	NEW: tab.NewEditHandler(handler).Register()
func (ts *tabSection) NewField(handler *anyHandler) *tabSection {
	// Log deprecation warning
	if ts.tui != nil && ts.tui.LogToFile != nil {
		ts.tui.LogToFile("WARNING: NewField is deprecated. Use NewEditHandler/NewExecutionHandler/NewDisplayHandler instead")
	}

	f := &field{
		handler:    handler,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}

	// AUTO-REGISTER: FieldHandlers are automatically registered for writing
	ts.registerAnyHandler(handler)

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

func (f *field) Value() string {
	if f.handler != nil {
		return f.handler.Value()
	}
	return ""
}

// GetHandlerForTest returns the handler for testing purposes
func (f *field) GetHandlerForTest() *anyHandler {
	return f.handler
}

func (f *field) Editable() bool {
	if f.handler != nil {
		return f.handler.Editable()
	}
	return false
}

// READONLY FIELD CONVENTION:
// - FieldHandler with Label() == "" (exactly empty string) indicates readonly/info display
// - Uses fieldReadOnlyStyle (highlight background + clear text)
// - No keyboard interaction allowed (no cursor, no Enter response)
// - Message content displayed without timestamp for cleaner visual
// - Navigation between fields works, but no interaction within readonly content
func (f *field) isDisplayOnly() bool {
	if f.handler == nil {
		return false
	}
	return f.handler.handlerType == handlerTypeDisplay
}

// NUEVO: Detección para execution con footer expandido
func (f *field) isExecutionHandler() bool {
	if f.handler == nil {
		return false
	}
	return f.handler.handlerType == handlerTypeExecution
}

// NUEVO: Detección para handlers que usan footer expandido (Display + Execution)
func (f *field) usesExpandedFooter() bool {
	return f.isDisplayOnly() || f.isExecutionHandler()
}

// NUEVO: Método para mostrar contenido en la sección principal
func (f *field) getDisplayContent() string {
	if f.isDisplayOnly() && f.handler != nil {
		if f.handler.contentFunc != nil {
			return f.handler.contentFunc() // Content() se muestra en la sección principal
		}
	}
	return ""
}

// NUEVO: Método para footer expandido - Name() usa espacio de label + value
func (f *field) getExpandedFooterLabel() string {
	if f.usesExpandedFooter() && f.handler != nil {
		if f.isDisplayOnly() && f.handler.nameFunc != nil {
			// Display handlers show Name() in footer
			return f.handler.nameFunc()
		} else if f.isExecutionHandler() && f.handler.valueFunc != nil {
			// Execution handlers show Value() in footer for better UX
			return f.handler.valueFunc()
		}
	}
	return ""
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

		// SOLUCIÓN: Usar detección centralizada como Writers
		msgType := messagetype.DetectMessageType(content)
		f.parentTab.tui.sendMessageWithHandler(content, msgType, f.parentTab, handlerName, f.asyncState.operationID)
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

		// SOLUCIÓN: Usar detección centralizada como sendProgressMessage
		msgType := messagetype.DetectMessageType(content)
		f.parentTab.tui.sendMessageWithHandler(content, msgType, f.parentTab, handlerName, operationID)
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

		// SOLUCIÓN: Usar detección centralizada como sendProgressMessage
		msgType := messagetype.DetectMessageType(content)
		f.parentTab.tui.sendMessageWithHandler(content, msgType, f.parentTab, handlerName, operationID)
	}
}

// executeAsyncChange executes the handler's Change method asynchronously
func (f *field) executeAsyncChange(valueToSave any) {
	if f.handler == nil || f.asyncState == nil {
		return
	}

	// In test mode, execute synchronously for predictable test behavior
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.isTestMode() {
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
	} else {
		// Log when id is nil for debugging
		if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.LogToFile != nil {
			f.parentTab.tui.LogToFile("Warning: Cannot generate operation ID, unixid not initialized")
		}
	}
	f.asyncState.startTime = time.Now()

	// Create progress callback for handler
	progressCallback := func(message string) {
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
		f.handler.Change(currentValue.(string), progressCallback)
		result := f.handler.Value() // Obtener valor actualizado
		resultChan <- struct {
			result string
			err    error
		}{result, nil}
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
			switch f.handler.handlerType {
			case handlerTypeEdit:
				f.sendSuccessMessage(res.result)
			case handlerTypeExecution:
				// Only send if handler explicitly implements Value()
				if _, ok := f.handler.origHandler.(interface{ Value() string }); ok {
					f.sendSuccessMessage(res.result)
				}
				// Other handler types: do not send success message
			}
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

// executeChangeSyncWithValue executes the handler's Change method synchronously with pre-captured value
func (f *field) executeChangeSyncWithValue(valueToSave any) {
	if f.handler == nil {
		return
	}

	// In sync test mode, we don't generate operation IDs or send messages to avoid race conditions
	// Use the pre-captured value directly

	// Create empty progress callback for sync test execution
	progressCallback := func(message string) {
		// In sync test mode, we don't send messages to avoid race conditions
	}

	f.handler.Change(valueToSave.(string), progressCallback)
	// In test mode, we don't send messages to UI to avoid race conditions
	// The test can verify the handler's internal state directly
}

// executeChangeSyncWithTracking executes the handler's Change method synchronously but maintains operation ID tracking
// This is specifically for testing operation ID reuse functionality
func (f *field) executeChangeSyncWithTracking(valueToSave any) {
	if f.handler == nil {
		return
	}

	// Generate or reuse operation ID like in async mode
	var operationID string
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.id != nil {
		// Check if handler has existing operationID to reuse (for updates)
		if existingID := f.handler.GetLastOperationID(); existingID != "" {
			operationID = existingID
		} else {
			// Generate new ID for new operations
			operationID = f.parentTab.tui.id.GetNewID()
		}
	}

	// Create progress callback that sends messages with operation tracking
	handlerName := f.handler.Name()
	progressCallback := func(message string) {
		if f.parentTab != nil {
			msgType := messagetype.DetectMessageType(message)
			f.parentTab.tui.sendMessageWithHandler(message, msgType, f.parentTab, handlerName, operationID)
		}
	}

	// Execute handler
	f.handler.Change(valueToSave.(string), progressCallback)

	// Set operation ID on handler for tracking
	f.handler.SetLastOperationID(operationID)

	// Send success message
	result := f.handler.Value()
	if f.parentTab != nil {
		msgType := messagetype.DetectMessageType(result)
		f.parentTab.tui.sendMessageWithHandler(result, msgType, f.parentTab, handlerName, operationID)
	}
}

// handleEnter triggers async operation when user presses Enter
func (f *field) handleEnter() {
	if f.handler == nil {
		return
	}

	// NEW: Readonly fields don't respond to any keys
	if f.isDisplayOnly() {
		return
	}

	// Capture the current value BEFORE any state changes
	valueToSave := f.getCurrentValue()

	// In test mode, execute synchronously without goroutine
	if f.parentTab != nil && f.parentTab.tui != nil && f.parentTab.tui.isTestMode() {
		f.executeChangeSyncWithValue(valueToSave)
		return
	}

	// DevTUI handles async internally - user doesn't see this complexity
	go f.executeAsyncChange(valueToSave)
}
