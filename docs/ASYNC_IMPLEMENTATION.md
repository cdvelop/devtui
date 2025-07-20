# DevTUI Async Field Operations - Implementation Plan

## Objective
Transform DevTUI field operations from synchronous (UI-blocking) to asynchronous with transparent internal handling.

## Current Architecture Analysis

### Files to Modify
Based on analysis of current codebase:

#### Core Files:
- `/home/cesar/Dev/Pkg/Mine/devtui/field.go` - Add async state, modify field struct
- `/home/cesar/Dev/Pkg/Mine/devtui/tabSection.go` - Update NewField method signature  
- `/home/cesar/Dev/Pkg/Mine/devtui/userKeyboard.go` - Trigger async operations on Enter
- `/home/cesar/Dev/Pkg/Mine/devtui/update.go` - Handle async messages in bubbletea cycle
- `/home/cesar/Dev/Pkg/Mine/devtui/print.go` - Update tabContent struct for async messages

#### Message System:
- `/home/cesar/Dev/Pkg/Mine/devtui/print.go` - Extend existing tabContent with async fields
- `/home/cesar/Dev/Pkg/Mine/messagetype/messagetype.go` - Use existing message type detection
- `/home/cesar/Dev/Pkg/Mine/unixid/unixid.go` - Use existing ID generation

#### Example Usage:
- `/home/cesar/Dev/Pkg/Mine/devtui/cmd/main.go` - Update to new API
- Create `/home/cesar/Dev/Pkg/Mine/devtui/examples/` directory with handler examples

## API Design (CORRECTED)

### Current API (from README.md)
```go
tui.NewTabSection("Tab", "Description").
    NewField("Label", "value", true, changeFunc)

// Current changeFunc signature
func changeFunc(newValue any) (string, error) {
    // Synchronous operation that blocks UI
    return "success message", nil
}
```

### New API (Simplified Interface-based)
```go
tui.NewTabSection("Tab", "Description").
    NewField(&handler)

// Handler interface - replaces individual parameters
type FieldHandler interface {
    Label() string
    Value() string  
    Editable() bool
    Change(newValue any) (string, error)  // SAME signature as current changeFunc
    
    // NEW: Optional timeout configuration
    Timeout() time.Duration  // Return 0 for no timeout, or specific duration
    
    // REMOVED: ProgressMessage() - Not needed, Change() already returns message
}
```

### Key Design Principles
1. **Transparent Async**: DevTUI handles all async complexity internally
2. **Same Change Signature**: `Change(any) (string, error)` - no context, no channels
3. **Simple Interface**: Handler only provides metadata and change logic
4. **Internal Management**: DevTUI manages goroutines, contexts, channels internally

## Implementation Details

### Modified Structures

#### field struct (internal changes)
```go
// field.go - internal struct modifications
type field struct {
    // MODIFIED: Change from individual parameters to handler-based approach
    handler      FieldHandler    // NEW: Replaces name, value, editable, changeFunc
    
    // NEW: Internal async state (not exposed)
    asyncState   *internalAsyncState
    
    // NEW: Spinner for visual feedback (using charmbracelet/bubbles/spinner)
    spinner      spinner.Model
    
    // UNCHANGED: Existing fields from current implementation
    tempEditValue string // use for edit (already exists)
    index         int    // already exists  
    cursor        int    // cursor position in text value (already exists)
    
    // REMOVED: These will be replaced by handler methods
    // name       string
    // value      string  
    // editable   bool
    // changeFunc func(newValue any) (string, error)
}

// Internal async management (not public)
type internalAsyncState struct {
    isRunning    bool
    operationID  string
    cancel       context.CancelFunc
    startTime    time.Time
}
```

#### view.go modifications (field rendering with spinner)
```go
// MODIFIED: Field rendering to show spinner when async operation is running
func (f *field) renderField() string {
    label := f.handler.Label()
    value := f.handler.Value()
    
    var statusIndicator string
    if f.asyncState != nil && f.asyncState.isRunning {
        // Show animated spinner during async operations
        statusIndicator = f.spinner.View() + " "
    } else {
        // Normal field rendering
        statusIndicator = ""
    }
    
    // Render: [⠋ Spinning] Label: Value (in same line)
    return fmt.Sprintf("%s%s: %s", statusIndicator, label, value)
}
```

