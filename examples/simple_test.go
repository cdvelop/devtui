package main

import (
	"fmt"
	"sync"

	"github.com/cdvelop/devtui"
	"github.com/cdvelop/devtui/examples/handlers"
)

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Simple Test",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"TEST LOG:"}, messages...)...)
		},
	}

	tui := devtui.NewTUI(config)

	// Just one simple tab for testing
	dbHost := handlers.NewDatabaseHostHandler("localhost")
	fmt.Printf("Created handler - Label: %s, Value: %s\n", dbHost.Label(), dbHost.Value())

	tab := tui.NewTabSection("Test", "Simple test")
	tab.NewField(dbHost)

	fmt.Printf("Tab sections count: %d\n", len(tui.GetTabSections()))
	if len(tui.GetTabSections()) > 0 {
		fmt.Printf("First tab has %d fields\n", len(tui.GetTabSections()[0].GetFieldHandlers()))
	}

	fmt.Println("Starting simple test...")

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
