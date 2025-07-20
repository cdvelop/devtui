package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
	"github.com/cdvelop/devtui/examples/handlers"
)

func main() {
	config := &devtui.TuiConfig{
		AppName:       "DevTUI - Complete Async Demo",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#1A1A1A",
			Highlight:  "#00D7FF",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"DevTUI Demo:"}, messages...)...)
		},
	}

	tui := devtui.NewTUI(config)

	fmt.Printf("TUI created successfully\n")

	// === Tab 1: Database Configuration ===
	dbHost := handlers.NewDatabaseHostHandler("localhost")
	pgPort := handlers.NewDatabasePortHandler("postgres")
	mysqlPort := handlers.NewDatabasePortHandler("mysql")
	mongoPort := handlers.NewDatabasePortHandler("mongodb")

	fmt.Printf("Database handlers created:\n")
	fmt.Printf("  - %s: %s\n", dbHost.Label(), dbHost.Value())
	fmt.Printf("  - %s: %s\n", pgPort.Label(), pgPort.Value())
	fmt.Printf("  - %s: %s\n", mysqlPort.Label(), mysqlPort.Value())
	fmt.Printf("  - %s: %s\n", mongoPort.Label(), mongoPort.Value())

	tab1 := tui.NewTabSection("Database", "Database connection configuration")
	tab1.NewField(dbHost).
		NewField(pgPort).
		NewField(mysqlPort).
		NewField(mongoPort)

	fmt.Printf("Database tab created with fields\n")

	// === Tab 2: Health Monitoring ===
	apiHealth := handlers.NewHealthCheckHandler("API Server", "http://localhost:8080/health")
	dbHealth := handlers.NewHealthCheckHandler("Database", "http://localhost:5432/ping")
	cacheHealth := handlers.NewHealthCheckHandler("Redis Cache", "http://localhost:6379/ping")

	tui.NewTabSection("Health", "Service health monitoring").
		NewField(apiHealth).
		NewField(dbHealth).
		NewField(cacheHealth)

	// === Tab 3: Build System ===
	goProd := handlers.NewAdvancedBuildHandler("production", "go")
	goDev := handlers.NewAdvancedBuildHandler("development", "go")
	rustProd := handlers.NewAdvancedBuildHandler("production", "rust")
	nodeProd := handlers.NewAdvancedBuildHandler("production", "node")

	tui.NewTabSection("Build", "Multi-language build system").
		NewField(goProd).
		NewField(goDev).
		NewField(rustProd).
		NewField(nodeProd)

	// === Tab 4: CI/CD Pipeline ===
	dockerLinux := handlers.NewDockerBuildHandler("myapp:latest", "linux/amd64")
	dockerArm := handlers.NewDockerBuildHandler("myapp:latest", "linux/arm64")
	secScanDep := handlers.NewSecurityScanHandler("dependency")
	secScanSAST := handlers.NewSecurityScanHandler("sast")

	tui.NewTabSection("CI/CD", "Continuous Integration pipeline").
		NewField(dockerLinux).
		NewField(dockerArm).
		NewField(secScanDep).
		NewField(secScanSAST)

	// === Tab 5: Load Testing ===
	lightLoad := handlers.NewLoadTestHandler("Light", 5*time.Second)
	mediumLoad := handlers.NewLoadTestHandler("Medium", 15*time.Second)
	heavyLoad := handlers.NewLoadTestHandler("Heavy", 30*time.Second)
	stressLoad := handlers.NewLoadTestHandler("Stress", 60*time.Second)

	tui.NewTabSection("Load Test", "Performance and stress testing").
		NewField(lightLoad).
		NewField(mediumLoad).
		NewField(heavyLoad).
		NewField(stressLoad)

	fmt.Println("Starting DevTUI Complete Demo...")
	fmt.Println("Features demonstrated:")
	fmt.Println("  • Async field operations with spinners")
	fmt.Println("  • Configurable timeouts per operation")
	fmt.Println("  • Network validation and health checks")
	fmt.Println("  • Build system integration")
	fmt.Println("  • CI/CD pipeline operations")
	fmt.Println("  • Load testing with different durations")
	fmt.Println("  • Error handling and progress feedback")
	fmt.Println()

	// Start the TUI
	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
