package devtui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetFirstTestTabIndex returns the index of the first test tab
// This centralizes the index calculation to avoid test failures when tabs are added/removed
// Currently, NewTUI always adds SHORTCUTS tab at index 0, so test tabs start at index 1
func GetFirstTestTabIndex() int {
	return 1 // SHORTCUTS tab is always at index 0, so first test tab is at index 1
}

// GetSecondTestTabIndex returns the index of the second test tab
func GetSecondTestTabIndex() int {
	return GetFirstTestTabIndex() + 1 // Second test tab follows first test tab
}

// TestEditableHandler - Handler para campos editables (input fields)
type TestEditableHandler struct {
	mu           sync.RWMutex
	label        string
	currentValue string
	lastOpID     string
	updateMode   bool // Para controlar si actualiza mensajes existentes
}

func NewTestEditableHandler(label, value string) *TestEditableHandler {
	return &TestEditableHandler{
		label:        label,
		currentValue: value,
	}
}

func (h *TestEditableHandler) Label() string { return h.label }

func (h *TestEditableHandler) Value() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.currentValue
}

func (h *TestEditableHandler) Editable() bool         { return true }
func (h *TestEditableHandler) Timeout() time.Duration { return 0 }

func (h *TestEditableHandler) Change(newValue any, progress ...func(string)) (string, error) {
	strValue := newValue.(string)
	h.mu.Lock()
	h.currentValue = strValue
	h.mu.Unlock()
	return "Saved: " + strValue, nil
}

// WritingHandler methods
func (h *TestEditableHandler) Name() string { return h.label + "Handler" }

func (h *TestEditableHandler) SetLastOperationID(lastOpID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastOpID = lastOpID
}

func (h *TestEditableHandler) GetLastOperationID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestEditableHandler) SetUpdateMode(update bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.updateMode = update
}

// TestNonEditableHandler - Handler para botones de acción (action buttons)
type TestNonEditableHandler struct {
	mu         sync.RWMutex
	label      string
	actionText string
	lastOpID   string
	updateMode bool
}

func NewTestNonEditableHandler(label, actionText string) *TestNonEditableHandler {
	return &TestNonEditableHandler{
		label:      label,
		actionText: actionText,
	}
}

func (h *TestNonEditableHandler) Label() string { return h.label }

func (h *TestNonEditableHandler) Value() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.actionText
}

func (h *TestNonEditableHandler) Editable() bool         { return false }
func (h *TestNonEditableHandler) Timeout() time.Duration { return 0 }

func (h *TestNonEditableHandler) Change(newValue any, progress ...func(string)) (string, error) {
	h.mu.RLock()
	actionText := h.actionText
	h.mu.RUnlock()
	return "Action executed: " + actionText, nil
}

// WritingHandler methods
func (h *TestNonEditableHandler) Name() string { return h.label + "Handler" }

func (h *TestNonEditableHandler) SetLastOperationID(lastOpID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastOpID = lastOpID
}
func (h *TestNonEditableHandler) GetLastOperationID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestNonEditableHandler) SetUpdateMode(update bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.updateMode = update
}

// TestWriterHandler - Handler solo para escribir (no es field, para componentes externos)
type TestWriterHandler struct {
	mu         sync.RWMutex
	name       string
	lastOpID   string
	updateMode bool
}

func NewTestWriterHandler(name string) *TestWriterHandler {
	return &TestWriterHandler{name: name}
}

// Solo implementa WritingHandler (no FieldHandler)
func (h *TestWriterHandler) Name() string { return h.name }

func (h *TestWriterHandler) SetLastOperationID(lastOpID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastOpID = lastOpID
}

func (h *TestWriterHandler) GetLastOperationID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestWriterHandler) SetUpdateMode(update bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.updateMode = update
}

// NewTestFieldHandler creates a basic test handler - compatibility function
func NewTestFieldHandler(label, value string, editable bool, changeFunc func(newValue any) (string, error)) FieldHandler {
	if editable {
		handler := NewTestEditableHandler(label, value)
		return handler
	} else {
		handler := NewTestNonEditableHandler(label, value)
		return handler
	}
}

// PortTestHandler - Handler específico para tests de puerto con validación
type PortTestHandler struct {
	mu          sync.RWMutex
	currentPort string
	lastOpID    string
	updateMode  bool
}

func NewPortTestHandler(initialPort string) *PortTestHandler {
	return &PortTestHandler{currentPort: initialPort}
}

func (h *PortTestHandler) Label() string { return "Port" }

