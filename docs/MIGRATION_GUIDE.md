# Migration Guide

## Updating Existing Handlers

If you have existing handlers, you need to add WritingHandler interface methods to comply with the new API.

### Before (Old API):
```go
type MyHandler struct {
    currentValue string
}

func (h *MyHandler) Label() string { return "My Field" }
func (h *MyHandler) Value() string { return h.currentValue }
func (h *MyHandler) Editable() bool { return true }
func (h *MyHandler) Timeout() time.Duration { return 5 * time.Second }
func (h *MyHandler) Change(newValue any) (string, error) {
    h.currentValue = newValue.(string)
    return "Value updated successfully", nil
}
```

### After (New API):
```go
type MyHandler struct {
    currentValue string
    lastOpID     string  // NEW: For WritingHandler interface
}

// NEW: WritingHandler implementation
func (h *MyHandler) Name() string { return "MyHandler" }
func (h *MyHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *MyHandler) GetLastOperationID() string { return h.lastOpID }

// Existing FieldHandler implementation (unchanged)
func (h *MyHandler) Label() string { return "My Field" }
func (h *MyHandler) Value() string { return h.currentValue }
func (h *MyHandler) Editable() bool { return true }
func (h *MyHandler) Timeout() time.Duration { return 5 * time.Second }
func (h *MyHandler) Change(newValue any) (string, error) {
    h.currentValue = newValue.(string)
    return "Value updated successfully", nil
}
```

## Key Changes Required

1. **Add `lastOpID string` field** to your handler structs
2. **Implement `Name()` method** returning a unique handler identifier
3. **Implement `SetLastOperationID(id string)`** method for operation tracking
4. **Implement `GetLastOperationID() string`** method for operation retrieval
5. All existing `Label()`, `Value()`, `Editable()`, `Timeout()`, and `Change()` methods remain **unchanged**

## Breaking Changes

### FieldHandler Interface
The `FieldHandler` interface now embeds `WritingHandler`:

```go
// OLD: Simple FieldHandler
type FieldHandler interface {
    Label() string
    Value() string
    Editable() bool
    Change(newValue any) (string, error)
    Timeout() time.Duration
}

// NEW: FieldHandler with embedded WritingHandler
type FieldHandler interface {
    WritingHandler                     // Embedded interface
    Label() string
    Value() string
    Editable() bool
    Change(newValue any) (string, error)
    Timeout() time.Duration
}

type WritingHandler interface {
    Name() string
    SetLastOperationID(id string)
    GetLastOperationID() string
}
```

## Message System Changes

### Enhanced Message Format
Messages now include handler identification:
- **Old format**: `12:34:56 Operation completed successfully`
- **New format**: `12:34:56 [HandlerName] Operation completed successfully`

### Operation ID Tracking
Handlers can now control message updates vs creating new messages:
- Messages with the same operation ID update existing content
- New operation IDs create new messages
- Automatic operation ID assignment for transparent operation tracking

## Benefits of Migration

### For Developers:
1. **Better Debugging**: Messages clearly identify their source handlers
2. **Message Correlation**: Related operations can be tracked and updated
3. **Concurrent Safety**: Multiple handlers can write to the same tab safely
4. **io.Writer Support**: Handlers can use standard Go writing patterns

### For Users:
1. **Clearer Interface**: Messages show which component generated them
2. **Progressive Updates**: Long operations can update their status in place
3. **Better Organization**: Message sources are always identifiable

## Migration Checklist

- [ ] Add `lastOpID string` field to all handler structs
- [ ] Implement `Name()` method for each handler (use unique names)
- [ ] Implement `SetLastOperationID(id string)` method
- [ ] Implement `GetLastOperationID() string` method
- [ ] Test that all existing functionality continues to work
- [ ] Verify message source identification is working correctly
- [ ] Update any custom interfaces or mocks to include WritingHandler methods

## Backward Compatibility

The migration maintains full backward compatibility for:
- ✅ All existing `Label()`, `Value()`, `Editable()`, `Change()`, `Timeout()` methods
- ✅ Message generation and display functionality  
- ✅ Async operation handling and timeouts
- ✅ Keyboard navigation and field editing
- ✅ Tab management and content organization

Only **new methods** need to be implemented - no existing code needs to be changed.
