# DevTUI: Refactor Progress Callbacks to Channels

## Objective
Refactor `Change` and `Execute` methods from callback-based `func(msgs ...any)` to channel-based `chan<- string` for consistency with external tools and better streaming support.

## Current vs Target Signatures

**Current:**
```go
Change(newValue string, progress func(msgs ...any))
Execute(progress func(msgs ...any))
```

**Target:**
```go
Change(newValue string, progress chan<- string)
Execute(progress chan<- string)
```

## Rationale
1. **Consistency**: TinyWasm already uses `chan<- string` in `Change()` and MCP tools
2. **Streaming**: True asynchronous progress updates for long operations
3. **Simplicity**: Single string type instead of variadic any
4. **Idiomatic Go**: Channels are standard for goroutine communication

## Core Files to Modify

### 1. `interfaces.go` (3 interfaces)
```go
type HandlerEdit interface {
    Name() string
    Label() string
    Value() string
    Change(newValue string, progress chan<- string) // Updated
}

type HandlerExecution interface {
    Name() string
    Label() string
    Execute(progress chan<- string) // Updated
}

type HandlerInteractive interface {
    Name() string
    Label() string
    Value() string
    Change(newValue string, progress chan<- string) // Updated
    WaitingForUser() bool
}
```

### 2. `anyHandler.go` (struct + methods)
```go
type anyHandler struct {
    // ... existing fields ...
    changeFunc  func(string, chan<- string) // Updated
    executeFunc func(chan<- string)         // Updated
}

func (a *anyHandler) Change(newValue string, progress chan<- string) {
    if a.changeFunc != nil {
        a.changeFunc(newValue, progress)
    }
}

func (a *anyHandler) Execute(progress chan<- string) {
    if a.executeFunc != nil {
        a.executeFunc(progress)
    }
}
```

### 3. `anyHandler.go` - Factory Methods
Update `NewEditHandler`, `NewExecutionHandler`, `NewInteractiveHandler`:
```go
// In NewExecutionHandler
changeFunc: func(_ string, progress chan<- string) {
    h.Execute(progress)
},
```

### 4. `field.go` - Progress Callbacks
Update all `progressCallback` functions (4 locations):
- `triggerContentDisplay()` - line ~145
- `executeAsyncChange()` - line ~279  
- `executeChangeSyncWithValue()` - line ~364
- `executeChangeSyncWithTracking()` - line ~396

**Pattern:**
```go
// Create channel
progressChan := make(chan string, 10)
messages := []string{}

// Collect in goroutine
done := make(chan struct{})
go func() {
    for msg := range progressChan {
        messages = append(messages, msg)
        // Process message immediately if needed
        f.sendMessage(msg)
    }
    close(done)
}()

// Execute handler
f.handler.Change(valueToSave.(string), progressChan)
close(progressChan)
<-done
```

### 5. `shortcuts.go`
```go
func (h *shortcutsInteractiveHandler) Change(newValue string, progress chan<- string) {
    defer close(progress) // Always close when done
    
    if newValue == "" && !h.needsLanguageInput {
        progress <- h.generateHelpContent()
        return
    }
    
    h.lang = OutLang(newValue)
    h.needsLanguageInput = false
    progress <- h.generateHelpContent()
}
```

## Test Files to Update

**Edit Handlers:**
- `new_api_test.go` - `testEditHandler`
- `handler_test.go` - `TestEditableHandler`, `TestNonEditableHandler`, `PortTestHandler`, etc.
- `operation_id_reuse_test.go` - `TestOperationIDHandler`, `TestNewOperationHandler`
- `handler_value_update_test.go` - Any edit handlers
- `tabSection_move_to_end_test.go` - `testTracker`

**Execution Handlers:**
- `new_api_test.go` - `testRunHandler`
- `handler_test.go` - Execution test handlers
- `execution_footer_bug_test.go` - `ExecHandler`
- `race_condition_test.go` - `RaceConditionHandler`

**Interactive Handlers:**
- `chat_handler_test.go` - Chat handlers
- Any content display tests

**Example conversion:**
```go
// Before
func (h *TestHandler) Change(v string, progress func(msgs ...any)) {
    progress("Starting", "update")
    h.value = v
    progress("Done:", v)
}

// After  
func (h *TestHandler) Change(v string, progress chan<- string) {
    progress <- "Starting update"
    h.value = v
    progress <- fmt.Sprintf("Done: %s", v)
}
```

## External Packages

**devbrowser/tui.go:**
```go
func (h *DevBrowser) Execute(progress chan<- string) {
    if h.isOpen {
        progress <- "Closing..."
        if err := h.CloseBrowser(); err != nil {
            progress <- fmt.Sprintf("Close error: %v", err)
        } else {
            progress <- "Closed."
        }
    } else {
        progress <- "Opening..."
        h.OpenBrowser()
    }
}
```

**example/ directory:**
- `HandlerEdit.go` - `DatabaseHandler`
- `HandlerExecution.go` - `BackupHandler`
- `HandlerInteractive.go` - `SimpleChatHandler`

## Implementation Order

1. Update `interfaces.go` (3 interfaces)
2. Update `anyHandler.go` (struct, methods, factories)
3. Update `field.go` (4 progressCallback patterns)
4. Update `shortcuts.go` (1 handler)
5. Update all test files (~15 files)
6. Update example handlers (3 files)
7. Update external package: `devbrowser/tui.go`
8. Run: `go test ./...`

## Critical Notes

- **⚠️ HANDLERS MUST NOT CLOSE THE CHANNEL**: DevTUI owns the progress channel and is responsible for closing it. Handlers should ONLY send messages, never close the channel.
- **Buffered channels**: DevTUI uses size 10 to prevent blocking
- **Format first**: Use `fmt.Sprintf` to format before sending
- **Collect async**: DevTUI uses goroutine to collect messages without blocking
- **No defer close in handlers**: Handlers must never use `defer close(progress)` - this causes "close of closed channel" panics

## Breaking Changes
⚠️ **BREAKING CHANGE** - All handler implementations must update signatures.

## Success Criteria
- [ ] All interfaces updated
- [ ] All internal implementations updated  
- [ ] All test handlers updated
- [ ] External packages updated
- [ ] `go test ./...` passes
- [ ] No compilation errors
