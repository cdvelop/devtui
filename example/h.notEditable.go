package example

import (
	"fmt"
	"time"
)

// BuildActionHandler - Action button example (non-editable)
type BuildActionHandler struct {
	buildType string
	lastOpID  string
}

func NewBuildActionHandler(buildType string) *BuildActionHandler {
	return &BuildActionHandler{buildType: buildType}
}

// WritingHandler implementation
func (h *BuildActionHandler) Name() string                 { return fmt.Sprintf("Build_%s", h.buildType) }
func (h *BuildActionHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *BuildActionHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *BuildActionHandler) Label() string          { return fmt.Sprintf("Build %s", h.buildType) }
func (h *BuildActionHandler) Value() string          { return "Press Enter to build" }
func (h *BuildActionHandler) Editable() bool         { return false }
func (h *BuildActionHandler) Timeout() time.Duration { return 10 * time.Second }

func (h *BuildActionHandler) Change(newValue any, progress ...func(string)) (string, error) {
	// Use progress callback for long-running operations
	if len(progress) > 0 {
		progressCallback := progress[0]
		progressCallback(fmt.Sprintf("Initiating %s build...", h.buildType))
		time.Sleep(1 * time.Second)
		progressCallback("Compiling source code...")
		time.Sleep(2 * time.Second)
		progressCallback("Build complete")
	}

	return fmt.Sprintf("%s build completed successfully", h.buildType), nil
}