#### update.go modifications (spinner animation)
```go
// MODIFIED: Handle spinner updates in bubbletea cycle
func (m *DevTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    switch msg := msg.(type) {
    case spinner.TickMsg:
        // Update spinners for all running async operations
        for _, tab := range m.tabSections {
            for _, field := range tab.fields {
                if field.asyncState.isRunning {
                    var cmd tea.Cmd
                    field.spinner, cmd = field.spinner.Update(msg)
                    if cmd != nil {
                        cmds = append(cmds, cmd)
                    }
                }
            }
        }
    
    // ... existing update logic
    
    case tabContent:
        // Handle async message updates
        // ... existing tabContent handling
    }
    
    return m, tea.Batch(cmds...)
}
```

#### tabContent struct (extended)
```go
// print.go - extend existing message structure
type tabContent struct {
    // UNCHANGED: Existing fields
    Id         string            // unixid timestamp-based (same field name as current)
    Content    string            // message text (same field name as current)
    Type       messagetype.Type  // Auto-detected (same field name as current)
    tabSection *tabSection       // reference to tab (same field name as current)
    
    // NEW: Async fields (always present, nil when not async)
    operationID *string     // nil for sync messages, value for async operations
    isProgress  bool        // true if this is a progress update
    isComplete  bool        // true if async operation completed
}
```

#### print.go modifications (extend existing methods)
```go
// MODIFIED: newContent method to accept optional ID for async operations
func (h *DevTUI) newContent(content string, mt messagetype.Type, tabSection *tabSection, operationID ...string) tabContent {
    var id string
    var opID *string
    
    if len(operationID) > 0 && operationID[0] != "" {
        // Use provided operation ID for async operations
        id = operationID[0]
        opID = &operationID[0]
    } else {
        // Generate new ID for regular operations (current behavior)
        if h.id != nil {
            id = h.id.GetNewID()
        } else {
            id = "temp-id"
            h.LogToFile("Warning: unixid not initialized, using fallback ID")
        }
        opID = nil // Not an async operation
    }

    return tabContent{
        Id:          id,
        Content:     content,
        Type:        mt,
        tabSection:  tabSection,
        operationID: opID,
        isProgress:  false,  // Will be set by specific async methods
        isComplete:  false,  // Will be set by specific async methods
    }
}

// MODIFIED: sendMessage to support optional operation ID
func (d *DevTUI) sendMessage(content string, mt messagetype.Type, tabSection *tabSection, operationID ...string) {
    tabSection.addNewContent(mt, content)
    newContent := d.newContent(content, mt, tabSection, operationID...)
    d.tabContentsChan <- newContent
}

// NEW: Internal async message methods (FINAL IMPLEMENTATION)
func (f *field) sendProgressMessage(content string) {
    f.parentTab.tui.sendMessage(content, messagetype.Info, f.parentTab, f.asyncState.operationID)
}

func (f *field) sendErrorMessage(content string) {
    f.parentTab.tui.sendMessage(content, messagetype.Error, f.parentTab, f.asyncState.operationID)
}

func (f *field) sendSuccessMessage(content string) {
    f.parentTab.tui.sendMessage(content, messagetype.Success, f.parentTab, f.asyncState.operationID)
}

// REMOVED: getParentTUI() method - not needed with direct parentTab reference
```

### Internal Async Flow (Transparent to User)

