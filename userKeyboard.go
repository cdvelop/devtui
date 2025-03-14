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
		switch msg.Type {
		case tea.KeyEnter: // Guardar cambios o ejecutar acción
			msg, err := currentField.FieldValueChange(currentField.Value)
			if err != nil {
				h.addTerminalPrint(Error, fmt.Sprintf("Error: %v %v", currentField.Label, err))
			}

			h.editingConfigOpen(false, currentField, msg)
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyEsc: // Al presionar ESC, descartamos los cambios y salimos del modo edición
			h.editingConfigOpen(false, currentField, "Exited config mode")
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyLeft: // Mover el cursor a la izquierda dentro del texto
			if currentField.cursor > 0 {
				currentField.cursor--
			}

		case tea.KeyRight: // Mover el cursor a la derecha dentro del texto
			if currentField.cursor < len(currentField.Value) {
				currentField.cursor++
			}

		case tea.KeyBackspace: // Borrar carácter a la izquierda
			if currentField.cursor > 0 {
				currentField.Value = currentField.Value[:currentField.cursor-1] + currentField.Value[currentField.cursor:]
				currentField.cursor--
			}

		default:
			// Soportar entrada de caracteres (runes)
			if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
				currentField.Value = currentField.Value[:currentField.cursor] + string(msg.Runes) + currentField.Value[currentField.cursor:]
				currentField.cursor += len(msg.Runes)
			}
		}
	} else { // Si el campo no es editable, solo ejecutar la acción
		switch msg.Type {
		case tea.KeyEnter:
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
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyEsc: // Permitir también salir con ESC para campos no editables
			h.editingConfigOpen(false, currentField, "Exited config mode")
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil
		}
	}

	return true, nil
}

// handleNormalModeKeyboard handles keyboard input in normal mode (not editing config)
func (h *DevTUI) handleNormalModeKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	currentTab := &h.tabSections[h.activeTab]
	totalFields := len(currentTab.FieldHandlers)

	switch msg.Type {
	case tea.KeyUp, tea.KeyDown:
		// Las teclas arriba y abajo ya no modifican el campo activo
		// Solo controlarán el desplazamiento del viewport
		// No hacemos nada aquí para permitir que el manejo del viewport siga su curso normal

	case tea.KeyLeft: // Navegar al campo anterior (ciclo continuo)
		if totalFields > 0 {
			currentTab.indexActiveEditField = (currentTab.indexActiveEditField - 1 + totalFields) % totalFields
		}

	case tea.KeyRight: // Navegar al campo siguiente (ciclo continuo)
		if totalFields > 0 {
			currentTab.indexActiveEditField = (currentTab.indexActiveEditField + 1) % totalFields
		}

	case tea.KeyTab: // cambiar tabSection
		h.activeTab = (h.activeTab + 1) % len(h.tabSections)

		// Comprobar si debe entrar automáticamente en modo edición
		h.checkAutoEditMode()
		h.updateViewport()

	case tea.KeyShiftTab: // cambiar tabSection
		h.activeTab = (h.activeTab - 1 + len(h.tabSections)) % len(h.tabSections)

		// Comprobar si debe entrar automáticamente en modo edición
		h.checkAutoEditMode()
		h.updateViewport()

	case tea.KeyCtrlL:
		// h.tabSections[h.activeTab].tabContents = []tabContent{}

	case tea.KeyEnter: //Enter para entrar en modo edición, ejecuta la acción directamente si el campo no es editable
		if totalFields > 0 {
			field := &currentTab.FieldHandlers[currentTab.indexActiveEditField]
			if !field.Editable {
				msgType := OK
				content, err := field.FieldValueChange(field.Value)
				if err != nil {
					msgType = Error
					content = fmt.Sprintf("%s %s %s", field.Label, content, err.Error())
				}
				field.Value = content
				h.addTerminalPrint(msgType, content)
			} else {
				// Para campos editables, activar modo de edición explícitamente
				h.tabEditingConfig = true
				h.editingConfigOpen(true, field, "Entered config editing mode press 'esc' to exit")
			}
			h.updateViewport()
		}

	case tea.KeyCtrlC:
		close(h.ExitChan) // Cerrar el canal para señalizar a todas las goroutines
		return false, tea.Quit
	}

	return true, nil
}

// checkAutoEditMode verifica si debe entrar automáticamente en modo edición
// cuando hay un solo campo y este es editable
func (h *DevTUI) checkAutoEditMode() {
	currentTab := &h.tabSections[h.activeTab]

	// Entrar automáticamente en modo edición si hay un solo campo editable
	if len(currentTab.FieldHandlers) == 1 && currentTab.FieldHandlers[0].Editable {
		h.tabEditingConfig = true
		currentTab.indexActiveEditField = 0
	} else {
		// Si hay múltiples campos, no entrar en modo edición automáticamente
		h.tabEditingConfig = false
	}
}
