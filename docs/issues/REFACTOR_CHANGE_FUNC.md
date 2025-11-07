# DevTUI: Refactor Change Method to Use Channel-Based Progress

## Objective
Refactor the `Change` method signature across DevTUI library from callback-based `func(msgs ...any)` to channel-based `chan<- string` for consistency with external tools and better streaming support.

## Current Signature
```go
// interfaces.go - HandlerEdit & HandlerInteractive
Change(newValue string, progress func(msgs ...any))

// anyHandler.go - Internal implementation
changeFunc func(string, func(msgs ...any))
```

## Target Signature
```go
// interfaces.go - HandlerEdit & HandlerInteractive
Change(newValue string, progress chan<- string)

// anyHandler.go - Internal implementation
changeFunc func(string, chan<- string)
```

## Rationale
1. **Consistency**: External tools (TinyWasm, GoServer) will use channels for progress streaming
2. **Streaming**: Channels enable true asynchronous progress updates
3. **Simplicity**: Single message type (string) instead of variadic any
4. **Idiomatic Go**: Channels are Go's standard for communication between goroutines

## Files to Modify

### 1. `/interfaces.go`
Update interface definitions:
```go
type HandlerEdit interface {
    Name() string
    Label() string
    Value() string
    Change(newValue string, progress chan<- string) // Changed from func(msgs ...any)
}

type HandlerInteractive interface {
    Name() string
    Label() string
    Value() string
    Change(newValue string, progress chan<- string) // Changed from func(msgs ...any)
    WaitingForUser() bool
}
```

### 2. `/anyHandler.go`
Update internal function pointer and method:
```go
type anyHandler struct {
    // ... existing fields ...
    changeFunc func(string, chan<- string) // Changed from func(string, func(msgs ...any))
    // ... rest of fields ...
}

func (a *anyHandler) Change(newValue string, progress chan<- string) {
    if a.changeFunc != nil {
        a.changeFunc(newValue, progress)
    }
}
```

### 3. `/handlerRegistration.go`
Update handler creation for Edit handlers:
```go
func (ts *tabSection) registerEditHandler(h HandlerEdit, timeout time.Duration, color string) {
    handler := &anyHandler{
        // ... existing fields ...
        changeFunc: func(newValue string, progress chan<- string) {
            h.Change(newValue, progress)
        },
        // ... rest of fields ...
    }
    // ... registration logic ...
}
```

Update for Interactive handlers similarly.

### 4. `/field.go`
Update `executeChangeSyncWithTracking` method to use channel instead of callback:
```go
func (f *field) executeChangeSyncWithTracking(newValue string) {
    // Create buffered channel
    progressChan := make(chan string, 10)
    messages := []string{}
    
    // Collect messages in goroutine
    done := make(chan bool)
    go func() {
        for msg := range progressChan {
            messages = append(messages, msg)
        }
        done <- true
    }()
    
    // Execute change (sends messages to channel)
    f.handler.Change(newValue, progressChan)
    close(progressChan)
    
    // Wait for collection
    <-done
    
    // Process collected messages
    for _, msg := range messages {
        // ... existing message processing logic ...
    }
}
```

### 5. `/shortcuts.go`
Update `shortcutsInteractiveHandler` implementation:
```go
func (h *shortcutsInteractiveHandler) Change(newValue string, progress chan<- string) {
    if newValue == "" && !h.needsLanguageInput {
        progress <- h.generateHelpContent()
        return
    }
    
    lang := OutLang(newValue)
    h.lang = lang
    h.needsLanguageInput = false
    
    progress <- h.generateHelpContent()
}
```

### 6. Tests (All)
Update all test handlers that implement HandlerEdit or HandlerInteractive:

**Files to update:**
- `/new_api_test.go` - `testEditHandler`
- `/handler_value_update_test.go` - `ThreadSafePortTestHandler`
- `/operation_id_reuse_test.go` - `TestOperationIDHandler`, `TestNewOperationHandler`
- `/content_handler_test.go` - Any interactive handlers
- `/chat_handler_test.go` - Chat handlers
- `/cursor_behavior_test.go` - Test handlers
- `/field_editing_bug_test.go` - Bug test handlers
- `/empty_field_enter_test.go` - Test handlers

Example test update:
```go
// Before
func (h *testEditHandler) Change(newValue string, progress func(msgs ...any)) {
    h.value = newValue
    progress("Value updated to:", newValue)
}

// After
func (h *testEditHandler) Change(newValue string, progress chan<- string) {
    h.value = newValue
    progress <- fmt.Sprintf("Value updated to: %s", newValue)
}
```

## Implementation Steps

1. **Update interfaces.go** - Change method signatures
2. **Update anyHandler.go** - Change function pointer type and method
3. **Update handlerRegistration.go** - Update handler wrappers
4. **Update field.go** - Refactor executeChangeSyncWithTracking to use channels
5. **Update shortcuts.go** - Update shortcutsInteractiveHandler
6. **Update all tests** - Change all test handler implementations
7. **Run tests** - Verify all tests pass: `go test ./...`

## Breaking Changes
⚠️ **This is a BREAKING CHANGE** - All implementations of HandlerEdit and HandlerInteractive must update their Change method signature.

## Success Criteria
- [ ] All interfaces updated with new signature
- [ ] All internal implementations updated
- [ ] All test handlers updated
- [ ] All tests pass: `go test ./...`
- [ ] No compilation errors
- [ ] Channel properly closed after use
- [ ] Messages collected and processed correctly

## Notes
- Use buffered channels (size 10) to avoid blocking
- Always close channel after sending messages
- Collect messages in separate goroutine to avoid deadlock
- Format messages before sending (use `fmt.Sprintf` if needed)
- Preserve existing message tracking and timeout behavior
