package devtui

import (
	"testing"
)

// Test simulating the exact demo behavior
func TestDemoMessageDuplication(t *testing.T) {
	t.Run("Simulate real demo usage to detect duplication", func(t *testing.T) {
		h := DefaultTUIForTest(func(messages ...any) {})

		// Initialize like demo
		h.viewport.Width = 80
		h.viewport.Height = 24
		h.activeTab = 0
		h.tabSections[0].indexActiveEditField = 0

		field := h.tabSections[0].fieldHandlers[0]

		// Clear messages like demo start
		h.tabSections[0].tabContents = nil

		t.Logf("=== SIMULATING DEMO BEHAVIOR ===")

		// 1. User navigates to shortcuts field (triggers auto-display)
		h.checkAndTriggerInteractiveContent()
		afterAutoDisplay := len(h.tabSections[0].tabContents)
		t.Logf("After auto-display: %d messages", afterAutoDisplay)

		// 2. User presses Enter to edit (normal demo flow)
		// This might trigger another content display
		field.triggerContentDisplay()
		afterManualTrigger := len(h.tabSections[0].tabContents)
		t.Logf("After manual trigger: %d messages", afterManualTrigger)

		// 3. User types "ES" and presses Enter (demo user input flow)
		field.executeChangeSyncWithTracking("ES")
		afterFirstChange := len(h.tabSections[0].tabContents)
		t.Logf("After first change (ES): %d messages", afterFirstChange)

		// Print all messages to see what's happening
		for i, msg := range h.tabSections[0].tabContents {
			var opID string
			if msg.operationID != nil {
				opID = *msg.operationID
			}
			t.Logf("Message %d: OpID='%s', Content='%s'", i, opID, msg.Content)
		}

		// 4. User changes to "FR" - should UPDATE same message
		field.executeChangeSyncWithTracking("FR")
		afterSecondChange := len(h.tabSections[0].tabContents)
		t.Logf("After second change (FR): %d messages", afterSecondChange)

		// Print all messages again
		for i, msg := range h.tabSections[0].tabContents {
			var opID string
			if msg.operationID != nil {
				opID = *msg.operationID
			}
			t.Logf("Final Message %d: OpID='%s', Content='%s'", i, opID, msg.Content)
		}

		// CRITICAL: If messages with same OpID exist, that's duplication!
		opIDCounts := make(map[string]int)
		for _, msg := range h.tabSections[0].tabContents {
			if msg.operationID != nil {
				opIDCounts[*msg.operationID]++
			}
		}

		for opID, count := range opIDCounts {
			if count > 1 {
				t.Errorf("DUPLICATION DETECTED: OpID '%s' appears %d times!", opID, count)
				t.Error("This is the problem you're seeing in the demo!")
			}
		}
	})
}
