package example

import "time"

type DatabaseHandler struct {
	ConnectionString string
}

func (h *DatabaseHandler) Name() string  { return "DatabaseConfig" }
func (h *DatabaseHandler) Label() string { return "Database Connection" }
func (h *DatabaseHandler) Value() string { return h.ConnectionString }
func (h *DatabaseHandler) Change(newValue string, progress func(msgs ...any)) {
	if progress != nil {
		progress("Validating", "Connection", newValue)
		time.Sleep(500 * time.Millisecond)
		progress("Testing database connectivity...", newValue)
		time.Sleep(500 * time.Millisecond)
		progress("Connection", "Database", "configured", "successfully", newValue)
	}
	h.ConnectionString = newValue
}
