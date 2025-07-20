package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// SimpleTextHandler - Example of an editable field
type SimpleTextHandler struct {
	currentValue string
}

func (h *SimpleTextHandler) Label() string { return "Text Field" }
func (h *SimpleTextHandler) Value() string { return h.currentValue }
func (h *SimpleTextHandler) Editable() bool { return true }
func (h *SimpleTextHandler) Timeout() time.Duration { return 0 } // No timeout for simple validation

func (h *SimpleTextHandler) Change(newValue any) (string, error) {
	text := strings.TrimSpace(newValue.(string))
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}
	
	h.currentValue = text
	return fmt.Sprintf("Text updated to: %s", text), nil
}

// PortHandler - Example with validation and timeout
type PortHandler struct {
	currentPort string
}

func (h *PortHandler) Label() string { return "Port" }
func (h *PortHandler) Value() string { return h.currentPort }
func (h *PortHandler) Editable() bool { return true }
func (h *PortHandler) Timeout() time.Duration { return 3 * time.Second } // 3 second timeout

func (h *PortHandler) Change(newValue any) (string, error) {
	portStr := newValue.(string)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("port must be between 1 and 65535")
	}
	
	h.currentPort = portStr
	return fmt.Sprintf("Port configured: %d", port), nil
}

// BuildHandler - Example of a non-editable action button with async operation
type BuildHandler struct {
	projectPath string
}

func (h *BuildHandler) Label() string { return "Build Project" }
func (h *BuildHandler) Value() string { return "Press Enter to build" }
func (h *BuildHandler) Editable() bool { return false }
func (h *BuildHandler) Timeout() time.Duration { return 10 * time.Second } // 10 second timeout for build

func (h *BuildHandler) Change(newValue any) (string, error) {
	// Simulate long running build operation
	time.Sleep(2 * time.Second)
	
	return "Build completed successfully", nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Handler Example",
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

	// Create handlers
	textHandler := &SimpleTextHandler{currentValue: "initial text"}
	portHandler := &PortHandler{currentPort: "8080"}
	buildHandler := &BuildHandler{projectPath: "./"}

	// Use new handler-based API
	tui.NewTabSection("Configuration", "Edit configuration values").
		NewField(textHandler).
		NewField(portHandler)
		
	tui.NewTabSection("Actions", "Build operations").
		NewField(buildHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
