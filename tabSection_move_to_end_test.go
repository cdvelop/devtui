package devtui

import (
	"testing"
	"time"

	. "github.com/cdvelop/tinystring"
)

type testTracker struct {
	lastOpID string
}

func (t *testTracker) Name() string                                       { return "TestTracker" }
func (t *testTracker) Label() string                                      { return "TrackerLabel" }
func (t *testTracker) Value() string                                      { return "" }
func (t *testTracker) Change(newValue string, progress func(msgs ...any)) {}
func (t *testTracker) GetLastOperationID() string                         { return t.lastOpID }
func (t *testTracker) SetLastOperationID(id string)                       { t.lastOpID = id }

func TestMessageTrackerMoveToEnd(t *testing.T) {
	config := &TuiConfig{
		ExitChan:  make(chan bool),
		Color:     &ColorStyle{},
		LogToFile: func(messages ...any) {},
	}
	tui := NewTUI(config)
	tui.SetTestMode(true)
	tab := tui.NewTabSection("TRACKER", "")

	tracker := &testTracker{}
	tab.AddEditHandler(tracker, 5*time.Second)

	// Add a normal message
	tab.addNewContent(Msg.Info, "Normal message")
	// Add a tracker message (first time, operationID empty)
	updated, _ := tab.updateOrAddContentWithHandler(Msg.Info, "Tracker message 1", tracker.Name(), "op-1")
	tracker.SetLastOperationID("op-1") // simulate tracker storing op id after first message
	if updated {
		t.Fatal("First tracker message should not be an update")
	}
	if tab.tabContents[len(tab.tabContents)-1].Content != "Tracker message 1" {
		t.Fatal("Tracker message 1 should be at the end")
	}

	// Add another normal message
	tab.addNewContent(Msg.Info, "Another normal message")
	if tab.tabContents[len(tab.tabContents)-1].Content != "Another normal message" {
		t.Fatal("Another normal message should be at the end")
	}

	// Update tracker message (should move to end)
	updated, _ = tab.updateOrAddContentWithHandler(Msg.Info, "Tracker message UPDATED", tracker.Name(), tracker.GetLastOperationID())
	if !updated {
		t.Fatal("Tracker message update should return updated=true")
	}
	if tab.tabContents[len(tab.tabContents)-1].Content != "Tracker message UPDATED" {
		t.Fatalf("Tracker message should be moved to end after update, got '%s' at end", tab.tabContents[len(tab.tabContents)-1].Content)
	}

	// Ensure only one tracker message exists
	trackerCount := 0
	for _, c := range tab.tabContents {
		if c.handlerName == tracker.Name() {
			trackerCount++
		}
	}
	if trackerCount != 1 {
		t.Fatalf("Expected only one tracker message, found %d", trackerCount)
	}
}
