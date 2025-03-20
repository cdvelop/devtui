package devtui

import (
	"sync"
	"testing"
	"time"

	"github.com/cdvelop/messagetype"
)

// RunAsyncFieldTest is a helper function for testing async field functionality
func RunAsyncFieldTest(t *testing.T, tui *DevTUI, tabIndex, fieldIndex int, testValue string) []Message {
	// Get the tab and field
	tab := &tui.tabSections[tabIndex]
	field := &tab.FieldHandlers[fieldIndex]

	if !field.IsAsync {
		t.Fatalf("Field is not configured as async")
	}

	// Create a channel to collect messages
	msgChan := make(chan Message, 10)
	tui.SetAsyncMessageChannel(msgChan)

	// Process the field value change in a goroutine
	done := make(chan bool)
	go func() {
		tui.ProcessFieldValueChange(field, testValue)
		// Señal para indicar que procesamiento terminó
		done <- true
	}()

	// Collect messages with a shorter timeout to avoid test hanging
	messages := collectAsyncMessages(t, msgChan, done, 2*time.Second)

	return messages
}

// collectAsyncMessages collects async messages from a channel with a timeout
func collectAsyncMessages(t *testing.T,
	msgChan chan Message,
	done chan bool,
	timeout time.Duration) []Message {

	var messages []Message
	timeoutChan := time.After(timeout)

	collecting := true
	for collecting {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				// Canal cerrado
				t.Log("Message channel closed")
				collecting = false
				break
			}
			messages = append(messages, msg)

		case <-done:
			// Procesamiento terminado, esperamos un poco más por mensajes finales
			t.Log("Processing done signal received")
			time.Sleep(100 * time.Millisecond)

			// Recolectar mensajes pendientes sin bloquear
			for {
				select {
				case msg, ok := <-msgChan:
					if !ok {
						collecting = false
						break
					}
					messages = append(messages, msg)
				default:
					// No hay más mensajes
					collecting = false
					break
				}

				if !collecting {
					break
				}
			}

			collecting = false

		case <-timeoutChan:
			t.Logf("Timeout reached after collecting %d messages", len(messages))
			collecting = false
		}
	}

	return messages
}

// drainChannel recolecta todos los mensajes disponibles sin bloquear
func drainChannel(ch chan Message) []Message {
	var messages []Message

	// Intentar leer mensajes sin bloquear
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return messages
			}
			messages = append(messages, msg)
		default:
			// No hay más mensajes disponibles
			return messages
		}
	}
}

// CollectAsyncMessages runs a function that triggers async messages and collects them with a timeout
func CollectAsyncMessages(t *testing.T, triggerFn func(), timeout time.Duration) []Message {
	// Create a channel to receive messages
	msgChan := make(chan Message, 10)

	// Run the trigger function that should send messages to the channel
	go triggerFn()

	// Collect messages with timeout
	var messages []Message
	timeoutChan := time.After(timeout)
	collecting := true

	for collecting {
		select {
		case msg := <-msgChan:
			messages = append(messages, msg)
		case <-timeoutChan:
			collecting = false
		}
	}

	return messages
}

// VerifyAsyncMessages checks that messages follow the expected pattern
func VerifyAsyncMessages(t *testing.T, messages []Message, expectedCount int) {
	if len(messages) != expectedCount {
		t.Errorf("Expected %d messages, got %d", expectedCount, len(messages))
	}

	// Check progress messages
	for i := 0; i < len(messages)-1 && i < expectedCount-1; i++ {
		if messages[i].Type != messagetype.Info {
			t.Errorf("Message %d should be of type Info, got %v", i+1, messages[i].Type)
		}
	}

	// Check completion message
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		if lastMsg.Type != messagetype.OK {
			t.Errorf("Final message should be of type OK, got %v", lastMsg.Type)
		}
	}
}

// MockAsyncProcessor creates a mock async processor function for testing
func MockAsyncProcessor(steps int, delay time.Duration) func(string, chan<- Message) {
	return func(value string, msgChan chan<- Message) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Send progress messages
			for i := 0; i < steps; i++ {
				msgChan <- Message{
					Content: "Step " + string(rune('A'+i)) + " for " + value,
					Type:    messagetype.Info,
				}
				time.Sleep(delay)
			}

			// Send completion
			msgChan <- Message{
				Content: "Completed processing " + value,
				Type:    messagetype.OK,
			}
		}()
	}
}
