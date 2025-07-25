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

// NewEditHandlerTracking creates a new EditHandlerBuilder with tracking support
// Note: The builder automatically detects if handler implements MessageTracker interface
func (ts *tabSection) NewEditHandlerTracking(handler HandlerEditTracker) *editHandlerBuilder {
	return &editHandlerBuilder{
		tabSection: ts,
		handler:    handler, // HandlerEditTracker extends HandlerEdit
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
func (ts *tabSection) NewExecutionHandlerTracking(handler HandlerExecutionTracker) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler, // HandlerExecutionTracker extends HandlerExecution
		timeout:    0,       // Default: synchronous
	}
}

// RegisterHandlerDisplay registers a HandlerDisplay directly (no builder pattern).
func (ts *tabSection) RegisterHandlerDisplay(handler HandlerDisplay) *tabSection {
	anyH := newDisplayHandler(handler)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}
