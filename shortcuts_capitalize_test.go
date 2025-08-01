package devtui

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestCapitalizeWithMultilineTranslation(t *testing.T) {
	tests := []struct {
		name        string
		appName     string
		lang        string
		expected    string
		description string
	}{
		{
			name:        "Simple multiline with Capitalize",
			appName:     "TestApp",
			lang:        "EN",
			expected:    "Testapp Shortcuts Keyboard (\"En\"):\n\nTabs:\n  • Tab/Shift+Tab  - Switch Tabs\n\nFields:\n  • Left/Right     - Navigate Fields\n  • Enter          - Edit/Execute\n  • Esc            - Cancel\n\nLanguage Supported: En, Es, Zh, Hi, Ar, Pt, Fr, De, Ru",
			description: "Test that Capitalize preserves multiline structure",
		},
		{
			name:        "Spanish translation multiline",
			appName:     "TestApp",
			lang:        "ES",
			expected:    "Testapp Atajos Teclado (\"Es\"):\n\nPestañas:\n  • Tab/Shift+Tab  - Cambiar Pestañas\n\nCampos:\n  • Left/Right     - Navegar Campos\n  • Enter          - Editar/Ejecutar\n  • Esc            - Cancelar\n\nIdioma Soportado: En, Es, Zh, Hi, Ar, Pt, Fr, De, Ru",
			description: "Test Spanish translation with multiline format preservation",
		},
		{
			name:        "Complex format preservation",
			appName:     "DevTUI",
			lang:        "EN",
			expected:    "Devtui Shortcuts Keyboard (\"En\"):\n\nTabs:\n  • Tab/Shift+Tab  - Switch Tabs\n\nFields:\n  • Left/Right     - Navigate Fields\n  • Enter          - Edit/Execute\n  • Esc            - Cancel\n\nText Edit:\n  • Left/Right     - Move Cursor\n  • Backspace      - Create Space\n  • Space/Letters  - Insert Character\n\nLanguage Supported: En, Es, Zh, Hi, Ar, Pt, Fr, De, Ru",
			description: "Test complex multiline structure with indentation and bullets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the generateHelpContent method but simplified
			result := generateSimplifiedHelpContent(tt.appName, tt.lang)

			// Check if the result maintains proper formatting
			if result != tt.expected {
				t.Errorf("Test %s failed.\nExpected: %q\nGot:      %q", tt.name, tt.expected, result)
			}

			// Additional checks for format preservation
			if !containsNewlines(result) {
				t.Errorf("Test %s: Result should contain newlines for proper formatting", tt.name)
			}

			if !containsBulletPoints(result) {
				t.Errorf("Test %s: Result should contain bullet points (•)", tt.name)
			}

			if !containsIndentation(result) {
				t.Errorf("Test %s: Result should preserve indentation spaces", tt.name)
			}
		})
	}
}

// generateSimplifiedHelpContent simulates the corrected method without Capitalize
func generateSimplifiedHelpContent(appName, lang string) string {
	// Fixed version: NOT using Capitalize() to preserve multiline formatting
	return T(appName, D.Shortcuts, D.Keyboard, `("`+lang+`"):

Tabs:
  • Tab/Shift+Tab  -`, D.Switch, ` tabs

`, D.Fields, `:
  • Left/Right     - Navigate fields
  • Enter          - Edit/Execute
  • Esc            -`, D.Cancel, `

Text Edit:
  • Left/Right     -`, D.Move, `cursor
  • Backspace      -`, D.Create, D.Space, `
  • Space/Letters  -`, D.Insert, D.Character, `

`, D.Language, D.Supported, `: EN, ES, ZH, HI, AR, PT, FR, DE, RU`).String()
}

// Helper functions to verify format preservation
func containsNewlines(s string) bool {
	return Contains(s, "\n")
}

func containsBulletPoints(s string) bool {
	return Contains(s, "•")
}

func containsIndentation(s string) bool {
	return Contains(s, "  •") // Two spaces before bullet point
}

// TestCapitalizeFormatPreservation tests specific formatting issues
func TestCapitalizeFormatPreservation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Preserve newlines and indentation",
			input:    "hello world\n  • item one\n  • item two",
			expected: "Hello World\n  • Item One\n  • Item Two",
		},
		{
			name:     "Preserve bullet points formatting",
			input:    "section:\n  • first item\n  • second item",
			expected: "Section:\n  • First Item\n  • Second Item",
		},
		{
			name:     "Complex multiline with various symbols",
			input:    "tabs:\n  • tab/shift+tab  - switch tabs\n\nfields:\n  • left/right     - navigate",
			expected: "Tabs:\n  • Tab/Shift+Tab  - Switch Tabs\n\nFields:\n  • Left/Right     - Navigate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Capitalize().String()
			if result != tt.expected {
				t.Errorf("Capitalize() formatting test failed.\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, tt.expected, result)
			}
		})
	}
}

// TestTranslationWithCapitalize tests the actual problem scenario
func TestTranslationWithCapitalize(t *testing.T) {
	// Test the exact pattern used in generateHelpContent
	tests := []struct {
		name     string
		lang     string
		expected string
	}{
		{
			name:     "English with proper formatting",
			lang:     "EN",
			expected: "Test Shortcuts Keyboard (\"En\"):\n\nTabs:\n  • Tab/Shift+Tab  - Switch Tabs",
		},
		{
			name:     "Spanish with proper formatting",
			lang:     "ES",
			expected: "Test Atajos Teclado (\"Es\"):\n\nPestañas:\n  • Tab/Shift+Tab  - Cambiar Pestañas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the exact problem pattern from generateHelpContent
			result := T("Test", D.Shortcuts, D.Keyboard, `("`+tt.lang+`"):

Tabs:
  • Tab/Shift+Tab  -`, D.Switch, ` tabs`).Capitalize().String()

			// Check if the basic structure is preserved
			if !Contains(result, "\n") {
				t.Errorf("Newlines should be preserved in result: %q", result)
			}

			if !Contains(result, "•") {
				t.Errorf("Bullet points should be preserved in result: %q", result)
			}

			// Log the actual result for debugging
			t.Logf("Language: %s\nResult: %q", tt.lang, result)
		})
	}
}
