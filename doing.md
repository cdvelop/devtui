## TabSection Refactoring

### Current Implementation
The `fieldHandler` struct currently has two methods for handling value changes:
- Synchronous: Using `ChangeValue` function that returns a message and error immediately
- Asynchronous: Using `AsyncFieldValueChange` and `IsAsync` flag that sends multiple messages over time

This dual approach increases complexity and creates inconsistent patterns for handling updates.

### Proposed Changes
Remove the separate async fields (`AsyncFieldValueChange` and `IsAsync`) and implement a unified asynchronous approach:

1. Modify `ChangeValue` to work asynchronously by default
2. Use a message ID system to update existing messages rather than creating new ones
3. Each field change operation will create/update a single `tuiMessage` entry

### Implementation Details 
```go
func (t *DevTUI) NewTabSection(title string, fhAdapters... fieldHandlerAdapter) *tabSection {
    // create a new tab section with the given title and field handlers
    newSection := &tabSection{
        title:                title,
        FieldHandlers:        make([]fieldHandler, 0, len(FieldHandlers)),
        tuiMessages:          make([]tuiMessage, 0),
        indexActiveEditField: 0, //  selected field index
        tui:                  t,
    }
    
    // build fieldHandler from FieldHandlers
    for i, fha := range fhAdapters {
        stdHandler := fieldHandler{
            fieldHandlerAdapter:  fha,
            tempEditValue: "",
            index:         i,
            cursor:        len(fha.Value()), // position cursor at end of text by default
        }
        newSection.FieldHandlers = append(newSection.FieldHandlers, stdHandler)
    }
    
    // Set the index of the tab section (usually managed by DevTUI internally)
    newSection.index = len(t.tabSections)
    
    // Add to the DevTUI sections
    t.tabSections = append(t.tabSections, newSection)
    
    return newSection
}

type fieldHandlerAdapter interface{
	Name() string // eg: "Server Port"
    Value() string //initial value or after change eg: "8080"
    Editable() bool // if no editable eject the action changed from 8080 to 9090"
	ChangeValue(newValue string) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port 
}

type fieldHandler struct {
    fieldHandlerAdapter // handler field interface
	tempEditValue string                                                // use for edit
	// Async handler that can send multiple messages over time
	AsyncFieldValueChange func(newValue string, messageChan chan<- tuiMessage)
	IsAsync               bool // Flag indicating if this handler uses async processing
	//internal use
	index  int
	cursor int // cursor position in text value
}

type tabSection struct {
	index                int                   // index of the tab
	title                string                // eg: "BUILD", "TEST"
	FieldHandlers        []fieldHandler // Field actions configured for the section
	tuiMessages          []tuiMessage          // message contents
	indexActiveEditField int                   // Índice del campo de configuración seleccionado
	tui                  *DevTUI
}

// tuiMessage represents a message sent asynchronously or synchronously in the TUI
type tuiMessage struct {
	id         string // Unique ID for the message
	Content    string
	Type       messagetype.Type
	tabSection *tabSection
}

func (t *tabSection) newTuiMessage(content string, mt messagetype.Type) tuiMessage {

	return tuiMessage{
		id:         h.id.GetNewID(),
		Content:    content,
		Type:       mt,
		tabSection: t,
	}
}


//*** tui_message.go
package devtui

// MessageUpdate represents a status update for a tuiMessage
type MessageUpdate struct {
    Content string
    Type    messagetype.Type
}

// MessageID represents a unique identifier for tracking message updates
type MessageID string

// messageTracker keeps track of messages and their IDs
type messageTracker struct {
    messages map[MessageID]*tuiMessage
    idGen    IDGenerator
}

// NewMessageTracker creates a new message tracker
func NewMessageTracker() *messageTracker {
    return &messageTracker{
        messages: make(map[MessageID]*tuiMessage),
        idGen:    NewIDGenerator(),
    }
}

// TrackMessage adds a message to be tracked
func (mt *messageTracker) TrackMessage(msg *tuiMessage) MessageID {
    id := MessageID(mt.idGen.GetNewID())
    msg.id = string(id)
    mt.messages[id] = msg
    return id
}

// UpdateMessage updates a tracked message with new content
func (mt *messageTracker) UpdateMessage(id MessageID, update MessageUpdate) bool {
    msg, exists := mt.messages[id]
    if !exists {
        return false
    }
    
    msg.Content = update.Content
    msg.Type = update.Type
    return true
}


//***** fiel_handler.go

package devtui
// fieldHandlerAdapter defines the interface for field handlers
type fieldHandlerAdapter interface {
    Name() string                                    // eg: "Server Port"
    Value() string                                  // initial value or after change eg: "8080"
    Editable() bool                                 // if not editable execute action directly
    // ChangeValue now returns a channel that will send message updates
    ChangeValue(newValue string) <-chan MessageUpdate // Returns a channel for updates
}

// fieldHandler embeds the adapter and adds UI-specific fields
type fieldHandler struct {
    fieldHandlerAdapter                           // handler field interface
    tempEditValue       string                    // used for edit mode
    index               int                       // position in the list
    cursor              int                       // cursor position in text value
    currentMessageID    MessageID                 // ID of the current message being updated
}

// ExecuteValueChange handles the value change and returns the initial message
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

//***** update.go
package devtui

// Update the AsyncMessageMsg handling in Update function
func (h *DevTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ...existing code...
    
    switch msg := msg.(type) {
    // ...existing code...
    
    case AsyncMessageMsg:
        asyncMsg := tuiMessage(msg)
        
        // Check if this is an update to an existing message
        if asyncMsg.id != "" {
            // Find and update the existing message
            for i, existingMsg := range h.tabSections[h.activeTab].tuiMessages {
                if existingMsg.id == asyncMsg.id {
                    h.tabSections[h.activeTab].tuiMessages[i].Content = asyncMsg.Content
                    h.tabSections[h.activeTab].tuiMessages[i].Type = asyncMsg.Type
                    break
                }
            }
        } else {
            // If no ID, treat as a new message
            h.sendMessage(asyncMsg.Content, asyncMsg.Type, asyncMsg.tabSection)
        }
        
        // Continue listening for more async messages
        cmds = append(cmds, h.listenForAsyncMessages(h.asyncMessageChan))
    
    // ...existing code...
    }
    
    // ...existing code...
}


// example usage
package main
// Example implementation of fieldHandlerAdapter
type ServerPortHandler struct {
    currentValue string
}

func (s *ServerPortHandler) Name() string {
    return "Server Port"
}

func (s *ServerPortHandler) Value() string {
    return s.currentValue
}

func (s *ServerPortHandler) Editable() bool {
    return true
}

func (s *ServerPortHandler) ChangeValue(newValue string) <-chan MessageUpdate {
    updates := make(chan MessageUpdate)
    
    go func() {
        defer close(updates)
        
        // Send "processing" update
        updates <- MessageUpdate{
            Content: "Changing port from " + s.currentValue + " to " + newValue + "...",
            Type:    messagetype.Info,
        }
        
        // Simulate some work
        time.Sleep(1 * time.Second)
        
        // Check if port is valid
        portNum, err := strconv.Atoi(newValue)
        if err != nil || portNum < 1 || portNum > 65535 {
            updates <- MessageUpdate{
                Content: "Invalid port number: " + newValue,
                Type:    messagetype.Error,
            }
            return
        }
        
        // Update the value
        s.currentValue = newValue
        
        // Send success message
        updates <- MessageUpdate{
            Content: "Port successfully changed to " + newValue,
            Type:    messagetype.Success,
        }
    }()
    
    return updates
}

```
