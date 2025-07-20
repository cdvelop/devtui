package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// SlowTaskHandler - Example with long running operation and timeout
type SlowTaskHandler struct {
	taskName string
}

func (h *SlowTaskHandler) Label() string { return h.taskName }
func (h *SlowTaskHandler) Value() string { return "Press Enter to start slow task" }
func (h *SlowTaskHandler) Editable() bool { return false }
func (h *SlowTaskHandler) Timeout() time.Duration { return 10 * time.Second }

func (h *SlowTaskHandler) Change(newValue any) (string, error) {
	// Simulate long running operation
	time.Sleep(3 * time.Second)
	return fmt.Sprintf("%s completed successfully", h.taskName), nil
}

// QuickValidationHandler - Example with fast validation
type QuickValidationHandler struct {
	currentValue string
}

func (h *QuickValidationHandler) Label() string { return "Quick Validation" }
func (h *QuickValidationHandler) Value() string { return h.currentValue }
func (h *QuickValidationHandler) Editable() bool { return true }
func (h *QuickValidationHandler) Timeout() time.Duration { return 2 * time.Second }

func (h *QuickValidationHandler) Change(newValue any) (string, error) {
	text := strings.TrimSpace(newValue.(string))
	if text == "" {
		return "", fmt.Errorf("value cannot be empty")
	}
	if len(text) < 3 {
		return "", fmt.Errorf("value must be at least 3 characters")
	}
	
	// Simulate some processing
	time.Sleep(1 * time.Second)
	
	h.currentValue = text
	return fmt.Sprintf("Validated and saved: %s", text), nil
}

// NumberHandler - Example with number validation and spinner
type NumberHandler struct {
	currentValue string
}

func (h *NumberHandler) Label() string { return "Number Field" }
func (h *NumberHandler) Value() string { return h.currentValue }
func (h *NumberHandler) Editable() bool { return true }
func (h *NumberHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *NumberHandler) Change(newValue any) (string, error) {
	numStr := strings.TrimSpace(newValue.(string))
	if numStr == "" {
		return "", fmt.Errorf("number cannot be empty")
	}
	
	// Simulate validation processing
	time.Sleep(2 * time.Second)
	
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("must be a valid number")
	}
	if num < 0 {
		return "", fmt.Errorf("number must be positive")
	}
	
	h.currentValue = numStr
	return fmt.Sprintf("Number validated: %d", num), nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Async Spinner Test",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000", 
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"DevTUI Log:"}, messages...)...)
		},
	}

	tui := devtui.NewTUI(config)

	// Create handlers with different async behaviors
	quickValidation := &QuickValidationHandler{currentValue: "initial"}
	numberField := &NumberHandler{currentValue: "42"}
	slowTask1 := &SlowTaskHandler{taskName: "Build Project"}
	slowTask2 := &SlowTaskHandler{taskName: "Deploy App"}

	// Configuration tab with editable fields
	tui.NewTabSection("Config", "Configuration with async validation").
		NewField(quickValidation).
		NewField(numberField)
		
	// Actions tab with long-running operations  
	tui.NewTabSection("Actions", "Long-running async operations").
		NewField(slowTask1).
		NewField(slowTask2)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
