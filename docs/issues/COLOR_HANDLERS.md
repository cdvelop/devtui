# Handler Color Identification Enhancement

## Problem Statement

Currently, DevTUI applies a fixed color scheme to all handler names in message formatting through the `formatHandlerName()` function. This implementation uses the `infoStyle` (defined with `palette.Info` color) uniformly for all handler types, making it difficult for users to visually distinguish between different types of handlers in the message stream.

### Current Implementation Issues

1. **No Visual Differentiation**: All handler names `[HandlerName]` use the same color (`palette.Info` - typically cyan `#00FFFF`)
2. **Poor UX**: Users cannot quickly identify the source type of messages (Edit vs Execution vs Logger vs Interactive)  
3. **Missed Color Palette Potential**: The existing `ColorPalette` struct provides semantic colors (`Success`, `Warning`, `Error`, `Primary`, `Secondary`) that could enhance handler identification
4. **Scalability Concerns**: As projects grow with more handlers, visual identification becomes increasingly important

### Current Color Application

```go
func (t *DevTUI) formatHandlerName(handlerName string) string {
    if handlerName == "" {
        return ""
    }
    // Fixed style applied to ALL handler types
    styledName := t.infoStyle.Render(Fmt("[%s]", handlerName))
    return styledName + " "
}
```

## Proposed Solution

### Automatic Type-Based Color Mapping (Breaking Change)

**Description**: Automatically assign colors based on handler type using predefined semantic mapping. This requires breaking API changes but provides consistent visual identification.

**Implementation**:
```go
// Add Color field to registration - BREAKING CHANGE
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration, color string) *tabSection
func (ts *tabSection) AddExecutionHandler(handler HandlerExecution, timeout time.Duration, color string) *tabSection
func (ts *tabSection) AddInteractiveHandler(handler HandlerInteractive, timeout time.Duration, color string) *tabSection
func (ts *tabSection) NewLogger(name string, enableTracking bool, color string) func(message ...any)

// Predefined semantic mapping
var DefaultHandlerColors = map[handlerType]string{
    handlerTypeEdit:        palette.Primary,   // Edit operations - main color
    handlerTypeExecution:   palette.Success,   // Actions/execution - green
    handlerTypeInteractive: palette.Info,      // Interactive content - cyan  
    handlerTypeWriter:      palette.Secondary, // Logging - muted
    handlerTypeDisplay:     palette.Muted,     // Read-only info - subtle
}

func (t *DevTUI) formatHandlerName(handlerName string, handlerType handlerType, customColor string) string {
    if handlerName == "" {
        return ""
    }
    
    color := customColor
    if color == "" {
        color = DefaultHandlerColors[handlerType] // Use semantic default
    }
    
    style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))
    styledName := style.Render(Fmt("[%s]", handlerName))
    return styledName + " "
}
```

**Pros**:
- ✅ Consistent visual identity - same handler types always use same colors by default
- ✅ Semantic meaning - colors have logical association with handler purpose
- ✅ Clean API - single method per handler type with required color parameter
- ✅ Override capability - custom colors still possible when needed
- ✅ Better UX - immediate visual differentiation of handler types

**Cons**:
- ❌ Breaking change - all existing registration calls must be updated
- ❌ Forced decision - developers must specify color even if they don't care
- ❌ Migration effort - existing codebases need updating
- ❌ Potentially opinionated - default color assignments might not suit all use cases

## Implementation Details

### CRITICAL IMPLEMENTATION NOTES

⚠️ **IMPORTANT**: The implementation requires multiple coordinated changes across several files. Follow this exact sequence to avoid compilation errors:

1. **File: `field.go`** - Add `handlerColor` field to `anyHandler` struct FIRST
2. **File: `print.go`** - Update `formatHandlerName` signature and `tabContent` struct
3. **File: `handlerRegistration.go`** - Update all registration methods signatures
4. **File: All factory functions** - Update `newEditHandler`, `newExecutionHandler`, etc.

### Required Struct Changes

#### 1. anyHandler Struct Modification (field.go)
```go
type anyHandler struct {
    handlerType handlerType
    timeout     time.Duration
    lastOpID    string
    mu          sync.RWMutex
    origHandler interface{}
    handlerColor string // NEW: Store handler-specific color - ADD THIS FIELD
    
    // ... rest of existing fields
}
```

