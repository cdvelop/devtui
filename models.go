package devtui

import "github.com/cdvelop/messagetype"

// Interfaz centralizada para manejo de campos
type fieldHandlerAdapter interface {
	Name() string
	Value() string
	Editable() bool
	ChangeValue(newValue string) <-chan MessageUpdate
}

// Estructura unificada para handlers de campos
type fieldHandler struct {
	fieldHandlerAdapter
	tempEditValue    string
	index            int
	cursor           int
	currentMessageID MessageID
}

// ExecuteValueChange handles the value change and returns the initial message ID
func (fh *fieldHandler) ExecuteValueChange(newValue string, tabSection *tabSection) MessageID {
	// Create initial "processing" message
	initialMsg := tabSection.newTuiMessage("Processing: "+fh.Name()+" change...", messagetype.Info)
	msgID := tabSection.tui.messageTracker.TrackMessage(&initialMsg)
	fh.currentMessageID = msgID

	// Add the message to the tab section
	tabSection.tuiMessages = append(tabSection.tuiMessages, initialMsg)

	// Start a goroutine to handle updates
	go func() {
		updateChan := fh.ChangeValue(newValue)

		// Process updates as they come in
		for update := range updateChan {
			tabSection.tui.asyncMessageChan <- tuiMessage{
				id:         string(msgID),
				Content:    update.Content,
				Type:       update.Type,
				tabSection: tabSection,
			}
		}
	}()

	return msgID
}

// Representación de una sección de pestaña en la UI
type tabSection struct {
	index                int
	title                string
	fieldHandlers        []fieldHandler
	tuiMessages          []tuiMessage
	indexActiveEditField int
	tui                  *DevTUI
}

// Mensaje para el sistema de UI
type tuiMessage struct {
	id         string
	Content    string
	Type       messagetype.Type
	tabSection *tabSection
}

// Actualización de estado para mensajes
type MessageUpdate struct {
	Content string
	Type    messagetype.Type
}

// Identificador único para mensajes
type MessageID string
