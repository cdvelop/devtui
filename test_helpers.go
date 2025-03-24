package devtui

import (
	"github.com/cdvelop/messagetype"
)

// Mock fieldHandlerAdapter for testing
type mockFieldHandler struct {
	name        string
	value       string
	editable    bool
	changeValue func(newValue string) <-chan MessageUpdate
}

func (m *mockFieldHandler) Name() string {
	return m.name
}

func (m *mockFieldHandler) Value() string {
	return m.value
}

func (m *mockFieldHandler) Editable() bool {
	return m.editable
}

func (m *mockFieldHandler) ChangeValue(newValue string) <-chan MessageUpdate {
	if m.changeValue != nil {
		return m.changeValue(newValue)
	}
	// Default behavior: return a channel that immediately sends a success message
	updateChan := make(chan MessageUpdate, 1)
	updateChan <- MessageUpdate{Content: "Value changed", Type: messagetype.Success}
	close(updateChan)
	return updateChan
}
