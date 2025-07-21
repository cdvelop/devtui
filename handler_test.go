package devtui

import (
	"errors"
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

// TestEditableHandler - Manejador editable básico para todos los tests
type TestEditableHandler struct {
	currentValue string
}

func (h *TestEditableHandler) Label() string          { return "Editable Field" }
func (h *TestEditableHandler) Value() string          { return h.currentValue }
func (h *TestEditableHandler) Editable() bool         { return true }
func (h *TestEditableHandler) Timeout() time.Duration { return 0 }

func (h *TestEditableHandler) Change(newValue any) (string, error) {
	strValue := newValue.(string)
	h.currentValue = strValue
	return "Saved: " + strValue, nil
}

// TestNonEditableHandler - Manejador no editable básico para todos los tests
type TestNonEditableHandler struct{}

func (h *TestNonEditableHandler) Label() string          { return "Non-Editable Field" }
func (h *TestNonEditableHandler) Value() string          { return "action button" }
func (h *TestNonEditableHandler) Editable() bool         { return false }
func (h *TestNonEditableHandler) Timeout() time.Duration { return 0 }

func (h *TestNonEditableHandler) Change(newValue any) (string, error) {
	return "Action executed", nil
}

// PortTestHandler - Manejador específico para tests de puerto (centralizado aquí)
type PortTestHandler struct {
	currentPort string
	mu          sync.RWMutex // Para proteger currentPort de race conditions
}

func (h *PortTestHandler) Label() string { return "Server Port" }
func (h *PortTestHandler) Value() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.currentPort
}
func (h *PortTestHandler) Editable() bool         { return true }
func (h *PortTestHandler) Timeout() time.Duration { return 0 }

func (h *PortTestHandler) Change(newValue any) (string, error) {
	newPort := newValue.(string)

	// Simple validation - reject obviously invalid ports for testing
	if newPort == "99999" {
		return "", errors.New("port out of range")
	}

	// Accept valid ports - protegido por mutex
	h.mu.Lock()
	h.currentPort = newPort
	h.mu.Unlock()

	return "Port updated to " + newPort, nil
}

// NewTestFieldHandler creates a new test handler with basic functionality
// Esta función mantiene compatibilidad con tests existentes
func NewTestFieldHandler(label, value string, editable bool, changeFunc func(newValue any) (string, error)) FieldHandler {
	// Siempre usar CustomTestHandler para permitir configuración completa
	return &CustomTestHandler{
		label:      label,
		value:      value,
		editable:   editable,
		changeFunc: changeFunc,
	}
}

// CustomTestHandler - Handler personalizable para casos específicos
type CustomTestHandler struct {
	label      string
	value      string
	editable   bool
	changeFunc func(newValue any) (string, error)
}

func (h *CustomTestHandler) Label() string          { return h.label }
func (h *CustomTestHandler) Value() string          { return h.value }
func (h *CustomTestHandler) Editable() bool         { return h.editable }
func (h *CustomTestHandler) Timeout() time.Duration { return 0 }

func (h *CustomTestHandler) Change(newValue any) (string, error) {
	if h.changeFunc != nil {
		result, err := h.changeFunc(newValue)
		if err == nil {
			inputStr := newValue.(string)

			// Special handling for empty values and transformations
			if inputStr == "" {
				// If input is empty and result doesn't look like a status message,
				// treat result as the new field value (for default value transformations)
				if result != "" && !strings.Contains(result, "Saved") && !strings.Contains(result, "Error") {
					h.value = result
				} else {
					// Empty input with status message - field becomes empty
					h.value = ""
				}
			} else {
				// Non-empty input: use input as the new value
				h.value = inputStr
			}
		}
		return result, err
	}
	// Default behavior
	h.value = newValue.(string)
	return h.value, nil
}

// SetLabel allows updating the label for testing
func (h *CustomTestHandler) SetLabel(label string) {
	h.label = label
}

// SetValue allows updating the value for testing (simulates external changes)
func (h *CustomTestHandler) SetValue(value string) {
	h.value = value
}

// SetEditable allows changing the editable state for testing
func (h *CustomTestHandler) SetEditable(editable bool) {
	h.editable = editable
}

// Aliases para compatibilidad con tests existentes
type TestFieldHandler = CustomTestHandler
type TestField1Handler = TestEditableHandler

// DefaultTUIForTest creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
// LogToFile is optional - if not provided, will use a no-op logger
func DefaultTUIForTest(LogToFile ...func(messages ...any)) *DevTUI {
	// Default no-op logger if none provided
	var logFunc func(messages ...any)
	if len(LogToFile) > 0 && LogToFile[0] != nil {
		logFunc = LogToFile[0]
	} else {
		logFunc = func(messages ...any) {
			// No-op logger for tests
		}
	}

	// Create basic tabSections for testing
	tmpTUI := &DevTUI{TuiConfig: &TuiConfig{}}

	// Tab 1: Con manejadores
	tab1 := tmpTUI.NewTabSection("Tab 1", "Tab with handlers")
	editableHandler := &TestEditableHandler{currentValue: "initial test value"}
	nonEditableHandler := &TestNonEditableHandler{}

	tab1.NewField(editableHandler).
		NewField(nonEditableHandler)
	tab1.SetIndex(GetFirstTestTabIndex()) // Index 1 (SHORTCUTS is 0)
	tab1.SetActiveEditField(0)

	// Tab 2: Sin manejadores (para tests simples)
	tab2 := tmpTUI.NewTabSection("Tab 2", "Empty tab")
	tab2.SetIndex(GetSecondTestTabIndex()) // Index 2

	tabSections := []*tabSection{tab1, tab2}

	// Initialize the UI with TestMode enabled for synchronous execution
	h := NewTUI(&TuiConfig{
		TabIndexStart: 0,               // Start with the first tab
		ExitChan:      make(chan bool), // Channel to signal exit
		TestMode:      true,            // Enable test mode for synchronous execution
		Color:         nil,             // Use default colors
		LogToFile:     logFunc,
	}).AddTabSections(tabSections...)

	return h
}
