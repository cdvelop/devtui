package devtui

// HandlerDisplay defines the interface for read-only information display handlers.
// These handlers show static or dynamic content without user interaction.
type HandlerDisplay interface {
	Name() string    // Full text to display in footer (handler responsible for content) eg. "System Status Information Display"
	Content() string // Display content (e.g., "help\n1-..\n2-...", "executing deploy wait...")
}

// HandlerEdit defines the interface for interactive fields that accept user input.
// These handlers allow users to modify values through text input.
type HandlerEdit interface {
	Name() string                                       // Identificador para logging: "ServerPort", "DatabaseURL"
	Label() string                                      // Field label (e.g., "Server Port", "Host Configuration")
	Value() string                                      // Current/initial value (e.g., "8080", "localhost")
	Change(newValue string, progress func(msgs ...any)) // Nueva firma: sin error, sin variádico, string
}

// HandlerEditTracker combines HandlerEdit with MessageTracker for advanced edit handlers
// that need message tracking capabilities.
type HandlerEditTracker interface {
	HandlerEdit
	MessageTracker
}

// HandlerExecution defines the interface for action buttons that execute operations.
// These handlers trigger business logic when activated by the user.
type HandlerExecution interface {
	Name() string                       // Identificador para logging: "DeployProd", "BuildProject"
	Label() string                      // Button label (e.g., "Deploy to Production", "Build Project")
	Execute(progress func(msgs ...any)) // Nueva firma: sin error, sin variádico
}

// HandlerExecutionTracker combines HandlerExecution with MessageTracker for advanced execution handlers
// that need message tracking capabilities.
type HandlerExecutionTracker interface {
	HandlerExecution
	MessageTracker
}

// HandlerWriter defines the interface for basic writers that create new lines for each write.
// These writers are suitable for simple logging or output display.
type HandlerWriter interface {
	Name() string // Writer identifier (e.g., "webBuilder", "ApplicationLog")
}

// HandlerWriterTracker defines the interface for advanced writers that can update existing lines.
// These writers support message tracking and can modify previously written content.
type HandlerWriterTracker interface {
	Name() string
	MessageTracker
}

// MessageTracker provides optional interface for message tracking control.
// Handlers can implement this to control message updates and operation tracking.
type MessageTracker interface {
	GetLastOperationID() string
	SetLastOperationID(id string)
}
