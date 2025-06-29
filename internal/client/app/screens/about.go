package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type aboutModel struct {
	versions *models.Versions
}

func NewAbout(versions *models.Versions) tea.Model {
	return &aboutModel{versions: versions}
}

func (m aboutModel) Init() tea.Cmd {
	return nil
}

func (m aboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, commands.BackToMenu
	}
	return m, nil
}

func (m aboutModel) View() string {
	return fmt.Sprintf(
		"Client version: %s\nClient build date: %s\nServer version: %s\nServer build date: %s\n\nPress any key to return.\n",
		m.versions.Client.BuildVersion,
		m.versions.Client.BuildDate,
		m.versions.Server.BuildVersion,
		m.versions.Server.BuildDate,
	)
}
