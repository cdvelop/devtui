# Advanced Examples

This document contains comprehensive examples demonstrating advanced DevTUI features and patterns.

## Multiple Field Types with Different Behaviors

### Network Configuration Handlers

```go
// Network configuration handlers
type HostHandler struct {
    currentHost string
    lastOpID    string  // For WritingHandler interface
}

// WritingHandler implementation
func (h *HostHandler) Name() string { return "HostHandler" }
func (h *HostHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *HostHandler) GetLastOperationID() string { return h.lastOpID }

// FieldHandler implementation
func (h *HostHandler) Label() string { return "Host" }
func (h *HostHandler) Value() string { return h.currentHost }
func (h *HostHandler) Editable() bool { return true }
func (h *HostHandler) Timeout() time.Duration { return 5 * time.Second }
func (h *HostHandler) Change(newValue any) (string, error) {
    host := strings.TrimSpace(newValue.(string))
    if host == "" {
        return "", fmt.Errorf("host cannot be empty")
    }
    
    // Network validation - async with timeout
    time.Sleep(1 * time.Second)
    h.currentHost = host
    return fmt.Sprintf("Host configured: %s", host), nil
}

type PortHandler struct {
    currentPort string
    lastOpID    string  // For WritingHandler interface
}

// WritingHandler implementation
func (h *PortHandler) Name() string { return "PortHandler" }
func (h *PortHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *PortHandler) GetLastOperationID() string { return h.lastOpID }

// FieldHandler implementation
func (h *PortHandler) Label() string { return "Port" }
func (h *PortHandler) Value() string { return h.currentPort }
func (h *PortHandler) Editable() bool { return true }
func (h *PortHandler) Timeout() time.Duration { return 0 } // No timeout
func (h *PortHandler) Change(newValue any) (string, error) {
    portStr := newValue.(string)
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return "", fmt.Errorf("port must be a number")
    }
    if port < 1 || port > 65535 {
        return "", fmt.Errorf("port must be between 1 and 65535")
    }
    h.currentPort = portStr
    return fmt.Sprintf("Port set to: %d", port), nil
}

// Usage
hostHandler := &HostHandler{currentHost: "localhost"}
portHandler := &PortHandler{currentPort: "8080"}

tui.NewTabSection("Network", "Server configuration").
    NewField(hostHandler).
    NewField(portHandler)
```

## Long-Running Operations with Progress Feedback

### CI/CD Pipeline Handlers

```go
// CI/CD Pipeline handler with timeout
type DeploymentHandler struct {
    environment string
    version     string
    lastOpID    string  // For WritingHandler interface
}

// WritingHandler implementation
func (h *DeploymentHandler) Name() string { return "DeploymentHandler" }
func (h *DeploymentHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *DeploymentHandler) GetLastOperationID() string { return h.lastOpID }

// FieldHandler implementation
func (h *DeploymentHandler) Label() string { 
    return fmt.Sprintf("Deploy v%s", h.version) 
}
func (h *DeploymentHandler) Value() string { 
    return fmt.Sprintf("Deploy to %s", h.environment) 
}
func (h *DeploymentHandler) Editable() bool { return false }
func (h *DeploymentHandler) Timeout() time.Duration { return 2 * time.Minute }
func (h *DeploymentHandler) Change(newValue any) (string, error) {
    // Long deployment process - DevTUI shows spinner automatically
    time.Sleep(5 * time.Second) // Simulate deployment
    return fmt.Sprintf("Successfully deployed v%s to %s", h.version, h.environment), nil
}

// Health check handler with shorter timeout
type HealthCheckHandler struct {
    serviceName string
    lastOpID    string  // For WritingHandler interface
}

// WritingHandler implementation
func (h *HealthCheckHandler) Name() string { return "HealthCheckHandler" }
func (h *HealthCheckHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *HealthCheckHandler) GetLastOperationID() string { return h.lastOpID }

// FieldHandler implementation
func (h *HealthCheckHandler) Label() string { return "Health Check" }
func (h *HealthCheckHandler) Value() string { 
    return fmt.Sprintf("Check %s status", h.serviceName) 
}
func (h *HealthCheckHandler) Editable() bool { return false }
func (h *HealthCheckHandler) Timeout() time.Duration { return 10 * time.Second }
func (h *HealthCheckHandler) Change(newValue any) (string, error) {
    // Simulate health check API call
    time.Sleep(2 * time.Second)
    return fmt.Sprintf("%s is healthy and responding", h.serviceName), nil
}

// Usage
deployHandler := &DeploymentHandler{environment: "production", version: "1.2.3"}
healthHandler := &HealthCheckHandler{serviceName: "API Service"}

tui.NewTabSection("Operations", "Deployment and monitoring").
    NewField(deployHandler).
    NewField(healthHandler)
```

## Multiple Tabs with Different Purposes

### Complete Application Example

