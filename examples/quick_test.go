package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// Simple handlers for quick testing

type QuickTestHandler struct {
	label string
	value string
	delay time.Duration
}

func (h *QuickTestHandler) Label() string          { return h.label }
func (h *QuickTestHandler) Value() string          { return h.value }
func (h *QuickTestHandler) Editable() bool         { return false }
func (h *QuickTestHandler) Timeout() time.Duration { return h.delay + (5 * time.Second) }

func (h *QuickTestHandler) Change(newValue any) (string, error) {
	time.Sleep(h.delay)
	return fmt.Sprintf("Operation completed in %v", h.delay), nil
}

type EditableTestHandler struct {
	label        string
	currentValue string
}

func (h *EditableTestHandler) Label() string          { return h.label }
func (h *EditableTestHandler) Value() string          { return h.currentValue }
func (h *EditableTestHandler) Editable() bool         { return true }
func (h *EditableTestHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *EditableTestHandler) Change(newValue any) (string, error) {
	newVal := strings.TrimSpace(newValue.(string))
	if newVal == "" {
		return "", fmt.Errorf("%s cannot be empty", h.label)
	}

	// Simulate validation
	time.Sleep(500 * time.Millisecond)

	h.currentValue = newVal
	return fmt.Sprintf("%s updated to: %s", h.label, newVal), nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Quick Async Test",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#1A1A1A",
			Highlight:  "#00FF88",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"Test Log:"}, messages...)...)
		},
	}

	tui := devtui.NewTUI(config)

	// Quick test operations
	tui.NewTabSection("Quick Test", "Fast async operations for testing").
		NewField(&QuickTestHandler{"Fast Op", "Click to run", 1 * time.Second}).
		NewField(&QuickTestHandler{"Medium Op", "Click to run", 3 * time.Second}).
		NewField(&QuickTestHandler{"Slow Op", "Click to run", 5 * time.Second}).
		NewField(&EditableTestHandler{"Server Name", "localhost"}).
		NewField(&EditableTestHandler{"Port", "8080"})

	fmt.Println("🚀 DevTUI Quick Async Test")
	fmt.Println("✨ Test spinner animations")
	fmt.Println("⚡ Test async operations")
	fmt.Println("🔧 Test editable fields")
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
