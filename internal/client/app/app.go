package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/screens"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type appModel struct {
	svc           interfaces.Service
	screen        tea.Model
	versions      models.Versions
	width         int
	errMsg        string
	connected     bool
	authenticated bool
	hasConflicts  bool
}

func New(svc interfaces.Service, clientBuildVersion, clientBuildDate string) tea.Model {
	return appModel{
		svc:    svc,
		screen: screens.NewMenu(svc, screens.MenuGuest),
		versions: models.Versions{Client: models.VersionInfo{
			BuildVersion: models.NewOptString(clientBuildVersion),
			BuildDate:    models.NewOptString(clientBuildDate),
		}},
	}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(commands.FetchVersions(m.svc), tea.WindowSize())
}
