# DevTUI Shortcut Keys Implementation Plan

## Executive Summary

This document outlines the comprehensive implementation plan for adding global shortcut key functionality to DevTUI. The feature will allow `HandlerEdit` implementations to optionally provide shortcuts that can be triggered from any tab, automatically navigating to the handler's tab and invoking its `Change()` method with the shortcut value.

## Current Architecture Analysis

### Existing Keyboard Handling Flow
1. **Entry Point**: `handleKeyboard()` in `userKeyboard.go` - main dispatcher
2. **Mode Detection**: Splits into `handleEditingConfigKeyboard()` vs `handleNormalModeKeyboard()`
3. **Normal Mode**: Handles navigation (Tab/Shift+Tab, Left/Right, Up/Down, PageUp/PageDown)
4. **Edit Mode**: Handles text input, cursor movement, Enter/Esc
5. **Global Keys**: Ctrl+C for exit

### Current Handler Registration Flow
1. `AddEditHandler()` in `handlerRegistration.go` creates `anyHandler` wrapper
2. Creates `field` struct with handler and adds to `tabSection.fieldHandlers`
3. Handler interfaces are well-defined in `interfaces.go`

### DevTUI Core Structure
```go
type DevTUI struct {
    tabSections       []*tabSection
    activeTab         int
    editModeActivated bool
    // ... other fields
}

type tabSection struct {
    fieldHandlers []*field
    indexActiveEditField int
    // ... other fields
}
```

## Proposed Implementation

### 1. New Interface Definition

**File**: `interfaces.go`
```go
// HandlerEditWithShortcuts extends HandlerEdit with shortcut key support
type HandlerEditWithShortcuts interface {
    HandlerEdit
    Shortcuts() []string // Returns shortcut keys (e.g., ["c", "d", "p"])
}
```

**Alternative Approach**: Optional interface detection
```go
// Optional interface that HandlerEdit can implement
type ShortcutProvider interface {
    Shortcuts() map[string]string // Returns shortcut keys with descriptions (e.g., {"c": "coding mode", "d": "debug mode", "p": "production mode"})
}
```

### 2. Shortcut Registry Structure

**File**: `shortcuts_registry.go` (new file)
```go
package devtui

import "sync"

// ShortcutEntry represents a registered shortcut
type ShortcutEntry struct {
    Key         string        // The shortcut key (e.g., "c", "d", "p")
    Description string        // Human-readable description (e.g., "coding mode", "debug mode")
    TabIndex    int          // Index of the tab containing the handler
    FieldIndex  int          // Index of the field within the tab
    HandlerName string       // Handler name for identification
    Value       string       // Value to pass to Change()
}

// ShortcutRegistry manages global shortcut keys
type ShortcutRegistry struct {
    mu        sync.RWMutex
    shortcuts map[string]*ShortcutEntry // key -> entry
}

func newShortcutRegistry() *ShortcutRegistry {
    return &ShortcutRegistry{
        shortcuts: make(map[string]*ShortcutEntry),
    }
}

func (sr *ShortcutRegistry) Register(key string, entry *ShortcutEntry) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    sr.shortcuts[key] = entry
}

func (sr *ShortcutRegistry) Get(key string) (*ShortcutEntry, bool) {
    sr.mu.RLock()
    defer sr.mu.RUnlock()
    entry, exists := sr.shortcuts[key]
    return entry, exists
}

func (sr *ShortcutRegistry) Unregister(key string) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    delete(sr.shortcuts, key)
}

func (sr *ShortcutRegistry) List() []string {
    sr.mu.RLock()
    defer sr.mu.RUnlock()
    keys := make([]string, 0, len(sr.shortcuts))
    for k := range sr.shortcuts {
        keys = append(keys, k)
    }
    return keys
}
```

### 3. DevTUI Integration

**File**: `init.go` - Add shortcut registry to DevTUI struct
```go
type DevTUI struct {
    // ... existing fields
    shortcutRegistry *ShortcutRegistry
}

func NewTUI(c *TuiConfig) *DevTUI {
    tui := &DevTUI{
        // ... existing initialization
        shortcutRegistry: newShortcutRegistry(),
    }
    // ... rest of initialization
    return tui
}
```

