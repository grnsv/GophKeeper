package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type editModel struct {
	svc        interfaces.Service
	record     *models.Record
	screen     tea.Model
	bodyHeight int
}

func NewEdit(svc interfaces.Service, record *models.Record) (tea.Model, error) {
	m := editModel{svc: svc, record: record}
	if m.record == nil {
		m.record = &models.Record{}
	}
	var err error
	m.screen, err = getScreenByType(m.record)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getScreenByType(record *models.Record) (tea.Model, error) {
	switch record.Type {
	case models.RecordTypeCredentials:
		return NewEditCredentials(record.Data)
	case models.RecordTypeText:
		return NewEditText(record.Data)
	case models.RecordTypeBinary:
		return NewEditBinary(record.Data)
	case models.RecordTypeCard:
		return NewEditCard(record.Data)
	default:
		return NewEditType(), nil
	}
}

func (m editModel) Init() tea.Cmd {
	return tea.Batch(m.screen.Init(), tea.WindowSize())
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, commands.BackToMenu
		}
	case types.RecordTypeSelectedMsg:
		m.record.Type = msg.RecordType
		screen, err := getScreenByType(m.record)
		if err != nil {
			return m, commands.Error(err)
		}
		return m.changeScreen(screen)
	case types.DataMsg:
		m.record.Data = msg.Data
		return m, tea.Batch(commands.BackToMenu, commands.SaveRecord(m.svc, m.record))
	case tea.WindowSizeMsg:
		m.bodyHeight = styles.CalcBodyHeight(msg.Height)
		msg.Height = m.bodyHeight
		var cmd tea.Cmd
		m.screen, cmd = m.screen.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.screen, cmd = m.screen.Update(msg)

	return m, cmd
}

func (m editModel) changeScreen(newScreen tea.Model) (tea.Model, tea.Cmd) {
	m.screen = newScreen
	return m, tea.Batch(m.screen.Init(), tea.WindowSize())
}

func (m editModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Height(m.bodyHeight).Render(m.screen.View()),
		styles.FooterStyle.Render("Press Esc to return to the menu."),
	)
}
