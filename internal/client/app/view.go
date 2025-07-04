package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
)

func (m appModel) View() string {
	headerStyle := styles.HeaderStyle.Width(m.width)
	return lipgloss.JoinVertical(lipgloss.Top,
		headerStyle.Render(m.renderHeader(m.width-headerStyle.GetHorizontalPadding())),
		styles.ErrorStyle.Width(m.width).Render(m.renderError()),
		styles.BodyStyle.Width(m.width).Render(m.screen.View()),
	)
}

func (m appModel) renderHeader(width int) string {
	title := "GophKeeper"

	var connBadge string
	if m.connected {
		connBadge = styles.GreenBadgeStyle.Render("ONLINE")
	} else {
		connBadge = styles.RedBadgeStyle.Render("OFFLINE")
	}

	var conflictBadge string
	if m.hasConflicts {
		conflictBadge = styles.RedBadgeStyle.Render("CONFLICT")
	} else {
		conflictBadge = styles.GreenBadgeStyle.Render("NO CONFLICTS")
	}

	badges := lipgloss.JoinHorizontal(lipgloss.Top, connBadge, styles.HeaderBackground.Render(" "), conflictBadge)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		lipgloss.PlaceHorizontal(width-lipgloss.Width(title), lipgloss.Right, badges),
	)
}

func (m appModel) renderError() string {
	if m.errMsg == "" {
		return ""
	}

	return styles.ErrorStyle.Render("⚠ " + m.errMsg)
}