### 4. Handler Registration Enhancement

**File**: `handlerRegistration.go`
```go
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration) *tabSection {
    // ... existing code to create anyHandler and field
    
    // NEW: Check for shortcut support
    ts.registerShortcutsIfSupported(handler, len(ts.fieldHandlers)-1)
    
    return ts
}

// NEW: Helper method to register shortcuts
func (ts *tabSection) registerShortcutsIfSupported(handler HandlerEdit, fieldIndex int) {
    // Check if handler implements shortcut interface
    if shortcutProvider, hasShortcuts := handler.(ShortcutProvider); hasShortcuts {
        shortcuts := shortcutProvider.Shortcuts()
        for key, description := range shortcuts {
            entry := &ShortcutEntry{
                Key:         key,
                Description: description,
                TabIndex:    ts.index,
                FieldIndex:  fieldIndex,
                HandlerName: handler.Name(),
                Value:       key, // Use the key as the value by default
            }
            ts.tui.shortcutRegistry.Register(key, entry)
        }
    }
}
```

### 5. Keyboard Handling Enhancement

**File**: `userKeyboard.go` - Modify `handleNormalModeKeyboard()`
```go
func (h *DevTUI) handleNormalModeKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
    // ... existing switch cases for navigation

    // NEW: Handle single character shortcuts
    case tea.KeyRunes:
        if len(msg.Runes) == 1 {
            key := string(msg.Runes[0])
            if entry, exists := h.shortcutRegistry.Get(key); exists {
                return h.executeShortcut(entry)
            }
        }

    // ... rest of existing cases
}

// NEW: Execute shortcut action
func (h *DevTUI) executeShortcut(entry *ShortcutEntry) (bool, tea.Cmd) {
    // Validate indexes are still valid
    if entry.TabIndex >= len(h.tabSections) {
        if h.Logger != nil {
            h.Logger("Shortcut error: invalid tab index", entry.TabIndex)
        }
        return true, nil
    }

    targetTab := h.tabSections[entry.TabIndex]
    if entry.FieldIndex >= len(targetTab.fieldHandlers) {
        if h.Logger != nil {
            h.Logger("Shortcut error: invalid field index", entry.FieldIndex)
        }
        return true, nil
    }

    targetField := targetTab.fieldHandlers[entry.FieldIndex]
    
    // Navigate to target tab if not already there
    if h.activeTab != entry.TabIndex {
        h.activeTab = entry.TabIndex
    }
    
    // Set active field
    targetTab.indexActiveEditField = entry.FieldIndex
    
    // Execute the Change method with shortcut value
    if targetField.handler != nil {
        progress := func(msgs ...any) {
            // REUTILIZA el método unificado propuesto en field.sendMessage()
            // En lugar de crear nuevos métodos duplicados
            targetField.sendMessage(msgs...)
        }
        targetField.handler.Change(entry.Value, progress)
    }
    
    // Update viewport to show changes
    h.updateViewport()
    
    return false, nil // Stop further processing
}

// NOTA ARQUITECTURAL: Método de mensajería consolidado
// Este método reemplazará los 3 métodos actuales duplicados en field.go:
// - sendProgressMessage()
// - sendErrorMessage() 
// - sendSuccessMessage()
//
// La implementación del shortcut system NO creará nuevos métodos duplicados,
// sino que aprovechará la refactorización propuesta de field.sendMessage()
```

## Implementation Questions & Decisions

### Q1: Refactorización de Métodos de Mensajería (CRÍTICO)
**Problema**: Actualmente `field.go` tiene 3 métodos casi idénticos:
- `sendProgressMessage(content string)`
- `sendErrorMessage(content string)` 
- `sendSuccessMessage(content string)`

