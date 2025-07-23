package devtui_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cdvelop/devtui"
	"github.com/cdvelop/devtui/example"
)

// TestShortcutsDisplayBug reproduce el bug específico donde el contenido de shortcuts aparece verde
func TestShortcutsDisplayBug(t *testing.T) {
	t.Run("Shortcuts content should display without green color and without timestamp", func(t *testing.T) {
		// Use the same configuration as manual testing
		config := example.CreateTestConfig()
		config.AppName = "DevTUI - Shortcuts Display Test"
		config.TestMode = false // Use real behavior to reproduce the bug

		tui := devtui.NewTUI(config)

		// Add the same handlers as in cmd/main.go for consistency
		example.SetupHandlersAndTabs(tui)

		t.Logf("=== REPRODUCING SHORTCUTS DISPLAY BUG ===")

		// The SHORTCUTS tab is automatically created by NewTUI at index 0
		// Give time for auto-display if it happens (like in init.go)
		time.Sleep(100 * time.Millisecond)

		// Check the content view for green text or timestamps
		contentView := tui.ContentView()
		t.Logf("Content view:\n%s", contentView)

		// CRITICAL TESTS: The content should NOT be green and should NOT have timestamps

		// Test 1: Check if content contains ANSI color codes for green
		if strings.Contains(contentView, "\x1b[32m") || strings.Contains(contentView, "\x1b[92m") {
			t.Errorf("BUG CONFIRMED: Shortcuts content contains green color codes")
		}

		// Test 2: Check if content contains timestamp patterns (like "11:28:08")
		timestampPatterns := []string{"10:", "11:", "12:", "13:", "14:", "15:", "16:", "17:", "18:", "19:", "20:", "21:", "22:", "23:", "00:", "01:", "02:", "03:", "04:", "05:", "06:", "07:", "08:", "09:"}
		hasTimestamp := false
		for _, pattern := range timestampPatterns {
			if strings.Contains(contentView, pattern) {
				t.Errorf("BUG CONFIRMED: Shortcuts content contains timestamp pattern '%s' - should be clean for readonly fields", pattern)
				hasTimestamp = true
				break
			}
		}

		// Test 3: Check if content contains handler name in brackets
		if strings.Contains(contentView, "[Shortcuts]") {
			t.Errorf("BUG CONFIRMED: Shortcuts content contains handler name - should be clean for readonly fields")
		}

		// Test 4: The content should be the clean shortcuts text
		expectedContent := "Keyboard Navigation Commands:"
		if !strings.Contains(contentView, expectedContent) {
			t.Errorf("Expected shortcuts content not found in view")
		}

		// Summary
		if !hasTimestamp && !strings.Contains(contentView, "\x1b[32m") && !strings.Contains(contentView, "[Shortcuts]") {
			t.Log("✅ SUCCESS: Shortcuts content displays correctly without green color, timestamps, or handler names")
		}
	})
}
