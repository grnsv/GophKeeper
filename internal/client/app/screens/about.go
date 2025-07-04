package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type aboutModel struct {
	versions   models.Versions
	bodyHeight int
}

func NewAbout(versions models.Versions) tea.Model {
	return aboutModel{versions: versions}
}

func (m aboutModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m aboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m, commands.BackToMenu

	case tea.WindowSizeMsg:
		m.bodyHeight = styles.CalcBodyHeight(msg.Height)
		return m, nil

	}

	return m, nil
}

func (m aboutModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Height(m.bodyHeight).Render(fmt.Sprintf(
			"Client version: %s\nClient build date: %s\nServer version: %s\nServer build date: %s",
			m.versions.Client.BuildVersion,
			m.versions.Client.BuildDate,
			m.versions.Server.BuildVersion,
			m.versions.Server.BuildDate,
		)),
		styles.FooterStyle.Render("Press any key to return."),
	)
}
