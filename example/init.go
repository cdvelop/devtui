package example

import (
	"fmt"

	"github.com/cdvelop/devtui"
)

// CreateTestConfig creates the standard configuration used for both manual testing and automated tests
func CreateTestConfig() *devtui.TuiConfig {
	return &devtui.TuiConfig{
		AppName:       "DevTUI",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"DevTUI Log:"}, messages...)...)
		},
	}
}

// SetupHandlersAndTabs configures the TUI with the standard handlers for testing
func SetupHandlersAndTabs(tui *devtui.DevTUI) {
	// Create only the essential handlers - 3 different types
	welcomeHandler := NewWelcomeHandler()               // Readonly (empty label)
	hostHandler := NewHostConfigHandler("localhost")    // Editable field
	buildHandler := NewBuildActionHandler("Production") // Action button

	// Configure tabs with handler-based API
	tui.NewTabSection("Welcome", "DevTUI Demo Features").
		NewField(welcomeHandler)

	tui.NewTabSection("Server", "Server configuration").
		NewField(hostHandler)

	tui.NewTabSection("Build", "Build operations").
		NewField(buildHandler)
}
