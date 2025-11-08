package devtui

import (
	"fmt"
	"testing"
	"time"
)

// TestTimeoutChannelClosePanic replicates the "close of closed channel" panic
// that occurs when a handler operation times out while still trying to send
// progress messages.
func TestTimeoutChannelClosePanic(t *testing.T) {
	t.Run("Handler that exceeds timeout should not cause panic", func(t *testing.T) {
		// Create a handler that takes longer than the timeout
		slowHandler := &SlowTestHandler{
			delay: 600 * time.Millisecond, // Longer than 500ms timeout
		}

		tui := DefaultTUIForTest()
		tui.testMode = true // Enable test mode

		tab := tui.NewTabSection("Slow Operations", "Testing timeout")

		// Add execution handler with 500ms timeout using new API
		tui.AddHandler(slowHandler, 500*time.Millisecond, "", tab)

		// Manually create field and trigger operation
		ts := tab.(*tabSection)
		if len(ts.fieldHandlers) == 0 {
			t.Fatal("Expected at least one field handler")
		}

		field := ts.fieldHandlers[0]

		// Initialize async state
		field.asyncState = &internalAsyncState{}

		// This should trigger the timeout scenario
		// The handler will try to send messages after timeout occurs
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occurred: %v", r)
			}
		}()

		// Disable test mode temporarily to test real async behavior
		tui.testMode = false
		field.handleEnter()

		// Wait for timeout to occur and goroutine to finish
		time.Sleep(800 * time.Millisecond)

		t.Log("Test completed without panic")
	})

	t.Run("Handler that sends messages after timeout should be safe", func(t *testing.T) {
		// Handler that sends multiple progress messages slowly
		verboseSlowHandler := &VerboseSlowTestHandler{
			delay:        100 * time.Millisecond, // Total: 500ms
			messageCount: 5,
		}

		tui := DefaultTUIForTest()
		tab := tui.NewTabSection("Verbose Operations", "Testing verbose timeout")

		// Add execution handler with 400ms timeout (will timeout before completion)
		tui.AddHandler(verboseSlowHandler, 400*time.Millisecond, "", tab)

		ts := tab.(*tabSection)
		field := ts.fieldHandlers[0]
		field.asyncState = &internalAsyncState{}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occurred: %v", r)
			}
		}()

		tui.testMode = false
		field.handleEnter()

		// Wait for timeout and handler to finish trying to send messages
		time.Sleep(700 * time.Millisecond)

		t.Log("Test completed without panic - handler continued sending after timeout")
	})
}

// SlowTestHandler simulates a handler that takes too long to complete
// This handler does NOT check if channel is closed before sending
type SlowTestHandler struct {
	delay time.Duration
}

func (h *SlowTestHandler) Name() string  { return "SlowHandler" }
func (h *SlowTestHandler) Label() string { return "Slow Operation" }
func (h *SlowTestHandler) Execute(progress chan<- string) {
	// Handler should NOT need to recover - devtui should handle gracefully
	progress <- "Starting slow operation..."
	time.Sleep(h.delay)
	// Even if timeout occurs, this should not crash the application
	progress <- "Operation completed"
}

// VerboseSlowTestHandler sends multiple messages with delays
type VerboseSlowTestHandler struct {
	delay        time.Duration
	messageCount int
}

func (h *VerboseSlowTestHandler) Name() string  { return "VerboseSlowHandler" }
func (h *VerboseSlowTestHandler) Label() string { return "Verbose Slow Operation" }
func (h *VerboseSlowTestHandler) Execute(progress chan<- string) {
	for i := 0; i < h.messageCount; i++ {
		select {
		case progress <- fmt.Sprintf("Progress message %d/%d", i+1, h.messageCount):
			time.Sleep(h.delay)
		default:
			// Channel might be closed, exit gracefully
			return
		}
	}
}
