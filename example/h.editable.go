package example

import (
	"fmt"
	"strings"
	"time"
)

// HostConfigHandler - Editable field example
type HostConfigHandler struct {
	currentHost string
	lastOpID    string
}

func NewHostConfigHandler(initialHost string) *HostConfigHandler {
	return &HostConfigHandler{currentHost: initialHost}
}

// WritingHandler implementation
func (h *HostConfigHandler) Name() string                 { return "HostConfigHandler" }
func (h *HostConfigHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *HostConfigHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *HostConfigHandler) Label() string          { return "Host" }
func (h *HostConfigHandler) Value() string          { return h.currentHost }
func (h *HostConfigHandler) Editable() bool         { return true }
func (h *HostConfigHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *HostConfigHandler) Change(newValue any, progress ...func(string)) (string, error) {
	host := strings.TrimSpace(newValue.(string))
	if host == "" {
		return "", fmt.Errorf("host cannot be empty")
	}

	// Use progress callback for real-time updates
	if len(progress) > 0 {
		progressCallback := progress[0]
		progressCallback("Validating host configuration...")
		time.Sleep(500 * time.Millisecond)
		progressCallback("Checking network connectivity...")
		time.Sleep(500 * time.Millisecond)
		progressCallback("Host validation complete")
	}

	h.currentHost = host
	return fmt.Sprintf("Host configured: %s", host), nil
}
