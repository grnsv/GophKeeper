package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
)

func (m appModel) View() string {
	var b strings.Builder
	b.WriteString(renderHeader(m.width, m.offline, m.hasConflicts))
	b.WriteString("\n")
	b.WriteString(renderError(m.err))
	b.WriteString("\n")
	b.WriteString(m.screen.View())

	return b.String()
}

func renderHeader(width int, offline, hasConflicts bool) string {
	title := styles.TitleStyle.Render("GophKeeper")

	var connBadge string
	if offline {
		connBadge = styles.BoldRed.Render("[OFFLINE]")
	} else {
		connBadge = styles.BoldGreen.Render("[ONLINE]")
	}

	var conflictBadge string
	if hasConflicts {
		conflictBadge = styles.BoldRed.Render("[CONFLICT]")
	} else {
		conflictBadge = styles.BoldGreen.Render("[NO CONFLICTS]")
	}

	badges := lipgloss.JoinHorizontal(lipgloss.Top, connBadge, " ", conflictBadge)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		lipgloss.PlaceHorizontal(width-lipgloss.Width(title), lipgloss.Right, badges),
	)
}

func renderError(err error) string {
	if err == nil {
		return ""
	}

	return styles.ErrorStyle.Render("Error: " + err.Error())
}
