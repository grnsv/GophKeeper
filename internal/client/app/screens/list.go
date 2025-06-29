package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
}

func NewList() tea.Model {
	return &listModel{}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m listModel) View() string {
	return ""
}