```go
// userKeyboard.go - When user presses Enter on field
func (f *field) handleEnter() {
    if f.handler == nil {
        return // fallback for old API during transition
    }
    
    // DevTUI handles async internally - user doesn't see this complexity
    go f.executeAsyncChange()
}

func (f *field) executeAsyncChange() {
    // Create internal context with timeout from handler
    timeout := f.handler.Timeout()
    var ctx context.Context
    var cancel context.CancelFunc
    
    if timeout > 0 {
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        ctx, cancel = context.WithCancel(context.Background())
    }
    
    f.asyncState.cancel = cancel
    f.asyncState.isRunning = true
    
    // Generate ONE operation ID for the entire async operation
    f.asyncState.operationID = f.parentTab.tui.id.GetNewID()
    
    // Initialize and start spinner (DevTUI only handles animation)
    f.spinner = spinner.New()
    f.spinner.Spinner = spinner.Dot  // Simple animated dot: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
    f.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69")) // Default spinner style
    
    // Spinner will be updated in the main bubbletea Update() cycle
    // No need to manually trigger here
    
    // Get current value based on field type
    currentValue := f.getCurrentValue()
    
    // Execute user's Change method with context monitoring
    resultChan := make(chan struct{
        result string
        err    error
    }, 1)
    
    go func() {
        result, err := f.handler.Change(currentValue)
        resultChan <- struct{
            result string
            err    error
        }{result, err}
    }()
    
    // Wait for completion or timeout
    select {
    case res := <-resultChan:
        // Operation completed normally
        f.asyncState.isRunning = false
        
        if res.err != nil {
            // Handler decides error message content
            f.sendErrorMessage(res.err.Error())
        } else {
            // Handler decides success message content
            f.sendSuccessMessage(res.result)
        }
        
    case <-ctx.Done():
        // Operation timed out
        f.asyncState.isRunning = false
        
        if ctx.Err() == context.DeadlineExceeded {
            f.sendErrorMessage(fmt.Sprintf("Operation timed out after %v", timeout))
        } else {
            f.sendErrorMessage("Operation was cancelled")
        }
    }
    
    cancel() // Clean up context
    // Spinner will automatically stop when isRunning = false
}

// NEW: Get current value for Change() method (matches existing field logic)
func (f *field) getCurrentValue() any {
    if f.handler.Editable() {
        // For editable fields, return the edited text (tempEditValue or current value)
        // This matches current field behavior with tempEditValue
        if f.tempEditValue != "" {
            return f.tempEditValue
        }
        return f.handler.Value()
    } else {
        // For non-editable fields (action buttons), return the original value
        return f.handler.Value()
    }
}
```

## Example Implementations

### Server Configuration Handler (Editable Field)
```go
// examples/server_config.go
type HostHandler struct {
    currentHost string
}

func (h *HostHandler) Label() string { return "Host" }
func (h *HostHandler) Value() string { return h.currentHost }
func (h *HostHandler) Editable() bool { return true }
func (h *HostHandler) Timeout() time.Duration { return 5 * time.Second } // 5s timeout for network validation

// SAME simple signature as current changeFunc - handles ALL messaging
func (h *HostHandler) Change(newValue any) (string, error) {
    host := strings.TrimSpace(newValue.(string))
    if host == "" {
        return "", fmt.Errorf("host cannot be empty")
    }
    if !isValidHost(host) {
        return "", fmt.Errorf("invalid host format")
    }
    
    h.currentHost = host
    return fmt.Sprintf("Host configured: %s", host), nil // Change() handles success message
}

type PortHandler struct {
    currentPort string
}

func (h *PortHandler) Label() string { return "Port" }
func (h *PortHandler) Value() string { return h.currentPort }
func (h *PortHandler) Editable() bool { return true }
func (h *PortHandler) Timeout() time.Duration { return 0 } // No timeout for simple validation

func (h *PortHandler) Change(newValue any) (string, error) {
    portStr := newValue.(string)
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return "", fmt.Errorf("port must be a number") // Change() handles error message
    }
    if port < 1 || port > 65535 {
        return "", fmt.Errorf("port must be between 1 and 65535") // Change() handles validation message
    }
    
    h.currentPort = portStr
    return fmt.Sprintf("Port set to: %d", port), nil // Change() handles success message
}
```

### Build Action Handler (Non-editable Field)
```go
// examples/build_actions.go
type BuildHandler struct {
    projectPath string
}

func (h *BuildHandler) Label() string { return "Build Project" }
func (h *BuildHandler) Value() string { return "Click to build" }
func (h *BuildHandler) Editable() bool { return false }
func (h *BuildHandler) Timeout() time.Duration { return 30 * time.Second } // 30s timeout for build

// Same simple signature - DevTUI handles async internally
func (h *BuildHandler) Change(newValue any) (string, error) {
    // User writes normal code - no async complexity
    // This will run in background automatically
    
    err := exec.Command("go", "build", h.projectPath).Run()
    if err != nil {
        return "", fmt.Errorf("build failed: %v", err) // Change() handles error message
    }
    
    return "Build completed successfully", nil // Change() handles success message
}

type DeployHandler struct {
    environment string
}

func (h *DeployHandler) Label() string { return "Deploy" }
func (h *DeployHandler) Value() string { return "Press Enter to deploy" }
func (h *DeployHandler) Editable() bool { return false }
func (h *DeployHandler) Timeout() time.Duration { return 2 * time.Minute } // 2 minutes for deployment

func (h *DeployHandler) Change(newValue any) (string, error) {
    // Long operation - will show progress automatically
    time.Sleep(5 * time.Second) // Simulates deployment
    
    return fmt.Sprintf("Deployed to %s successfully", h.environment), nil // Change() handles result message
}
```

