package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// === Network Configuration Handlers ===

type HostHandler struct {
	currentHost string
}

func (h *HostHandler) Label() string          { return "Host" }
func (h *HostHandler) Value() string          { return h.currentHost }
func (h *HostHandler) Editable() bool         { return true }
func (h *HostHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *HostHandler) Change(newValue any) (string, error) {
	host := strings.TrimSpace(newValue.(string))
	if host == "" {
		return "", fmt.Errorf("host cannot be empty")
	}

	// Simulate network validation
	time.Sleep(1 * time.Second)

	// Basic host validation
	if !strings.Contains(host, ".") && host != "localhost" {
		return "", fmt.Errorf("invalid host format")
	}

	h.currentHost = host
	return fmt.Sprintf("Host configured: %s", host), nil
}

type PortHandler struct {
	currentPort string
}

func (h *PortHandler) Label() string          { return "Port" }
func (h *PortHandler) Value() string          { return h.currentPort }
func (h *PortHandler) Editable() bool         { return true }
func (h *PortHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *PortHandler) Change(newValue any) (string, error) {
	portStr := strings.TrimSpace(newValue.(string))
	if portStr == "" {
		return "", fmt.Errorf("port cannot be empty")
	}

	// Simulate port validation
	time.Sleep(500 * time.Millisecond)

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

type ConnectionTestHandler struct {
	host string
	port string
}

func NewConnectionTestHandler(hostHandler *HostHandler, portHandler *PortHandler) *ConnectionTestHandler {
	return &ConnectionTestHandler{
		host: hostHandler.currentHost,
		port: portHandler.currentPort,
	}
}

func (h *ConnectionTestHandler) Label() string          { return "Test Connection" }
func (h *ConnectionTestHandler) Value() string          { return "Press Enter to test" }
func (h *ConnectionTestHandler) Editable() bool         { return false }
func (h *ConnectionTestHandler) Timeout() time.Duration { return 10 * time.Second }

func (h *ConnectionTestHandler) Change(newValue any) (string, error) {
	if h.host == "" || h.port == "" {
		return "", fmt.Errorf("configure host and port first")
	}

	// Simulate connection test
	target := fmt.Sprintf("http://%s:%s", h.host, h.port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return "", fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	return fmt.Sprintf("Connection successful! Status: %s", resp.Status), nil
}

// === Build System Handlers ===

type BuildHandler struct {
	projectPath string
	buildType   string
}

func NewBuildHandler(buildType string) *BuildHandler {
	return &BuildHandler{
		projectPath: "./",
		buildType:   buildType,
	}
}

func (h *BuildHandler) Label() string          { return fmt.Sprintf("Build %s", h.buildType) }
func (h *BuildHandler) Value() string          { return "Press Enter to build" }
func (h *BuildHandler) Editable() bool         { return false }
func (h *BuildHandler) Timeout() time.Duration { return 30 * time.Second }

func (h *BuildHandler) Change(newValue any) (string, error) {
	var cmd *exec.Cmd

	switch h.buildType {
	case "Production":
		cmd = exec.Command("go", "build", "-ldflags", "-s -w", ".")
	case "Development":
		cmd = exec.Command("go", "build", ".")
	case "Test":
		cmd = exec.Command("go", "test", "./...")
	default:
		return "", fmt.Errorf("unknown build type: %s", h.buildType)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}

	return fmt.Sprintf("%s build completed successfully", h.buildType), nil
}

type CleanHandler struct{}

func (h *CleanHandler) Label() string          { return "Clean Build" }
func (h *CleanHandler) Value() string          { return "Press Enter to clean" }
func (h *CleanHandler) Editable() bool         { return false }
func (h *CleanHandler) Timeout() time.Duration { return 10 * time.Second }

func (h *CleanHandler) Change(newValue any) (string, error) {
	cmd := exec.Command("go", "clean", "-cache", "-modcache", "-testcache")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("clean failed: %v\nOutput: %s", err, string(output))
	}

	return "Build cache cleaned successfully", nil
}

// === Deployment Handlers ===

type DeployHandler struct {
	environment string
}

func NewDeployHandler(env string) *DeployHandler {
	return &DeployHandler{environment: env}
}

func (h *DeployHandler) Label() string  { return fmt.Sprintf("Deploy to %s", h.environment) }
func (h *DeployHandler) Value() string  { return "Press Enter to deploy" }
func (h *DeployHandler) Editable() bool { return false }
func (h *DeployHandler) Timeout() time.Duration {
	// Different timeouts based on environment
	switch h.environment {
	case "Staging":
		return 1 * time.Minute
	case "Production":
		return 3 * time.Minute
	default:
		return 30 * time.Second
	}
}

func (h *DeployHandler) Change(newValue any) (string, error) {
	// Simulate deployment process
	steps := []string{
		"Preparing deployment package",
		"Uploading to %s environment",
		"Running deployment scripts",
		"Verifying deployment",
		"Updating load balancer",
	}

	for _, step := range steps {
		// Each step takes some time
		time.Sleep(time.Duration(len(steps)) * 200 * time.Millisecond)
	}

	return fmt.Sprintf("Successfully deployed to %s environment", h.environment), nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Complete Dev Environment",
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

	// Network Configuration Tab
	hostHandler := &HostHandler{currentHost: "localhost"}
	portHandler := &PortHandler{currentPort: "8080"}
	connectionTest := NewConnectionTestHandler(hostHandler, portHandler)

	tui.NewTabSection("Network", "Network configuration and testing").
		NewField(hostHandler).
		NewField(portHandler).
		NewField(connectionTest)

	// Build System Tab
	prodBuild := NewBuildHandler("Production")
	devBuild := NewBuildHandler("Development")
	testBuild := NewBuildHandler("Test")
	cleanHandler := &CleanHandler{}

	tui.NewTabSection("Build", "Build system operations").
		NewField(prodBuild).
		NewField(devBuild).
		NewField(testBuild).
		NewField(cleanHandler)

	// Deployment Tab
	stagingDeploy := NewDeployHandler("Staging")
	prodDeploy := NewDeployHandler("Production")

	tui.NewTabSection("Deploy", "Deployment operations").
		NewField(stagingDeploy).
		NewField(prodDeploy)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
