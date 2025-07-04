package styles

import "github.com/charmbracelet/lipgloss"

const (
	HeaderHeight  = 3
	ErrorHeight   = 1
	FooterHeight  = 2
	MinBodyHeight = 13
	InputWidth    = 32

	Brand   = lipgloss.Color("4")
	Accent  = lipgloss.Color("99")
	Success = lipgloss.Color("2")
	Error   = lipgloss.Color("1")

	TextPrimary = lipgloss.Color("7")
	TextSubtle  = lipgloss.Color("240")
	TextMuted   = lipgloss.Color("241")

	ButtonBg = lipgloss.Color("57")
)

var (
	NoStyle    = lipgloss.NewStyle()
	HelpStyle  = lipgloss.NewStyle().Foreground(TextMuted)
	ErrorStyle = lipgloss.NewStyle().Foreground(Error).Bold(true).Height(ErrorHeight).Padding(0, 1)

	FocusedStyle = lipgloss.NewStyle().Foreground(Accent)
	BlurredStyle = lipgloss.NewStyle().Foreground(TextSubtle)
	CursorStyle  = FocusedStyle

	HeaderBackground = lipgloss.NewStyle().Background(Brand)
	HeaderStyle      = HeaderBackground.Foreground(TextPrimary).Bold(true).Height(HeaderHeight).Padding(1)
	BodyStyle        = lipgloss.NewStyle().Padding(0, 1)
	FooterStyle      = HelpStyle.Height(FooterHeight).PaddingBottom(1)

	BaseBadgeStyle  = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	RedBadgeStyle   = BaseBadgeStyle.Background(Error)
	GreenBadgeStyle = BaseBadgeStyle.Background(Success)

	BorderStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Width(InputWidth).Padding(0, 1)
	FocusedBorderStyle = BorderStyle.BorderForeground(Accent)
	BlurredBorderStyle = BorderStyle.BorderForeground(TextSubtle)
	InputTextStyle     = lipgloss.NewStyle().Foreground(Accent)

	ButtonStyle        = lipgloss.NewStyle().Background(ButtonBg).Foreground(TextPrimary).Padding(0, 1)
	FocusedButtonStyle = ButtonStyle.Background(Accent)
)

func CalcBodyHeight(height int) int {
	return max(MinBodyHeight, height-HeaderHeight-ErrorHeight-FooterHeight)
}
