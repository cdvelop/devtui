package devtui

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestManualCrashReproduction tries to reproduce the exact crash scenario
// that occurs when manually executing the Backup handler
func TestManualCrashReproduction(t *testing.T) {
	t.Run("Reproduce exact main.go scenario with manual navigation", func(t *testing.T) {
		// Recreate the exact same setup as main.go
		tui := NewTUI(&TuiConfig{
			AppName:  "New API Demo",
			ExitChan: make(chan bool),
			TestMode: true, // Enable test mode to avoid full UI
			LogToFile: func(messages ...any) {
				t.Logf("TUI Log: %v", messages)
			},
		})

		// Exact same setup as main.go
		dashboard := tui.NewTabSection("Dashboard", "System Overview")
		statusHandler := &RealStatusHandler{}
		dashboard.NewDisplayHandler(statusHandler).Register()

		config := tui.NewTabSection("Config", "System Configuration")
		dbHandler := &RealDatabaseHandler{connectionString: "postgres://localhost:5432/mydb"}
		config.NewEditHandler(dbHandler).WithTimeout(2 * time.Millisecond)

		ops := tui.NewTabSection("Operations", "System Operations")
		backupHandler := &RealBackupHandler{}
		ops.NewExecutionHandler(backupHandler).WithTimeout(5 * time.Millisecond)

		logs := tui.NewTabSection("Logs", "System Logs")
		systemWriter := logs.RegisterHandlerWriter(&RealSystemLogWriter{})
		systemWriter.Write([]byte("System initialized"))
		systemWriter.Write([]byte("API demo started"))

		opWriter := logs.RegisterHandlerTrackerWriter(&RealOperationLogWriter{})
		opWriter.Write([]byte("Operation tracking enabled"))

		// Initialize viewport like in real execution
		tui.viewport.Width = 80
		tui.viewport.Height = 24
		tui.ready = true

		// Test navigation exactly as user would do manually
		t.Logf("Starting manual navigation simulation...")

		// User starts on SHORTCUTS tab (index 0), navigates to Operations tab (index 3)
		// Simulate Tab key presses to get to Operations
		tui.activeTab = 0 // Start at SHORTCUTS

		// Tab -> Dashboard (1)
		continueParsing, _ := tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})
		if !continueParsing {
			t.Fatal("Navigation to Dashboard failed")
		}

		// Tab -> Config (2)
		continueParsing, _ = tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})
		if !continueParsing {
			t.Fatal("Navigation to Config failed")
		}

		// Tab -> Operations (3) - This is where the crash might occur
		continueParsing, _ = tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})
		if !continueParsing {
			t.Fatal("Navigation to Operations failed")
		}

		// Verify we're on Operations tab
		if tui.activeTab != 3 {
			t.Fatalf("Expected to be on Operations tab (3), but on tab %d", tui.activeTab)
		}

		// Verify the backup handler field exists
		opsTab := tui.tabSections[3]
		if len(opsTab.FieldHandlers()) == 0 {
			t.Fatal("Operations tab should have backup handler field")
		}

		// Ensure we're focused on the backup handler field
		opsTab.indexActiveEditField = 0

		t.Logf("About to execute Backup handler - this is where crash might occur...")

		// This is the critical moment - execute the Backup handler
		continueParsing, cmd := tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		// Check for crash indicators
		if !continueParsing {
			t.Error("CRASH DETECTED: HandleKeyboard returned false, indicating UI stopped parsing")
		}

		if len(tui.tabSections) == 0 {
			t.Error("CRASH DETECTED: Tab sections disappeared")
		}

		if tui.tabSections[3] == nil {
			t.Error("CRASH DETECTED: Operations tab became nil")
		}

		// Try to continue navigation after handler execution
		continueParsing, _ = tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})
		if !continueParsing {
			t.Error("CRASH DETECTED: Navigation failed after handler execution")
		}

		// Check handler execution
		if !backupHandler.wasExecuted {
			t.Error("Backup handler was not executed")
		}

		// Test command handling if any
		if cmd != nil {
			t.Logf("Command returned: %T", cmd)
		}

		t.Logf("Manual crash reproduction test completed")
	})

	t.Run("Test unixid initialization failure - should handle gracefully now", func(t *testing.T) {
		// This test verifies that when unixid fails to initialize,
		// the system now handles it gracefully instead of panicking

		// Create a TUI instance where unixid failed to initialize
		tui := &DevTUI{
			TuiConfig: &TuiConfig{
				AppName:  "Crash Test",
				ExitChan: make(chan bool),
				TestMode: false, // DISABLE test mode to use async execution path
				LogToFile: func(messages ...any) {
					t.Logf("TUI Log: %v", messages)
				},
			},
			tabSections:     []*tabSection{},
			activeTab:       0,
			tabContentsChan: make(chan tabContent, 100),
			tuiStyle:        newTuiStyle(nil),
			id:              nil, // This simulates unixid initialization failure
		}

		// Add a tab and handler
		ops := tui.NewTabSection("Operations", "Test Operations")
		backupHandler := &RealBackupHandler{}
		ops.NewExecutionHandler(backupHandler).WithTimeout(5 * time.Millisecond)

		// Ensure the field has asyncState initialized
		if len(ops.FieldHandlers()) > 0 {
			field := ops.FieldHandlers()[0]
			if field.asyncState == nil {
				field.asyncState = &internalAsyncState{}
			}
		}

		tui.viewport.Width = 80
		tui.viewport.Height = 24
		tui.ready = true
		tui.activeTab = 0 // First (and only) tab - no SHORTCUTS added manually

		// Execute the handler - this should now handle gracefully instead of panicking
		t.Logf("Calling HandleKeyboard(Enter) to trigger async execution...")
		continueParsing, _ := tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		if !continueParsing {
			t.Error("HandleKeyboard should continue parsing even with nil unixid")
		}

		// Wait for async execution to complete
		t.Logf("Waiting for async execution to complete gracefully...")
		time.Sleep(100 * time.Millisecond)
		time.Sleep(1 * time.Second)

		// Verify that the application continued to work despite nil unixid
		if !backupHandler.wasExecuted {
			t.Error("Handler should have been executed despite nil unixid")
		}

		// Verify interface is still stable
		if len(tui.tabSections) == 0 {
			t.Error("Interface should remain stable despite nil unixid")
		}

		t.Logf("SUCCESS: unixid initialization failure now handled gracefully")
	})

	t.Run("Test rapid repeated execution like real user behavior", func(t *testing.T) {
		// This simulates a user rapidly pressing Enter multiple times
		tui := NewTUI(&TuiConfig{
			AppName:   "Crash Test",
			ExitChan:  make(chan bool),
			TestMode:  true,
			LogToFile: func(messages ...any) {},
		})

		ops := tui.NewTabSection("Operations", "System Operations")
		backupHandler := &RealBackupHandler{}
		ops.NewExecutionHandler(backupHandler).WithTimeout(5 * time.Millisecond)

		tui.viewport.Width = 80
		tui.viewport.Height = 24
		tui.ready = true
		tui.activeTab = 1 // Operations tab (after SHORTCUTS)

		// Rapid fire execution - like impatient user
		for i := 0; i < 10; i++ {
			continueParsing, _ := tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})
			if !continueParsing {
				t.Errorf("CRASH on iteration %d: HandleKeyboard stopped parsing", i)
				break
			}
			time.Sleep(1 * time.Millisecond) // Very rapid execution
		}

		if len(tui.tabSections) == 0 {
			t.Error("CRASH DETECTED: Rapid execution caused tab sections to disappear")
		}
	})

	t.Run("Test during async execution interruption", func(t *testing.T) {
		// This tests what happens if user navigates while async operation is running
		tui := NewTUI(&TuiConfig{
			AppName:   "Async Test",
			ExitChan:  make(chan bool),
			TestMode:  false, // Disable test mode to allow real async
			LogToFile: func(messages ...any) {},
		})

		ops := tui.NewTabSection("Operations", "System Operations")
		longHandler := &RealLongBackupHandler{} // Takes 3+ seconds
		ops.NewExecutionHandler(longHandler).WithTimeout(5 * time.Millisecond)

		tui.viewport.Width = 80
		tui.viewport.Height = 24
		tui.ready = true
		tui.activeTab = 1

		// Start async execution
		continueParsing, _ := tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})
		if !continueParsing {
			t.Fatal("Failed to start async execution")
		}

		// Immediately try to navigate while async operation is running
		// This might cause the crash
		time.Sleep(5 * time.Millisecond) // Let async start

		for i := 0; i < 5; i++ {
			continueParsing, _ = tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyTab})
			if !continueParsing {
				t.Errorf("CRASH during async execution: Navigation failed on iteration %d", i)
			}
			time.Sleep(5 * time.Millisecond)
		}

		// Wait for async to complete
		time.Sleep(5 * time.Millisecond)

		if len(tui.tabSections) == 0 {
			t.Error("CRASH DETECTED: Async execution + navigation caused crash")
		}
	})
}

