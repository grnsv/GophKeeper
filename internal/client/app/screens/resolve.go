package screens

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type resolveModel struct {
	svc            interfaces.Service
	localRecord    *models.Record
	serverRecord   *models.Record
	localViewport  viewport.Model
	serverViewport viewport.Model
	focusIndex     int
}

func NewResolve(svc interfaces.Service, localRecord *models.Record) (tea.Model, error) {
	return resolveModel{
		svc:            svc,
		localRecord:    localRecord,
		serverRecord:   &models.Record{},
		localViewport:  viewport.New(0, 0),
		serverViewport: viewport.New(0, 0),
	}, nil
}

func (m resolveModel) Init() tea.Cmd {
	return tea.Batch(commands.PullRecord(m.svc, m.localRecord.ID), tea.WindowSize())
}

func (m resolveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, commands.BackToMenu
		case tea.KeyLeft, tea.KeyRight, tea.KeyTab, tea.KeyShiftTab:
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, nil
		case tea.KeyEnter:
			record := m.localRecord
			if m.focusIndex == 1 {
				record = m.serverRecord
			}
			record.Version = max(m.localRecord.Version, m.serverRecord.Version)
			return m, tea.Batch(commands.BackToMenu, commands.SaveRecord(m.svc, record))
		}
	case types.RecordPulledMsg:
		m.serverRecord = msg.Record
		m.updateViewports()
		return m, nil
	case tea.WindowSizeMsg:
		viewportHeight := styles.CalcBodyHeight(msg.Height) - 1
		viewportWidth := msg.Width/2 - 2
		m.localViewport.Width = viewportWidth
		m.localViewport.Height = viewportHeight
		m.serverViewport.Width = viewportWidth
		m.serverViewport.Height = viewportHeight
		m.updateViewports()
		return m, nil
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.localViewport, cmd = m.localViewport.Update(msg)
	} else {
		m.serverViewport, cmd = m.serverViewport.Update(msg)
	}

	return m, cmd
}

func (m *resolveModel) updateViewports() {
	m.localViewport.SetContent(styles.NoStyle.Width(m.localViewport.Width).Render(string(m.localRecord.Data)))
	m.serverViewport.SetContent(styles.NoStyle.Width(m.serverViewport.Width).Render(string(m.serverRecord.Data)))
}

func (m resolveModel) View() string {
	if m.focusIndex == 0 {
		m.localViewport.Style = styles.FocusedBorderStyle.Width(m.localViewport.Width)
		m.serverViewport.Style = styles.BlurredBorderStyle.Width(m.serverViewport.Width)
	} else {
		m.localViewport.Style = styles.BlurredBorderStyle.Width(m.localViewport.Width)
		m.serverViewport.Style = styles.FocusedBorderStyle.Width(m.serverViewport.Width)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		"Time to choose...",
		lipgloss.JoinHorizontal(lipgloss.Top, m.localViewport.View(), " ", m.serverViewport.View()),
		styles.FooterStyle.Render("Press Esc to return to the menu."),
	)
}
