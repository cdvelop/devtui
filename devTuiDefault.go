package devtui

import (
	"time"

	"github.com/cdvelop/messagetype"
)

const defaultTabName = "DEFAULT"

// Standard field implementation for test purposes
type standardField struct {
	name       string
	value      string
	editable   bool
	changeFunc func(string) <-chan MessageUpdate
}

func (f *standardField) Name() string {
	return f.name
}

func (f *standardField) Value() string {
	return f.value
}

func (f *standardField) Editable() bool {
	return f.editable
}

func (f *standardField) ChangeValue(newValue string) <-chan MessageUpdate {
	if f.changeFunc != nil {
		return f.changeFunc(newValue)
	}

	// Default implementation
	updates := make(chan MessageUpdate)
	go func() {
		defer close(updates)
		// Update the field value
		f.value = newValue
		// Send a success message
		updates <- MessageUpdate{
			Content: "Changed " + f.name + " to " + newValue,
			Type:    messagetype.Success,
		}
	}()
	return updates
}

// NewDefaultTUI creates a DevTUI instance with basic default configuration
// useful for unit tests and for quick initialization in real applications
func DefaultTUIForTest(LogToFile func(messageErr any)) *DevTUI {
	// Create a new DevTUI instance
	devtui := NewTUI(&TuiConfig{
		TabIndexStart: 0,               // Start with the first tab
		ExitChan:      make(chan bool), // Channel to signal exit
		Color:         nil,             // Use default colors
		LogToFile:     LogToFile,
	})

	// Create Tab 1 with basic fields
	tab1 := devtui.NewTabSection("Tab 1",
		&standardField{
			name:     "Field 1 (Editable)",
			value:    "initial test value",
			editable: true,
			changeFunc: func(newValue string) <-chan MessageUpdate {
				updates := make(chan MessageUpdate)
				go func() {
					defer close(updates)
					updates <- MessageUpdate{
						Content: "Saved value: " + newValue,
						Type:    messagetype.Success,
					}
				}()
				return updates
			},
		},
		&standardField{
			name:     "Field 2 (Non-Editable)",
			value:    "special action",
			editable: false,
			changeFunc: func(newValue string) <-chan MessageUpdate {
				updates := make(chan MessageUpdate)
				go func() {
					defer close(updates)
					updates <- MessageUpdate{
						Content: "Action executed",
						Type:    messagetype.Success,
					}
				}()
				return updates
			},
		},
	)

	// Create Tab 2 with an async field
	tab2 := devtui.NewTabSection("Tab 2",
		&standardField{
			name:     "Field 1",
			value:    "tab 2 value 1",
			editable: true,
			changeFunc: func(newValue string) <-chan MessageUpdate {
				updates := make(chan MessageUpdate)
				go func() {
					defer close(updates)
					updates <- MessageUpdate{
						Content: "Tab 2 saved: " + newValue,
						Type:    messagetype.Success,
					}
				}()
				return updates
			},
		},
		&standardField{
			name:     "Async Operation",
			value:    "Start",
			editable: true,
			changeFunc: func(newValue string) <-chan MessageUpdate {
				updates := make(chan MessageUpdate)

				go func() {
					defer close(updates)

					// Simulate a long-running operation
					for i := range 5 {
						// Send progress messages
						updates <- MessageUpdate{
							Content: "Processing step " + string(rune('A'+i)) + " for value: " + newValue,
							Type:    messagetype.Info,
						}
						time.Sleep(time.Millisecond * 100) // Shortened for tests
					}

					// Send completion message
					updates <- MessageUpdate{
						Content: "Operation completed successfully for: " + newValue,
						Type:    messagetype.Success,
					}
				}()

				return updates
			},
		},
	)

	// Ensure both tabs are properly initialized
	_ = tab1
	_ = tab2

	return devtui
}

// prepareForTesting configures a DevTUI instance for use in unit tests
func prepareForTesting() *DevTUI {
	// Create a logger that doesn't do anything during tests
	testLogger := func(messageErr any) {
		// In test mode, we don't need to log
		// This is a no-op logger for tests
	}

	// Get default TUI instance
	h := DefaultTUIForTest(testLogger)

	return h
}
