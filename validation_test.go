package devtui

import (
	"testing"
	"time"
)

func TestValidateTabSection_Nil(t *testing.T) {
	tui := NewTUI(&TuiConfig{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil tabSection")
		}
	}()

	tui.AddHandler(&validationTestHandler{}, time.Second, "", nil)
}

func TestValidateTabSection_WrongType(t *testing.T) {
	tui := NewTUI(&TuiConfig{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for wrong type")
		}
	}()

	tui.AddHandler(&validationTestHandler{}, time.Second, "", "not a tabSection")
}

func TestValidateTabSection_WrongDevTUI(t *testing.T) {
	tui1 := NewTUI(&TuiConfig{})
	tui2 := NewTUI(&TuiConfig{})

	tab := tui1.NewTabSection("TEST", "test")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for tabSection from different DevTUI")
		}
	}()

	tui2.AddHandler(&validationTestHandler{}, time.Second, "", tab)
}

func TestValidateTabSection_Success(t *testing.T) {
	tui := NewTUI(&TuiConfig{})
	tab := tui.NewTabSection("TEST", "test")

	// Should not panic
	tui.AddHandler(&validationTestDisplayHandler{name: "test"}, 0, "", tab)
	log := tui.AddLogger("test", true, "", tab)

	if log == nil {
		t.Error("Expected logger function, got nil")
	}
}

// validationTestHandler is a minimal handler for testing purposes
type validationTestHandler struct{}

func (h *validationTestHandler) Name() string {
	return "test"
}

func (h *validationTestHandler) Value() string {
	return "test value"
}

// validationTestDisplayHandler is a minimal display handler for testing purposes
type validationTestDisplayHandler struct {
	name string
}

func (h *validationTestDisplayHandler) Name() string {
	return h.name
}

func (h *validationTestDisplayHandler) Value() string {
	return "display value"
}