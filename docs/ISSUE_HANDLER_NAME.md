# WritingHandler Interface Implementation Plan

## Executive Summary

**Objective**: Implement a `WritingHandler` interface that provides message source identification and operation ID management for DevTUI field handlers, enabling precise control over message placement and updates.

**Integration**: Optional interface pattern that extends existing `FieldHandler` functionality without breaking changes, using `io.Writer` as the primary writing mechanism.

---

## Requirements Analysis

### Current Problems
1. **Ambiguous Message Source**: Messages don't identify which handler/component sent them
2. **No Message Update Control**: Handlers cannot control whether to create new messages or update existing ones
3. **Public Print Method Issues**: `DevTUI.Print()` allows writing to arbitrary tabs causing confusion

### Proposed Solution
1. **WritingHandler Interface**: Optional interface for handlers that need writing control
2. **Message Source Identification**: Include handler name in formatted messages  
3. **Operation ID Management**: Handler-controlled message placement and updates
4. **io.Writer Standardization**: Use `tabSection.Write()` as primary writing mechanism

---

## Technical Specification

### 1. WritingHandler Interface Definition

```go
// WritingHandler interface provides message source identification and operation ID management
// ALL handlers must implement this interface for message source control
type WritingHandler interface {
    Name() string                    // Handler identifier (e.g., "TinyWasm", "MainServer")
    SetLastOperationID(lastOpID string)  // DevTUI calls this after processing each message
    GetLastOperationID() string      // Handler returns ID for message updates, "" for new messages
}
```

### 2. FieldHandler Integration

**REQUIRED**: All FieldHandler implementations must implement WritingHandler:

```go
// FieldHandler now INCLUDES WritingHandler - all handlers must implement both
type FieldHandler interface {
    Label() string                       
    Value() string                      
    Editable() bool                     
    Change(newValue any) (string, error)
    Timeout() time.Duration
    
    // REQUIRED: WritingHandler methods
    WritingHandler
}

// Example implementation
type MyHandler struct {
    lastOpID string
    needsUpdate bool
}

func (h *MyHandler) Label() string { return "My Field" }
func (h *MyHandler) Value() string { return "value" }
func (h *MyHandler) Editable() bool { return true }
func (h *MyHandler) Change(newValue any) (string, error) { return "changed", nil }
func (h *MyHandler) Timeout() time.Duration { return 0 }

// REQUIRED WritingHandler methods
func (h *MyHandler) Name() string { return "MyHandler" }
func (h *MyHandler) SetLastOperationID(lastOpID string) { h.lastOpID = lastOpID }
func (h *MyHandler) GetLastOperationID() string { 
    if h.needsUpdate {
        return h.lastOpID  // Update existing message
    }
    return ""  // Create new message
}
```

### 3. Message Format Enhancement

Update `formatMessage` to include handler identification:

```go
func (t *DevTUI) formatMessage(msg tabContent) string {
    var timeStr string
    if t.id != nil {
        timeStr = t.timeStyle.Render(t.id.UnixNanoToTime(msg.Id))
    } else {
        timeStr = t.timeStyle.Render("--:--:--")
    }

    // NEW: Include handler name in message format
    var handlerName string
    if msg.handlerName != "" {
        handlerName = fmt.Sprintf("[%s] ", msg.handlerName)
    }

    // Apply message type styling
    switch msg.Type {
    case messagetype.Error:
        msg.Content = t.errStyle.Render(msg.Content)
    case messagetype.Warning:
        msg.Content = t.warnStyle.Render(msg.Content)
    case messagetype.Info:
        msg.Content = t.infoStyle.Render(msg.Content)
    case messagetype.Success:
        msg.Content = t.okStyle.Render(msg.Content)
    }

    return fmt.Sprintf("%s %s%s", timeStr, handlerName, msg.Content)
}
```

### 4. Enhanced tabContent Structure

```go
type tabContent struct {
    Id          string
    Content     string
    Type        messagetype.Type
    tabSection  *tabSection
    operationID *string  // Existing async operation ID
    isProgress  bool     // Existing async flags
    isComplete  bool
    
    // NEW: Handler identification
    handlerName string   // Handler name for message source identification
}
```

### 5. tabSection.Write() Integration - CRITICAL QUESTION

**PROBLEM**: How do we identify which handler is writing when using `io.Writer`?

