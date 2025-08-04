package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/devtui"

	example "github.com/cdvelop/devtui/example"
)

func main() {
	tui := devtui.NewTUI(&devtui.TuiConfig{
		AppName:  "Demo",
		ExitChan: make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(messages...) // Replace with actual logging implementation
		},
	})

	// Method chaining with optional timeout configuration
	// New API dramatically simplifies handler implementation

	// Dashboard tab with DisplayHandlers (read-only information)
	dashboard := tui.NewTabSection("Dashboard", "System Overview")
	dashboard.AddDisplayHandler(&example.StatusHandler{})

	// Configuration tab with EditHandlers (interactive fields)
	config := tui.NewTabSection("Config", "System Configuration")
	config.AddEditHandler(&example.DatabaseHandler{ConnectionString: "postgres://localhost:5432/mydb"}, 2*time.Second)
	config.AddExecutionHandler(&example.BackupHandler{}, 5*time.Second)

	// NEW: Chat tab with InteractiveHandler - Demonstrates interactive content management
	chat := tui.NewTabSection("Chat", "AI Chat Assistant")
	chatHandler := &example.SimpleChatHandler{
		Messages:           make([]example.ChatMessage, 0),
		WaitingForUserFlag: false, // Start showing content, not waiting for input
		IsProcessing:       false, // Not processing initially
	}
	chat.AddInteractiveHandler(chatHandler, 3*time.Second)

	// Logging tab with Writers
	logs := tui.NewTabSection("Logs", "System Logs")

	// Basic writer (always creates new lines)
	systemWriter := logs.NewWriter("SystemLogWriter", false)
	systemWriter.Write([]byte("System initialized"))
	systemWriter.Write([]byte("API demo started"))
	systemWriter.Write([]byte("Chat interface enabled"))

	// Generate multiple log entries to test scrolling (30 total)
	go func() {
		for i := 1; i <= 30; i++ {
			time.Sleep(3 * time.Second) // Simulate processing delay
			systemWriter.Write([]byte(fmt.Sprintf("System log entry #%d - Processing data batch", i)))
		}
	}()

	// Advanced writer (can update existing messages with tracking)
	opWriter := logs.NewWriter("OperationLogWriter", true)
	opWriter.Write([]byte("Operation tracking enabled"))

	// Generate more tracking entries to test Page Up/Page Down navigation
	go func() {
		for i := 1; i <= 50; i++ {
			time.Sleep(3 * time.Second) // Simulate processing delay
			opWriter.Write([]byte(fmt.Sprintf("Operation #%d - Background task completed successfully", i)))
		}
	}()

	// Different timeout configurations:
	// - Synchronous (default): .Register() or timeout = 0
	// - Asynchronous with timeout: .WithTimeout(duration)
	// - Example timeouts: 100*time.Millisecond, 2*time.Second, 1*time.Minute
	// - Tip: Keep timeouts reasonable (2-10 seconds) for good UX

	// Handler Types Summary:
	// • HandlerDisplay: Name() + Content() - Shows immediate content
	// • HandlerEdit: Name() + Label() + Value() + Change() - Interactive fields
	// • HandlerExecution: Name() + Label() + Execute() - Action buttons
	// • HandlerInteractive: Name() + Label() + Value() + Change() + WaitingForUser() - Interactive content
	// • HandlerWriter: Name() - Basic logging (new lines)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
