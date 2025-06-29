package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editTypeModel struct {
}

func NewEditType() tea.Model {
	return &editTypeModel{}
}

func (m editTypeModel) Init() tea.Cmd {
	return nil
}

func (m editTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editTypeModel) View() string {
	return ""
}
