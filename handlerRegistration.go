package devtui

// NEW: Handler registration methods using builder pattern

// AddEditHandler creates a new EditHandlerBuilder for method chaining
func (ts *tabSection) AddEditHandler(handler HandlerEdit) *editHandlerBuilder {
	return &editHandlerBuilder{
		tabSection: ts,
		handler:    handler,
		timeout:    0, // Default: synchronous
	}
}

// AddEditHandlerTracking creates a new EditHandlerBuilder with tracking support
// Note: The builder automatically detects if handler implements MessageTracker interface
func (ts *tabSection) AddEditHandlerTracking(handler HandlerEditTracker) *editHandlerBuilder {
	return &editHandlerBuilder{
		tabSection: ts,
		handler:    handler, // HandlerEditTracker extends HandlerEdit
		timeout:    0,       // Default: synchronous
	}
}

// AddExecutionHandler creates a new RunHandlerBuilder for method chaining
func (ts *tabSection) AddExecutionHandler(handler HandlerExecution) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler,
		timeout:    0, // Default: synchronous
	}
}

// AddExecutionHandlerTracking creates a new ExecutionHandlerBuilder with tracking support
// Note: The builder automatically detects if handler implements MessageTracker interface
func (ts *tabSection) AddExecutionHandlerTracking(handler HandlerExecutionTracker) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler, // HandlerExecutionTracker extends HandlerExecution
		timeout:    0,       // Default: synchronous
	}
}

// AddHandlerDisplay registers a HandlerDisplay directly (no builder pattern).
func (ts *tabSection) AddHandlerDisplay(handler HandlerDisplay) *tabSection {
	anyH := newDisplayHandler(handler)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}
