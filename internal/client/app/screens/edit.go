package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type editModel struct {
	record *models.Record
	screen tea.Model
}

func NewEdit(record *models.Record) tea.Model {
	m := editModel{record: record}
	if m.record == nil {
		m.record = &models.Record{}
	}
	m.screen = getScreenByType(m.record)
	return m
}

func getScreenByType(record *models.Record) tea.Model {
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
		return NewEditType()
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
		return m.changeScreen(getScreenByType(m.record))
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
