package devtui

import (
	. "github.com/cdvelop/tinystring"
)

// Print sends a normal Label or error to the tui in current tab
func (h *DevTUI) Print(messages ...any) {
	message, msgType := Translate(messages...).StringType()
	h.sendMessage(message, msgType, h.tabSections[h.activeTab])
}

// sendMessage envía un mensaje al tui por el canal de mensajes
func (d *DevTUI) sendMessage(content string, mt MessageType, tabSection *tabSection, operationID ...string) {
	var opID string
	if len(operationID) > 0 {
		opID = operationID[0]
	}
	newContent := d.createTabContent(content, mt, tabSection, "", opID)

	// Agregar contenido directamente al slice
	tabSection.mu.Lock()
	tabSection.tabContents = append(tabSection.tabContents, newContent)
	tabSection.mu.Unlock()

	d.tabContentsChan <- newContent
}

// NEW: sendMessageWithHandler sends a message with handler identification
func (d *DevTUI) sendMessageWithHandler(content string, mt MessageType, tabSection *tabSection, handlerName string, operationID string) {
	// Use update or add function that handles operationID reuse
	_, newContent := tabSection.updateOrAddContentWithHandler(mt, content, handlerName, operationID)

	// Always send to channel to trigger UI update, regardless of whether content was updated or added new
	d.tabContentsChan <- newContent

	// Call SetLastOperationID on the handler after processing
	// First try writing handlers, then field handlers
	var targetHandler *anyHandler
	if handler := tabSection.getWritingHandler(handlerName); handler != nil {
		targetHandler = handler
	} else {
		// Search in field handlers
		for _, field := range tabSection.FieldHandlers() {
			if field.handler != nil && field.handler.Name() == handlerName {
				targetHandler = field.handler
				break
			}
		}
	}

	if targetHandler != nil {
		targetHandler.SetLastOperationID(newContent.Id)
	} else {
		// DEBUG: Log when handler is not found (temporary for debugging)
		if tabSection.tui != nil && tabSection.tui.LogToFile != nil {
			tabSection.tui.LogToFile(Fmt("DEBUG: Handler not found for '%s'. Available field handlers:", handlerName))
			for i, field := range tabSection.FieldHandlers() {
				if field.handler != nil {
					tabSection.tui.LogToFile(Fmt("  [%d] %s", i, field.handler.Name()))
				}
			}
		}
	}
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tabContent) string {
	// Check if message comes from a readonly field handler (HandlerDisplay)
	if msg.handlerName != "" && t.isReadOnlyHandler(msg.handlerName) {
		// For readonly fields: no timestamp, cleaner visual content, no special coloring
		return msg.Content
	}

	// Apply message type styling to content (unified for all handler types)
	styledContent := t.applyMessageTypeStyle(msg.Content, msg.Type)

	// Generate timestamp (unified for all handler types that need it)
	timeStr := t.generateTimestamp(msg.Timestamp)

	// Check if message comes from interactive handler - clean format with timestamp only
	if msg.handlerName != "" && t.isInteractiveHandler(msg.handlerName) {
		// Interactive handlers: timestamp + content (no handler name for cleaner UX)
		return Fmt("%s %s", timeStr, styledContent)
	}

	// Default format for other handlers (Edit, Execution, Writers)
	handlerName := t.formatHandlerName(msg.handlerName)
	return Fmt("%s %s%s", timeStr, handlerName, styledContent)
}

// Helper methods to reduce code duplication

func (t *DevTUI) applyMessageTypeStyle(content string, msgType MessageType) string {
	switch msgType {
	case Msg.Error:
		return t.errStyle.Render(content)
	case Msg.Warning:
		return t.warnStyle.Render(content)
	case Msg.Info:
		return t.infoStyle.Render(content)
	case Msg.Success:
		return t.okStyle.Render(content)
	default:
		return content
	}
}

func (t *DevTUI) generateTimestamp(timestamp string) string {
	if t.id != nil {
		return t.timeStyle.Render(t.id.UnixNanoToTime(timestamp))
	}
	return t.timeStyle.Render("--:--:--")
}

func (t *DevTUI) formatHandlerName(handlerName string) string {
	if handlerName == "" {
		return ""
	}
	// Aplicar estilo completo a [handlerName] como una unidad
	styledName := t.infoStyle.Render(Fmt("[%s]", handlerName))
	return styledName + " "
}

// Helper to detect readonly handlers
func (t *DevTUI) isReadOnlyHandler(handlerName string) bool {
	// Check if handler has empty label (readonly convention)
	for _, tab := range t.tabSections {
		if handler := tab.getWritingHandler(handlerName); handler != nil {
			// Check if it's a display handler (readonly)
			return handler.handlerType == handlerTypeDisplay
		}
	}
	return false
}

// NEW: Helper to detect interactive handlers
func (t *DevTUI) isInteractiveHandler(handlerName string) bool {
	for _, tab := range t.tabSections {
		for _, field := range tab.FieldHandlers() {
			if field.handler != nil && field.handler.Name() == handlerName {
				return field.handler.handlerType == handlerTypeInteractive
			}
		}
	}
	return false
}

// createTabContent creates tabContent with unified logic (replaces newContent and newContentWithHandler)
func (h *DevTUI) createTabContent(content string, mt MessageType, tabSection *tabSection, handlerName string, operationID string) tabContent {
	// Timestamp SIEMPRE nuevo usando GetNewID - Handle gracefully if unixid failed to initialize
	var timestamp string
	if h.id != nil {
		timestamp = h.id.GetNewID()
	} else {
		errMsg := "error: unixid not initialized, using fallback timestamp for content: " + content
		// Log the issue before using fallback
		if h.LogToFile != nil {
			h.LogToFile(errMsg)
		}
		panic(errMsg) // Panic to ensure we catch this critical issue
		// Graceful fallback when unixid initialization failed
	}

	var id string
	var opID *string

	// Lógica unificada para ID
	if operationID != "" {
		id = operationID
		opID = &operationID
	} else {
		// Usar el mismo timestamp como ID para operaciones nuevas
		id = timestamp
		opID = nil
	}

	return tabContent{
		Id:          id,
		Timestamp:   timestamp, // NUEVO campo
		Content:     content,
		Type:        mt,
		tabSection:  tabSection,
		operationID: opID,
		isProgress:  false,
		isComplete:  false,
		handlerName: handlerName,
	}
}
