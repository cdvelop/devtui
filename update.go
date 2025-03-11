package devtui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
						h.addTerminalPrint(ErrorMsg, fmt.Sprintf("Error updating field: %v %v", currentField.Name, err))
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
					h.addTerminalPrint(OkMsg, "Exited config editing mode")
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

					msgType := OkMsg
					// content eg: "Browser Opened"
					content, err := currentField.FieldValueChange(currentField.Value)
					if err != nil {
						msgType = ErrorMsg
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
				h.TabSections[h.activeTab].tabContents = []tabContent{}
			case "enter":
				if h.tabEditingConfig {
					h.tabEditingConfig = false
					h.addTerminalPrint(WarnMsg, "Exited config editing mode")
				} else {
					h.tabEditingConfig = true
					h.addTerminalPrint(WarnMsg, "Entered config editing mode")
				}

				h.updateViewport()
			case "ctrl+c":
				close(h.ExitChan) // Cerrar el canal para señalizar a todas las goroutines
				return h, tea.Quit
			default:

			}
		}
	case channelMsg:
		h.TabSections[h.activeTab].tabContents = append(h.TabSections[h.activeTab].tabContents, tabContent(msg))
		cmds = append(cmds, h.listenToMessages())

		h.updateViewport()

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
