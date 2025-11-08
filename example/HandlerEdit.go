package example

import "time"

type DatabaseHandler struct {
	ConnectionString string
	LastAction       string
}

func (h *DatabaseHandler) Name() string  { return "DatabaseConfig" }
func (h *DatabaseHandler) Label() string { return "Database Connection" }
func (h *DatabaseHandler) Value() string { return h.ConnectionString }

func (h *DatabaseHandler) Change(newValue string, progress chan<- string) {
	switch newValue {
	case "t":
		h.LastAction = "test"
		if progress != nil {
			progress <- "Testing database connection..."
			time.Sleep(500 * time.Millisecond)
			progress <- "Connection test completed successfully"
		}
	case "b":
		h.LastAction = "backup"
		if progress != nil {
			progress <- "Starting database backup..."
			time.Sleep(1000 * time.Millisecond)
			progress <- "Database backup completed successfully"
		}
	default:
		// Regular connection string update
		if progress != nil {
			progress <- "Validating Connection " + newValue
			time.Sleep(500 * time.Millisecond)
			progress <- "Testing database connectivity... " + newValue
			time.Sleep(500 * time.Millisecond)
			progress <- "Connection Database configured successfully " + newValue
		}
		h.ConnectionString = newValue
	}
}

// NEW: Add shortcut support
func (h *DatabaseHandler) Shortcuts() []map[string]string {
	return []map[string]string{
		{"t": "test connection"},
		{"b": "backup database"},
	}
}
