package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editBinaryModel struct {
}

func NewEditBinary() tea.Model {
	return &editBinaryModel{}
}

func (m editBinaryModel) Init() tea.Cmd {
	return nil
}

func (m editBinaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editBinaryModel) View() string {
	return ""
}
