package main

import (
	"sync"

	"github.com/cdvelop/devtui"
	"github.com/cdvelop/devtui/example"
)

func main() {
	// Use common configuration
	config := example.CreateTestConfig()
	tui := devtui.NewTUI(config)

	// Setup handlers and tabs using common setup
	example.SetupHandlersAndTabs(tui)

	// Start the UI
	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
