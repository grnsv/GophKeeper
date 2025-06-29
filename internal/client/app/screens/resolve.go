package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type resolveModel struct {
}

func NewResolve() tea.Model {
	return &resolveModel{}
}

func (m resolveModel) Init() tea.Cmd {
	return nil
}

func (m resolveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m resolveModel) View() string {
	return ""
}
