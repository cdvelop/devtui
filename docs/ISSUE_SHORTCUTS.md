# DevTUI Shortcut Keys Interface

## Executive Summary

**Objective**: Implement an optional shortcut interface for DevTUI FieldHandler to enable quick keyboard shortcuts for specific fields (e.g., TinyWasm compiler modes: c/d/p).

**Integration**: Optional interface pattern that doesn't affect existing handlers, providing enhanced UX for power users.

---

## Technical Requirements

### 1. Optional ShortcutHandler Interface

```go
// Optional interface - only implemented by fields that need shortcuts
type ShortcutHandler interface {
    Shortcuts() map[string]string // key -> description
    HandleShortcut(key string) error
}
```

### 2. DevTUI Integration Pattern

```go
// In field handling logic
if shortcutHandler, ok := field.handler.(ShortcutHandler); ok {
    shortcuts := shortcutHandler.Shortcuts()
    // Display shortcuts in field help
    // Handle shortcut key presses
    if err := shortcutHandler.HandleShortcut(pressedKey); err != nil {
        // Handle shortcut error
    }
}
```

### 3. Example Implementation (TinyWasm)

```go
// TinyWasm implements ShortcutHandler optionally
func (w *TinyWasm) Shortcuts() map[string]string {
    return map[string]string{
        w.Config.CodingShortcut:     "Coding mode (fast compilation)",
        w.Config.DebuggingShortcut:  "Debug mode (TinyGo with debug symbols)",
        w.Config.ProductionShortcut: "Production mode (optimized binary)",
    }
}

func (w *TinyWasm) HandleShortcut(key string) error {
    // Reuse existing Change method
    _, err := w.Change(key)
    return err
}
```

### 4. UI Display Enhancement

**Field Label Enhancement**:
```
Current: "Build Mode: c, d, p"
Enhanced: "Build Mode: c, d, p (shortcuts: c=coding, d=debug, p=production)"
```

### 5. Benefits

- ✅ **Optional**: Doesn't affect other handlers
- ✅ **Simple**: Minimal implementation complexity
- ✅ **Great UX**: Quick mode switching for power users
- ✅ **Self-documenting**: Shortcuts display their purpose
- ✅ **Consistent**: Uses existing Change() validation logic

---

## Implementation Phases

### Phase 1: Interface Definition
1. Add `ShortcutHandler` interface to DevTUI
2. Update field handling logic to detect interface
3. Basic shortcut key processing

### Phase 2: UI Integration
1. Enhanced field label display with shortcuts
2. Shortcut help text generation
3. Key press handling and routing

### Phase 3: Testing & Polish
1. Unit tests for shortcut interface
2. Integration tests with TinyWasm
3. Performance testing for key handling

---

## Future Applications

- **Server configuration**: Port shortcuts (8080, 3000, 5000)
- **Environment modes**: dev/staging/prod shortcuts
- **Build targets**: Multiple architecture shortcuts
- **File operations**: Common path shortcuts

---

**Priority**: Medium - UX enhancement  
**Estimated Effort**: 1 day development + testing  
**Dependencies**: None - optional interface pattern