```go
package main

import (
    "fmt"
    "os/exec"
    "strings"
    "time"
    "github.com/cdvelop/devtui"
)

// Database configuration handler
type DatabaseURLHandler struct {
    currentURL string
    lastOpID   string
}

func (h *DatabaseURLHandler) Name() string { return "DatabaseURLHandler" }
func (h *DatabaseURLHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *DatabaseURLHandler) GetLastOperationID() string { return h.lastOpID }
func (h *DatabaseURLHandler) Label() string { return "Database URL" }
func (h *DatabaseURLHandler) Value() string { return h.currentURL }
func (h *DatabaseURLHandler) Editable() bool { return true }
func (h *DatabaseURLHandler) Timeout() time.Duration { return 10 * time.Second }
func (h *DatabaseURLHandler) Change(newValue any) (string, error) {
    url := strings.TrimSpace(newValue.(string))
    if url == "" {
        return "", fmt.Errorf("database URL cannot be empty")
    }
    
    // Simulate database connection test
    time.Sleep(2 * time.Second)
    
    h.currentURL = url
    return "Database connection verified successfully", nil
}

// API Key handler
type APIKeyHandler struct {
    currentKey string
    lastOpID   string
}

func (h *APIKeyHandler) Name() string { return "APIKeyHandler" }
func (h *APIKeyHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *APIKeyHandler) GetLastOperationID() string { return h.lastOpID }
func (h *APIKeyHandler) Label() string { return "API Key" }
func (h *APIKeyHandler) Value() string {
    if len(h.currentKey) > 8 {
        return h.currentKey[:8] + "..."
    }
    return h.currentKey
}
func (h *APIKeyHandler) Editable() bool { return true }
func (h *APIKeyHandler) Timeout() time.Duration { return 5 * time.Second }
func (h *APIKeyHandler) Change(newValue any) (string, error) {
    key := strings.TrimSpace(newValue.(string))
    if len(key) < 16 {
        return "", fmt.Errorf("API key must be at least 16 characters")
    }
    
    h.currentKey = key
    return "API key configured successfully", nil
}

// Build project handler
type BuildProjectHandler struct {
    projectPath string
    lastOpID    string
}

func (h *BuildProjectHandler) Name() string { return "BuildProjectHandler" }
func (h *BuildProjectHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *BuildProjectHandler) GetLastOperationID() string { return h.lastOpID }
func (h *BuildProjectHandler) Label() string { return "Build Project" }
func (h *BuildProjectHandler) Value() string { return "Press Enter to build" }
func (h *BuildProjectHandler) Editable() bool { return false }
func (h *BuildProjectHandler) Timeout() time.Duration { return 60 * time.Second }
func (h *BuildProjectHandler) Change(newValue any) (string, error) {
    // Long-running build operation
    cmd := exec.Command("go", "build", h.projectPath)
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("build failed: %v", err)
    }
    return "Build completed successfully", nil
}

// Test suite handler
type TestSuiteHandler struct {
    testPath string
    lastOpID string
}

func (h *TestSuiteHandler) Name() string { return "TestSuiteHandler" }
func (h *TestSuiteHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *TestSuiteHandler) GetLastOperationID() string { return h.lastOpID }
func (h *TestSuiteHandler) Label() string { return "Run Tests" }
func (h *TestSuiteHandler) Value() string { return "Press Enter to test" }
func (h *TestSuiteHandler) Editable() bool { return false }
func (h *TestSuiteHandler) Timeout() time.Duration { return 30 * time.Second }
func (h *TestSuiteHandler) Change(newValue any) (string, error) {
    // Run test suite
    cmd := exec.Command("go", "test", h.testPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("tests failed: %s", string(output))
    }
    return "All tests passed successfully", nil
}

// View logs handler
type ViewLogsHandler struct {
    logLevel string
    lastOpID string
}

func (h *ViewLogsHandler) Name() string { return "ViewLogsHandler" }
func (h *ViewLogsHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *ViewLogsHandler) GetLastOperationID() string { return h.lastOpID }
func (h *ViewLogsHandler) Label() string { return "Log Level" }
func (h *ViewLogsHandler) Value() string { return h.logLevel }
func (h *ViewLogsHandler) Editable() bool { return true }
func (h *ViewLogsHandler) Timeout() time.Duration { return 0 }
func (h *ViewLogsHandler) Change(newValue any) (string, error) {
    level := strings.ToUpper(strings.TrimSpace(newValue.(string)))
    validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
    
    for _, validLevel := range validLevels {
        if level == validLevel {
            h.logLevel = level
            return fmt.Sprintf("Log level set to %s", level), nil
        }
    }
    
    return "", fmt.Errorf("invalid log level. Must be one of: %v", validLevels)
}

func main() {
    config := &devtui.TuiConfig{
        AppName:       "DevApp",
        TabIndexStart: 0,
        ExitChan:      make(chan bool),
    }

    tui := devtui.NewTUI(config)

    // Create handlers for different functional areas
    databaseHandler := &DatabaseURLHandler{currentURL: "postgresql://localhost:5432/mydb"}
    apiKeyHandler := &APIKeyHandler{currentKey: ""}

    buildHandler := &BuildProjectHandler{projectPath: "./"}
    testHandler := &TestSuiteHandler{testPath: "./tests"}

    healthHandler := &HealthCheckHandler{serviceName: "API Service"}
    logHandler := &ViewLogsHandler{logLevel: "INFO"}

    // Organize into logical tabs
    tui.NewTabSection("Config", "Application settings").
        NewField(databaseHandler).
        NewField(apiKeyHandler)
        
    tui.NewTabSection("Build", "Development operations").  
        NewField(buildHandler).
        NewField(testHandler)
        
    tui.NewTabSection("Monitor", "System monitoring").
        NewField(healthHandler).
        NewField(logHandler)

    var wg sync.WaitGroup
    wg.Add(1)
    go tui.Start(&wg)
    wg.Wait()
}
```

