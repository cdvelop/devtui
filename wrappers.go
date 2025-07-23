package devtui

import "time"

// Wrapper handlers that adapt new interfaces to FieldHandler

// displayFieldHandler wraps HandlerDisplay to implement FieldHandler
type displayFieldHandler struct {
	display HandlerDisplay
	timeout time.Duration
}

func (d *displayFieldHandler) Label() string  { return d.display.Label() }
func (d *displayFieldHandler) Value() string  { return d.display.Content() }
func (d *displayFieldHandler) Editable() bool { return false }
func (d *displayFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	return "", nil
}
func (d *displayFieldHandler) Timeout() time.Duration             { return d.timeout }
func (d *displayFieldHandler) Name() string                       { return d.display.Label() }
func (d *displayFieldHandler) SetLastOperationID(lastOpID string) {}
func (d *displayFieldHandler) GetLastOperationID() string         { return "" }

// editFieldHandler wraps HandlerEdit to implement FieldHandler
type editFieldHandler struct {
	edit     HandlerEdit
	timeout  time.Duration
	lastOpID string
}

func (e *editFieldHandler) Label() string  { return e.edit.Label() }
func (e *editFieldHandler) Value() string  { return e.edit.Value() }
func (e *editFieldHandler) Editable() bool { return true }
func (e *editFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	err := e.edit.Change(newValue, progress...)
	if err != nil {
		return "", err
	}
	return "", nil // Success message handled by progress callback
}
func (e *editFieldHandler) Timeout() time.Duration             { return e.timeout }
func (e *editFieldHandler) Name() string                       { return e.edit.Label() }
func (e *editFieldHandler) SetLastOperationID(lastOpID string) { e.lastOpID = lastOpID }
func (e *editFieldHandler) GetLastOperationID() string         { return e.lastOpID }

// runFieldHandler wraps HandlerExecution to implement FieldHandler
type runFieldHandler struct {
	run      HandlerExecution
	timeout  time.Duration
	lastOpID string
}

func (r *runFieldHandler) Label() string  { return r.run.Label() }
func (r *runFieldHandler) Value() string  { return "Execute" }
func (r *runFieldHandler) Editable() bool { return false }
func (r *runFieldHandler) Change(newValue any, progress ...func(string)) (string, error) {
	err := r.run.Execute(progress...)
	if err != nil {
		return "", err
	}
	return "", nil // Success message handled by progress callback
}
func (r *runFieldHandler) Timeout() time.Duration             { return r.timeout }
func (r *runFieldHandler) Name() string                       { return r.run.Label() }
func (r *runFieldHandler) SetLastOperationID(lastOpID string) { r.lastOpID = lastOpID }
func (r *runFieldHandler) GetLastOperationID() string         { return r.lastOpID }