#### 2. tabContent Struct Modification (print.go)
```go
type tabContent struct {
    Id          string
    Timestamp   string
    Content     string
    Type        MessageType
    tabSection  *tabSection
    operationID *string
    isProgress  bool
    isComplete  bool
    handlerName string
    handlerColor string // NEW: Add this field to pass color to formatMessage
}
```

### API Changes Required

All existing registration methods will require a color parameter:

```go
// BEFORE (current API)
config.AddEditHandler(&DatabaseHandler{}, 2*time.Second)
config.AddExecutionHandler(&BackupHandler{}, 5*time.Second) 
chat.AddInteractiveHandler(chatHandler, 3*time.Second)
logs.NewLogger("SystemLogWriter", false)

// AFTER (new API with required color parameter)
config.AddEditHandler(&DatabaseHandler{}, 2*time.Second, "#00ADD8")
config.AddExecutionHandler(&BackupHandler{}, 5*time.Second, "#00AA00")
chat.AddInteractiveHandler(chatHandler, 3*time.Second, "#0088FF")
logs.NewLogger("SystemLogWriter", false, "#666666")

// AFTER (using default Primary color with empty string)
config.AddEditHandler(&StandardHandler{}, 2*time.Second, "")
config.AddExecutionHandler(&DefaultAction{}, 5*time.Second, "")
logs.NewLogger("DefaultLogger", false, "")
```

### Color Parameter Usage

The color parameter is required but can be an empty string to use the default `palette.Primary` color:

```go
// Using custom colors (explicit color specification)
config.AddEditHandler(&DatabaseHandler{}, 2*time.Second, "#00ADD8")
config.AddExecutionHandler(&BackupHandler{}, 5*time.Second, "#00AA00")
config.AddInteractiveHandler(chatHandler, 3*time.Second, "#0088FF")
logs.NewLogger("SystemLogWriter", false, "#666666")

// Using default Primary color (empty string)
config.AddEditHandler(&StandardHandler{}, 2*time.Second, "")
config.AddExecutionHandler(&DefaultAction{}, 5*time.Second, "")
logs.NewLogger("DefaultLogger", false, "")
```

### Internal Implementation Changes

#### CRITICAL: Update formatHandlerName Function (print.go)
```go
// CURRENT SIGNATURE (must be changed)
func (t *DevTUI) formatHandlerName(handlerName string) string

// NEW SIGNATURE (required)
func (t *DevTUI) formatHandlerName(handlerName string, handlerColor string) string {
    if handlerName == "" {
        return ""
    }
    
    // Use Primary color if no specific color provided
    color := handlerColor
    if color == "" {
        color = t.Primary // Use palette.Primary as default
    }
    
    // SIMPLE COLOR VALIDATION: Check for minimum hex format (#RRGGBB = 7 chars)
    if color != "" && (len(color) != 7 || color[0] != '#') {
        // Invalid color format - use Primary and log warning
        if t.Logger != nil {
            t.Logger("warning: invalid color format '", color, "' for handler '", handlerName, "' - using Primary color instead")
        }
        color = t.Primary
    }
    
    // Create style with handler-specific color as background
    style := lipgloss.NewStyle().
        Bold(true).
        Background(lipgloss.Color(color)).
        Foreground(lipgloss.Color(t.Foreground)) // Use foreground for text contrast
    
    styledName := style.Render(Fmt("[%s]", handlerName))
    return styledName + " "
}
```

#### CRITICAL: Update formatMessage Function Call (print.go)
```go
// CURRENT CALL (must be changed)
handlerName := t.formatHandlerName(msg.handlerName)

// NEW CALL (required)
handlerName := t.formatHandlerName(msg.handlerName, msg.handlerColor)
```

#### CRITICAL: Update createTabContent Function (print.go)
```go
// Must add handlerColor parameter and set field
func (h *DevTUI) createTabContent(content string, mt MessageType, tabSection *tabSection, handlerName string, operationID string, handlerColor string) tabContent {
    // ... existing timestamp and ID logic
    
    return tabContent{
        Id:          id,
        Timestamp:   timestamp,
        Content:     content,
        Type:        mt,
        tabSection:  tabSection,
        operationID: opID,
        isProgress:  false,
        isComplete:  false,
        handlerName: handlerName,
        handlerColor: handlerColor, // NEW: Set the color field
    }
}
```

