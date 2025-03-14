package devtui

import (
	"fmt"

	. "github.com/cdvelop/messagetype"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleKeyboard processes keyboard input and updates the model state
// returns whether the update function should continue processing or return early
func (h *DevTUI) HandleKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	if h.tabEditingConfig { // EDITING CONFIG IN SECTION
		return h.handleEditingConfigKeyboard(msg)
	} else {
		return h.handleNormalModeKeyboard(msg)
	}
}

// handleEditingConfigKeyboard handles keyboard input while in config editing mode
func (h *DevTUI) handleEditingConfigKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	currentTab := &h.tabSections[h.activeTab]
	currentField := &h.tabSections[h.activeTab].FieldHandlers[currentTab.indexActiveEditField]

	if currentField.Editable { // Si el campo es editable, permitir la edición
		switch msg.String() {
		case "enter": // Al presionar ENTER, guardamos los cambios o ejecutamos la acción
			msg, err := currentField.FieldValueChange(currentField.Value)
			if err != nil {
				h.addTerminalPrint(Error, fmt.Sprintf("Error: %v %v", currentField.Label, err))
			}

			h.editingConfigOpen(false, currentField, msg)
			return false, nil

		case "esc": // Al presionar ESC, descartamos los cambios
			h.editingConfigOpen(false, currentField, "Exited config mode")
			return false, nil

		case "left": // Mover el cursor a la izquierda
			if currentField.cursor > 0 {
				currentField.cursor--
			}

		case "right": // Mover el cursor a la derecha
			if currentField.cursor < len(currentField.Value) {
				currentField.cursor++
			}

		default:
			if msg.String() == "backspace" && currentField.cursor > 0 {
				currentField.Value = currentField.Value[:currentField.cursor-1] + currentField.Value[currentField.cursor:]
				currentField.cursor--
			} else if len(msg.String()) == 1 {
				currentField.Value = currentField.Value[:currentField.cursor] + msg.String() + currentField.Value[currentField.cursor:]
				currentField.cursor++
			}
		}
	} else { // Si el campo no es editable, solo ejecutar la acción
		switch msg.String() {
		case "enter":
			msgType := OK
			// content eg: "Browser Opened"
			content, err := currentField.FieldValueChange(currentField.Value)
			if err != nil {
				msgType = Error
				content = fmt.Sprintf("%s %s %s", currentField.Label, content, err.Error())
			}
			currentField.Value = content
			h.addTerminalPrint(msgType, content)
			h.tabEditingConfig = false
		}
	}

	return true, nil
}

// handleNormalModeKeyboard handles keyboard input in normal mode (not editing config)
func (h *DevTUI) handleNormalModeKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch msg.String() {
	case "up": // Mover hacia arriba el indice del campo activo
		currentTab := &h.tabSections[h.activeTab]
		if currentTab.indexActiveEditField > 0 {
			currentTab.indexActiveEditField--
		}

	case "down": // Mover hacia abajo el indice del campo activo
		currentTab := &h.tabSections[h.activeTab]
		if currentTab.indexActiveEditField < len(h.tabSections[0].FieldHandlers)-1 {
			currentTab.indexActiveEditField++
		}

	case "tab": // change tabSection
		h.activeTab = (h.activeTab + 1) % len(h.tabSections)
		h.editingConfigOpen(false, nil, "")
		h.updateViewport()

	case "shift+tab": // change tabSection
		h.editingConfigOpen(false, nil, "")
		h.activeTab = (h.activeTab - 1 + len(h.tabSections)) % len(h.tabSections)
		h.updateViewport()

	case "ctrl+l":
		// h.tabSections[h.activeTab].tabContents = []tabContent{}

	case "enter":
		h.editingConfigOpen(true, nil, "Entered config editing mode press 'esc' to exit")
		h.updateViewport()

	case "ctrl+c":
		close(h.ExitChan) // Cerrar el canal para señalizar a todas las goroutines
		return false, tea.Quit
	}

	return true, nil
}