func (h *PortTestHandler) Value() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.currentPort
}

func (h *PortTestHandler) Editable() bool         { return true }
func (h *PortTestHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *PortTestHandler) Change(newValue any, progress ...func(string)) (string, error) {
	portStr := strings.TrimSpace(newValue.(string))
	if portStr == "" {
		return "", fmt.Errorf("port cannot be empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("port must be between 1 and 65535")
	}

	h.mu.Lock()
	h.currentPort = portStr
	h.mu.Unlock()

	return fmt.Sprintf("Port configured: %d", port), nil
}

// WritingHandler methods
func (h *PortTestHandler) Name() string { return "PortHandler" }

func (h *PortTestHandler) SetLastOperationID(lastOpID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastOpID = lastOpID
}

func (h *PortTestHandler) GetLastOperationID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *PortTestHandler) SetUpdateMode(update bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.updateMode = update
}

// TestErrorHandler - Handler que siempre genera errores para testing
type TestErrorHandler struct {
	label      string
	value      string
	lastOpID   string
	updateMode bool
}

func NewTestErrorHandler(label, value string) *TestErrorHandler {
	return &TestErrorHandler{
		label: label,
		value: value,
	}
}

func (h *TestErrorHandler) Label() string          { return h.label }
func (h *TestErrorHandler) Value() string          { return h.value }
func (h *TestErrorHandler) Editable() bool         { return true }
func (h *TestErrorHandler) Timeout() time.Duration { return 0 }

func (h *TestErrorHandler) Change(newValue any, progress ...func(string)) (string, error) {
	return "", fmt.Errorf("simulated error occurred")
}

// WritingHandler methods
func (h *TestErrorHandler) Name() string                       { return h.label + "ErrorHandler" }
func (h *TestErrorHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *TestErrorHandler) GetLastOperationID() string {
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestErrorHandler) SetUpdateMode(update bool) {
	h.updateMode = update
}

// TestRequiredFieldHandler - Handler que rechaza valores vacíos
type TestRequiredFieldHandler struct {
	label        string
	currentValue string
	lastOpID     string
	updateMode   bool
}

func NewTestRequiredFieldHandler(label, initialValue string) *TestRequiredFieldHandler {
	return &TestRequiredFieldHandler{
		label:        label,
		currentValue: initialValue,
	}
}

func (h *TestRequiredFieldHandler) Label() string          { return h.label }
func (h *TestRequiredFieldHandler) Value() string          { return h.currentValue }
func (h *TestRequiredFieldHandler) Editable() bool         { return true }
func (h *TestRequiredFieldHandler) Timeout() time.Duration { return 0 }

func (h *TestRequiredFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	strValue := newValue.(string)
	if strValue == "" {
		return "", fmt.Errorf("Field cannot be empty")
	}
	h.currentValue = strValue
	return "Accepted: " + strValue, nil
}

// WritingHandler methods
func (h *TestRequiredFieldHandler) Name() string                       { return h.label + "RequiredHandler" }
func (h *TestRequiredFieldHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *TestRequiredFieldHandler) GetLastOperationID() string {
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestRequiredFieldHandler) SetUpdateMode(update bool) {
	h.updateMode = update
}

// TestOptionalFieldHandler - Handler que acepta valores vacíos
type TestOptionalFieldHandler struct {
	label        string
	currentValue string
	lastOpID     string
	updateMode   bool
}

func NewTestOptionalFieldHandler(label, initialValue string) *TestOptionalFieldHandler {
	return &TestOptionalFieldHandler{
		label:        label,
		currentValue: initialValue,
	}
}

func (h *TestOptionalFieldHandler) Label() string          { return h.label }
func (h *TestOptionalFieldHandler) Value() string          { return h.currentValue }
func (h *TestOptionalFieldHandler) Editable() bool         { return true }
func (h *TestOptionalFieldHandler) Timeout() time.Duration { return 0 }

func (h *TestOptionalFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	if strValue == "" {
		h.currentValue = "Default Value" // Para el test que espera esta transformación
		return "Default Value", nil
	}
	return "Updated: " + strValue, nil
}

// WritingHandler methods
func (h *TestOptionalFieldHandler) Name() string                       { return h.label + "OptionalHandler" }
func (h *TestOptionalFieldHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *TestOptionalFieldHandler) GetLastOperationID() string {
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestOptionalFieldHandler) SetUpdateMode(update bool) {
	h.updateMode = update
}

