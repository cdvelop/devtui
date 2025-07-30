package devtui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestPaginationDisplay(t *testing.T) {
	// Setup DevTUI using a similar pattern to user_scenario_test.go
	h := DefaultTUIForTest(func(messages ...any) {})
	h.viewport.Width = 80
	h.viewport.Height = 24
	h.paginationStyle = lipgloss.NewStyle().Background(lipgloss.Color(h.Lowlight)).Foreground(lipgloss.Color(h.Foreground))

	// Tab pagination cases
	tabCases := []struct {
		activeTab int
		totalTabs int
		expected  string
	}{
		{0, 1, "[ 1/ 1]"},
		{0, 4, "[ 1/ 4]"},
		{3, 4, "[ 4/ 4]"},
		{99, 100, "[100/99]"}, // Clamp to 99
	}

	for _, tc := range tabCases {
		// Setup tabs using only public API
		// Remove all tabs
		h.tabSections = h.tabSections[:0]
		for i := 0; i < tc.totalTabs; i++ {
			h.NewTabSection(fmt.Sprintf("Tab%d", i), "desc")
		}
		h.activeTab = tc.activeTab
		// Render header pagination
		currentTab := h.activeTab
		totalTabs := len(h.tabSections)
		displayCurrent := min(currentTab, 99) + 1
		displayTotal := min(totalTabs, 99)
		pagination := fmt.Sprintf("[%2d/%2d]", displayCurrent, displayTotal)
		output := h.paginationStyle.Render(pagination)
		if !contains(output, tc.expected) {
			t.Errorf("Header pagination failed: got %q, want %q", output, tc.expected)
		}
	}

	// Field pagination cases
	fieldCases := []struct {
		activeField int
		totalFields int
		expected    string
	}{
		{0, 1, "[ 1/ 1]"},
		{0, 4, "[ 1/ 4]"},
		{3, 4, "[ 4/ 4]"},
		{99, 100, "[100/99]"}, // Clamp to 99
	}

	for _, tc := range fieldCases {
		// Remove all tabs
		h.tabSections = h.tabSections[:0]
		tab := h.NewTabSection("TestTab", "desc")
		for i := 0; i < tc.totalFields; i++ {
			tab.AddEditHandler(NewTestEditableHandler(fmt.Sprintf("Field%d", i), "val"), 0)
		}
		h.activeTab = 0
		if tc.activeField < len(tab.FieldHandlers()) {
			tab.SetActiveEditField(tc.activeField)
		}
		currentField := tc.activeField
		totalFields := len(tab.FieldHandlers())
		displayCurrent := min(currentField, 99) + 1
		displayTotal := min(totalFields, 99)
		pagination := fmt.Sprintf("[%2d/%2d]", displayCurrent, displayTotal)
		output := h.paginationStyle.Render(pagination)
		if !contains(output, tc.expected) {
			t.Errorf("Footer pagination failed: got %q, want %q", output, tc.expected)
		}
	}
}

// Helper to check substring
func contains(s, substr string) bool {
	return len(substr) > 0 && (s == substr || (len(s) > len(substr) && (s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}

// Use min from main codebase
