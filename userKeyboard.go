package devtui

import (
	"fmt"

	"slices"

	"github.com/cdvelop/messagetype"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleKeyboard processes keyboard input and updates the model state
// returns whether the update function should continue processing or return early
func (h *DevTUI) HandleKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	if h.editModeActivated { // EDITING CONFIG IN SECTION
		return h.handleEditingConfigKeyboard(msg)
	} else {
		return h.handleNormalModeKeyboard(msg)
	}
}

// handleEditingConfigKeyboard handles keyboard input while in config editing mode
func (h *DevTUI) handleEditingConfigKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	currentTab := h.tabSections[h.activeTab]
	fieldHandlers := currentTab.FieldHandlers()
	currentField := fieldHandlers[currentTab.indexActiveEditField]

	if currentField.Editable() { // Si el campo es editable, permitir la edición
		// Calcular el ancho máximo disponible para el texto
		// Esto sigue la misma lógica que en footerInput.go
		_, availableTextWidth := h.calculateInputWidths(currentField.Name())

		switch msg.Type {
		case tea.KeyEnter: // Guardar cambios o ejecutar acción
			// Verificar si hubo cambios (incluyendo borrar el contenido)
			if currentField.tempEditValue != currentField.Value() {
				if currentField.changeFunc != nil {
					// Llamar changeFunc con el nuevo valor ANTES de establecerlo
					msg, err := currentField.changeFunc(currentField.tempEditValue)
					if err != nil {
						// Si hay un error, mostrarlo en la pestaña actual
						currentTab.addNewContent(messagetype.Error, fmt.Sprintf("%v %v", currentField.Name(), err))
					} else {
						// Solo aplicar el resultado de changeFunc si no hay error
						currentField.SetValue(msg)
					}
					h.editingConfigOpen(false, currentField, msg)
				} else {
					// Si no hay changeFunc, simplemente aplicar el valor directamente
					currentField.SetValue(currentField.tempEditValue)
					h.editingConfigOpen(false, currentField, "")
				}
			} else {
				// Si no hubo cambios, solo salimos del modo edición sin mostrar mensajes
				h.editingConfigOpen(false, currentField, "")
			}

			currentField.tempEditValue = "" // Limpiar el valor temporal
			h.updateViewport()              // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyEsc: // Al presionar ESC, descartamos los cambios y salimos del modo edición
			currentField.tempEditValue = "" // Limpiar el valor temporal
			h.editingConfigOpen(false, currentField, "")
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyLeft: // Mover el cursor a la izquierda dentro del texto
			if currentField.cursor > 0 {
				currentField.cursor--
			}

		case tea.KeyRight: // Mover el cursor a la derecha dentro del texto
			value := currentField.Value()
			if currentField.tempEditValue != "" {
				value = currentField.tempEditValue
			}
			if currentField.cursor < len([]rune(value)) {
				currentField.cursor++
			}

		case tea.KeyBackspace: // Borrar carácter a la izquierda
			if currentField.cursor > 0 {
				// Si aún no hay valor temporal, copiar el valor original solo la primera vez
				if currentField.tempEditValue == "" {
					currentField.tempEditValue = currentField.Value()
				}

				// Convert to runes to handle multi-byte characters correctly
				runes := []rune(currentField.tempEditValue)
				if currentField.cursor <= len(runes) {
					newRunes := slices.Delete(runes, currentField.cursor-1, currentField.cursor)
					currentField.tempEditValue = string(newRunes)
					currentField.cursor--
				}
			}

		case tea.KeySpace: // Manejar la tecla espacio como un carácter especial
			// Si aún no hay valor temporal, NO copiar el valor original automáticamente
			if currentField.tempEditValue == "" {
				currentField.tempEditValue = ""
			}

			runes := []rune(currentField.tempEditValue)
			if currentField.cursor > len(runes) {
				currentField.cursor = len(runes)
			}

			// Verificar si agregar un espacio excedería el ancho disponible
			if len(runes)+1 < availableTextWidth {
				// Insert the space at cursor position
				newRunes := make([]rune, 0, len(runes)+1)
				newRunes = append(newRunes, runes[:currentField.cursor]...)
				newRunes = append(newRunes, ' ') // Agregar el espacio
				newRunes = append(newRunes, runes[currentField.cursor:]...)
				currentField.tempEditValue = string(newRunes)
				currentField.cursor++
			}

		case tea.KeyRunes:
			// Handle normal character input - convert everything to runes for proper handling
			if len(msg.Runes) > 0 {
				// Si aún no hay valor temporal, NO copiar el valor original automáticamente
				// Solo inicializar como string vacío si está vacío
				if currentField.tempEditValue == "" {
					currentField.tempEditValue = ""
				}

				runes := []rune(currentField.tempEditValue)
				if currentField.cursor > len(runes) {
					currentField.cursor = len(runes)
				}

				// Verificar si agregar los nuevos caracteres excedería el ancho disponible
				if len(runes)+len(msg.Runes) < availableTextWidth {
					// Insert the new runes at cursor position
					newRunes := make([]rune, 0, len(runes)+len(msg.Runes))
					newRunes = append(newRunes, runes[:currentField.cursor]...)
					newRunes = append(newRunes, msg.Runes...)
					newRunes = append(newRunes, runes[currentField.cursor:]...)
					currentField.tempEditValue = string(newRunes)
					currentField.cursor += len(msg.Runes)
				}
				// Si excede el ancho, simplemente no agregar los caracteres
			}
		}
	} else { // Si el campo no es editable, solo ejecutar la acción
		switch msg.Type {
		case tea.KeyEnter:
			msgType := messagetype.Success
			// content eg: "Browser Opened"
			if currentField.changeFunc != nil {
				content, err := currentField.changeFunc(currentField.Value())
				if err != nil {
					msgType = messagetype.Error
					content = fmt.Sprintf("%s %s %s", currentField.Name(), content, err.Error())
				}
				currentField.SetValue(content)
				currentTab.addNewContent(msgType, content)
			}
			h.editModeActivated = false
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil

		case tea.KeyEsc: // Permitir también salir con ESC para campos no editables
			h.editingConfigOpen(false, currentField, "")
			h.updateViewport() // Asegurar que se actualice la vista para mostrar el mensaje
			return false, nil
		}
	}

	return true, nil
}