**Solución Propuesta**: Consolidar en un método unificado:
```go
// field.go - Método unificado que reemplaza los 3 existentes
func (f *field) sendMessage(msgs ...any) {
    if f.parentTab == nil || f.parentTab.tui == nil || len(msgs) == 0 {
        return
    }

    // Manejo inteligente de operationID
    var operationID string
    if f.asyncState != nil && f.asyncState.operationID != "" {
        operationID = f.asyncState.operationID
    }

    // Handler name extraction
    handlerName := ""
    if f.handler != nil {
        handlerName = f.handler.Name()
    }

    // Content() method optimization
    if f.hasContentMethod() {
        f.parentTab.tui.updateViewport()
        return
    }

    // Unified message processing
    message := tinystring.Translate(msgs...).String()
    msgType := messagetype.DetectMessageType(message)
    f.parentTab.tui.sendMessageWithHandler(message, msgType, f.parentTab, handlerName, operationID)
}
```

**Beneficios**:
- Elimina duplicación de código
- Consistencia en el manejo de mensajes
- Facilita el mantenimiento futuro
- Los shortcuts reutilizan infraestructura existente

### Q2: Should shortcuts be case-sensitive?
**Recommendation**: Yes, to maximize available shortcuts (both 'c' and 'C' could be different shortcuts).

### Q2: How should shortcut conflicts be handled?
**Options**:
1. **First-registered wins** (simple, predictable)

**Recommendation**: First-registered wins with warning log for conflicts.

### Q3: Should shortcuts work in edit mode?
**Recommendation**: No, only in normal mode to avoid conflicts with text input.

### Q4: Should the shortcut value be customizable per key?
**Current Design**: Uses the key itself as the value (e.g., "c" key sends "c" to Change())
**Recommendation**: Start with key-as-value for simplicity, extend later if needed.

### Q5: Should shortcuts be displayed in the UI?
**Options**:
1. Show in existing shortcuts handler
2. Show in handler labels (e.g., "Compiler (c)")
3. Both

**Decision**: Extended existing shortcuts handler to list registered shortcuts with descriptions.

**Interface Change**: Updated `ShortcutProvider` interface from `Shortcuts() []string` to `Shortcuts() map[string]string` where:
- **Key**: The shortcut character (e.g., "t", "b")
- **Value**: Human-readable description (e.g., "test connection", "backup database")

This allows the shortcuts to be displayed in the UI with meaningful descriptions.

**Implementation**: Modify `shortcuts.go` to display registered shortcuts:
```go
// NEW: Add to shortcutsInteractiveHandler
func (h *shortcutsInteractiveHandler) generateHelpContent() string {
    content := Translate(h.appName, D.Shortcuts, D.Keyboard, `:

`, D.Content, D.Tab, `:
  • Tab/Shift+Tab  -`, D.Switch, D.Content, `

`, D.Fields, `:
  • `, D.Arrow, D.Left, `/`, D.Right, `     -`, D.Switch, D.Field, `
  • Enter          				-`, D.Edit, `/`, D.Execute, `
  • Esc            				-`, D.Cancel, `

`, D.Edit, D.Text, `:
  • `, D.Arrow, D.Left, `/`, D.Right, `   -`, D.Move, `cursor
  • Backspace      			-`, D.Create, D.Space, `

Viewport:
  • `, D.Arrow, D.Up, "/", D.Down, `    - Scroll`, D.Line, D.Text, `
  • PgUp/PgDown    		- Scroll`, D.Page, `
  • Mouse Wheel    		- Scroll`, D.Page, `

Scroll `, D.Status, D.Icons, `:
  •  ■  - `, D.All, D.Content, D.Visible, `
  •  ▼  - `, D.Can, `scroll`, D.Down, `
  •  ▲  - `, D.Can, `scroll`, D.Up, `
  • ▼ ▲ - `, D.Can, `scroll`, D.Down, `/`, D.Up, `

`, D.Quit, `:
  • Ctrl+C         - `, D.Quit, `
