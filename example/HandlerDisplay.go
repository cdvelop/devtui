package example

import (
	. "github.com/cdvelop/tinystring"
)

type StatusHandler struct{}

func (h *StatusHandler) Name() string { return T(D.Information, D.Status, D.System) }
func (h *StatusHandler) Content() string {
	return "Status: Running\nPID: 12345\nUptime: 2h 30m\nMemory: 45MB\nCPU: 12%"
}