### Usage Example (Replacing Current README Example)
```go
// cmd/main.go - Updated to use handler-based API
func main() {
    config := &devtui.TuiConfig{
        AppName:       "MyApp",
        TabIndexStart: 0,
        ExitChan:      make(chan bool),
        Color: &devtui.ColorStyle{
            Foreground: "#F4F4F4",
            Background: "#000000", 
            Highlight:  "#FF6600",
            Lowlight:   "#666666",
        },
        LogToFile: func(messages ...any) {
            fmt.Println(append([]any{"DevTUI Log:"}, messages...)...)
        },
    }

    tui := devtui.NewTUI(config)

    // Create handlers
    hostHandler := &HostHandler{currentHost: "localhost"}
    portHandler := &PortHandler{currentPort: "8080"}
    buildHandler := &BuildHandler{projectPath: "./"}
    deployHandler := &DeployHandler{environment: "production"}

    // NEW API - Clean and simple
    tui.NewTabSection("Server", "Server configuration").
        NewField(hostHandler).      // Instead of: NewField("Host", "localhost", true, func...)
        NewField(portHandler)       // Instead of: NewField("Port", "8080", true, func...)
        
    tui.NewTabSection("Actions", "Available operations").
        NewField(buildHandler).     // Instead of: NewField("Build Project", "Click to build", false, func...)
        NewField(deployHandler)     // Instead of: NewField("Deploy", "Click to deploy", false, func...)

    var wg sync.WaitGroup
    wg.Add(1)
    go tui.Start(&wg)
    wg.Wait()
}
```

## Implementation Plan

### Phase 1: Core Infrastructure
1. **field.go**: 
   - Add FieldHandler interface
   - Add internal async state management
   - Modify field struct to use handler
   - Add `getCurrentValue()` method for editable vs non-editable fields
   - Add spinner initialization and async state management

2. **tabSection.go**:
   - Update NewField method: `NewField(handler FieldHandler) *tabSection`
   - Remove old parameter-based method
   - Ensure fields have reference to parent tabSection for message routing

3. **print.go**:
   - Extend tabContent with async fields
   - Add internal progress message methods
   - Modify newContent and sendMessage to support operation IDs

4. **userKeyboard.go**:
   - Update Enter key handling to trigger async operations
   - Add internal goroutine management
   - Handle field focus and editing state

5. **init.go or devTuiDefault.go**:
   - Add spinner styles to DevTUI struct
   - Initialize spinner-related fields in TUI configuration

### Phase 2: Message Integration  
6. **update.go**:
   - Integrate async message updates into bubbletea cycle
   - Handle spinner.TickMsg for all active async operations
   - Update spinner animations for running operations

7. **view.go**:
   - Add progress indicators for running operations
   - Show spinner in field rendering when asyncState.isRunning = true
   - Ensure proper field layout with spinner integration