`).String()

    // NEW: Add registered shortcuts section
    if h.tui != nil && h.tui.shortcutRegistry != nil {
        shortcuts := h.getRegisteredShortcuts()
        if len(shortcuts) > 0 {
            content += "\n\nRegistered Shortcuts:\n"
            for key, description := range shortcuts {
                content += fmt.Sprintf("  • %s - %s\n", key, description)
            }
        }
    }

    content += "\n" + Translate(D.Language, D.Supported, `: en, es, zh, hi, ar, pt, fr, de, ru`).String()
    return content
}

// NEW: Helper to get registered shortcuts with descriptions
func (h *shortcutsInteractiveHandler) getRegisteredShortcuts() map[string]string {
    shortcuts := make(map[string]string)
    if h.tui != nil && h.tui.shortcutRegistry != nil {
        for key, entry := range h.tui.shortcutRegistry.shortcuts {
            shortcuts[key] = entry.Description
        }
    }
    return shortcuts
}
```

## Testing Strategy

### Unit Tests Required

1. **ShortcutRegistry Tests** (`shortcuts_registry_test.go`)
   - Registration/unregistration
   - Concurrent access
   - Conflict handling

2. **Integration Tests** (`shortcuts_integration_test.go`)
   - End-to-end shortcut execution
   - Tab navigation
   - Handler method invocation

3. **Keyboard Tests** (`shortcuts_keyboard_test.go`)
   - Shortcut key detection
   - Normal vs edit mode behavior
   - Invalid shortcut handling

> **Note:** For integration and keyboard tests, you can use [`real_user_scenario_test.go`](../real_user_scenario_test.go) as a reference. This file simulates real user scenarios and verifies the complete flow of editing and confirming values in handlers. Adapt the structure of these tests to validate shortcut behavior, including automatic navigation to the corresponding tab and invocation of the `Change()` method with the shortcut value.


### Test Handler Example
```go
type TestShortcutHandler struct {
    label string
    value string
    lastValue string
}

func (h *TestShortcutHandler) Shortcuts() map[string]string {
    return map[string]string{
        "c": "coding mode",
        "d": "debug mode", 
        "p": "production mode",
    }
}

func (h *TestShortcutHandler) Change(newValue string, progress func(msgs ...any)) {
    h.lastValue = newValue
    if progress != nil {
        progress("Mode changed to:", newValue)
    }
}
```

## Files to Create/Modify

### New Files
1. `shortcuts_registry.go` - Shortcut registry implementation
2. `shortcuts_registry_test.go` - Registry unit tests
3. `shortcuts_integration_test.go` - Integration tests
4. `shortcuts_keyboard_test.go` - Keyboard handling tests

### Modified Files
1. `interfaces.go` - Add `ShortcutProvider` interface with `map[string]string` signature
2. `init.go` - Add registry to DevTUI struct
3. `handlerRegistration.go` - Add shortcut registration logic with description support
4. `userKeyboard.go` - Add shortcut handling in normal mode
5. `shortcuts.go` - Extend `generateHelpContent()` to show registered shortcuts with descriptions
6. `example/HandlerEdit.go` - Add shortcut support to `DatabaseHandler` for testing
7. `README.md` - Document shortcut feature

## Implementation Phases

### Phase 1: Core Infrastructure
1. Create `ShortcutRegistry` and interface
2. Add registry to `DevTUI` struct
3. Basic registration in `AddEditHandler`

### Phase 2: Keyboard Integration
1. Add shortcut detection in `handleNormalModeKeyboard`
2. Implement `executeShortcut` method
3. Add progress message routing

### Phase 3: Testing & Documentation
1. Write comprehensive tests
2. Update documentation
3. Update shortcuts display

### Phase 4: Polish & Edge Cases
1. Handle dynamic handler removal
2. Add conflict resolution
3. Performance optimization

## Potential Issues & Mitigations

### Issue 1: Handler Lifecycle
**Problem**: What happens if a handler is removed but shortcuts remain registered?
**Mitigation**: 
- Add cleanup in handler removal (if such method exists)
- Validate indexes before execution
- Add error logging for invalid shortcuts

### Issue 2: Performance Impact
**Problem**: Registry lookup on every keypress
**Mitigation**: 
- Use fast map lookup (O(1))
- Only check single-character keys
- Profile and optimize if needed

