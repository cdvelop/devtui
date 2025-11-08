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
	Name() string                                       // Identifier for logging: "ServerPort", "DatabaseURL"
	Label() string                                      // Field label (e.g., "Server Port", "Host Configuration")
	Value() string                                      // Current/initial value (e.g., "8080", "localhost")
	Change(newValue string, progress chan<- string)
}

// HandlerExecution defines the interface for action buttons that execute operations.
// These handlers trigger business logic when activated by the user.
type HandlerExecution interface {
	Name() string      // Identifier for logging: "DeployProd", "BuildProject"
	Label() string     // Button label (e.g., "Deploy to Production", "Build Project")
	Execute(progress chan<- string)
}

// HandlerLogger defines the interface for basic writers that create new lines for each write.
// These writers are suitable for simple logging or output display.
type HandlerLogger interface {
	Name() string // Writer identifier (e.g., "webBuilder", "ApplicationLog")
}

// HandlerLoggerTracker defines the interface for advanced writers that can update existing lines.
// These writers support message tracking and can modify previously written content.
// HandlerLoggerTracker defines the interface for advanced writers that can update existing lines.
// These writers support message tracking and can modify previously written content.
// This interface has been removed to limit public exposure.

// HandlerInteractive defines the interface for interactive content handlers.
// These handlers combine content display with user interaction capabilities.
// All content display is handled through progress() for consistency.
type HandlerInteractive interface {
	Name() string                                // Identifier for logging: "ChatBot", "ConfigWizard"
	Label() string                               // Field label (updates dynamically)
	Value() string                               // Current input value
	Change(newValue string, progress chan<- string)  // Handle user input + content display via progress
	WaitingForUser() bool                        // Should edit mode be auto-activated?
}

// MessageTracker provides optional interface for message tracking control.
// Handlers can implement this to control message updates and operation tracking.
type MessageTracker interface {
	GetLastOperationID() string
	SetLastOperationID(id string)
}

// ShortcutProvider defines the optional interface for handlers that provide global shortcuts.
// HandlerEdit implementations can implement this interface to enable global shortcut keys.
type ShortcutProvider interface {
	Shortcuts() []map[string]string // Returns ordered list of single-entry maps with shortcut->description, preserving registration order
}
