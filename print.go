package devtui

import (
	"fmt"

	"github.com/cdvelop/messagetype"
)

// Print sends a normal Label or error to the tui in current tab
func (h *DevTUI) Print(messages ...any) {
	msgType := messagetype.DetectMessageType(messages...)
	h.sendMessage(joinMessages(messages...), msgType, h.tabSections[h.activeTab])
}

func joinMessages(messages ...any) (Label string) {
	var space string
	for _, m := range messages {
		Label += space + fmt.Sprint(m)
		space = " "
	}
	return
}

// sendMessage envía un mensaje al tui por el canal de mensajes
func (d *DevTUI) sendMessage(content string, mt messagetype.Type, tabSection *tabSection, operationID ...string) {
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
func (d *DevTUI) sendMessageWithHandler(content string, mt messagetype.Type, tabSection *tabSection, handlerName string, operationID string) {
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
			tabSection.tui.LogToFile(fmt.Sprintf("DEBUG: Handler not found for '%s'. Available field handlers:", handlerName))
			for i, field := range tabSection.FieldHandlers() {
				if field.handler != nil {
					tabSection.tui.LogToFile(fmt.Sprintf("  [%d] %s", i, field.handler.Name()))
				}
			}
		}
	}
}

// formatMessage formatea un mensaje según su tipo
func (t *DevTUI) formatMessage(msg tabContent) string {
	// Check if message comes from a readonly field handler
	if msg.handlerName != "" && t.isReadOnlyHandler(msg.handlerName) {
		// For readonly fields: no timestamp, cleaner visual content, no special coloring
		return msg.Content
	}

	var timeStr string
	if t.id != nil {
		timeStr = t.timeStyle.Render(t.id.UnixNanoToTime(msg.Timestamp))
	} else {
		// When unixid is not initialized, use a simple timestamp format
		timeStr = t.timeStyle.Render("--:--:--")
	}

	// NEW: Include handler name in message format with brackets united
	var handlerName string
	if msg.handlerName != "" {
		// Aplicar estilo completo a [handlerName] como una unidad
		styledName := t.infoStyle.Render(fmt.Sprintf("[%s]", msg.handlerName))
		handlerName = styledName + " "
	}

	// timeStr := t.timeStyle.Render(msg.Time.Format("15:04:05"))
	// content := fmt.Sprintf("[%s] %s", timeStr, msg.Content)

	switch msg.Type {
	case messagetype.Error:
		msg.Content = t.errStyle.Render(msg.Content)
	case messagetype.Warning:
		msg.Content = t.warnStyle.Render(msg.Content)
	case messagetype.Info:
		msg.Content = t.infoStyle.Render(msg.Content)
	case messagetype.Success:
		msg.Content = t.okStyle.Render(msg.Content)
		// default:
		// 	msg.Content= msg.Content
	}

	return fmt.Sprintf("%s %s%s", timeStr, handlerName, msg.Content)
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

// createTabContent creates tabContent with unified logic (replaces newContent and newContentWithHandler)
func (h *DevTUI) createTabContent(content string, mt messagetype.Type, tabSection *tabSection, handlerName string, operationID string) tabContent {
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
