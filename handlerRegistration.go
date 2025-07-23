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

// NewRunHandler creates a new RunHandlerBuilder for method chaining
func (ts *tabSection) NewRunHandler(handler HandlerExecution) *executionHandlerBuilder {
	return &executionHandlerBuilder{
		tabSection: ts,
		handler:    handler,
		timeout:    0, // Default: synchronous
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
