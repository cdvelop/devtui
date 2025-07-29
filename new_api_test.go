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

type testWriterBasic struct{}

func (w *testWriterBasic) Name() string { return "TestWriter" }

type testWriterTracker struct {
	lastOpID string
}

func (w *testWriterTracker) Name() string                 { return "TestTrackerWriter" }
func (w *testWriterTracker) GetLastOperationID() string   { return w.lastOpID }
func (w *testWriterTracker) SetLastOperationID(id string) { w.lastOpID = id }

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
	tab.AddHandlerDisplay(&testDisplayHandler{})

	// Test HandlerEdit registration with and without timeout
	tab.AddEditHandler(&testEditHandler{value: "initial"}).Register()                 // Sync
	tab.AddEditHandler(&testEditHandler{value: "async"}).WithTimeout(5 * time.Second) // Async

	// Test HandlerExecution registration with and without timeout
	tab.AddExecutionHandler(&testRunHandler{}).Register()                    // Sync
	tab.AddExecutionHandler(&testRunHandler{}).WithTimeout(10 * time.Second) // Async

	// Test Writer registration
	basicWriter := tab.RegisterHandlerWriter(&testWriterBasic{})
	trackerWriter := tab.RegisterHandlerWriter(&testWriterTracker{})

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
	if !fields[1].Editable() {
		t.Error("Second field should be editable")
	}
	if !fields[2].Editable() {
		t.Error("Third field should be editable")
	}

	// Fourth and fifth fields should be HandlerExecution (not editable, but not display-only)
	if fields[3].Editable() {
		t.Error("Fourth field should not be editable")
	}
	if fields[3].isDisplayOnly() {
		t.Error("Fourth field should not be display-only")
	}

	// Test writers
	if basicWriter == nil {
		t.Error("Basic writer should not be nil")
	}
	if trackerWriter == nil {
		t.Error("Tracker writer should not be nil")
	}

	// Test writing to basic writer
	n, err := basicWriter.Write([]byte("test message"))
	if err != nil {
		t.Errorf("Basic writer failed: %v", err)
	}
	if n != 12 { // "test message" length
		t.Errorf("Expected 12 bytes written, got %d", n)
	}

	// Test writing to tracker writer
	n, err = trackerWriter.Write([]byte("tracked message"))
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