**Current Challenge**: When external packages call `fmt.Fprintf(tabWriter, "message")`, we need to know which handler is sending the message to get the handler name and operation ID.

**Suggested Solutions**:

#### Option A: Handler Context Registry (RECOMMENDED)
Add a registry system to track active handlers per tabSection:

```go
type tabSection struct {
    // Existing fields...
    index         int
    title         string   
    fieldHandlers []*field
    sectionFooter string
    tabContents   []tabContent
    indexActiveEditField int
    tui           *DevTUI
    mu            sync.RWMutex
    
    // NEW: Handler context registry
    activeHandlers map[string]FieldHandler // handler name -> handler instance
    defaultHandler FieldHandler            // fallback handler for unknown writes
}

// NEW: Handler registration methods
func (ts *tabSection) RegisterHandler(handler FieldHandler) {
    if ts.activeHandlers == nil {
        ts.activeHandlers = make(map[string]FieldHandler)
    }
    ts.activeHandlers[handler.Name()] = handler
}

func (ts *tabSection) SetActiveHandler(handlerName string) {
    // Set which handler is currently active for io.Writer operations
    if handler, exists := ts.activeHandlers[handlerName]; exists {
        ts.defaultHandler = handler
    }
}

func (ts *tabSection) Write(p []byte) (n int, err error) {
    msg := strings.TrimSpace(string(p))
    if msg != "" {
        msgType := messagetype.DetectMessageType(msg)
        
        var handlerName string
        var operationID string
        
        // Use default handler if available
        if ts.defaultHandler != nil {
            handlerName = ts.defaultHandler.Name()
            operationID = ts.defaultHandler.GetLastOperationID()
        }
        
        ts.tui.sendMessageWithHandler(msg, msgType, ts, handlerName, operationID)
        
        if msgType == messagetype.Error {
            ts.tui.LogToFile(msg)
        }
    }
    return len(p), nil
}
```

#### Option B: Thread-Local Storage Pattern
Use goroutine-local storage to track current handler:

```go
var handlerContext = make(map[int64]FieldHandler) // goroutine ID -> handler
var handlerMutex sync.RWMutex

func (ts *tabSection) SetWritingContext(handler FieldHandler) {
    handlerMutex.Lock()
    defer handlerMutex.Unlock()
    
    // Store handler for current goroutine
    goroutineID := getGoroutineID() // helper function to get goroutine ID
    handlerContext[goroutineID] = handler
}

func (ts *tabSection) GetCurrentHandler() FieldHandler {
    handlerMutex.RLock()
    defer handlerMutex.RUnlock()
    
    goroutineID := getGoroutineID()
    return handlerContext[goroutineID]
}

func (ts *tabSection) Write(p []byte) (n int, err error) {
    msg := strings.TrimSpace(string(p))
    if msg != "" {
        msgType := messagetype.DetectMessageType(msg)
        
        var handlerName string
        var operationID string
        
        if handler := ts.GetCurrentHandler(); handler != nil {
            handlerName = handler.Name()
            operationID = handler.GetLastOperationID()
        }
        
        ts.tui.sendMessageWithHandler(msg, msgType, ts, handlerName, operationID)
        
        if msgType == messagetype.Error {
            ts.tui.LogToFile(msg)
        }
    }
    return len(p), nil
}
```

#### Option C: Wrapper Writer Pattern
Create handler-specific writers:

```go
type HandlerWriter struct {
    tabSection *tabSection
    handler    FieldHandler
}

func (hw *HandlerWriter) Write(p []byte) (n int, err error) {
    msg := strings.TrimSpace(string(p))
    if msg != "" {
        msgType := messagetype.DetectMessageType(msg)
        handlerName := hw.handler.Name()
        operationID := hw.handler.GetLastOperationID()
        
        hw.tabSection.tui.sendMessageWithHandler(msg, msgType, hw.tabSection, handlerName, operationID)
        
        if msgType == messagetype.Error {
            hw.tabSection.tui.LogToFile(msg)
        }
    }
    return len(p), nil
}

// Usage in external packages:
func (handler *MyHandler) Change(newValue any) (string, error) {
    // Get handler-specific writer
    writer := tabSection.GetHandlerWriter(handler)
    fmt.Fprintf(writer, "Processing...")
    
    return "Completed", nil
}
```

**Question for Decision**: Which approach do you prefer for handler identification in `io.Writer`?

---

