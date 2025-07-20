package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// Example handlers demonstrating the new async API

type HostConfigHandler struct {
	currentHost string
}

func (h *HostConfigHandler) Label() string          { return "Host" }
func (h *HostConfigHandler) Value() string          { return h.currentHost }
func (h *HostConfigHandler) Editable() bool         { return true }
func (h *HostConfigHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *HostConfigHandler) Change(newValue any) (string, error) {
	host := strings.TrimSpace(newValue.(string))
	if host == "" {
		return "", fmt.Errorf("host cannot be empty")
	}

	// Simulate network validation
	time.Sleep(1 * time.Second)

	h.currentHost = host
	return fmt.Sprintf("Host configured: %s", host), nil
}

type PortConfigHandler struct {
	currentPort string
}

func (h *PortConfigHandler) Label() string          { return "Port" }
func (h *PortConfigHandler) Value() string          { return h.currentPort }
func (h *PortConfigHandler) Editable() bool         { return true }
func (h *PortConfigHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *PortConfigHandler) Change(newValue any) (string, error) {
	portStr := strings.TrimSpace(newValue.(string))
	if portStr == "" {
		return "", fmt.Errorf("port cannot be empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("port must be between 1 and 65535")
	}

	h.currentPort = portStr
	return fmt.Sprintf("Port configured: %d", port), nil
}

type BuildActionHandler struct {
	buildType string
}

func NewBuildActionHandler(buildType string) *BuildActionHandler {
	return &BuildActionHandler{buildType: buildType}
}

func (h *BuildActionHandler) Label() string          { return fmt.Sprintf("Build %s", h.buildType) }
func (h *BuildActionHandler) Value() string          { return "Press Enter to build" }
func (h *BuildActionHandler) Editable() bool         { return false }
func (h *BuildActionHandler) Timeout() time.Duration { return 30 * time.Second }

func (h *BuildActionHandler) Change(newValue any) (string, error) {
	// Simulate build process
	var buildDuration time.Duration
	if h.buildType == "Production" {
		buildDuration = 4 * time.Second
	} else {
		buildDuration = 2 * time.Second
	}

	time.Sleep(buildDuration)
	return fmt.Sprintf("%s build completed successfully", h.buildType), nil
}

type DeployActionHandler struct {
	environment string
}

func NewDeployActionHandler(env string) *DeployActionHandler {
	return &DeployActionHandler{environment: env}
}

func (h *DeployActionHandler) Label() string          { return fmt.Sprintf("Deploy to %s", h.environment) }
func (h *DeployActionHandler) Value() string          { return "Press Enter to deploy" }
func (h *DeployActionHandler) Editable() bool         { return false }
func (h *DeployActionHandler) Timeout() time.Duration { return 45 * time.Second }

func (h *DeployActionHandler) Change(newValue any) (string, error) {
	// Simulate deployment process
	var deployDuration time.Duration
	switch h.environment {
	case "Staging":
		deployDuration = 3 * time.Second
	case "Production":
		deployDuration = 8 * time.Second
	default:
		deployDuration = 2 * time.Second
	}

	time.Sleep(deployDuration)
	return fmt.Sprintf("Successfully deployed to %s", h.environment), nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "DevTUI - New Async API Demo",
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

	tui := devtui.NewTUI(config)

	// Create handlers
	hostHandler := &HostConfigHandler{currentHost: "localhost"}
	portHandler := &PortConfigHandler{currentPort: "8080"}
	prodBuild := NewBuildActionHandler("Production")
	devBuild := NewBuildActionHandler("Development")
	stagingDeploy := NewDeployActionHandler("Staging")
	prodDeploy := NewDeployActionHandler("Production")

	// Configure tabs with new handler-based API
	tui.NewTabSection("Server", "Server configuration").
		NewField(hostHandler).
		NewField(portHandler)

	tui.NewTabSection("Build", "Build operations").
		NewField(prodBuild).
		NewField(devBuild)

	tui.NewTabSection("Deploy", "Deployment operations").
		NewField(stagingDeploy).
		NewField(prodDeploy)

	fmt.Println("Starting DevTUI with New Async API...")
	fmt.Println("Features:")
	fmt.Println("  • Async operations with spinners")
	fmt.Println("  • Configurable timeouts")
	fmt.Println("  • Error handling")
	fmt.Println("  • Progress feedback")
	fmt.Println()

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.Start(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()

	fmt.Println("Aplicación finalizada correctamente")
}
