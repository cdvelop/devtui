package devtui

import (
	"io"
	"time"
)

// AddDisplayHandler registers a HandlerDisplay directly
func (ts *tabSection) AddDisplayHandler(handler HandlerDisplay) *tabSection {
	anyH := newDisplayHandler(handler)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

// AddEditHandler registers a HandlerEdit with mandatory timeout
func (ts *tabSection) AddEditHandler(handler HandlerEdit, timeout time.Duration) *tabSection {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := newEditHandler(handler, timeout, tracker)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)

	// Auto-register handler for writing if it implements HandlerWriterTracker
	if _, ok := handler.(HandlerWriterTracker); ok {
		if writerHandler, ok := handler.(HandlerWriter); ok {
			ts.RegisterWriterHandler(writerHandler)
		}
	}

	return ts
}

// AddEditHandlerTracking registers a HandlerEditTracker with mandatory timeout
func (ts *tabSection) AddEditHandlerTracking(handler HandlerEditTracker, timeout time.Duration) *tabSection {
	return ts.AddEditHandler(handler, timeout) // HandlerEditTracker extends HandlerEdit
}

// AddExecutionHandler registers a HandlerExecution with mandatory timeout
func (ts *tabSection) AddExecutionHandler(handler HandlerExecution, timeout time.Duration) *tabSection {
	anyH := newExecutionHandler(handler, timeout)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

// AddExecutionHandlerTracking registers a HandlerExecutionTracker with mandatory timeout
func (ts *tabSection) AddExecutionHandlerTracking(handler HandlerExecutionTracker, timeout time.Duration) *tabSection {
	return ts.AddExecutionHandler(handler, timeout) // HandlerExecutionTracker extends HandlerExecution
}

// RegisterWriterHandler registers a writer handler and returns io.Writer (kept from existing API)
func (ts *tabSection) RegisterWriterHandler(handler HandlerWriter) io.Writer {
	return ts.RegisterHandlerWriter(handler) // Delegate to existing implementation
}

// AddInteractiveHandler registers a HandlerInteractive with mandatory timeout
func (ts *tabSection) AddInteractiveHandler(handler HandlerInteractive, timeout time.Duration) *tabSection {
	var tracker MessageTracker
	if t, ok := handler.(MessageTracker); ok {
		tracker = t
	}

	anyH := newInteractiveHandler(handler, timeout, tracker)
	f := &field{
		handler:    anyH,
		parentTab:  ts,
		asyncState: &internalAsyncState{},
	}
	ts.addFields(f)
	return ts
}

// AddInteractiveHandlerTracking registers a HandlerInteractiveTracker with mandatory timeout
func (ts *tabSection) AddInteractiveHandlerTracking(handler HandlerInteractiveTracker, timeout time.Duration) *tabSection {
	return ts.AddInteractiveHandler(handler, timeout) // HandlerInteractiveTracker extends HandlerInteractive
}
