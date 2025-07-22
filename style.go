package devtui

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorStyle struct {
	Foreground string // eg: #F4F4F4
	Background string // eg: #000000
	Highlight  string // eg: #FF6600
	Lowlight   string // eg: #666666
}

const (
	infoColor = "#00FF00" // Verde brillante
	warnColor = "#FFFF00" // Amarillo brillante
	errColor  = "#FF0000" // Rojo brillante
)

type tuiStyle struct {
	*ColorStyle

	contentBorder    lipgloss.Border
	headerTitleStyle lipgloss.Style
	labelWidth       int // Ancho estándar para etiquetas
	labelStyle       lipgloss.Style

	footerInfoStyle lipgloss.Style

	fieldLineStyle     lipgloss.Style
	fieldSelectedStyle lipgloss.Style
	fieldEditingStyle  lipgloss.Style

	textContentStyle  lipgloss.Style
	lineHeadFootStyle lipgloss.Style // header right and footer left line

	// Estilos globales mensajes
	okStyle   lipgloss.Style
	errStyle  lipgloss.Style
	warnStyle lipgloss.Style
	infoStyle lipgloss.Style
	normStyle lipgloss.NoColor
	timeStyle lipgloss.Style
}

func newTuiStyle(cs *ColorStyle) *tuiStyle {
	// check if color is nil
	if cs == nil {
		cs = &ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		}
	}

	t := &tuiStyle{
		ColorStyle: cs,
		labelWidth: 15, // Definir un ancho estándar en caracteres para etiquetas
	}

	t.labelStyle = lipgloss.NewStyle().
		Width(t.labelWidth).
		Align(lipgloss.Left).
		Padding(0, 0)

	// El borde del contenido necesita conectarse con las pestañas
	t.contentBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	t.headerTitleStyle = lipgloss.NewStyle().
		Padding(0, 1).
		BorderForeground(lipgloss.Color(t.Highlight)).
		Background(lipgloss.Color(t.Highlight)).
		Foreground(lipgloss.Color(t.Foreground))

	t.footerInfoStyle = t.headerTitleStyle

	t.fieldLineStyle = lipgloss.NewStyle().
		Padding(0, 2)

	t.fieldSelectedStyle = t.fieldLineStyle
	t.fieldSelectedStyle = t.fieldSelectedStyle.
		Bold(true).
		Background(lipgloss.Color(t.Highlight)).
		Foreground(lipgloss.Color(t.Foreground))

	t.fieldEditingStyle = t.fieldSelectedStyle.
		Foreground(lipgloss.Color(t.Background))

	// Estilo para los mensajes - VISUAL UPGRADE: Padding interno para mejor legibilidad
	t.textContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Foreground)).
		PaddingLeft(1).
		PaddingRight(1)

	t.lineHeadFootStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Highlight))

	// Inicializar los estilos que antes eran globales
	t.okStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(infoColor)) // Verde brillante

	t.errStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(errColor)) // Rojo brillante

	t.warnStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(warnColor)) // Amarillo brillante

	t.infoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(infoColor)) // verde brillante

	t.normStyle = lipgloss.NoColor{}

	t.timeStyle = lipgloss.NewStyle().Foreground(
		lipgloss.Color(t.Lowlight),
	)

	return t
}