// Real handlers that match main.go exactly

type RealStatusHandler struct{}

func (h *RealStatusHandler) Name() string  { return "SystemStatus" }
func (h *RealStatusHandler) Label() string { return "System Status" }
func (h *RealStatusHandler) Content() string {
	return "Status: Running\nPID: 12345\nUptime: 2h 30m\nMemory: 45MB\nCPU: 12%"
}

type RealDatabaseHandler struct {
	connectionString string
}

func (h *RealDatabaseHandler) Name() string  { return "DatabaseConfig" }
func (h *RealDatabaseHandler) Label() string { return "Database Connection" }
func (h *RealDatabaseHandler) Value() string { return h.connectionString }
func (h *RealDatabaseHandler) Change(newValue any, progress ...func(string)) error {
	if len(progress) > 0 {
		progress[0]("Validating connection string...")
		time.Sleep(20 * time.Millisecond)
		progress[0]("Testing database connectivity...")
		time.Sleep(5 * time.Millisecond)
		progress[0]("Database connection configured successfully")
	}
	h.connectionString = newValue.(string)
	return nil
}

type RealBackupHandler struct {
	wasExecuted bool
}

func (h *RealBackupHandler) Name() string  { return "SystemBackup" }
func (h *RealBackupHandler) Label() string { return "Create System Backup" }
func (h *RealBackupHandler) Execute(progress ...func(string)) error {
	h.wasExecuted = true
	// ALWAYS call progress to trigger the crash
	if len(progress) > 0 {
		progress[0]("Preparing backup...")
		time.Sleep(10 * time.Millisecond)
		progress[0]("Backing up database...")
		time.Sleep(50 * time.Millisecond)
		progress[0]("Backing up files...")
		time.Sleep(30 * time.Millisecond)
		progress[0]("Backup completed successfully")
	} else {
		// Even if no progress callback, this should not happen in normal execution
		return fmt.Errorf("no progress callback provided")
	}
	return nil
}

type RealLongBackupHandler struct {
	wasExecuted bool
}

func (h *RealLongBackupHandler) Name() string  { return "LongBackup" }
func (h *RealLongBackupHandler) Label() string { return "Long System Backup" }
func (h *RealLongBackupHandler) Execute(progress ...func(string)) error {
	h.wasExecuted = true
	if len(progress) > 0 {
		progress[0]("Starting very long backup...")
		time.Sleep(10 * time.Millisecond)
		progress[0]("Processing files...")
		time.Sleep(20 * time.Millisecond)
		progress[0]("Finalizing backup...")
		time.Sleep(10 * time.Millisecond)
		progress[0]("Long backup completed")
	}
	return nil
}

type RealSystemLogWriter struct{}

func (w *RealSystemLogWriter) Name() string { return "SystemLog" }

type RealOperationLogWriter struct {
	lastOpID string
}

func (w *RealOperationLogWriter) Name() string                 { return "OperationLog" }
func (w *RealOperationLogWriter) GetLastOperationID() string   { return w.lastOpID }
func (w *RealOperationLogWriter) SetLastOperationID(id string) { w.lastOpID = id }
