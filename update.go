package devtui

import (
	"fmt"
	"strings"
	"time"

	. "github.com/cdvelop/messagetype"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (cf *SectionField) SetCursorAtEnd() {
	cf.cursor = len(cf.Value)
}

// listenToMessages crea un comando para escuchar mensajes del canal
func (h *DevTUI) listenToMessages() tea.Cmd {
	return func() tea.Msg {
		msg := <-h.tabContentsChan
		return channelMsg(msg)
	}
}

// tickEverySecond crea un comando para actualizar el tiempo
func (h *DevTUI) tickEverySecond() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update maneja las actualizaciones del estado
func (h *DevTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg: // Al presionar una tecla
		if h.tabEditingConfig { // EDITING CONFIG IN SECTION

			currentTab := &h.TabSections[h.activeTab]

			currentField := &h.TabSections[h.activeTab].SectionFields[currentTab.indexActiveEditField]

			if currentField.Editable { // Si el campo es editable, permitir la edición

				switch msg.String() {
				case "enter": // Al presionar ENTER, guardamos los cambios o ejecutamos la acción
					if _, err := currentField.FieldValueChange(currentField.Value); err != nil {
						h.addTerminalPrint(Error, fmt.Sprintf("Error updating field: %v %v", currentField.Name, err))
					}
					// return the cursor to its position in the field
					currentField.SetCursorAtEnd()
					h.tabEditingConfig = false
					return h, nil
				case "esc": // Al presionar ESC, descartamos los cambios
					currentField := &h.TabSections[h.activeTab].SectionFields[currentTab.indexActiveEditField]
					// currentField.Value = GetConfigFields()[currentTab.indexActiveEditField].value // Restaurar valor original

					// volvemos el cursor a su posición
					currentField.SetCursorAtEnd()

					h.tabEditingConfig = false
					h.addTerminalPrint(OK, "Exited config editing mode")
					return h, nil
				case "left": // Mover el cursor a la izquierda
					currentField := &h.TabSections[h.activeTab].SectionFields[currentTab.indexActiveEditField]
					if currentField.cursor > 0 {
						currentField.cursor--
					}
				case "right": // Mover el cursor a la derecha
					currentField := &h.TabSections[h.activeTab].SectionFields[currentTab.indexActiveEditField]
					if currentField.cursor < len(currentField.Value) {
						currentField.cursor++
					}
				default:
					currentField := &h.TabSections[h.activeTab].SectionFields[currentTab.indexActiveEditField]
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

		} else {

			switch msg.String() {
			case "up": // Mover hacia arriba el indice del campo activo
				currentTab := &h.TabSections[h.activeTab]

				if currentTab.indexActiveEditField > 0 {
					currentTab.indexActiveEditField--
				}
			case "down": // Mover hacia abajo el indice del campo activo
				currentTab := &h.TabSections[h.activeTab]
				if currentTab.indexActiveEditField < len(h.TabSections[0].SectionFields)-1 {
					currentTab.indexActiveEditField++
				}

			case "tab": // change tabSection
				h.activeTab = (h.activeTab + 1) % len(h.TabSections)
				h.cancelEditingConfig(true)
				h.updateViewport()
			case "shift+tab": // change tabSection
				h.cancelEditingConfig(true)
				h.activeTab = (h.activeTab - 1 + len(h.TabSections)) % len(h.TabSections)
				h.updateViewport()
			case "ctrl+l":
				// h.TabSections[h.activeTab].tabContents = []tabContent{}
			case "enter":
				if h.tabEditingConfig {
					h.tabEditingConfig = false
					h.addTerminalPrint(Warning, "Exited config editing mode")
				} else {
					h.tabEditingConfig = true
					h.addTerminalPrint(Warning, "Entered config editing mode")
				}

				h.updateViewport()
			case "ctrl+c":
				close(h.ExitChan) // Cerrar el canal para señalizar a todas las goroutines
				return h, tea.Quit
			default:

			}
		}
	case channelMsg:
		// Start listening for new messages again after processing the current one
		cmds = append(cmds, h.listenToMessages())

		// Convert the channel message to a tabContent type
		tc := tabContent(msg)

		// Only update the viewport if the message belongs to the currently active tab
		if tc.tabSection.index == h.activeTab {
			h.updateViewport()
		}

	case tea.WindowSizeMsg:

		headerHeight := lipgloss.Height(h.headerView())
		footerHeight := lipgloss.Height(h.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !h.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			h.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			h.viewport.YPosition = headerHeight
			h.viewport.SetContent(h.ContentView())
			h.ready = true
		} else {
			h.viewport.Width = msg.Width
			h.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tickMsg:
		h.currentTime = time.Now().Format("15:04:05")
		cmds = append(cmds, h.tickEverySecond())
	}
	// Handle keyboard and mouse events in the viewport
	h.viewport, cmd = h.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return h, tea.Batch(cmds...)
}

// Write implementa io.Writer para capturar la salida de otros procesos
func (ts *TabSection) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		// Detectar automáticamente el tipo de mensaje
		msgType := DetectMessageType(msg)

		ts.tui.sendMessage(msg, msgType, ts)
		// Si es un error, escribirlo en el archivo de log
		if msgType == Error {
			ts.tui.LogToFile(msg)
		}

	}
	return len(p), nil
}

func (h *DevTUI) updateViewport() {
	h.viewport.SetContent(h.ContentView())
	h.viewport.GotoBottom()
}

func (h *DevTUI) cancelEditingConfig(cancel bool) {
	if cancel {
		h.tabEditingConfig = false
		h.addTerminalPrint(OK, "Exited config editing mode")
	} else {
		h.tabEditingConfig = true
		h.addTerminalPrint(Warning, "Entered config editing mode")
	}
}

// Add this helper function
func (h *DevTUI) addTerminalPrint(msgType MessageType, content string) {
	h.TabSections[h.activeTab].tabContents = append(
		h.TabSections[h.activeTab].tabContents,
		tabContent{
			Type:    msgType,
			Content: content,
			Time:    time.Now(),
		},
	)
}
