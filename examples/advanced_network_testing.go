package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdvelop/devtui"
)

// === Advanced Network Configuration Handlers ===

type DatabaseHostHandler struct {
	currentHost string
}

func (h *DatabaseHostHandler) Label() string          { return "Database Host" }
func (h *DatabaseHostHandler) Value() string          { return h.currentHost }
func (h *DatabaseHostHandler) Editable() bool         { return true }
func (h *DatabaseHostHandler) Timeout() time.Duration { return 8 * time.Second }

func (h *DatabaseHostHandler) Change(newValue any) (string, error) {
	host := strings.TrimSpace(newValue.(string))
	if host == "" {
		return "", fmt.Errorf("database host cannot be empty")
	}

	// Simulate database connection validation
	time.Sleep(2 * time.Second) // Simulate network lookup and validation

	if !strings.Contains(host, ".") && host != "localhost" {
		return "", fmt.Errorf("invalid database host format")
	}

	h.currentHost = host
	return fmt.Sprintf("Database host configured: %s", host), nil
}

type DatabasePortHandler struct {
	currentPort string
	dbType      string
}

func NewDatabasePortHandler(dbType string) *DatabasePortHandler {
	defaultPorts := map[string]string{
		"postgres": "5432",
		"mysql":    "3306",
		"mongodb":  "27017",
		"redis":    "6379",
	}

	return &DatabasePortHandler{
		currentPort: defaultPorts[dbType],
		dbType:      dbType,
	}
}

func (h *DatabasePortHandler) Label() string          { return fmt.Sprintf("%s Port", h.dbType) }
func (h *DatabasePortHandler) Value() string          { return h.currentPort }
func (h *DatabasePortHandler) Editable() bool         { return true }
func (h *DatabasePortHandler) Timeout() time.Duration { return 3 * time.Second }

func (h *DatabasePortHandler) Change(newValue any) (string, error) {
	portStr := strings.TrimSpace(newValue.(string))
	if portStr == "" {
		return "", fmt.Errorf("port cannot be empty")
	}

	// Simulate port validation and availability check
	time.Sleep(800 * time.Millisecond)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("port must be between 1 and 65535")
	}

	h.currentPort = portStr
	return fmt.Sprintf("%s port configured: %d", h.dbType, port), nil
}

type HealthCheckHandler struct {
	endpoint string
	name     string
}

func NewHealthCheckHandler(name, endpoint string) *HealthCheckHandler {
	return &HealthCheckHandler{
		endpoint: endpoint,
		name:     name,
	}
}

func (h *HealthCheckHandler) Label() string          { return fmt.Sprintf("Health Check - %s", h.name) }
func (h *HealthCheckHandler) Value() string          { return "Press Enter to check" }
func (h *HealthCheckHandler) Editable() bool         { return false }
func (h *HealthCheckHandler) Timeout() time.Duration { return 15 * time.Second }

func (h *HealthCheckHandler) Change(newValue any) (string, error) {
	if h.endpoint == "" {
		return "", fmt.Errorf("endpoint not configured")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(h.endpoint)
	if err != nil {
		return "", fmt.Errorf("%s health check failed: %v", h.name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s health check failed with status: %s", h.name, resp.Status)
	}

	return fmt.Sprintf("%s is healthy! Response: %s", h.name, resp.Status), nil
}

// === Load Testing Handlers ===

type LoadTestHandler struct {
	testType string
	duration time.Duration
}

func NewLoadTestHandler(testType string, duration time.Duration) *LoadTestHandler {
	return &LoadTestHandler{
		testType: testType,
		duration: duration,
	}
}

func (h *LoadTestHandler) Label() string          { return fmt.Sprintf("Load Test - %s", h.testType) }
func (h *LoadTestHandler) Value() string          { return fmt.Sprintf("Run %v test", h.duration) }
func (h *LoadTestHandler) Editable() bool         { return false }
func (h *LoadTestHandler) Timeout() time.Duration { return h.duration + (10 * time.Second) }

func (h *LoadTestHandler) Change(newValue any) (string, error) {
	// Simulate load test
	startTime := time.Now()

	for elapsed := time.Duration(0); elapsed < h.duration; elapsed = time.Since(startTime) {
		// Simulate load test progress
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Sprintf("%s load test completed in %v", h.testType, time.Since(startTime)), nil
}

// === Environment Variable Handler ===

type EnvVarHandler struct {
	varName      string
	currentValue string
	isRequired   bool
}

func NewEnvVarHandler(name, defaultValue string, required bool) *EnvVarHandler {
	return &EnvVarHandler{
		varName:      name,
		currentValue: defaultValue,
		isRequired:   required,
	}
}

func (h *EnvVarHandler) Label() string {
	if h.isRequired {
		return fmt.Sprintf("%s *", h.varName)
	}
	return h.varName
}
func (h *EnvVarHandler) Value() string          { return h.currentValue }
func (h *EnvVarHandler) Editable() bool         { return true }
func (h *EnvVarHandler) Timeout() time.Duration { return 2 * time.Second }

func (h *EnvVarHandler) Change(newValue any) (string, error) {
	value := strings.TrimSpace(newValue.(string))

	if h.isRequired && value == "" {
		return "", fmt.Errorf("%s is required and cannot be empty", h.varName)
	}

	// Simulate environment variable validation
	time.Sleep(300 * time.Millisecond)

	h.currentValue = value
	return fmt.Sprintf("Environment variable %s updated", h.varName), nil
}

func main() {
	config := &devtui.TuiConfig{
		AppName:       "Advanced Network Testing",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#1A1A1A",
			Highlight:  "#00D7FF",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"Network Test Log:"}, messages...)...)
		},
	}

	tui := devtui.NewTUI(config)

	// Database Configuration Tab
	dbHost := &DatabaseHostHandler{currentHost: "localhost"}
	pgPort := NewDatabasePortHandler("postgres")
	mysqlPort := NewDatabasePortHandler("mysql")

	tui.NewTabSection("Database", "Database connection configuration").
		NewField(dbHost).
		NewField(pgPort).
		NewField(mysqlPort)

	// Health Checks Tab
	apiHealth := NewHealthCheckHandler("API", "http://localhost:8080/health")
	dbHealth := NewHealthCheckHandler("Database", "http://localhost:5432/ping")
	cacheHealth := NewHealthCheckHandler("Cache", "http://localhost:6379/ping")

	tui.NewTabSection("Health", "Service health monitoring").
		NewField(apiHealth).
		NewField(dbHealth).
		NewField(cacheHealth)

	// Load Testing Tab
	lightLoad := NewLoadTestHandler("Light", 5*time.Second)
	mediumLoad := NewLoadTestHandler("Medium", 15*time.Second)
	heavyLoad := NewLoadTestHandler("Heavy", 30*time.Second)

	tui.NewTabSection("Load Test", "Performance testing tools").
		NewField(lightLoad).
		NewField(mediumLoad).
		NewField(heavyLoad)

	// Environment Configuration Tab
	apiKey := NewEnvVarHandler("API_KEY", "", true)
	debug := NewEnvVarHandler("DEBUG", "false", false)
	logLevel := NewEnvVarHandler("LOG_LEVEL", "info", false)

	tui.NewTabSection("Environment", "Environment variable configuration").
		NewField(apiKey).
		NewField(debug).
		NewField(logLevel)

	var wg sync.WaitGroup
	wg.Add(1)
	go tui.Start(&wg)
	wg.Wait()
}
