package main

import (
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// Example showcasing all new handler types with minimal implementation

// 1. HandlerDisplay - Read-only information display (3 methods)
type StatusHandler struct{}

func (h *StatusHandler) Name() string  { return "SystemStatus" }
func (h *StatusHandler) Label() string { return "System Status" }
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
func (h *DatabaseHandler) Change(newValue any, progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Validating connection string...")
		time.Sleep(200 * time.Millisecond)
		progress[0]("Testing database connectivity...")
		time.Sleep(400 * time.Millisecond)
		progress[0]("Database connection configured successfully")
	}
	h.connectionString = newValue.(string)
	return nil
}

// 3. HandlerExecution - Action buttons (3 methods) with tracking
type BackupHandler struct {
	lastOpID string
}

func (h *BackupHandler) Name() string  { return "SystemBackup" }
func (h *BackupHandler) Label() string { return "Create System Backup" }
func (h *BackupHandler) Execute(progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Preparing backup...")
		time.Sleep(200 * time.Millisecond)
		progress[0]("Backing up database...")
		time.Sleep(500 * time.Millisecond)
		progress[0]("Backing up files...")
		time.Sleep(300 * time.Millisecond)
		progress[0]("Backup completed successfully")
	}
	return nil
}

// MessageTracker implementation for operation tracking
func (h *BackupHandler) GetLastOperationID() string   { return h.lastOpID }
func (h *BackupHandler) SetLastOperationID(id string) { h.lastOpID = id }

// 4. HandlerWriter - Simple logging (1 method)
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
		AppName:   "New API Demo",
		ExitChan:  make(chan bool),
		LogToFile: func(messages ...any) {},
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

	// Advanced writer (can update existing messages with tracking)
	opWriter := logs.RegisterHandlerTrackerWriter(&OperationLogWriter{})
	opWriter.Write([]byte("Operation tracking enabled"))

	// Different timeout configurations:
	// - Synchronous (default): .Register() or timeout = 0
	// - Asynchronous with timeout: .WithTimeout(duration)
	// - Example timeouts: 100*time.Millisecond, 2*time.Second, 1*time.Minute
	// - Tip: Keep timeouts reasonable (2-10 seconds) for good UX

	// Handler Types Summary:
	// • HandlerDisplay: Name() + Label() + Content() - Shows immediate content
	// • HandlerEdit: Name() + Label() + Value() + Change() - Interactive fields
	// • HandlerExecution: Name() + Label() + Execute() - Action buttons
	// • HandlerWriter: Name() - Basic logging (new lines)
	// • HandlerTrackerWriter: Name() + MessageTracker - Advanced logging (can update)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
