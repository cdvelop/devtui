package devtui

import "testing"

func TestPadHandlerName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{"short name", "hi", 12, "     hi     "},
		{"exact length", "handler", 12, "  handler   "},
		{"longer name", "verylongname", 12, "verylongname"},
		{"empty name", "", 12, "            "},
		{"single char", "a", 5, "  a  "},
		{"width one short", "ab", 1, "a"},
		{"width one exact", "a", 1, "a"},
		{"odd padding", "test", 7, " test  "}, // 7-4=3, left=1, right=2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padHandlerName(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("padHandlerName(%q, %d) = %q; want %q", tt.input, tt.width, result, tt.expected)
			}
		})
	}
}
