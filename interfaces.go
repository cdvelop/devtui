package devtui

import "time"

// HandlerDisplay defines the interface for read-only information display handlers.
// These handlers show static or dynamic content without user interaction.
type HandlerDisplay interface {
	Label() string   // Display label (e.g., "Help", "Status")
	Content() string // Display content (e.g., "help\n1-..\n2-...", "executing deploy wait...")
}

// HandlerEdit defines the interface for interactive fields that accept user input.
// These handlers allow users to modify values through text input.
type HandlerEdit interface {
	Label() string // Field label (e.g., "Server Port", "Host Configuration")
	Value() string // Current/initial value (e.g., "8080", "localhost")
	Change(newValue any, progress ...func(string)) error
}

// HandlerExecution defines the interface for action buttons that execute operations.
// These handlers trigger business logic when activated by the user.
type HandlerExecution interface {
	Label() string // Button label (e.g., "Deploy to Production", "Build Project")
	Execute(progress ...func(string)) error
}

// HandlerWriter defines the interface for basic writers that create new lines for each write.
// These writers are suitable for simple logging or output display.
type HandlerWriter interface {
	Name() string // Writer identifier (e.g., "webBuilder", "ApplicationLog")
}

// HandlerTrackerWriter defines the interface for advanced writers that can update existing lines.
// These writers support message tracking and can modify previously written content.
type HandlerTrackerWriter interface {
	Name() string
	MessageTracker
}

// MessageTracker provides optional interface for message tracking control.
// Handlers can implement this to control message updates and operation tracking.
type MessageTracker interface {
	GetLastOperationID() string
	SetLastOperationID(id string)
}

// EditHandlerTracker combines HandlerEdit with MessageTracker for advanced edit handlers
// that need message tracking capabilities.
type EditHandlerTracker interface {
	HandlerEdit
	MessageTracker
}

// handlerWithTimeout represents a handler with an associated timeout configuration.
type handlerWithTimeout struct {
	Handler any
	Timeout time.Duration
}
