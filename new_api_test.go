package devtui

import (
	"testing"
	"time"
)

// Test handlers implementing the new interfaces

type testDisplayHandler struct{}

func (h *testDisplayHandler) Label() string   { return "Test Display" }
func (h *testDisplayHandler) Content() string { return "This is display content" }

type testEditHandler struct {
	value string
}

func (h *testEditHandler) Label() string { return "Test Edit" }
func (h *testEditHandler) Value() string { return h.value }
func (h *testEditHandler) Change(newValue any, progress ...func(string)) error {
	h.value = newValue.(string)
	return nil
}

type testRunHandler struct{}

func (h *testRunHandler) Label() string { return "Test Run" }
func (h *testRunHandler) Execute(progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Operation completed")
	}
	return nil
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

	// Test DisplayHandler registration
	tab.NewDisplayHandler(&testDisplayHandler{}).Register()

	// Test EditHandler registration with and without timeout
	tab.NewEditHandler(&testEditHandler{value: "initial"}).Register()                 // Sync
	tab.NewEditHandler(&testEditHandler{value: "async"}).WithTimeout(5 * time.Second) // Async

	// Test ExecutionHandler registration with and without timeout
	tab.NewRunHandler(&testRunHandler{}).Register()                    // Sync
	tab.NewRunHandler(&testRunHandler{}).WithTimeout(10 * time.Second) // Async

	// Test Writer registration
	basicWriter := tab.NewWriterHandler(&testWriterBasic{}).Register()
	trackerWriter := tab.NewWriterHandler(&testWriterTracker{}).Register()

	// Verify field count (5 fields registered)
	if len(tab.fieldHandlers) != 5 {
		t.Errorf("Expected 5 fields, got %d", len(tab.fieldHandlers))
	}

	// Test field types
	fields := tab.fieldHandlers

	// First field should be DisplayHandler (read-only)
	if !fields[0].isDisplayOnly() {
		t.Error("First field should be display-only")
	}

	// Second and third fields should be EditHandler (editable)
	if !fields[1].Editable() {
		t.Error("Second field should be editable")
	}
	if !fields[2].Editable() {
		t.Error("Third field should be editable")
	}

	// Fourth and fifth fields should be ExecutionHandler (not editable, but not display-only)
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

func TestFieldWrappers(t *testing.T) {
	// Test displayFieldHandler wrapper (private type, tested through public API)
	display := &testDisplayHandler{}
	wrapper := &displayFieldHandler{display: display, timeout: 0}

	if wrapper.Label() != "Test Display" {
		t.Error("displayFieldHandler.Label() mismatch")
	}
	if wrapper.Value() != "This is display content" {
		t.Error("displayFieldHandler.Value() mismatch")
	}
	if wrapper.Editable() {
		t.Error("displayFieldHandler should not be editable")
	}

	// Test editFieldHandler wrapper (private type, tested through public API)
	edit := &testEditHandler{value: "test"}
	editWrapper := &editFieldHandler{edit: edit, timeout: 5 * time.Second}

	if editWrapper.Label() != "Test Edit" {
		t.Error("editFieldHandler.Label() mismatch")
	}
	if editWrapper.Value() != "test" {
		t.Error("editFieldHandler.Value() mismatch")
	}
	if !editWrapper.Editable() {
		t.Error("editFieldHandler should be editable")
	}
	if editWrapper.Timeout() != 5*time.Second {
		t.Error("editFieldHandler.Timeout() mismatch")
	}

	// Test runFieldHandler wrapper (private type, tested through public API)
	run := &testRunHandler{}
	runWrapper := &runFieldHandler{run: run, timeout: 10 * time.Second}

	if runWrapper.Label() != "Test Run" {
		t.Error("runFieldHandler.Label() mismatch")
	}
	if runWrapper.Value() != "Execute" {
		t.Error("runFieldHandler.Value() should be 'Execute'")
	}
	if runWrapper.Editable() {
		t.Error("runFieldHandler should not be editable")
	}
	if runWrapper.Timeout() != 10*time.Second {
		t.Error("runFieldHandler.Timeout() mismatch")
	}
}

func TestBasicWriterAdapter(t *testing.T) {
	basic := &testWriterBasic{}
	adapter := &basicWriterAdapter{basic: basic}

	if adapter.Name() != "TestWriter" {
		t.Error("basicWriterAdapter.Name() mismatch")
	}

	// Test operation ID behavior (always returns empty for basic writers)
	adapter.SetLastOperationID("test-id")
	if adapter.GetLastOperationID() != "" {
		t.Error("basicWriterAdapter should always return empty operation ID")
	}
}
