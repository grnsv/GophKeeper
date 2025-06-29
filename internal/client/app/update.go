package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/screens"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = max(80, msg.Width)
		return m, nil

	case types.FetchVersionsMsg:
		m.versions.Server = msg.ServerVersion
		m.offline = msg.Offline
		m.err = msg.Err
		return m, nil

	case types.MenuSelectedMsg:
		switch msg.Item {
		case "Login":
			return m.changeScreen(screens.NewAuth(m.svc, screens.AuthModeLogin))
		case "Register":
			return m.changeScreen(screens.NewAuth(m.svc, screens.AuthModeRegister))
		case "About":
			return m.changeScreen(screens.NewAbout(m.versions))
		// case "Show":
		// 	m.screen = screens.NewList(m.svc)
		case "Add":
			return m.changeScreen(screens.NewEdit(nil))
		case "Sync":
			return m, tea.Batch(m.trySync(), commands.BackToMenu)
		}
		return m, nil

	case types.BackToMenuMsg:
		mode := screens.MenuGuest
		if m.isAuthenticated {
			mode = screens.MenuAuth
		}
		return m.changeScreen(screens.NewMenu(m.svc, mode))

	case types.AuthMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.isAuthenticated = true
		return m, tea.Batch(commands.SyncTick(), m.trySync(), commands.BackToMenu)

	case types.SyncTickMsg:
		return m, tea.Batch(commands.SyncTick(), m.trySync())

	case types.SyncMsg:
		if msg.Err != nil {
			m.err = msg.Err
		}
		return m, nil

	}

	newScreen, cmd := m.screen.Update(msg)
	m.screen = newScreen

	return m, cmd
}

func (m appModel) changeScreen(newScreen tea.Model) (tea.Model, tea.Cmd) {
	m.screen = newScreen
	return m, m.screen.Init()
}

func (m appModel) trySync() tea.Cmd {
	if m.offline {
		return commands.TrySync(m.svc)
	}
	return commands.Sync(m.svc)
}
