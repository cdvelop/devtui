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
	lastOpID    string
}

// WritingHandler implementation
func (h *HostConfigHandler) Name() string                 { return "HostConfigHandler" }
func (h *HostConfigHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *HostConfigHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *HostConfigHandler) Label() string          { return "Host" }
func (h *HostConfigHandler) Value() string          { return h.currentHost }
func (h *HostConfigHandler) Editable() bool         { return true }
func (h *HostConfigHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *HostConfigHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
	host := strings.TrimSpace(newValue.(string))
	if host == "" {
		return "", fmt.Errorf("host cannot be empty")
	}

	// Use progress callback if available
	if len(progress) > 0 {
		progressCallback := progress[0]
		progressCallback("Validating host configuration...")
		time.Sleep(500 * time.Millisecond)
		progressCallback("Checking network connectivity...", 50.0)
		time.Sleep(500 * time.Millisecond)
		progressCallback("Host validation complete", 100.0)
	} else {
		// Fallback - simulate network validation
		time.Sleep(1 * time.Second)
	}

	h.currentHost = host
	return fmt.Sprintf("Host configured: %s", host), nil
}

type PortConfigHandler struct {
	currentPort string
	lastOpID    string
}

// WritingHandler implementation
func (h *PortConfigHandler) Name() string                 { return "PortConfig" }
func (h *PortConfigHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *PortConfigHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *PortConfigHandler) Label() string          { return "Port" }
func (h *PortConfigHandler) Value() string          { return h.currentPort }
func (h *PortConfigHandler) Editable() bool         { return true }
func (h *PortConfigHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *PortConfigHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
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

	// Use progress callback if available
	if len(progress) > 0 {
		progressCallback := progress[0]
		progressCallback("Validating port number...")
		time.Sleep(300 * time.Millisecond)
		progressCallback("Checking port availability...", 60.0)
		time.Sleep(400 * time.Millisecond)
		progressCallback("Port validation complete", 100.0)
	}

	h.currentPort = portStr
	return fmt.Sprintf("Port configured: %d", port), nil
}

type BuildActionHandler struct {
	buildType string
	lastOpID  string
}

func NewBuildActionHandler(buildType string) *BuildActionHandler {
	return &BuildActionHandler{buildType: buildType}
}

// WritingHandler implementation
func (h *BuildActionHandler) Name() string                 { return fmt.Sprintf("Build_%s", h.buildType) }
func (h *BuildActionHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *BuildActionHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *BuildActionHandler) Label() string          { return fmt.Sprintf("Build %s", h.buildType) }
func (h *BuildActionHandler) Value() string          { return "Press Enter to build" }
func (h *BuildActionHandler) Editable() bool         { return false }
func (h *BuildActionHandler) Timeout() time.Duration { return 30 * time.Second }

func (h *BuildActionHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
	// Use progress callback if available
	if len(progress) > 0 {
		progressCallback := progress[0]

		progressCallback(fmt.Sprintf("Initiating %s build...", h.buildType))
		time.Sleep(500 * time.Millisecond)

		progressCallback("Checking dependencies...", 25.0)
		time.Sleep(1 * time.Second)

		progressCallback("Compiling source code...", 60.0)
		if h.buildType == "Production" {
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}

		progressCallback("Generating artifacts...", 85.0)
		time.Sleep(500 * time.Millisecond)

		progressCallback("Build finalization...", 100.0)
	} else {
		// Fallback - simulate build process
		var buildDuration time.Duration
		if h.buildType == "Production" {
			buildDuration = 4 * time.Second
		} else {
			buildDuration = 2 * time.Second
		}
		time.Sleep(buildDuration)
	}

	return fmt.Sprintf("%s build completed successfully", h.buildType), nil
}

type DeployActionHandler struct {
	environment string
	lastOpID    string
}

func NewDeployActionHandler(env string) *DeployActionHandler {
	return &DeployActionHandler{environment: env}
}

// WritingHandler implementation
func (h *DeployActionHandler) Name() string {
	return fmt.Sprintf("Deploy_%s", h.environment)
}
func (h *DeployActionHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *DeployActionHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *DeployActionHandler) Label() string          { return fmt.Sprintf("Deploy to %s", h.environment) }
func (h *DeployActionHandler) Value() string          { return "Press Enter to deploy" }
func (h *DeployActionHandler) Editable() bool         { return false }
func (h *DeployActionHandler) Timeout() time.Duration { return 45 * time.Second }

func (h *DeployActionHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
	// Use progress callback if available
	if len(progress) > 0 {
		progressCallback := progress[0]

		progressCallback(fmt.Sprintf("Initiating deployment to %s...", h.environment))
		time.Sleep(500 * time.Millisecond)

		progressCallback("Preparing deployment package...", 20.0)
		time.Sleep(1 * time.Second)

		progressCallback("Uploading to servers...", 40.0)
		if h.environment == "Production" {
			time.Sleep(3 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}

		progressCallback("Configuring services...", 70.0)
		time.Sleep(1 * time.Second)

		progressCallback("Running health checks...", 90.0)
		time.Sleep(800 * time.Millisecond)

		progressCallback("Deployment completed", 100.0)
	} else {
		// Fallback - simulate deployment process
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
	}

	return fmt.Sprintf("Successfully deployed to %s", h.environment), nil
}

type WelcomeHandler struct {
	lastOpID string
}

// WritingHandler implementation
func (h *WelcomeHandler) Name() string                 { return "WelcomeHandler" }
func (h *WelcomeHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *WelcomeHandler) GetLastOperationID() string   { return h.lastOpID }

// FieldHandler implementation
func (h *WelcomeHandler) Label() string          { return "DevTUI Features" }
func (h *WelcomeHandler) Value() string          { return "Press Enter to view features" }
func (h *WelcomeHandler) Editable() bool         { return false }
func (h *WelcomeHandler) Timeout() time.Duration { return 1 * time.Second }

func (h *WelcomeHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
	// Simple handler - no progress needed since it's instant
	return "DevTUI Features:\n• Async operations with dynamic progress messages\n• Configurable timeouts\n• Error handling\n• Real-time progress feedback\n• Handler-based architecture", nil
}

func main() {
	config := &devtui.TuiConfig{
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

	tui := devtui.NewTUI(config)

	// Create handlers
	welcomeHandler := &WelcomeHandler{}
	hostHandler := &HostConfigHandler{currentHost: "localhost"}
	portHandler := &PortConfigHandler{currentPort: "8080"}
	prodBuild := NewBuildActionHandler("Production")
	devBuild := NewBuildActionHandler("Development")
	stagingDeploy := NewDeployActionHandler("Staging")
	prodDeploy := NewDeployActionHandler("Production")

	// Configure tabs with new handler-based API
	tui.NewTabSection("Welcome", "DevTUI Demo Features").
		NewField(welcomeHandler)

	tui.NewTabSection("Server", "Server configuration").
		NewField(hostHandler).
		NewField(portHandler)

	tui.NewTabSection("Build", "Build operations").
		NewField(prodBuild).
		NewField(devBuild)

	tui.NewTabSection("Deploy", "Deployment operations").
		NewField(stagingDeploy).
		NewField(prodDeploy)

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.Start(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()
}