## IMPLEMENTATION DECISIONS - FINALIZED

### ✅ Question 1: Handler Identification in io.Writer
**DECISION**: **Option C - Wrapper Writer Pattern**

**Rationale**:
- ✅ **Explicit and Clear**: Each handler gets its own writer, no ambiguity
- ✅ **Thread-Safe**: No shared state between handlers
- ✅ **Simple to Implement**: Straightforward wrapper pattern
- ✅ **Type-Safe**: Compiler ensures handler is provided

### ✅ Question 2: Handler Registration Timing
**DECISION**: **Option A - Auto-registration in NewField()**

**Rationale**: 
- ✅ **Zero Configuration**: Works automatically when fields are added
- ✅ **No Extra Code**: Developers don't need to remember registration calls

### ✅ Question 3: Multiple Handlers per TabSection
**DECISION**: **Option A - Handler Context Registry with SetActiveWriter()**

**Rationale**:
- ✅ **Simple Implementation**: Easy to understand and debug  
- ✅ **Explicit Control**: Developers can set active handler when needed
- ✅ **Fallback Strategy**: Clear behavior when context is unclear

### ✅ Question 4: Migration Strategy for External Packages
**DECISION**: **Hybrid Approach - RegisterWritingHandler + HandlerWriter**

**Implementation**:
```go
// For FieldHandlers (auto-registered)
serverFieldHandler := &ServerFieldHandler{...}
sectionBuild.NewField(serverFieldHandler) // Auto-registers for writing

// For independent writers (manual registration)
watcherHandler := &WatcherWritingHandler{...}
watcherWriter := sectionBuild.RegisterWritingHandler(watcherHandler)

// Pass specific writers to external components
wasmHandler := tinywasm.New(&tinywasm.Config{
    Writer: watcherWriter, // Handler-specific writer
})
```

### ✅ Question 5: DevTUI.Print() Method
**DECISION**: **Option A - Remove completely (BREAKING CHANGE)**

**Rationale**:
- ✅ **Forces Migration**: All packages must use proper io.Writer pattern
- ✅ **Clear API**: No ambiguous methods in public interface
- ✅ **Explicit Targeting**: Messages must have clear source identification

---

## Implementation Strategy

### Phase 1: Core Interface Implementation (✅ COMPLETED)
1. ✅ **Update FieldHandler interface** to include WritingHandler methods (BREAKING CHANGE)
2. ✅ **Enhance tabContent structure** with `handlerName` field
3. ✅ **Update formatMessage method** to include handler names in format: `[HandlerName]`
4. ✅ **Remove DevTUI.Print()** method completely (BREAKING CHANGE)
5. ✅ **Implement HandlerWriter wrapper pattern**

### Phase 2: Integration and Handler Context (✅ COMPLETED)
1. ✅ **Implement HandlerWriter wrapper pattern** for handler-specific writers
2. ✅ **Add RegisterWritingHandler()** method to tabSection
3. ✅ **Enhance tabSection.Write()** to use WritingHandler data with activeWriter context
4. ✅ **Add SetLastOperationID callbacks** after message processing via sendMessageWithHandler
5. ✅ **Implement auto-registration** in NewField() for FieldHandlers

### Phase 3: Testing and Validation (✅ COMPLETED)
1. ✅ **Updated ALL existing test handlers** to implement new FieldHandler interface
2. ✅ **Created comprehensive tests** for WritingHandler functionality in `writing_handler_test.go`
3. ✅ **Created integration tests** for external writer registration and HandlerWriter
4. ✅ **Validated breaking changes** - all tests pass, no regressions detected
5. ✅ **Documented implementation** with architectural decisions and usage patterns

### Phase 4: Ready for External Package Migration (📋 PENDING)
1. ❓ **Update external packages** (godev, etc.) to use new io.Writer pattern
2. ❓ **Create migration examples** demonstrating new usage patterns in godev
3. ❓ **Update external documentation** with breaking change notices
4. ❓ **Provide migration utilities** if needed for complex transitions

---

## Usage Examples

### Basic WritingHandler Implementation

