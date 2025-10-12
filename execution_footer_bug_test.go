package devtui

import (
	"sync"
	"testing"
	"time"
)

// HandlerExecution without Value()
type ExecHandler struct{}

func (h *ExecHandler) Name() string  { return "TestExec" }
func (h *ExecHandler) Label() string { return "Action" }
func (h *ExecHandler) Execute(progress func(msgs ...any)) {
	progress("Step 1")
	time.Sleep(10 * time.Millisecond)
	progress("Step 2")
	time.Sleep(10 * time.Millisecond)
	progress("Final step")
}

// TestExecutionHandlerFooterBug verifies that only progress messages are shown for execution handlers without Value()
func TestExecutionHandlerFooterBug(t *testing.T) {
	var mu sync.Mutex
	var messages []string

	tui := NewTUI(&TuiConfig{
		AppName:  "TestApp",
		ExitChan: make(chan bool),
		Logger: func(msgs ...any) {
			mu.Lock()
			for _, m := range msgs {
				messages = append(messages, m.(string))
			}
			mu.Unlock()
		},
	})

	tab := tui.NewTabSection("Tab", "TestTab")
	tab.AddHandler(&ExecHandler{}, 50*time.Millisecond, "")

	// Simulate Enter key (async)
	fields := tab.fieldHandlers
	if len(fields) == 0 {
		t.Fatal("No field handlers registered")
	}
	field := fields[0]
	field.handleEnter()

	// Wait for async to finish
	time.Sleep(100 * time.Millisecond)

	// Check tabContents for correct messages
	tab.mu.RLock()
	defer tab.mu.RUnlock()
	var foundFinal bool
	for _, c := range tab.tabContents {
		if c.Content == "Final step" {
			foundFinal = true
		}
		if c.Content == "Action" {
			t.Errorf("BUG: Label() value 'Action' should not appear as a message after execution")
		}
	}
	if !foundFinal {
		t.Errorf("Final progress message not found in tabContents")
	}
}
