package app

import (
	"errors"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/screens"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

var netErr *net.OpError

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = max(80, msg.Width)
		newScreen, cmd := m.screen.Update(msg)
		m.screen = newScreen
		return m, cmd

	case types.FetchVersionsMsg:
		m.versions.Server = msg.ServerVersion
		return m.handleError(msg.Err)

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
			screen, err := screens.NewEdit(m.svc, nil)
			if err != nil {
				return m, commands.Error(err)
			}
			return m.changeScreen(screen)
		case "Sync":
			return m, tea.Batch(m.trySync(), commands.BackToMenu)
		}
		return m, nil

	case types.BackToMenuMsg:
		mode := screens.MenuGuest
		if m.authenticated {
			mode = screens.MenuAuth
		}
		return m.changeScreen(screens.NewMenu(m.svc, mode))

	case types.AuthMsg:
		if msg.Err != nil {
			return m.handleError(msg.Err)
		}
		m.connected = true
		m.authenticated = true
		return m, tea.Batch(commands.SyncTick(), m.trySync(), commands.BackToMenu)

	case types.SyncTickMsg:
		return m, tea.Batch(commands.SyncTick(), m.trySync())

	case types.SyncMsg:
		return m.handleError(msg.Err)

	case types.ErrMsg:
		return m.handleError(msg.Err)

	case types.ErrClearedMsg:
		m.errMsg = ""
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

func (m appModel) handleError(err error) (appModel, tea.Cmd) {
	if err == nil {
		m.connected = true
		return m, nil
	}

	m.errMsg = err.Error()
	if errors.As(err, &netErr) {
		m.connected = false
	}

	return m, commands.ClearErrorAfter(3 * time.Second)
}

func (m appModel) trySync() tea.Cmd {
	if m.connected {
		return commands.Sync(m.svc)
	}
	return commands.TrySync(m.svc)
}