```go
type DownloadHandler struct {
    downloadName string
    lastOpID     string
    isUpdating   bool
}

func (h *DownloadHandler) Name() string { return "Downloader" }
func (h *DownloadHandler) Label() string { return fmt.Sprintf("Download %s", h.downloadName) }
func (h *DownloadHandler) Value() string { return "Press Enter to download" }
func (h *DownloadHandler) Editable() bool { return false }
func (h *DownloadHandler) Timeout() time.Duration { return 30 * time.Second }

func (h *DownloadHandler) SetLastOperationID(lastOpID string) { 
    h.lastOpID = lastOpID 
}

func (h *DownloadHandler) GetLastOperationID() string { 
    if h.isUpdating {
        return h.lastOpID  // Update existing progress message
    }
    return ""  // Create new message
}

func (h *DownloadHandler) Change(newValue any) (string, error) {
    h.isUpdating = true
    defer func() { h.isUpdating = false }()
    
    // Write progress updates - these will update the same message
    fmt.Fprintf(tabWriter, "Downloading %s...", h.downloadName)
    time.Sleep(1 * time.Second)
    fmt.Fprintf(tabWriter, "Download progress: 50%%")
    time.Sleep(1 * time.Second) 
    fmt.Fprintf(tabWriter, "Download progress: 100%%")
    
    return fmt.Sprintf("Download of %s completed", h.downloadName), nil
}
```

### Handler Without WritingHandler (Uses Label)

```go
type SimpleHandler struct {
    value string
}

func (h *SimpleHandler) Label() string { return "Simple Field" }
func (h *SimpleHandler) Value() string { return h.value }
func (h *SimpleHandler) Editable() bool { return true }
func (h *SimpleHandler) Timeout() time.Duration { return 0 }

func (h *SimpleHandler) Change(newValue any) (string, error) {
    h.value = newValue.(string)
    // When this handler writes to tabWriter, messages will show:
    // "12:34:56 [Simple Field] Value updated successfully"
    return "Value updated successfully", nil
}
```

---

## Migration Path

### For External Packages (godev, etc.)
1. **Replace tui.Print() calls** with `io.Writer` usage:
   ```go
   // OLD: tui.Print("Server started")
   // NEW: fmt.Fprintf(tabWriter, "Server started")
   ```

2. **Optionally implement WritingHandler** for better control:
   ```go
   type ServerHandler struct {
       lastOpID string
   }
   
   func (h *ServerHandler) Name() string { return "MainServer" }
   // ... other WritingHandler methods
   ```

### Backward Compatibility
- **BREAKING CHANGE**: FieldHandler interface now requires WritingHandler methods
- **BREAKING CHANGE**: DevTUI.Print() method removed completely  
- **Migration Required**: All existing handlers must implement Name(), SetLastOperationID(), GetLastOperationID()
- **External Packages**: Must migrate from tui.Print() to io.Writer usage patterns

---

## Success Criteria

### Functional Requirements
1. ✅ **DECIDED**: All messages include handler identification in format `[HandlerName]`
2. ✅ **DECIDED**: All handlers must implement WritingHandler for message control  
3. ✅ **DECIDED**: io.Writer is the primary writing mechanism
4. ✅ **DECIDED**: WritingHandler is REQUIRED (embedded in FieldHandler)
5. ✅ **DECIDED**: DevTUI.Print() removed completely

### Implementation Decisions Needed
1. ❓ **PENDING**: Handler identification method in io.Writer (Options A/B/C)
2. ❓ **PENDING**: Handler registration approach (auto vs manual)  
3. ❓ **PENDING**: Multiple handlers per tabSection strategy
4. ❓ **PENDING**: External package migration pattern

### Performance Requirements  
1. ✅ No performance degradation from handler name lookup
2. ✅ Efficient operation ID management
3. ✅ Minimal memory overhead for WritingHandler support

### Usability Requirements
1. ✅ Clear documentation with examples
2. ✅ Easy migration path for existing handlers
3. ✅ Intuitive WritingHandler interface design
4. ✅ Consistent behavior across all handler types

---

## Risk Mitigation

### High Risk: Breaking Changes
- **Mitigation**: Use optional interface pattern, keep FieldHandler unchanged
- **Fallback**: Use field labels for non-WritingHandler identification

### Medium Risk: Performance Impact
- **Mitigation**: Cache handler name lookups, efficient context tracking
- **Monitoring**: Add benchmarks for message formatting performance

### Low Risk: Complex Handler Implementation  
- **Mitigation**: Provide clear examples and helper functions
- **Documentation**: Step-by-step migration guide for existing handlers

This implementation provides precise control over message source identification and placement while maintaining full backward compatibility with existing handlers.
