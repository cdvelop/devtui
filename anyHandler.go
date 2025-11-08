package devtui

import (
	"sync"
	"time"
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
	handlerTypeInteractive // NEW: Interactive content handler
)

// anyHandler - Estructura privada que unifica todos los handlers
type anyHandler struct {
	handlerType handlerType
	timeout     time.Duration // Solo edit/execution
	lastOpID    string        // Tracking interno
	mu          sync.RWMutex  // Protección para lastOpID

	origHandler any // Store original handler for type assertions

	handlerColor string // NEW: Handler-specific color for message formatting

	// Function pointers - solo los necesarios poblados
	nameFunc     func() string                      // Todos
	labelFunc    func() string                      // Display/Edit/Execution
	valueFunc    func() string                      // Edit/Display
	contentFunc  func() string                      // Display únicamente
	editableFunc func() bool                        // Por tipo
	editModeFunc func() bool                        // NEW: Auto edit mode activation
	changeFunc   func(string, chan<- string) // Edit/Execution (nueva firma)
	executeFunc  func(chan<- string)            // Execution únicamente (nueva firma)
	timeoutFunc  func() time.Duration               // Edit/Execution
	getOpIDFunc  func() string                      // Tracking
	setOpIDFunc  func(string)                       // Tracking
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

func (a *anyHandler) editable() bool {
	if a.editableFunc != nil {
		return a.editableFunc()
	}
	return false
}

func (a *anyHandler) Change(newValue string, progress chan<- string) {
	if a.changeFunc != nil {
		a.changeFunc(newValue, progress)
	}
}

func (a *anyHandler) Execute(progress chan<- string) {
	if a.executeFunc != nil {
		a.executeFunc(progress)
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

func (a *anyHandler) WaitingForUser() bool {
	if a.editModeFunc != nil {
		return a.editModeFunc()
	}
	return false
}

// ============================================================================
// Factory Methods
// ============================================================================

func NewEditHandler(h HandlerEdit, timeout time.Duration, tracker MessageTracker, color string) *anyHandler {
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
		handlerColor: color, // NEW: Store handler color
	}

	// NEW: Check if handler also implements Value() method (like TestNonEditableHandler)
	if valuer, ok := h.(interface{ Value() string }); ok {
		anyH.valueFunc = valuer.Value
	} else {
		anyH.valueFunc = h.Label // Fallback to Label
	}

	// REMOVED: Hybrid Content() detection - use HandlerInteractive instead

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

func NewDisplayHandler(h HandlerDisplay, color string) *anyHandler {
	return &anyHandler{
		handlerType:  handlerTypeDisplay,
		timeout:      0,         // Display no requiere timeout
		nameFunc:     h.Name,    // Solo Name()
		valueFunc:    h.Content, // Content como Value para compatibilidad interna
		contentFunc:  h.Content, // Solo Content()
		editableFunc: func() bool { return false },
		getOpIDFunc:  func() string { return "" },
		setOpIDFunc:  func(string) {},
		handlerColor: color, // NEW: Store handler color
	}
}

func NewExecutionHandler(h HandlerExecution, timeout time.Duration, color string) *anyHandler {
	anyH := &anyHandler{
		handlerType:  handlerTypeExecution,
		timeout:      timeout,
		nameFunc:     h.Name,
		labelFunc:    h.Label,
		editableFunc: func() bool { return false },
		executeFunc:  h.Execute,
		changeFunc: func(_ string, progress chan<- string) {
			h.Execute(progress)
		},
		timeoutFunc:  func() time.Duration { return timeout },
		origHandler:  h,
		handlerColor: color, // NEW: Store handler color
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

	// REMOVED: Hybrid Content() detection - use HandlerInteractive instead

	return anyH
}

func NewWriterHandler(h HandlerLogger, color string) *anyHandler {
	return &anyHandler{
		handlerType:  handlerTypeWriter,
		nameFunc:     h.Name,
		getOpIDFunc:  func() string { return "" }, // Siempre nuevas líneas
		setOpIDFunc:  func(string) {},
		handlerColor: color, // NEW: Store handler color
	}
}

func NewWriterTrackerHandler(h interface {
	Name() string
	GetLastOperationID() string
	SetLastOperationID(string)
}, color string) *anyHandler {
	return &anyHandler{
		handlerType:  handlerTypeTrackerWriter,
		nameFunc:     h.Name,
		getOpIDFunc:  h.GetLastOperationID,
		setOpIDFunc:  h.SetLastOperationID,
		handlerColor: color, // NEW: Store handler color
	}
}

func NewInteractiveHandler(h HandlerInteractive, timeout time.Duration, tracker MessageTracker, color string) *anyHandler {
	anyH := &anyHandler{
		handlerType: handlerTypeInteractive,
		timeout:     timeout,
		nameFunc:    h.Name,
		labelFunc:   h.Label,
		valueFunc:   h.Value,
		// NO contentFunc - interactive handlers use progress() only
		editableFunc: func() bool { return true },
		changeFunc:   h.Change,
		timeoutFunc:  func() time.Duration { return timeout },
		editModeFunc: h.WaitingForUser, // NEW: Auto edit mode detection
		origHandler:  h,
		handlerColor: color, // NEW: Store handler color
	}

	// Configure optional tracking
	if tracker != nil {
		anyH.getOpIDFunc = tracker.GetLastOperationID
		anyH.setOpIDFunc = tracker.SetLastOperationID
	} else {
		anyH.getOpIDFunc = func() string { return "" }
		anyH.setOpIDFunc = func(string) {}
	}

	return anyH
}