### Phase 3: Examples and Testing
8. **examples/** directory:
   - Create handler examples (build, deploy, config, etc.)
   - Show various use cases with different timeouts
   - Demonstrate editable vs non-editable field patterns

9. **cmd/main.go**:
   - Update to use new handler-based API
   - Demonstrate async capabilities
   - Show timeout configuration examples

10. **Testing**:
    - Add tests for handler interface
    - Test async operation flow
    - Test timeout behavior
    - Performance testing

## Critical Implementation Notes

### Field Value Handling
```go
// For editable fields: user input takes precedence
func (f *field) getCurrentValue() any {
    if f.handler.Editable() {
        if f.tempEditValue != "" {
            return f.tempEditValue  // User's current input
        }
        return f.handler.Value()    // Handler's current value
    } else {
        return f.handler.Value()    // Action buttons use handler value
    }
}
```

### Spinner Integration with BubbleTea
```go
// In DevTUI Update() method - handle spinner updates
case spinner.TickMsg:
    var cmds []tea.Cmd
    for _, tab := range m.tabSections {
        for _, field := range tab.fields {
            if field.asyncState != nil && field.asyncState.isRunning {
                var cmd tea.Cmd
                field.spinner, cmd = field.spinner.Update(msg)
                if cmd != nil {
                    cmds = append(cmds, cmd)
                }
            }
        }
    }
    return m, tea.Batch(cmds...)
```

### Message Routing
```go
// Fields need reference to parent structures for message sending
field.tabSection.tui.sendMessage(content, msgType, field.tabSection, operationID)
```

## Migration Strategy

### Current Field Structure Analysis
```go
// CURRENT field struct (field.go)
type field struct {
    name       string                                             
    value      string                                             
    editable   bool                                               
    changeFunc func(newValue any) (execMessage string, err error) 
    // internal use
    tempEditValue string
    index         int
    cursor        int
}
```

### Migration Steps

#### Step 1: Add Handler Support (Backward Compatible)
```go
// MODIFIED field struct - supports both old and new API
type field struct {
    // OLD API (keep for backward compatibility during transition)
    name       string                                             
    value      string                                             
    editable   bool                                               
    changeFunc func(newValue any) (execMessage string, err error) 
    
    // NEW API
    handler    FieldHandler        // NEW: Optional, nil if using old API
    asyncState *internalAsyncState // NEW: Only used for handler-based fields
    spinner    spinner.Model       // NEW: Only used for handler-based fields
    
    // UNCHANGED: existing internal fields
    tempEditValue string
    index         int
    cursor        int
}
```

#### Step 2: Update NewField Method
```go
// CURRENT NewField (tabSection.go)
func (ts *tabSection) NewField(name, value string, editable bool, changeFunc func(newValue any) (string, error)) *tabSection

// NEW NewField - add handler-based version
func (ts *tabSection) NewFieldHandler(handler FieldHandler) *tabSection {
    f := &field{
        handler:    handler,
        asyncState: &internalAsyncState{},
        spinner:    spinner.New(),
    }
    ts.addFields(f)
    return ts
}

// Keep old NewField for backward compatibility during transition
```

#### Step 3: Add Detection Logic
```go
// Helper methods to detect which API is being used
func (f *field) isHandlerBased() bool {
    return f.handler != nil
}

func (f *field) getLabel() string {
    if f.isHandlerBased() {
        return f.handler.Label()
    }
    return f.name // fallback to old API
}

func (f *field) getValue() string {
    if f.isHandlerBased() {
        return f.handler.Value()
    }
    return f.value // fallback to old API
}
```

### Critical Implementation Issues - RESOLVED

#### Issue 1: Field-to-TabSection Reference ✅ RESOLVED
**Problem**: Fields need to send messages through their parent tabSection, but don't have direct reference.

**SOLUTION CHOSEN**: Add `parentTab *tabSection` to field struct
- **Why**: Simplest and most direct approach
- **Implementation**: Add field reference during field creation in `addFields()`

#### Issue 2: Existing Method Names ✅ RESOLVED  
**Problem**: Current field.go methods conflict with handler interface:
- `field.Value()` exists, but handler also has `Value()`
- `field.Editable()` exists, but handler also has `Editable()`

**SOLUTION CHOSEN**: Clean replacement - remove old methods, use handler methods
- **Why**: No external users to break, cleaner design
- **Implementation**: Replace old methods entirely with handler-based approach

#### Issue 3: Migration Strategy ✅ RESOLVED
**DECISION**: No Backward Compatibility - Clean Implementation
- **Remove old API completely** ✅
- **Clean implementation without legacy baggage** ✅  
- **Focus on best design, not compatibility** ✅

### Final Implementation Decisions

#### Field Structure (FINAL)
```go
type field struct {
    // NEW: Handler-based approach (replaces name, value, editable, changeFunc)
    handler    FieldHandler        
    parentTab  *tabSection         // NEW: Direct reference to parent for message routing
    
    // NEW: Internal async state
    asyncState *internalAsyncState 
    spinner    spinner.Model       
    
    // UNCHANGED: Existing internal fields
    tempEditValue string
    index         int
    cursor        int
}
```

#### API Changes (FINAL)
```go
// OLD API - COMPLETELY REMOVED
// func (ts *tabSection) NewField(name, value string, editable bool, changeFunc func(newValue any) (string, error)) *tabSection

// NEW API - ONLY THIS EXISTS
func (ts *tabSection) NewField(handler FieldHandler) *tabSection
```

## Benefits of This Approach

1. **Transparent Async**: Users write simple synchronous Change methods
2. **Internal Complexity**: DevTUI handles all goroutines, contexts, channels
3. **Clean Interface**: Simple handler interface, no public async types
4. **Progress Feedback**: Automatic progress messages during operations
5. **Error Handling**: Standard Go error handling, async execution
6. **Maintainable**: Single responsibility - handlers handle logic, DevTUI handles execution

Ready to implement Phase 1?
