package devtui

import (
	"sync"
	"testing"
	"time"
)

// RaceConditionHandler simulates a handler that causes race conditions
type RaceConditionHandler struct {
	lastOpID string
	mu       sync.Mutex
}

func (h *RaceConditionHandler) Name() string  { return "RaceTest" }
func (h *RaceConditionHandler) Label() string { return "Race Condition Test" }
func (h *RaceConditionHandler) Execute(progress func(msgs ...any)) {
	// Simulate work that triggers progress updates
	if progress != nil {
		for i := 0; i < 10; i++ {
			progress("Step ", string(rune('0'+i)))
			time.Sleep(1 * time.Millisecond) // Small delay to increase race probability
		}
	}
}

func (h *RaceConditionHandler) GetLastOperationID() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.lastOpID
}

func (h *RaceConditionHandler) SetLastOperationID(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastOpID = id
}

func TestRaceConditionReproduction(t *testing.T) {
	// This test should reproduce the race condition
	// Run with: go test -race -run TestRaceConditionReproduction

	config := &TuiConfig{
		AppName:  "Race Test",
		ExitChan: make(chan bool),
	}

	tui := NewTUI(config)
	tab := tui.NewTabSection("Test", "Race Condition Test")

	// Create multiple handlers to increase concurrency
	handlers := make([]*RaceConditionHandler, 5)
	for i := 0; i < 5; i++ {
		handlers[i] = &RaceConditionHandler{}
		tab.AddExecutionHandler(handlers[i], 100*time.Millisecond)
	}

	// Simulate concurrent executions
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Each goroutine tries to execute handlers concurrently
			for j := 0; j < 10; j++ {
				handlerIndex := j % len(handlers)
				handler := handlers[handlerIndex]

				// Simulate the race condition by calling SetLastOperationID concurrently
				go func() {
					handler.SetLastOperationID("op-" + string(rune('0'+goroutineID)) + "-" + string(rune('0'+j)))
				}()

				// Also execute the handler to trigger the race condition path
				go func() {
					handler.Execute(func(msgs ...any) {
						// This progress callback will trigger SetLastOperationID internally
					})
				}()

				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentOperationIDAccess(t *testing.T) {
	// Test concurrent access to SetLastOperationID and GetLastOperationID
	// This focuses specifically on the race condition in anyHandler

	handler := &RaceConditionHandler{}

	// Create anyHandler through factory method (same as used in AddExecutionHandler)
	anyH := newExecutionHandler(handler, 100*time.Millisecond)

	var wg sync.WaitGroup
	numWriters := 50
	numReaders := 50

	// Writers - multiple goroutines setting operation ID
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				anyH.SetLastOperationID("writer-" + string(rune('0'+id)) + "-" + string(rune('0'+j)))
				time.Sleep(1 * time.Microsecond)
			}
		}(i)
	}

	// Readers - multiple goroutines reading operation ID
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = anyH.GetLastOperationID()
				time.Sleep(1 * time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
}

func TestAnyHandlerConcurrentAccess(t *testing.T) {
	// Direct test of anyHandler race condition
	handler := &RaceConditionHandler{}
	anyH := newExecutionHandler(handler, 100*time.Millisecond)

	var wg sync.WaitGroup
	numGoroutines := 100
	operationsPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				// This should trigger the race condition
				anyH.SetLastOperationID("test-" + string(rune('0'+goroutineID)))
				_ = anyH.GetLastOperationID()
			}
		}(i)
	}

	wg.Wait()
}

// TestChatHandlerRaceCondition reproduces the specific race condition found in SimpleChatHandler
// This test targets the exact race condition reported by the race detector:
// Write at generateAIResponse() vs Read at TestChatHandlerRealScenario
// DISABLED: This test intentionally creates race conditions and should not be run with -race
/*
func TestChatHandlerRaceCondition(t *testing.T) {
	t.Skip("Skipping intentional race condition test - use only for debugging race conditions")

	t.Log("Testing race condition in SimpleChatHandler fields")

	// Creating a local handler that reproduces the same race condition pattern
	type RacyChatHandler struct {
		IsProcessing       bool
		WaitingForUserFlag bool
		Messages           []string
		CurrentInput       string
	}

	handler := &RacyChatHandler{
		WaitingForUserFlag: true,
	}

	var wg sync.WaitGroup
	stopFlag := make(chan bool)

	// Goroutine 1: Simulates the test reading handler state (like TestChatHandlerRealScenario)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-stopFlag:
				return
			default:
				// These reads race with writes in the generateAIResponse equivalent
				_ = handler.IsProcessing       // Read that races with write
				_ = handler.WaitingForUserFlag // Read that races with write
				_ = len(handler.Messages)      // Read that races with write

				// Tiny delay to increase race probability
				time.Sleep(time.Microsecond)
			}
		}
	}()

	// Goroutine 2: Simulates generateAIResponse writing to the same fields
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			// Simulate the writes that happen in generateAIResponse
			handler.Messages = append(handler.Messages, "response")
			handler.IsProcessing = false      // Write that races with read
			handler.WaitingForUserFlag = true // Write that races with write

			time.Sleep(time.Microsecond * 10)
		}
	}()

	// Let the race condition develop
	time.Sleep(time.Millisecond * 100)

	// Stop the reader goroutine
	close(stopFlag)
	wg.Wait()

	t.Logf("Race condition test completed. Final state: Processing=%v, Messages=%d",
		handler.IsProcessing, len(handler.Messages))
}
*/
