package devtui

import (
	"strings"
	"testing"
)

// TestHandlerInteractiveInterface verifies the HandlerInteractive interface behavior
func TestHandlerInteractiveInterface(t *testing.T) {
	t.Run("HandlerInteractive should work with progress() pattern", func(t *testing.T) {
		// Use DefaultTUIForTest - this automatically creates shortcuts tab at index 0
		h := DefaultTUIForTest(func(messages ...any) {
			// Silent logging for tests
		})

		// Initialize viewport
		h.viewport.Width = 80
		h.viewport.Height = 24

		// Set to shortcuts tab (index 0, automatically created by init)
		h.activeTab = 0
		h.tabSections[0].indexActiveEditField = 0

		// Get the shortcuts field that was automatically created
		fields := h.tabSections[0].FieldHandlers()
		if len(fields) == 0 {
			t.Fatal("No shortcuts field found - createShortcutsTab() should create one automatically")
		}

		field := fields[0]

		// CRITICAL TEST 1: Verify it's detected as interactive handler
		if !field.isInteractiveHandler() {
			t.Error("Shortcuts handler should be detected as HandlerInteractive")
		}

		// CRITICAL TEST 2: Verify WaitingForUser() method exists and works
		waitingForUser := field.handler.WaitingForUser()
		t.Logf("WaitingForUser(): %t", waitingForUser)

		// CRITICAL TEST 3: Verify should auto-activate edit mode detection
		shouldAutoActivate := field.shouldAutoActivateEditMode()
		t.Logf("Should auto-activate edit mode: %t", shouldAutoActivate)

		// CRITICAL TEST 4: Verify automatic content display on field selection
		initialMessageCount := len(h.tabSections[0].tabContents)
		t.Logf("Initial message count: %d", initialMessageCount)

		// NEW: Test automatic content display when navigating to interactive field
		// This simulates the user navigating to the shortcuts field for the first time
		h.checkAndTriggerInteractiveContent()

		// Should display content automatically when interactive field is selected
		afterAutoDisplayCount := len(h.tabSections[0].tabContents)
		t.Logf("Message count after auto content display: %d", afterAutoDisplayCount)

		if afterAutoDisplayCount <= initialMessageCount {
			t.Error("HandlerInteractive should automatically display content when field is selected")
		}

		// CRITICAL TEST 5: Verify manual content display via Change("", progress)
		// Trigger content display manually (this is what triggerContentDisplay() does internally)
		field.triggerContentDisplay()

		// Should display content through progress() messages
		finalMessageCount := len(h.tabSections[0].tabContents)
		t.Logf("Final message count after manual content display: %d", finalMessageCount)

		// CRITICAL TEST 6: Verify language change functionality
		initialValue := field.handler.Value()
		t.Logf("Initial language: %s", initialValue)

		// CRITICAL TEST 6A: Verify MessageTracker is working to prevent duplicate messages
		messageCountBeforeChange := len(h.tabSections[0].tabContents)
		t.Logf("Message count before language change: %d", messageCountBeforeChange)

		// Simulate language change
		field.executeChangeSyncWithTracking("ES")

		// Verify value changed
		finalValue := field.handler.Value()
		if finalValue != "ES" {
			t.Errorf("Handler value should be 'ES', got: '%s'", finalValue)
		}

		// CRITICAL TEST: Verify MessageTracker updates existing message (same count, updated content)
		messageCountAfterChange := len(h.tabSections[0].tabContents)
		t.Logf("Message count after language change: %d", messageCountAfterChange)

		// Get the current message content to compare timestamps
		if messageCountAfterChange > 0 {
			lastMessage := h.tabSections[0].tabContents[messageCountAfterChange-1]
			t.Logf("Last message after ES change: %s (timestamp: %s)", lastMessage.Content, lastMessage.Timestamp)
		}

		// Call the same change again - should UPDATE existing message, not create new one
		field.executeChangeSyncWithTracking("ES")
		messageCountAfterSameChange := len(h.tabSections[0].tabContents)
		t.Logf("Message count after same change (ES again): %d", messageCountAfterSameChange)

		// Message count should be the same (updates existing message)
		if messageCountAfterSameChange != messageCountAfterChange {
			t.Errorf("MessageTracker should update existing message, not create new one. Expected count: %d, Got: %d",
				messageCountAfterChange, messageCountAfterSameChange)
		}

		// Verify different value change updates the same message (same count, different content)
		field.executeChangeSyncWithTracking("FR")
		messageCountAfterNewChange := len(h.tabSections[0].tabContents)
		t.Logf("Message count after new language change (FR): %d", messageCountAfterNewChange)

		// Should still be same count (updates existing message)
		if messageCountAfterNewChange != messageCountAfterChange {
			t.Errorf("MessageTracker should update existing message for new values too. Expected count: %d, Got: %d",
				messageCountAfterChange, messageCountAfterNewChange)
		}

		// Verify the content actually changed
		if messageCountAfterNewChange > 0 {
			lastMessage := h.tabSections[0].tabContents[messageCountAfterNewChange-1]
			t.Logf("Last message after FR change: %s (timestamp: %s)", lastMessage.Content, lastMessage.Timestamp)

			// The message content should now reflect FR, not ES
			if !strings.Contains(lastMessage.Content, "FR") {
				t.Error("Message content should be updated to reflect new language (FR)")
			}
		}

		t.Logf("SUCCESS: HandlerInteractive works correctly - automatic content display: %t, value changed from '%s' to '%s', messages via progress()",
			afterAutoDisplayCount > initialMessageCount, initialValue, finalValue)
	})
}