### Issue 3: Thread Safety
**Problem**: Concurrent access to registry during registration/execution
**Mitigation**: Use RWMutex in registry (already planned)

### Issue 4: Memory Leaks
**Problem**: Registry holding references to removed handlers
**Mitigation**: 
- Store minimal data in registry entries
- Implement cleanup methods
- Periodic registry validation

## Success Criteria

1. **Functional**: Shortcuts work as described in the original request
2. **Performance**: No noticeable impact on normal keyboard handling
3. **Robust**: Handles edge cases gracefully with proper error logging
4. **Testable**: Comprehensive test coverage (>90%)
5. **Maintainable**: Clean, well-documented code following project patterns
6. **Backward Compatible**: Existing handlers continue to work unchanged

## Example Usage

Once implemented, the feature would be used like this:

```go
type TinyWasmHandler struct {
    mode string
}

func (h *TinyWasmHandler) Name() string  { return "TinyWasm" }
func (h *TinyWasmHandler) Label() string { return "Compilation Mode" }
func (h *TinyWasmHandler) Value() string { return h.mode }

func (h *TinyWasmHandler) Change(newValue string, progress func(msgs ...any)) {
    switch newValue {
    case "c":
        h.mode = "coding"
        progress("Switched to coding mode")
    case "d":
        h.mode = "debug" 
        progress("Switched to debug mode")
    case "p":
        h.mode = "production"
        progress("Switched to production mode")
    }
}

// Enable shortcuts with descriptions
func (h *TinyWasmHandler) Shortcuts() map[string]string {
    return map[string]string{
        "c": "coding mode",
        "d": "debug mode",
        "p": "production mode",
    }
}

// Usage
tab := tui.NewTabSection("Build", "Build configuration")
tab.AddEditHandler(&TinyWasmHandler{mode: "coding"}, 5*time.Second)

// Now pressing 'c', 'd', or 'p' from any tab will switch to this handler
// and change the compilation mode accordingly
```

### DatabaseHandler Example for Testing

For testing the shortcut functionality in the demo interface (`example/demo/main.go`), here's how to extend the existing `DatabaseHandler`:

```go
// Update example/HandlerEdit.go
type DatabaseHandler struct {
    ConnectionString string
    LastAction       string
}

func (h *DatabaseHandler) Name() string  { return "DatabaseConfig" }
func (h *DatabaseHandler) Label() string { return "Database Connection" }
func (h *DatabaseHandler) Value() string { return h.ConnectionString }

func (h *DatabaseHandler) Change(newValue string, progress func(msgs ...any)) {
    switch newValue {
    case "t":
        h.LastAction = "test"
        if progress != nil {
            progress("Testing database connection...")
            time.Sleep(500 * time.Millisecond)
            progress("Connection test completed successfully")
        }
    case "b":
        h.LastAction = "backup"
        if progress != nil {
            progress("Starting database backup...")
            time.Sleep(1000 * time.Millisecond)
            progress("Database backup completed successfully")
        }
    default:
        // Regular connection string update
        if progress != nil {
            progress("Validating", "Connection", newValue)
            time.Sleep(500 * time.Millisecond)
            progress("Testing database connectivity...", newValue)
            time.Sleep(500 * time.Millisecond)
            progress("Connection", "Database", "configured", "successfully", newValue)
        }
        h.ConnectionString = newValue
    }
}

// NEW: Add shortcut support
func (h *DatabaseHandler) Shortcuts() map[string]string {
    return map[string]string{
        "t": "test connection",
        "b": "backup database",
    }
}
```

With this implementation:
- Pressing `t` from any tab will navigate to the DatabaseHandler and trigger a connection test
- Pressing `b` from any tab will navigate to the DatabaseHandler and trigger a database backup
- The shortcuts will be displayed in the SHORTCUTS tab automatically

This implementation provides a clean, efficient, and extensible solution for global shortcuts while maintaining DevTUI's existing architecture and patterns.
