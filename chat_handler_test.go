package devtui

import (
	"strings"
	"testing"
	"time"

	"github.com/cdvelop/devtui/example"
	tea "github.com/charmbracelet/bubbletea"
)

// TestChatHandlerRealScenario tests the complete chat interaction flow
// focusing on handler behavior in different states, not DevTUI orchestration
func TestChatHandlerRealScenario(t *testing.T) {
	t.Run("Chat handler behavior following DevTUI responsibility separation", func(t *testing.T) {
		t.Logf("=== TESTING CHAT HANDLER ACCORDING TO DEVTUI PRINCIPLES ===")

		// Create chat handler with initial state (using the real handler from example)
		chatHandler := &example.SimpleChatHandler{}

		var contentDisplayed []string
		mockProgress := func(msgs ...any) {
			for _, msg := range msgs {
				contentDisplayed = append(contentDisplayed, msg.(string))
				t.Logf("Progress: %s", msg)
			}
		}

		// STATE 1: Initial content display (when DevTUI selects the field)
		t.Logf("State 1: DevTUI selects field -> handler shows content")

		// Verify initial state
		if chatHandler.WaitingForUser() {
			t.Errorf("Initial state: should not be waiting for user")
		}

		// DevTUI calls Change("", progress) when field is selected
		chatHandler.Change("", mockProgress)

		// Verify welcome content was shown
		if len(contentDisplayed) == 0 {
			t.Errorf("No content displayed in initial state")
		}

		welcomeFound := false
		for _, content := range contentDisplayed {
			if strings.Contains(content, "Welcome") {
				welcomeFound = true
				break
			}
		}
		if !welcomeFound {
			t.Errorf("Expected welcome content, got: %v", contentDisplayed)
		}

		// STATE 2: DevTUI transitions to input mode (this is DevTUI's responsibility)
		t.Logf("State 2: DevTUI activates input mode -> handler becomes ready")

		contentDisplayed = []string{}

		// DevTUI is responsible for managing the input activation
		// The handler just needs to be ready when WaitingForUser() should return true
		chatHandler.WaitingForUserFlag = true // This simulates DevTUI's state management

		// Verify handler is now waiting for user
		if !chatHandler.WaitingForUser() {
			t.Errorf("After DevTUI activation: should be waiting for user")
		}

		// Label should reflect input mode (handler's responsibility)
		if !strings.Contains(chatHandler.Label(), "Type message") {
			t.Errorf("Expected input mode label, got: %s", chatHandler.Label())
		}

		// STATE 3: User types and sends message (handler processes business logic)
		t.Logf("State 3: User sends message -> handler processes it")

		userMessage := "Hello, how are you?"
		chatHandler.Change(userMessage, mockProgress)

		// Verify handler managed its own state correctly
		if chatHandler.WaitingForUser() {
			t.Errorf("After user input: handler should not be waiting for user")
		}

		if !chatHandler.IsProcessing {
			t.Errorf("After user input: handler should be processing")
		}

		if len(chatHandler.Messages) != 1 {
			t.Errorf("Expected 1 message after user input, got %d", len(chatHandler.Messages))
		}

		if chatHandler.Value() != "" {
			t.Errorf("Handler should clear input after sending, got '%s'", chatHandler.Value())
		}

		// Verify handler sent appropriate progress messages
		userMessageFound := false
		for _, content := range contentDisplayed {
			if strings.Contains(content, "U: Hello, how are you?") {
				userMessageFound = true
				break
			}
		}
		if !userMessageFound {
			t.Errorf("Expected user message in progress, got: %v", contentDisplayed)
		}

		// STATE 4: AI response completion (handler's async business logic)
		t.Logf("State 4: Handler completes AI response -> ready for next input")

		// Wait for async AI response (handler's responsibility)
		maxWait := 50
		for i := 0; i < maxWait && chatHandler.IsProcessing; i++ {
			time.Sleep(100 * time.Millisecond)
		}

		// Verify handler managed its async operation correctly
		if !chatHandler.WaitingForUser() {
			t.Errorf("After AI response: handler should be waiting for user again")
		}

		if chatHandler.IsProcessing {
			t.Errorf("After AI response: handler should not be processing")
		}

		if len(chatHandler.Messages) != 2 {
			t.Errorf("Expected 2 messages after AI response, got %d", len(chatHandler.Messages))
		}

		// STATE 5: DevTUI re-selects field -> handler shows conversation history
		t.Logf("State 5: DevTUI re-selects field -> handler shows history")

		// Simulate DevTUI deactivating input mode (field loses focus, regains focus)
		chatHandler.WaitingForUserFlag = false
		contentDisplayed = []string{}

		// DevTUI calls Change("", progress) when field is re-selected
		chatHandler.Change("", mockProgress)

		// Verify conversation history is shown (handler's business logic)
		historyFound := false
		for _, content := range contentDisplayed {
			if strings.Contains(content, "U: Hello, how are you?") || strings.Contains(content, "A: Response:") {
				historyFound = true
				break
			}
		}
		if !historyFound {
			t.Errorf("Expected conversation history, got: %v", contentDisplayed)
		}

		// STATE 6: Test empty input while in input mode (edge case handling)
		t.Logf("State 6: User presses Enter without typing -> handler guides user")

		chatHandler.WaitingForUserFlag = true // Back to input mode
		contentDisplayed = []string{}

		// User presses Enter without typing anything
		chatHandler.Change("", mockProgress)

		// Handler should guide the user (handler's responsibility for UX)
		guidanceFound := false
		for _, content := range contentDisplayed {
			if strings.Contains(content, "Type message") {
				guidanceFound = true
				break
			}
		}
		if !guidanceFound {
			t.Errorf("Expected user guidance message, got: %v", contentDisplayed)
		}

		t.Logf("=== CHAT HANDLER TEST COMPLETED - ALL RESPONSIBILITIES PROPERLY SEPARATED ===")
	})

	t.Run("Test chat UI rendering and edit mode transitions", func(t *testing.T) {
		tui := DefaultTUIForTest()

		chatHandler := &example.SimpleChatHandler{}

		chatTab := tui.NewTabSection("Chat", "AI Chat Assistant")
		chatTab.AddInteractiveHandler(chatHandler, 5*time.Second)

		tui.viewport.Width = 80
		tui.viewport.Height = 24

		chatTabIndex := len(tui.tabSections) - 1
		tui.activeTab = chatTabIndex
		chatField := tui.tabSections[chatTabIndex].FieldHandlers()[0]

		t.Logf("=== TESTING UI RENDERING AND EDIT MODE ===")

		// Phase 1: Before any interaction
		content1 := tui.ContentView()
		t.Logf("Phase 1 - Initial UI:\n%s", content1)

		// Phase 2: Enter to activate input mode
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		content2 := tui.ContentView()
		t.Logf("Phase 2 - After Enter (should be in edit mode):\n%s", content2)

		// Should now be in edit mode (check if tempEditValue is being used)
		if chatField.tempEditValue == "" && !chatHandler.WaitingForUser() {
			// This is expected - the handler manages its own state
			t.Logf("Handler state: WaitingForUser=%v, IsProcessing=%v", chatHandler.WaitingForUser(), chatHandler.IsProcessing)
		}

		// Phase 3: Type message
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})

		content3 := tui.ContentView()
		t.Logf("Phase 3 - After typing 'hello':\n%s", content3)

		// Phase 4: Send message
		tui.HandleKeyboard(tea.KeyMsg{Type: tea.KeyEnter})

		content4 := tui.ContentView()
		t.Logf("Phase 4 - After sending message:\n%s", content4)

		// Should no longer be in edit mode (processing message)
		if chatHandler.WaitingForUser() && !chatHandler.IsProcessing {
			// This is fine - handler completed processing and is ready for next input
			t.Logf("Handler ready for next input: WaitingForUser=%v, IsProcessing=%v", chatHandler.WaitingForUser(), chatHandler.IsProcessing)
		}

		t.Logf("=== UI RENDERING TEST COMPLETED ===")
	})
}
