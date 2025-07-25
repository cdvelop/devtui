package devtui

import (
	"time"
)

// editHandlerBuilder provides method chaining for HandlerEdit registration with optional timeout.
type editHandlerBuilder struct {
	tabSection *tabSection
	handler    HandlerEdit
	timeout    time.Duration
}

// WithTimeout configures the handler to run asynchronously with the specified timeout.
// timeout = 0 means synchronous execution (default behavior).
// timeout > 0 means asynchronous execution with the specified timeout.
func (b *editHandlerBuilder) WithTimeout(timeout time.Duration) *tabSection {
	b.timeout = timeout

	// Check if handler implements MessageTracker interface
	var tracker MessageTracker
	if t, ok := b.handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := newEditHandler(b.handler, b.timeout, tracker)

	f := &field{
		handler:    anyH,
		parentTab:  b.tabSection,
		asyncState: &internalAsyncState{},
	}

	b.tabSection.addFields(f)

	// Auto-register handler for writing if it implements HandlerWriterTracker interface (both HandlerWriter and MessageTracker)
	if _, ok := b.handler.(HandlerWriterTracker); ok {
		if writerHandler, ok := b.handler.(HandlerWriter); ok {
			b.tabSection.RegisterHandlerWriter(writerHandler)
		}
	}

	return b.tabSection
}

// Register finalizes the handler registration with default synchronous behavior (timeout = 0).
func (b *editHandlerBuilder) Register() *tabSection {
	return b.WithTimeout(0)
}

// executionHandlerBuilder provides method chaining for HandlerExecution registration with optional timeout.
type executionHandlerBuilder struct {
	tabSection *tabSection
	handler    HandlerExecution
	timeout    time.Duration
}

// WithTimeout configures the handler to run asynchronously with the specified timeout.
// timeout = 0 means synchronous execution (default behavior).
// timeout > 0 means asynchronous execution with the specified timeout.
func (b *executionHandlerBuilder) WithTimeout(timeout time.Duration) *tabSection {
	b.timeout = timeout
	anyH := newExecutionHandler(b.handler, b.timeout)

	f := &field{
		handler:    anyH,
		parentTab:  b.tabSection,
		asyncState: &internalAsyncState{},
	}

	b.tabSection.addFields(f)
	return b.tabSection
}

// Register finalizes the handler registration with default synchronous behavior (timeout = 0).
func (b *executionHandlerBuilder) Register() *tabSection {
	return b.WithTimeout(0)
}

// The writerHandlerBuilder struct and its Register method have been removed as part of the refactoring plan.
