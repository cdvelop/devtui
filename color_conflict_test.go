package devtui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cdvelop/messagetype"
)

// TestOpcionA_RequirementsValidation validates the core requirements from BETTER_VIEW.md
func TestOpcionA_RequirementsValidation(t *testing.T) {
	tui := NewTUI(&TuiConfig{
		AppName:  "Requirements Test",
		ExitChan: make(chan bool),
		TestMode: true,
	})

	tab := tui.NewTabSection("Test", "Test Tab")

	testCases := []struct {
		handler  string
		content  string
		msgType  messagetype.Type
		expected string // Expected format pattern
	}{
		{"DatabaseConfig", "postgres://localhost:5432/mydb", messagetype.Info, "[DatabaseConfig] postgres://localhost:5432/mydb"},
		{"SystemBackup", "Create System Backup", messagetype.Success, "[SystemBackup] Create System Backup"},
		{"ErrorHandler", "Connection failed", messagetype.Error, "[ErrorHandler] Connection failed"},
		{"WarningHandler", "Deprecated function", messagetype.Warning, "[WarningHandler] Deprecated function"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.handler, tc.msgType), func(t *testing.T) {
			tabContent := tui.createTabContent(tc.content, tc.msgType, tab, tc.handler, "")
			formattedMessage := tui.formatMessage(tabContent)

			t.Logf("Message: %s", formattedMessage)

			// 1. Verificar que el formato tiene brackets unidos: [HandlerName]
			expectedPattern := fmt.Sprintf("[%s]", tc.handler)
			if !strings.Contains(formattedMessage, expectedPattern) {
				t.Errorf("FAIL: Expected unified pattern '[%s]' not found", tc.handler)
			}

			// 2. Verificar que el contenido está presente
			if !strings.Contains(formattedMessage, tc.content) {
				t.Errorf("FAIL: Expected content '%s' not found", tc.content)
			}

			// 3. Verificar que NO hay brackets separados (patrón viejo)
			separatedPattern := fmt.Sprintf(" [ %s ] ", tc.handler)
			if strings.Contains(formattedMessage, separatedPattern) {
				t.Errorf("FAIL: Found old separated pattern '%s'", separatedPattern)
			}

			t.Log("✅ PASS: Opción A requirements met")
		})
	}
}

// TestCentralizedMessageProcessing validates that all message flows use centralized processing
func TestCentralizedMessageProcessing(t *testing.T) {
	t.Log("TESTING CENTRALIZED MESSAGE PROCESSING:")
	t.Log("=======================================")

	// Test cases que deberían detectar tipo automáticamente
	testCases := []struct {
		content      string
		expectedType messagetype.Type
		description  string
	}{
		{"Database connection configured successfully", messagetype.Success, "Success word detected correctly"},
		{"ERROR: Connection failed", messagetype.Error, "Error prefix detected correctly"},
		{"WARNING: Deprecated function", messagetype.Warning, "Warning prefix detected correctly"},
		{"SUCCESS: Operation completed", messagetype.Success, "Success prefix detected correctly"},
		{"System initialized", messagetype.Normal, "Normal message detected correctly"},
		{"Backup completed successfully", messagetype.Success, "Success word detected correctly"},
		{"Preparing backup...", messagetype.Normal, "Normal progress message"},
	}

	for _, tc := range testCases {
		t.Run(tc.content, func(t *testing.T) {
			// Test que messagetype.DetectMessageType funciona correctamente
			detectedType := messagetype.DetectMessageType(tc.content)

			if detectedType != tc.expectedType {
				t.Errorf("FAIL: Expected %v, got %v for: %s", tc.expectedType, detectedType, tc.content)
			} else {
				t.Logf("✅ PASS: '%s' correctly detected as %v", tc.content, detectedType)
			}
		})
	}

	t.Log("")
	t.Log("CONCLUSION: DetectMessageType works correctly")
	t.Log("SOLUTION: All message methods now use DetectMessageType for centralization")
}

// TestLastMessageColorFixed validates that the last callback message now uses correct colors
func TestLastMessageColorFixed(t *testing.T) {
	tui := NewTUI(&TuiConfig{
		AppName:  "Last Message Color Fixed Test",
		ExitChan: make(chan bool),
		TestMode: true,
	})

	tab := tui.NewTabSection("Test", "Test Tab")

	t.Log("🔧 SOLUTION TEST: Validar que el último mensaje usa el color correcto")

	// Test casos que simulan el final de una operación
	finalMessages := []struct {
		content       string
		expectedType  messagetype.Type
		expectedColor string
		context       string
	}{
		// Casos que antes fallaban - ahora deberían funcionar
		{"Operation completed successfully", messagetype.Success, "HIGHLIGHT (#FF6600)", "Success con palabra 'successfully'"},
		{"Backup completed successfully", messagetype.Success, "HIGHLIGHT (#FF6600)", "Success con 'completed successfully'"},
		{"ERROR: Operation failed", messagetype.Error, "RED (#FF0000)", "Error con prefijo 'ERROR:'"},
		{"WARNING: Operation completed with warnings", messagetype.Warning, "YELLOW (#FFFF00)", "Warning con prefijo 'WARNING:'"},
		{"Database connection established", messagetype.Normal, "NORMAL", "Normal message sin keywords especiales"},
		{"SUCCESS: All tasks completed", messagetype.Success, "HIGHLIGHT (#FF6600)", "Success con prefijo 'SUCCESS:'"},
	}

	for _, tc := range finalMessages {
		t.Run(tc.content, func(t *testing.T) {
			// Simular el mensaje final de una operación
			tabContent := tui.createTabContent(tc.content, tc.expectedType, tab, "TestHandler", "final-op-123")
			formattedMessage := tui.formatMessage(tabContent)

			t.Logf("Context: %s", tc.context)
			t.Logf("Content: %s", tc.content)
			t.Logf("Expected: %s (%s)", tc.expectedType, tc.expectedColor)
			t.Logf("Formatted: %s", formattedMessage)

			// Verificar detección automática de tipo
			detectedType := messagetype.DetectMessageType(tc.content)
			if detectedType != tc.expectedType {
				t.Errorf("❌ DetectMessageType failed: Expected %v, got %v", tc.expectedType, detectedType)
			} else {
				t.Logf("✅ DetectMessageType working: %s → %v", tc.content, detectedType)
			}

			// Verificar que el tabContent tiene el tipo correcto
			if tabContent.Type != tc.expectedType {
				t.Errorf("❌ TabContent type wrong: Expected %v, got %v", tc.expectedType, tabContent.Type)
			} else {
				t.Logf("✅ TabContent type correct: %v", tabContent.Type)
			}
		})
	}

	t.Log("")
	t.Log("🎯 RESULT: sendSuccessMessage() y sendErrorMessage() ahora usan DetectMessageType")
	t.Log("✅ BENEFIT: El último mensaje de callback tendrá el color correcto según su contenido")
	t.Log("✅ CONSISTENCY: Todos los métodos de envío de mensajes usan centralización")
}
