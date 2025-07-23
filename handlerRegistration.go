package devtui

// NEW: Handler registration methods using builder pattern

// NewEditHandler creates a new EditHandlerBuilder for method chaining
func (ts *tabSection) NewEditHandler(handler HandlerEdit) *editHandlerBuilder {
	return &editHandlerBuilder{
		tabSection: ts,
		handler:    handler,
		timeout:    0, // Default: synchronous
	}
}

// NewEditHandlerWithTracking creates a new EditHandlerBuilder with tracking support
// Note: The builder automatically detects if handler implements MessageTracker interface
func (ts *tabSection) NewEditHandlerWithTracking(handler EditHandlerTracker) *editHandlerBuilder {
	return &editHandlerBuilder{
		tabSection: ts,
		handler:    handler, // EditHandlerTracker extends HandlerEdit
		timeout:    0,       // Default: synchronous
	}
}

// NewExecutionHandler creates a new RunHandlerBuilder for method chaining
func (ts *tabSection) NewExecutionHandler(handler HandlerExecution) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler,
		timeout:    0, // Default: synchronous
	}
}

// NewExecutionHandlerTracking creates a new ExecutionHandlerBuilder with tracking support
// Note: The builder automatically detects if handler implements MessageTracker interface
func (ts *tabSection) NewExecutionHandlerTracking(handler ExecutionHandlerTracker) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler, // ExecutionHandlerTracker extends HandlerExecution
		timeout:    0,       // Default: synchronous
	}
}

// NewDisplayHandler creates a new DisplayHandlerBuilder for method chaining
func (ts *tabSection) NewDisplayHandler(handler HandlerDisplay) *displayHandlerBuilder {
	return &displayHandlerBuilder{
		tabSection: ts,
		handler:    handler,
	}
}

// NewWriterHandler creates a new WriterHandlerBuilder for method chaining
func (ts *tabSection) NewWriterHandler(handler any) *writerHandlerBuilder {
	return &writerHandlerBuilder{
		tabSection: ts,
		handler:    handler,
	}
}

// NewWriterHandlerTracking creates a new WriterHandlerBuilder with tracking support
// For handlers that implement both HandlerWriter and MessageTracker interfaces
func (ts *tabSection) NewWriterHandlerTracking(handler HandlerTrackerWriter) *writerHandlerBuilder {
	return &writerHandlerBuilder{
		tabSection: ts,
		handler:    handler, // HandlerTrackerWriter extends both interfaces
	}
}
