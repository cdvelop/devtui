# DevTUI

Interactive Terminal User Interface library for Go applications development (principal tui in [**GoDEV** App](https://github.com/cdvelop/godev))

## Features

- Tab-based interface organization
- Editable and non-editable fields
- Keyboard navigation (Tab, Shift+Tab, Left/Right arrows)
- Field validation and callbacks
- Customizable styles and colors

![devtui](tui.jpg)


## Basic Usage

```go
package main

import (
	"fmt"
	"github.com/cdvelop/devtui"
)

func main() {
	// Configuration
	config := &devtui.TuiConfig{
		AppName:  "MyApp", 
		ExitChan: make(chan bool),
		Config: devtui.Config{
			ForeGround:"#F4F4F4",
			Background:"#000000",
			Highlight: "#FF6600",
			Lowlight:  "#666666",
		},
		LogToFile: func(messageErr any) {
			// Implement your logging logic here
		},
	}

	// Create new TUI instance
	tui := devtui.NewTUI(config)

	// Add custom tabs if needed (optional)
	tui.AddTabSections(devtui.TabSection{
		Title: "Main",
		FieldHandlers: []devtui.Field{
			*devtui.NewField(
				"Username",
				"",
				true,
				func(newValue string) (string, error) {
					if len(newValue) < 5 {
						return "", fmt.Errorf("username must be at least 5 characters")
					}
					return newValue, nil
				},
			),
		},
	})

	// Start the TUI
	if err := tui.Start(); err != nil {
		panic(err)
	}
}
```

## Keyboard Shortcuts

| Key          | Action                                  |
|--------------|----------------------------------------|
| Tab          | Switch to next tab                     |
| Shift+Tab    | Switch to previous tab                 |
| Left/Right   | Navigate between fields in current tab |
| Enter        | Edit field or execute action           |
| Esc          | Cancel editing                         |
| Ctrl+C       | Quit application                       |

## Field Types

### Editable Fields
```go
*NewField(
	"Field Name", 
	"initial value", 
	true, 
	func(value string) (string, error) {
		// Validate or process the new value
		return value, nil
	},
)
```

### Non-Editable Fields (Action Buttons)
```go
*NewField(
	"Action Button", 
	"Click me", 
	false, 
	func(value string) (string, error) {
		// Execute action
		return "Action executed", nil
	},
)
```

## Dependencies

- [Charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- [Charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
