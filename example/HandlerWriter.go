package example

type SystemLogWriter struct{}

func (w *SystemLogWriter) Name() string { return "SystemLog" }
