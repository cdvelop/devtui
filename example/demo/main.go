package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// Example showcasing all new handler types with minimal implementation

// 1. HandlerDisplay - Read-only information display (2 methods)
type StatusHandler struct{}

func (h *StatusHandler) Name() string { return "System Status Information Display" }
func (h *StatusHandler) Content() string {
	return "Status: Running\nPID: 12345\nUptime: 2h 30m\nMemory: 45MB\nCPU: 12%"
}

// 2. HandlerEdit - Interactive input fields (4 methods)
type DatabaseHandler struct {
	connectionString string
}

func (h *DatabaseHandler) Name() string  { return "DatabaseConfig" }
func (h *DatabaseHandler) Label() string { return "Database Connection" }
func (h *DatabaseHandler) Value() string { return h.connectionString }
func (h *DatabaseHandler) Change(newValue string, progress func(string)) {
	if progress != nil {
		progress("Validating connection string...")
		time.Sleep(200 * time.Millisecond)
		progress("Testing database connectivity...")
		time.Sleep(400 * time.Millisecond)
		progress("Database connection configured successfully")
	}
	h.connectionString = newValue
}

// 3. HandlerExecution - Action buttons (3 methods) with tracking
type BackupHandler struct {
	lastOpID string
}

func (h *BackupHandler) Name() string  { return "SystemBackup" }
func (h *BackupHandler) Label() string { return "Create System Backup" }
func (h *BackupHandler) Execute(progress func(string)) {
	if progress != nil {
		progress("Preparing backup...")
		time.Sleep(200 * time.Millisecond)
		progress("Backing up database...")
		time.Sleep(500 * time.Millisecond)
		progress("Backing up files...")
		time.Sleep(300 * time.Millisecond)
		progress("Backup completed successfully")
	}
	// ...existing code...
}

// MessageTracker implementation for operation tracking
func (h *BackupHandler) GetLastOperationID() string   { return h.lastOpID }
func (h *BackupHandler) SetLastOperationID(id string) { h.lastOpID = id }

// 4. HandlerWriter - Simple logging (1 method)
type SystemLogWriter struct{}

func (w *SystemLogWriter) Name() string { return "SystemLog" }

// 5. HandlerWriterTracker - Advanced logging with message tracking (3 methods)
type OperationLogWriter struct {
	lastOpID string
}

func (w *OperationLogWriter) Name() string                 { return "OperationLog" }
func (w *OperationLogWriter) GetLastOperationID() string   { return w.lastOpID }
func (w *OperationLogWriter) SetLastOperationID(id string) { w.lastOpID = id }

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
	dashboard.NewDisplayHandler(&StatusHandler{}).Register()

	// Configuration tab with EditHandlers (interactive fields)
	config := tui.NewTabSection("Config", "System Configuration")
	config.NewEditHandler(&DatabaseHandler{connectionString: "postgres://localhost:5432/mydb"}).WithTimeout(2 * time.Second)

	// Operations tab with ExecutionHandlers (action buttons)
	ops := tui.NewTabSection("Operations", "System Operations")
	ops.NewExecutionHandlerTracking(&BackupHandler{}).WithTimeout(5 * time.Second)

	// Logging tab with Writers
	logs := tui.NewTabSection("Logs", "System Logs")

	// Basic writer (always creates new lines)
	systemWriter := logs.RegisterHandlerWriter(&SystemLogWriter{})
	systemWriter.Write([]byte("System initialized"))
	systemWriter.Write([]byte("API demo started"))

	// Generate multiple log entries to test scrolling (30 total)
	for i := 1; i <= 15; i++ {
		systemWriter.Write([]byte(fmt.Sprintf("System log entry #%d - Processing data batch", i)))
	}

	// Advanced writer (can update existing messages with tracking)
	opWriter := logs.RegisterHandlerWriter(&OperationLogWriter{})
	opWriter.Write([]byte("Operation tracking enabled"))

	// Generate more tracking entries to test Page Up/Page Down navigation
	for i := 1; i <= 13; i++ {
		opWriter.Write([]byte(fmt.Sprintf("Operation #%d - Background task completed successfully", i)))
	}

	// Different timeout configurations:
	// - Synchronous (default): .Register() or timeout = 0
	// - Asynchronous with timeout: .WithTimeout(duration)
	// - Example timeouts: 100*time.Millisecond, 2*time.Second, 1*time.Minute
	// - Tip: Keep timeouts reasonable (2-10 seconds) for good UX

	// Handler Types Summary:
	// • HandlerDisplay: Name() + Content() - Shows immediate content
	// • HandlerEdit: Name() + Label() + Value() + Change() - Interactive fields
	// • HandlerExecution: Name() + Label() + Execute() - Action buttons
	// • HandlerWriter: Name() - Basic logging (new lines)
	// • HandlerWriterTracker: Name() + MessageTracker - Advanced logging (can update)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
