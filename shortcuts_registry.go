package devtui

import "sync"

// ShortcutEntry represents a registered shortcut
type ShortcutEntry struct {
	Key         string // The shortcut key (e.g., "c", "d", "p")
	Description string // Human-readable description (e.g., "coding mode", "debug mode")
	TabIndex    int    // Index of the tab containing the handler
	FieldIndex  int    // Index of the field within the tab
	HandlerName string // Handler name for identification
	Value       string // Value to pass to Change()
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

// GetAll returns all registered shortcuts for UI display
func (sr *ShortcutRegistry) GetAll() map[string]*ShortcutEntry {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	// Return a copy to prevent external modification
	result := make(map[string]*ShortcutEntry)
	for k, v := range sr.shortcuts {
		result[k] = v
	}
	return result
}
