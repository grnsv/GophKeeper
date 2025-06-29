package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type EditMode int

const (
	EditModeCreate EditMode = iota
	EditModeUpdate
)

type editModel struct {
	record *models.Record
}

func NewEdit(mode EditMode) tea.Model {
	return &editModel{record: &models.Record{}}
}

func (m editModel) Init() tea.Cmd {
	return nil
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editModel) View() string {
	return ""
}
