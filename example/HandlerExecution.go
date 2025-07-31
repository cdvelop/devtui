package example

import "time"

type BackupHandler struct {
	lastOpID string
}

func (h *BackupHandler) Name() string  { return "SystemBackup" }
func (h *BackupHandler) Label() string { return "With Tracking" }
func (h *BackupHandler) Execute(progress func(msgs ...any)) {
	if progress != nil {
		progress("Preparing", "backup...", h.lastOpID)
		time.Sleep(500 * time.Millisecond)
		progress("BackingUp", "database...", h.lastOpID)
		time.Sleep(500 * time.Millisecond)
		progress("BackingUp", "Files", h.lastOpID)
		time.Sleep(500 * time.Millisecond)
		progress("Backup", "End", "OK", h.lastOpID)
	}
}

func (h *BackupHandler) GetLastOperationID() string   { return h.lastOpID }
func (h *BackupHandler) SetLastOperationID(id string) { h.lastOpID = id }
