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

	// Create and add custom tabs (recommended way)
	mainTab := tui.NewTabSection("Main", "")
	mainTab.AddFields(
		*devtui.NewField(
			"Username",
			"",
			true,
			func(newValue any) (string, error) {
				strValue := newValue.(string)
				if len(strValue) < 5 {
					return "", fmt.Errorf("username must be at least 5 characters")
				}
				return "Username updated to " + strValue, nil
			},
		),
	)
	tui.AddTabSections(mainTab)

	// Start the TUI
	if err := tui.Start(); err != nil {
		panic(err)
	}
}
```

## Keyboard Shortcuts

| Key          | Action                                 |
|--------------|----------------------------------------|
| Tab          | Switch to next tab                     |
| Shift+Tab    | Switch to previous tab                 |
| Left/Right   | Navigate between fields in current tab |
| Enter        | Edit field or execute action           |
| Esc          | Cancel editing                         |
| Ctrl+C       | Quit application                       |

## NewTabSection Method

```go
// NewTabSection creates a new TabSection with the given title and footer
// Example:

	tab := tui.NewTabSection("BUILD", "Press 't' to compile")
	// Preferred way to add fields (variadic parameters)
	tab.AddFields(
		*NewField(
			"Username",
			"",
			true,
			func(newValue any) (string, error) {
				strValue := newValue.(string)
				if len(strValue) < 5 {
					return "", fmt.Errorf("username must be at least 5 characters")
				}
				return "Username updated to " + strValue, nil
			},
		),
		*NewField(
			"Password",
			"",
			true,
			func(newValue any) (string, error) {
				strValue := newValue.(string)
				if len(strValue) < 8 {
					return "", fmt.Errorf("password must be at least 8 characters")
				}
				return "Password updated", nil
			},
		),
	)


//	 Get/Set title and footer
	currentTitle := tab.Title()
	tab.SetTitle("New Title")

	currentFooter := tab.Footer() 
	tab.SetFooter("New Footer")
```

## Field Types

### Editable Fields with Different Data Types

#### String Field
```go
*NewField(
	"Username", 
	"defaultUser", 
	true, 
	func(value any) (string, error) {
		strValue := value.(string)
		if len(strValue) < 3 {
			return "", fmt.Errorf("username must be at least 3 characters")
		}
		return "Username updated to " + strValue, nil
	},
)
```

#### Numeric Field (Port)
```go
*NewField(
	"Server Port", 
	"8080", 
	true, 
	func(value any) (string, error) {
		switch v := value.(type) {
		case string:
			if port, err := strconv.Atoi(v); err == nil && port > 0 && port < 65536 {
				return fmt.Sprintf("Port changed to %d", port), nil
			}
			return "", fmt.Errorf("invalid port number: %s", v)
		case int:
			if v > 0 && v < 65536 {
				return fmt.Sprintf("Port changed to %d", v), nil
			}
			return "", fmt.Errorf("port out of range: %d", v)
		default:
			return "", fmt.Errorf("unsupported type for port: %T", v)
		}
	},
)
```

#### Boolean Field
```go
*NewField(
	"Debug Mode", 
	"false", 
	true, 
	func(value any) (string, error) {
		switch v := value.(type) {
		case bool:
			return fmt.Sprintf("Debug mode set to %t", v), nil
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return fmt.Sprintf("Debug mode set to %t", b), nil
			}
			return "", fmt.Errorf("invalid boolean value: %s", v)
		default:
			return "", fmt.Errorf("unsupported type for boolean: %T", v)
		}
	},
)
```

### Non-Editable Fields (Action Buttons)
```go
*NewField(
	"Action Button", 
	"Click me", 
	false, 
	func(value any) (string, error) {
		// Execute action - value parameter is ignored for non-editable fields
		return "Action executed successfully", nil
	},
)
```

## errors
- [ ] al borrar todo un campo editable, al comenzar a escribir se copia el valor anterior al lado del cursor 

## Dependencies

- [Charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- [Charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