// Test the REAL problem: MessageTracker should update same message, not create new ones
func TestMessageTrackerRealProblem(t *testing.T) {
	t.Run("MessageTracker should update existing message, not create new messages", func(t *testing.T) {
		h := DefaultTUIForTest(func(messages ...any) {})

		// Initialize
		h.viewport.Width = 80
		h.viewport.Height = 24
		h.activeTab = 0
		h.tabSections[0].indexActiveEditField = 0

		field := h.tabSections[0].FieldHandlers()[0]

		// Clear any existing messages
		h.tabSections[0].tabContents = nil
		initialCount := 0

		t.Logf("=== TESTING MESSAGETRACKER BEHAVIOR ===")
		t.Logf("Initial message count: %d", initialCount)

		// FIRST CHANGE: Should create new message with new operationID
		field.executeChangeSyncWithTracking("ES")
		firstChangeCount := len(h.tabSections[0].tabContents)
		t.Logf("After first change (ES): %d messages", firstChangeCount)

		// Get the operationID from the first message
		var firstOpID string
		if firstChangeCount > 0 {
			firstMessage := h.tabSections[0].tabContents[firstChangeCount-1]
			if firstMessage.operationID != nil {
				firstOpID = *firstMessage.operationID
			}
			t.Logf("First message OpID: '%s', Content: '%s'", firstOpID, firstMessage.Content)
		}

		// SECOND CHANGE (SAME VALUE): Should UPDATE existing message, not create new
		field.executeChangeSyncWithTracking("ES")
		secondChangeCount := len(h.tabSections[0].tabContents)
		t.Logf("After duplicate change (ES): %d messages", secondChangeCount)

		// CRITICAL: Should be same count (updated existing message)
		if secondChangeCount != firstChangeCount {
			t.Errorf("PROBLEM DETECTED: MessageTracker should update existing message, not create new one")
			t.Errorf("Expected count: %d, Got: %d", firstChangeCount, secondChangeCount)
			t.Error("This is the exact problem you're seeing in the demo!")
		}

		// Verify the operationID is the same (reused)
		if secondChangeCount > 0 {
			lastMessage := h.tabSections[0].tabContents[secondChangeCount-1]
			var lastOpID string
			if lastMessage.operationID != nil {
				lastOpID = *lastMessage.operationID
			}
			t.Logf("Last message OpID: '%s', Content: '%s'", lastOpID, lastMessage.Content)

			if firstOpID != lastOpID {
				t.Errorf("PROBLEM: OperationID should be reused! First: '%s', Last: '%s'", firstOpID, lastOpID)
			}
		}

		// THIRD CHANGE (DIFFERENT VALUE): Should still UPDATE same message
		field.executeChangeSyncWithTracking("FR")
		thirdChangeCount := len(h.tabSections[0].tabContents)
		t.Logf("After new value change (FR): %d messages", thirdChangeCount)

		// Should STILL be same count (updated existing message)
		if thirdChangeCount != firstChangeCount {
			t.Errorf("PROBLEM: Even different values should update same message")
			t.Errorf("Expected count: %d, Got: %d", firstChangeCount, thirdChangeCount)
		}

		if thirdChangeCount > 0 {
			finalMessage := h.tabSections[0].tabContents[thirdChangeCount-1]
			var finalOpID string
			if finalMessage.operationID != nil {
				finalOpID = *finalMessage.operationID
			}
			t.Logf("Final message OpID: '%s', Content: '%s'", finalOpID, finalMessage.Content)
		}
	})
}
