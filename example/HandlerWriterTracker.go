package example

type OperationLogWriter struct {
	lastOpID string
}

func (w *OperationLogWriter) Name() string                 { return "OperationLog" }
func (w *OperationLogWriter) GetLastOperationID() string   { return w.lastOpID }
func (w *OperationLogWriter) SetLastOperationID(id string) { w.lastOpID = id }