// TestClearableFieldHandler - Handler que preserva valores vacíos tal como son
type TestClearableFieldHandler struct {
	label        string
	currentValue string
	lastOpID     string
	updateMode   bool
}

func NewTestClearableFieldHandler(label, initialValue string) *TestClearableFieldHandler {
	return &TestClearableFieldHandler{
		label:        label,
		currentValue: initialValue,
	}
}

func (h *TestClearableFieldHandler) Label() string          { return h.label }
func (h *TestClearableFieldHandler) Value() string          { return h.currentValue }
func (h *TestClearableFieldHandler) Editable() bool         { return true }
func (h *TestClearableFieldHandler) Timeout() time.Duration { return 0 }

func (h *TestClearableFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return strValue, nil // Return exactly what was input, including empty string
}

// WritingHandler methods
func (h *TestClearableFieldHandler) Name() string                       { return h.label + "ClearableHandler" }
func (h *TestClearableFieldHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *TestClearableFieldHandler) GetLastOperationID() string {
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestClearableFieldHandler) SetUpdateMode(update bool) {
	h.updateMode = update
}

// TestCapturingHandler - Handler que captura valores recibidos para testing
type TestCapturingHandler struct {
	label         string
	currentValue  string
	capturedValue *string // Puntero para capturar valores en tests
	lastOpID      string
	updateMode    bool
}

func NewTestCapturingHandler(label, initialValue string, capturedValue *string) *TestCapturingHandler {
	return &TestCapturingHandler{
		label:         label,
		currentValue:  initialValue,
		capturedValue: capturedValue,
	}
}

func (h *TestCapturingHandler) Label() string          { return h.label }
func (h *TestCapturingHandler) Value() string          { return h.currentValue }
func (h *TestCapturingHandler) Editable() bool         { return true }
func (h *TestCapturingHandler) Timeout() time.Duration { return 0 }

func (h *TestCapturingHandler) Change(newValue any, progress ...func(string)) (string, error) {
	strValue := newValue.(string)
	if h.capturedValue != nil {
		*h.capturedValue = strValue // Captura el valor para el test
	}
	if strValue == "" {
		h.currentValue = "Field was cleared" // Actualizar el valor interno también
		return "Field was cleared", nil
	}
	h.currentValue = strValue
	return "Field value: " + strValue, nil
}

// WritingHandler methods
func (h *TestCapturingHandler) Name() string                       { return h.label + "CapturingHandler" }
func (h *TestCapturingHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *TestCapturingHandler) GetLastOperationID() string {
	if h.updateMode {
		return h.lastOpID
	}
	return ""
}

// SetUpdateMode permite controlar si actualiza mensajes para tests
func (h *TestCapturingHandler) SetUpdateMode(update bool) {
	h.updateMode = update
}

// DefaultTUIForTest creates a DevTUI instance with configurable handlers
// Usage examples:
//   - DefaultTUIForTest() // Empty TUI, no handlers
//   - DefaultTUIForTest(handler1, handler2) // TUI with specified handlers
//   - DefaultTUIForTest(handler1, func(messages...any){}) // TUI with handlers + logger
func DefaultTUIForTest(handlersAndLogger ...any) *DevTUI {
	var logFunc func(messages ...any)
	var handlers []FieldHandler

	// Parse variadic arguments: handlers (FieldHandler) and optional logger (func)
	for _, arg := range handlersAndLogger {
		switch v := arg.(type) {
		case func(messages ...any):
			logFunc = v
		case FieldHandler:
			handlers = append(handlers, v)
		}
	}

	// Default no-op logger if none provided
	if logFunc == nil {
		logFunc = func(messages ...any) {
			// No-op logger for tests
		}
	}

	// Initialize the UI with TestMode enabled for synchronous execution
	h := NewTUI(&TuiConfig{
		TabIndexStart: 0,               // Start with the first tab
		ExitChan:      make(chan bool), // Channel to signal exit
		TestMode:      true,            // Enable test mode for synchronous execution
		Color:         nil,             // Use default colors
		LogToFile:     logFunc,
	})

	// Create test tab only if handlers are provided
	if len(handlers) > 0 {
		tab := h.NewTabSection("Test Tab", "Tab with test handlers")

		// Add all provided handlers to the tab
		for _, handler := range handlers {
			tab.NewField(handler)
		}

		tab.SetIndex(GetFirstTestTabIndex()) // Index 1 (SHORTCUTS is 0)
		tab.SetActiveEditField(0)
	}

	return h
}
