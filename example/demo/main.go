package main

import (
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// Example showcasing all new handler types with minimal implementation

// 1. HandlerDisplay - Read-only information display (2 methods - 75% reduction)
type StatusHandler struct{}

func (h *StatusHandler) Label() string { return "System Status" }
func (h *StatusHandler) Content() string {
	return "Status: Running\nPID: 12345\nUptime: 2h 30m\nMemory: 45MB\nCPU: 12%"
}

// 2. HandlerEdit - Interactive input fields (3 methods - 62.5% reduction)
type DatabaseHandler struct {
	connectionString string
}

func (h *DatabaseHandler) Label() string { return "Database Connection" }
func (h *DatabaseHandler) Value() string { return h.connectionString }
func (h *DatabaseHandler) Change(newValue any, progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Validating connection string...")
		time.Sleep(500 * time.Millisecond)
		progress[0]("Testing database connectivity...")
		time.Sleep(1 * time.Second)
		progress[0]("Database connection configured successfully")
	}
	h.connectionString = newValue.(string)
	return nil
}

// 3. HandlerExecution - Action buttons (2 methods - 75% reduction)
type BackupHandler struct{}

func (h *BackupHandler) Label() string { return "Create System Backup" }
func (h *BackupHandler) Execute(progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Preparing backup...")
		time.Sleep(1 * time.Second)
		progress[0]("Backing up database...")
		time.Sleep(2 * time.Second)
		progress[0]("Backing up files...")
		time.Sleep(1 * time.Second)
		progress[0]("Backup completed successfully")
	}
	return nil
}

// 4. HandlerWriter - Simple logging (1 method - 87.5% reduction)
type SystemLogWriter struct{}

func (w *SystemLogWriter) Name() string { return "SystemLog" }

// 5. HandlerTrackerWriter - Advanced logging with message tracking (3 methods)
type OperationLogWriter struct {
	lastOpID string
}

func (w *OperationLogWriter) Name() string                 { return "OperationLog" }
func (w *OperationLogWriter) GetLastOperationID() string   { return w.lastOpID }
func (w *OperationLogWriter) SetLastOperationID(id string) { w.lastOpID = id }

func main() {
	tui := devtui.NewTUI(&devtui.TuiConfig{
		AppName:  "New API Demo",
		ExitChan: make(chan bool),
	})

	// Method chaining with optional timeout configuration
	// Demonstrates 60-85% reduction in boilerplate code

	// Dashboard tab with DisplayHandlers (read-only information)
	dashboard := tui.NewTabSection("Dashboard", "System Overview")
	dashboard.NewDisplayHandler(&StatusHandler{}).Register()

	// Configuration tab with EditHandlers (interactive fields)
	config := tui.NewTabSection("Config", "System Configuration")
	config.NewEditHandler(&DatabaseHandler{connectionString: "postgres://localhost:5432/mydb"}).WithTimeout(3 * time.Second)

	// Operations tab with RunHandlers (action buttons)
	ops := tui.NewTabSection("Operations", "System Operations")
	ops.NewRunHandler(&BackupHandler{}).WithTimeout(15 * time.Second)

	// Logging tab with Writers
	logs := tui.NewTabSection("Logs", "System Logs")

	// Basic writer (always creates new lines)
	systemWriter := logs.NewWriterHandler(&SystemLogWriter{}).Register()
	systemWriter.Write([]byte("System initialized"))
	systemWriter.Write([]byte("API demo started"))

	// Advanced writer (can update existing messages)
	opWriter := logs.NewWriterHandler(&OperationLogWriter{}).Register()
	opWriter.Write([]byte("Operation tracking enabled"))

	// Different timeout configurations:
	// - Synchronous (default): timeout = 0
	// - Asynchronous with specific timeout: WithTimeout(duration)
	// - Millisecond precision for testing: WithTimeout(100*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
