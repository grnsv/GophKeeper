package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type editModel struct {
	svc    interfaces.Service
	record *models.Record
	screen tea.Model
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
		return NewEditText(record.Data), nil
	case models.RecordTypeBinary:
		return NewEditBinary(record.Data), nil
	case models.RecordTypeCard:
		return NewEditCard(record.Data), nil
	default:
		return NewEditType(), nil
	}
}

func (m editModel) Init() tea.Cmd {
	return m.screen.Init()
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
	}

	newScreen, cmd := m.screen.Update(msg)
	m.screen = newScreen

	return m, cmd
}

func (m editModel) changeScreen(newScreen tea.Model) (tea.Model, tea.Cmd) {
	m.screen = newScreen
	return m, m.screen.Init()
}

func (m editModel) View() string {
	return m.screen.View()
}
