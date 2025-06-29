package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	red   = lipgloss.Color("1")
	green = lipgloss.Color("2")

	BoldRed   = lipgloss.NewStyle().Foreground(red).Bold(true)
	BoldGreen = lipgloss.NewStyle().Foreground(green).Bold(true)

	TitleStyle = lipgloss.NewStyle().Bold(true)
	ErrorStyle = BoldRed

	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func FocusedButton(title string) string {
	return FocusedStyle.Render(fmt.Sprintf("[ %s ]", title))
}

func BlurredButton(title string) string {
	return fmt.Sprintf("[ %s ]", BlurredStyle.Render(title))
}