// handleNormalModeKeyboard handles keyboard input in normal mode (not editing config)
func (h *DevTUI) handleNormalModeKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
	currentTab := h.tabSections[h.activeTab]
	fieldHandlers := currentTab.FieldHandlers()
	totalFields := len(fieldHandlers)

	switch msg.Type {
	case tea.KeyUp, tea.KeyDown:
		// Las teclas arriba y abajo ya no modifican el campo activo
		// Solo controlarán el desplazamiento del viewport
		// No hacemos nada aquí para permitir que el manejo del viewport siga su curso normal

	case tea.KeyLeft: // Navegar al campo anterior (ciclo continuo)
		if totalFields > 0 {
			currentTab.indexActiveEditField = (currentTab.indexActiveEditField - 1 + totalFields) % totalFields
			h.updateViewport()
			return false, nil // Detener procesamiento adicional
		}

	case tea.KeyRight: // Navegar al campo siguiente (ciclo continuo)
		if totalFields > 0 {
			currentTab.indexActiveEditField = (currentTab.indexActiveEditField + 1) % totalFields
			h.updateViewport()
			return false, nil // Detener procesamiento adicional
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
			fieldHandlers := currentTab.FieldHandlers()
			field := fieldHandlers[currentTab.indexActiveEditField]
			if !field.Editable() {
				msgType := messagetype.Success
				if field.changeFunc != nil {
					content, err := field.changeFunc(field.Value())
					if err != nil {
						msgType = messagetype.Error
						content = fmt.Sprintf("%s %s %s", field.Name(), content, err.Error())
					}
					field.SetValue(content)
					currentTab.addNewContent(msgType, content)
				}
			} else {
				// Para campos editables, activar modo de edición explícitamente
				field.tempEditValue = field.Value()
				field.cursor = 0 // Asegurarnos de que el cursor comience al principio
				h.editModeActivated = true
				h.editingConfigOpen(true, field, "")
			}
			h.updateViewport()
		}

	case tea.KeyCtrlC:
		close(h.ExitChan) // Cerrar el canal para señalizar a todas las goroutines
		// Usar tea.Sequence para asegurar que ExitAltScreen se ejecute antes de Quit
		return false, tea.Sequence(tea.ExitAltScreen, tea.Quit)
	}

	return true, nil
}

// checkAutoEditMode verifica si debe entrar automáticamente en modo edición
// cuando hay un solo campo y este es editable
func (h *DevTUI) checkAutoEditMode() {
	currentTab := h.tabSections[h.activeTab]

	// Entrar automáticamente en modo edición si hay un solo campo editable
	fieldHandlers := currentTab.FieldHandlers()
	if len(fieldHandlers) == 1 && fieldHandlers[0].Editable() {
		h.editModeActivated = true
		currentTab.indexActiveEditField = 0
		// Inicializar tempEditValue y cursor
		field := fieldHandlers[0]
		field.tempEditValue = field.Value()
		field.cursor = 0
	} else {
		// Si hay múltiples campos, no entrar en modo edición automáticamente
		h.editModeActivated = false
	}
}
