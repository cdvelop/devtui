package devtui

import (
	"io"
	"time"
)

// editHandlerBuilder provides method chaining for EditHandler registration with optional timeout.
type editHandlerBuilder struct {
	tabSection *tabSection
	handler    EditHandler
	timeout    time.Duration
}

// WithTimeout configures the handler to run asynchronously with the specified timeout.
// timeout = 0 means synchronous execution (default behavior).
// timeout > 0 means asynchronous execution with the specified timeout.
func (b *editHandlerBuilder) WithTimeout(timeout time.Duration) *tabSection {
	b.timeout = timeout
	// Create a temporary wrapper to hold timeout info
	wrapper := &handlerWithTimeout{
		Handler: b.handler,
		Timeout: b.timeout,
	}
	b.tabSection.registerHandlerWithTimeout(wrapper)
	return b.tabSection
}

// Register finalizes the handler registration with default synchronous behavior (timeout = 0).
func (b *editHandlerBuilder) Register() *tabSection {
	return b.WithTimeout(0)
}

// executionHandlerBuilder provides method chaining for ExecutionHandler registration with optional timeout.
type executionHandlerBuilder struct {
	tabSection *tabSection
	handler    ExecutionHandler
	timeout    time.Duration
}

// WithTimeout configures the handler to run asynchronously with the specified timeout.
// timeout = 0 means synchronous execution (default behavior).
// timeout > 0 means asynchronous execution with the specified timeout.
func (b *executionHandlerBuilder) WithTimeout(timeout time.Duration) *tabSection {
	b.timeout = timeout
	// Create a temporary wrapper to hold timeout info
	wrapper := &handlerWithTimeout{
		Handler: b.handler,
		Timeout: b.timeout,
	}
	b.tabSection.registerHandlerWithTimeout(wrapper)
	return b.tabSection
}

// Register finalizes the handler registration with default synchronous behavior (timeout = 0).
func (b *executionHandlerBuilder) Register() *tabSection {
	return b.WithTimeout(0)
}

// displayHandlerBuilder provides method chaining for DisplayHandler registration.
// DisplayHandlers are always synchronous and don't support timeout configuration.
type displayHandlerBuilder struct {
	tabSection *tabSection
	handler    DisplayHandler
}

// Register finalizes the DisplayHandler registration.
func (b *displayHandlerBuilder) Register() *tabSection {
	b.tabSection.registerHandler(b.handler)
	return b.tabSection
}

// writerHandlerBuilder provides method chaining for Writer registration.
type writerHandlerBuilder struct {
	tabSection *tabSection
	handler    any // WriterBasic or WriterTracker
}

// Register finalizes the Writer registration and returns the io.Writer.
func (b *writerHandlerBuilder) Register() io.Writer {
	return b.tabSection.registerWriterHandler(b.handler)
}
