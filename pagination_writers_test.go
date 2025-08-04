package devtui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestPaginationWritersOnlyTab(t *testing.T) {
	h := DefaultTUIForTest(func(messages ...any) {})
	h.viewport.Width = 80
	h.viewport.Height = 24
	h.paginationStyle = lipgloss.NewStyle().Background(lipgloss.Color(h.Lowlight)).Foreground(lipgloss.Color(h.Foreground))

	// Create a tab with only writers, no field handlers
	h.tabSections = h.tabSections[:0]
	logs := h.NewTabSection("Logs", "System Logs")
	// Minimal SystemLogWriter for test

	h.activeTab = 0
	_ = logs.RegisterHandlerWriter(&SystemLogWriter{})

	h.activeTab = 0
	// Call the real footerView rendering logic
	output := h.footerView()
	expected := "1/ 1" // Look for the core pagination text without spacing
	// Use strings.Contains for ANSI-aware searching
	found := false
	for i := 0; i <= len(output)-len(expected); i++ {
		if output[i:i+len(expected)] == expected {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Writers-only tab pagination failed: got %q, want %q", output, expected)
	}
}

// Minimal SystemLogWriter for test
type SystemLogWriter struct{}

func (w *SystemLogWriter) Name() string { return "SystemLog" }