## Advanced Handler Patterns

### Progressive Status Updates

```go
type ProgressHandler struct {
    currentStatus string
    lastOpID      string
}

func (h *ProgressHandler) Name() string { return "ProgressHandler" }
func (h *ProgressHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *ProgressHandler) GetLastOperationID() string { return h.lastOpID }
func (h *ProgressHandler) Label() string { return "Process Data" }
func (h *ProgressHandler) Value() string { return "Start processing" }
func (h *ProgressHandler) Editable() bool { return false }
func (h *ProgressHandler) Timeout() time.Duration { return 30 * time.Second }

func (h *ProgressHandler) Change(newValue any) (string, error) {
    // Simulate multi-stage process with status updates
    stages := []string{
        "Initializing...",
        "Loading data...",
        "Processing records...",
        "Validating results...",
        "Finalizing...",
    }
    
    for i, stage := range stages {
        h.currentStatus = stage
        time.Sleep(1 * time.Second) // Simulate work
        
        // Each stage could potentially update the same message
        // using the operation ID tracking system
    }
    
    return "Processing completed successfully", nil
}
```

### Conditional Field Behavior

```go
type ConditionalHandler struct {
    mode     string
    lastOpID string
}

func (h *ConditionalHandler) Name() string { return "ConditionalHandler" }
func (h *ConditionalHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *ConditionalHandler) GetLastOperationID() string { return h.lastOpID }
func (h *ConditionalHandler) Label() string { return "Mode" }
func (h *ConditionalHandler) Value() string { return h.mode }
func (h *ConditionalHandler) Editable() bool {
    // Conditional editability based on current state
    return h.mode != "locked"
}
func (h *ConditionalHandler) Timeout() time.Duration { return 5 * time.Second }

func (h *ConditionalHandler) Change(newValue any) (string, error) {
    newMode := strings.TrimSpace(newValue.(string))
    validModes := []string{"development", "staging", "production", "locked"}
    
    for _, valid := range validModes {
        if newMode == valid {
            h.mode = newMode
            return fmt.Sprintf("Mode changed to %s", newMode), nil
        }
    }
    
    return "", fmt.Errorf("invalid mode. Valid modes: %v", validModes)
}
```

## io.Writer Integration Examples

### Direct Message Writing

```go
func demonstrateDirectWriting(tui *devtui.DevTUI) {
    // Create a tab
    tab := tui.NewTabSection("Messages", "Direct message examples")
    
    // Create a handler for writing
    handler := &MyWritingHandler{name: "DirectWriter"}
    
    // Register the handler and get a writer
    writer := tab.RegisterWritingHandler(handler)
    
    // Use standard io.Writer methods
    writer.Write([]byte("Starting process..."))
    writer.Write([]byte("Processing data..."))
    writer.Write([]byte("Operation completed successfully"))
    
    // Messages will appear with [DirectWriter] prefix automatically
}

type MyWritingHandler struct {
    name     string
    lastOpID string
}

func (h *MyWritingHandler) Name() string { return h.name }
func (h *MyWritingHandler) SetLastOperationID(id string) { h.lastOpID = id }
func (h *MyWritingHandler) GetLastOperationID() string { return h.lastOpID }
```

### Concurrent Writing

```go
func demonstrateConcurrentWriting(tui *devtui.DevTUI) {
    tab := tui.NewTabSection("Concurrent", "Multiple writers example")
    
    // Create multiple handlers
    handler1 := &MyWritingHandler{name: "Worker1"}
    handler2 := &MyWritingHandler{name: "Worker2"}
    
    writer1 := tab.RegisterWritingHandler(handler1)
    writer2 := tab.RegisterWritingHandler(handler2)
    
    // Concurrent writing
    go func() {
        writer1.Write([]byte("Worker 1 starting..."))
        time.Sleep(1 * time.Second)
        writer1.Write([]byte("Worker 1 completed"))
    }()
    
    go func() {
        writer2.Write([]byte("Worker 2 starting..."))
        time.Sleep(2 * time.Second)
        writer2.Write([]byte("Worker 2 completed"))
    }()
}
```

These examples demonstrate the full power and flexibility of the DevTUI handler system, from basic field validation to complex multi-stage operations with progress tracking and concurrent message writing.
