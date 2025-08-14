package devtui

import (
	"testing"
	"time"
)

// Test handlers implementing the new interfaces

type testDisplayHandler struct{}

func (h *testDisplayHandler) Name() string    { return "Test Display Handler" }
func (h *testDisplayHandler) Content() string { return "This is display content" }

type testEditHandler struct {
	value string
}

func (h *testEditHandler) Name() string  { return "TestEdit" }
func (h *testEditHandler) Label() string { return "Test Edit" }
func (h *testEditHandler) Value() string { return h.value }
func (h *testEditHandler) Change(newValue string, progress func(msgs ...any)) {
	h.value = newValue
	if progress != nil {
		progress("Changed")
	}
}

type testRunHandler struct{}

func (h *testRunHandler) Name() string  { return "TestRun" }
func (h *testRunHandler) Label() string { return "Test Run" }
func (h *testRunHandler) Execute(progress func(msgs ...any)) {
	if progress != nil {
		progress("Operation completed")
	}
}

type testLoggerBasic struct{}

func (w *testLoggerBasic) Name() string { return "TestWriter" }

type testLoggerTracker struct {
	lastOpID string
}

func (w *testLoggerTracker) Name() string                 { return "TestTrackerLogger" }
func (w *testLoggerTracker) GetLastOperationID() string   { return w.lastOpID }
func (w *testLoggerTracker) SetLastOperationID(id string) { w.lastOpID = id }

func TestNewAPIHandlers(t *testing.T) {
	// Create TUI
	exitChan := make(chan bool, 1)
	tui := NewTUI(&TuiConfig{
		AppName:  "Test New API",
		ExitChan: exitChan,
	})

	// Create tab section
	tab := tui.NewTabSection("Test", "Testing new API")

	// Test HandlerDisplay registration
	tab.AddDisplayHandler(&testDisplayHandler{})

	// Test HandlerEdit registration with and without timeout
	tab.AddEditHandler(&testEditHandler{value: "initial"}, 0)           // Sync
	tab.AddEditHandler(&testEditHandler{value: "async"}, 5*time.Second) // Async

	// Test HandlerExecution registration with and without timeout
	tab.AddExecutionHandler(&testRunHandler{}, 0)              // Sync
	tab.AddExecutionHandler(&testRunHandler{}, 10*time.Second) // Async

	// Test Writer registration
	basicLogger := tab.NewLogger("testLoggerBasic", false)
	trackerLogger := tab.NewLogger("testLoggerTracker", true)

	// Verify field count (5 fields registered)
	if len(tab.fieldHandlers) != 5 {
		t.Errorf("Expected 5 fields, got %d", len(tab.fieldHandlers))
	}

	// Test field types
	fields := tab.fieldHandlers

	// First field should be HandlerDisplay (read-only)
	if !fields[0].isDisplayOnly() {
		t.Error("First field should be display-only")
	}

	// Second and third fields should be HandlerEdit (editable)
	if !fields[1].editable() {
		t.Error("Second field should be editable")
	}
	if !fields[2].editable() {
		t.Error("Third field should be editable")
	}

	// Fourth and fifth fields should be HandlerExecution (not editable, but not display-only)
	if fields[3].editable() {
		t.Error("Fourth field should not be editable")
	}
	if fields[3].isDisplayOnly() {
		t.Error("Fourth field should not be display-only")
	}

	// Test writers
	if basicLogger == nil {
		t.Error("Basic writer should not be nil")
	}
	if trackerLogger == nil {
		t.Error("Tracker writer should not be nil")
	}

	// Test writing to basic writer
	n, err := basicLogger.Write([]byte("test message"))
	if err != nil {
		t.Errorf("Basic writer failed: %v", err)
	}
	if n != 12 { // "test message" length
		t.Errorf("Expected 12 bytes written, got %d", n)
	}

	// Test writing to tracker writer
	n, err = trackerLogger.Write([]byte("tracked message"))
	if err != nil {
		t.Errorf("Tracker writer failed: %v", err)
	}
	if n != 15 { // "tracked message" length
		t.Errorf("Expected 15 bytes written, got %d", n)
	}

	// Verify writing handlers were registered
	if len(tab.writingHandlers) != 2 {
		t.Errorf("Expected 2 writing handlers, got %d", len(tab.writingHandlers))
	}

	close(exitChan)
}