#### Registration Method Updates - All Files in handlerRegistration.go
```go
// CURRENT SIGNATURE (must be changed)
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration) *tabSection

// NEW SIGNATURE (required)
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration, color string) *tabSection {
    var tracker MessageTracker
    if t, ok := handler.(MessageTracker); ok {
        tracker = t
    }

    anyH := newEditHandler(handler, timeout, tracker)
    anyH.handlerColor = color // NEW: Store the color
    
    f := &field{
        handler:    anyH,
        parentTab:  ts,
        asyncState: &internalAsyncState{},
    }
    ts.addFields(f)
    return ts
}

// Apply same pattern to:
// - AddExecutionHandler
// - AddInteractiveHandler  
// - AddDisplayHandler
// - NewLogger
```

#### Factory Function Updates (field.go)
```go
// All factory functions must be updated to accept and store color:
func newEditHandler(h HandlerEdit, timeout time.Duration, tracker MessageTracker, color string) *anyHandler
func newExecutionHandler(h HandlerExecution, timeout time.Duration, color string) *anyHandler
func newDisplayHandler(h HandlerDisplay, color string) *anyHandler
func newInteractiveHandler(h HandlerInteractive, timeout time.Duration, tracker MessageTracker, color string) *anyHandler
func newWriterHandler(h HandlerLogger, color string) *anyHandler
func newTrackerWriterHandler(h interface{...}, color string) *anyHandler
```

### Message Creation Update Required

All calls to `sendMessageWithHandler` and `createTabContent` must be updated to include the handler color:

```go
// Find all calls like this and add color parameter:
d.sendMessageWithHandler(content, mt, tabSection, handlerName, operationID, handlerColor)
```

### Color Validation Implementation

#### Simple Hex Color Validation
DevTUI implements minimal color validation to ensure basic hex format compliance:

```go
// Color validation logic in formatHandlerName
func validateColor(color string, handlerName string, logger func(...any)) string {
    // Empty string is valid (uses Primary color)
    if color == "" {
        return ""
    }
    
    // Simple hex validation: must be exactly 7 characters (#RRGGBB)
    if len(color) != 7 || color[0] != '#' {
        // Log warning and return empty to use Primary color
        if logger != nil {
            logger("warning: invalid color format '", color, "' for handler '", handlerName, "' - using Primary color instead")
        }
        return "" // Return empty to trigger Primary color usage
    }
    
    return color // Valid format
}
```

