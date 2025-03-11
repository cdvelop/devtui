package devtui

import "github.com/charmbracelet/lipgloss"

type ColorStyle struct {
	ForeGround string // eg: #F4F4F4
	Background string // eg: #000000
	Highlight  string // eg: #FF6600
	Lowlight   string // eg: #666666
}

type tuiStyle struct {
	*ColorStyle

	contentBorder     lipgloss.Border
	headerTitleStyle  lipgloss.Style
	footerInfoStyle   lipgloss.Style
	textContentStyle  lipgloss.Style
	lineHeadFootStyle lipgloss.Style // header right and footer left line

	// Estilos que antes eran globales
	okStyle   lipgloss.Style
	errStyle  lipgloss.Style
	warnStyle lipgloss.Style
	infoStyle lipgloss.Style
	normStyle lipgloss.NoColor
	timeStyle lipgloss.Style
}

func newTuiStyle(cs *ColorStyle) *tuiStyle {
	t := &tuiStyle{
		ColorStyle: cs,
	}

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
		Foreground(lipgloss.Color(t.ForeGround))

	t.footerInfoStyle = t.headerTitleStyle

	// Estilo para los mensajes
	t.textContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.ForeGround)).
		PaddingLeft(0)

	t.lineHeadFootStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Highlight))

	// Inicializar los estilos que antes eran globales
	t.okStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")) // Verde brillante

	t.errStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000")) // Rojo brillante

	t.warnStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFF00")) // Amarillo brillante

	t.infoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.Background)) //

	t.normStyle = lipgloss.NoColor{}

	t.timeStyle = lipgloss.NewStyle().Foreground(
		lipgloss.Color(t.Lowlight),
	)

	return t
}