#### Validation Rules
- ✅ **Valid**: `""` (empty string, uses Primary color)
- ✅ **Valid**: `"#FF0000"` (7 characters, starts with #)
- ✅ **Valid**: `"#00ADD8"` (proper hex format)
- ❌ **Invalid**: `"red"` (named colors not supported)
- ❌ **Invalid**: `"FF0000"` (missing # prefix)
- ❌ **Invalid**: `"#FF00"` (too short)
- ❌ **Invalid**: `"#FF0000AA"` (too long)

#### Error Handling Behavior
1. **Invalid color detected** → Log warning via Logger
2. **Fallback to Primary color** → Use `palette.Primary` as background
3. **Continue execution** → No crashes or panics
4. **Warning format**: `"warning: invalid color format 'BADCOLOR' for handler 'HandlerName' - using Primary color instead"`

### Message Creation Update Required


## INCONSISTENCIAS DETECTADAS - REQUIEREN DECISIÓN

### 1. **Conflicto en Implementación Propuesta Inicial**
El documento original mostraba dos versiones diferentes de `formatHandlerName`:

**Versión A (en sección inicial):**
```go
func (t *DevTUI) formatHandlerName(handlerName string, handlerType handlerType, customColor string) string
```

**Versión B (en implementación):**
```go
func (t *DevTUI) formatHandlerName(handlerName string, handlerColor string) string
```

❓ **DECISIÓN REQUERIDA**: ¿Usar `handlerType` + `customColor` o solo `handlerColor`?

### 2. **Propagación del Color a través del Sistema**
Actualmente el color del handler NO se propaga hasta `formatMessage`. Se necesita:

1. Almacenar color en `anyHandler.handlerColor`
2. Pasar color a `tabContent.handlerColor` 
3. Usar color en `formatHandlerName(msg.handlerName, msg.handlerColor)`

❓ **DECISIÓN REQUERIDA**: ¿Confirmar esta propagación de color?

### 3. **AddDisplayHandler - Manejo Especial**
`AddDisplayHandler` no tiene timeout pero necesita color:

**Actual:**
```go
func (ts *tabSection) AddDisplayHandler(handler HandlerDisplay) *tabSection
```

**Propuesto:**
```go
func (ts *tabSection) AddDisplayHandler(handler HandlerDisplay, color string) *tabSection
```

❓ **DECISIÓN REQUERIDA**: ¿Confirmar que Display handlers también necesitan color?

### 4. **Manejo de Handlers ReadOnly e Interactive**
El código actual tiene lógica especial:

```go
// For readonly fields: no timestamp, cleaner visual content, no special coloring
if msg.handlerName != "" && t.isReadOnlyHandler(msg.handlerName) {
    return msg.Content // NO formatting applied
}

// Interactive handlers: timestamp + content (no handler name for cleaner UX)
if msg.handlerName != "" && t.isInteractiveHandler(msg.handlerName) {
    return Fmt("%s %s", timeStr, styledContent) // NO handler name
}
```

❓ **DECISIÓN REQUERIDA**: ¿Mantener estas excepciones o aplicar color a todos los handlers?

## Migration Strategy

### Phase 1: Add Color Parameter (Breaking Change)
- Update all registration method signatures to require color parameter (can be empty string for default)
- Update internal handler creation to store and use colors
- Use `palette.Primary` as background when no color specified

### Phase 2: Update Documentation and Examples
- Update all documentation with new API signatures
- Provide migration guide showing empty string usage for default color
- Update example code to demonstrate both custom colors and default usage


## Questions for Review

1. **Color Usage**: Do you prefer the approach where empty string defaults to `palette.Primary` background, or should we have a more explicit default mechanism?


3. **Color Validation**: Simple hex validation implemented - requires minimum 7 characters (#RRGGBB format). Invalid colors fallback to Primary with Logger warning.

4. **Background vs Foreground**: Should handler colors be applied as background (with foreground text for contrast) or as foreground color?

5. **Handler Identification**: Are there other visual indicators beyond color that could improve handler identification (icons, prefixes, etc.)?

6. **Error Handling**: Invalid color formats use Primary color fallback with Logger warning for notification.

## Benefits of This Approach

- **Clear Visual Distinction**: Each handler can have its own distinct background color in message streams
- **Flexible Color Usage**: Developers can specify custom colors or use the default Primary color
- **Clean API**: Single registration method per handler type with simple color parameter (empty string for default)
- **Consistent Default**: All handlers without custom colors use the same Primary background for consistency
- **Better UX**: Users can visually distinguish between different handlers, especially when custom colors are used for different purposes

## FILES THAT REQUIRE CHANGES (Implementation Checklist)

### 1. `/field.go` - Core Structure Changes
- [ ] Add `handlerColor string` to `anyHandler` struct
- [ ] Update all factory functions: `newEditHandler`, `newExecutionHandler`, etc.
- [ ] Add color parameter to all factory function signatures

### 2. `/print.go` - Message Formatting Changes  
- [ ] Add `handlerColor string` to `tabContent` struct
- [ ] Update `formatHandlerName` function signature
- [ ] Update `formatMessage` function call to `formatHandlerName`
- [ ] Update `createTabContent` function to accept and set color
- [ ] Update `sendMessageWithHandler` function to pass color

### 3. `/handlerRegistration.go` - API Changes
- [ ] Update `AddEditHandler` signature and implementation
- [ ] Update `AddExecutionHandler` signature and implementation  
- [ ] Update `AddInteractiveHandler` signature and implementation
- [ ] Update `AddDisplayHandler` signature and implementation
- [ ] Update `NewLogger` signature and implementation

### 4. `/example/` directory - Update Examples
- [ ] Update all example usage in demo files
- [ ] Update documentation examples
- [ ] Add color parameter to all registration calls

### 5. Test Files - Update Tests
- [ ] Update all test files that use registration methods
- [ ] Add tests for color functionality
- [ ] Test color fallback to Primary behavior
