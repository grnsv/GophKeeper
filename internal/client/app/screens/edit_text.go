package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editTextModel struct {
}

func NewEditText() tea.Model {
	return &editTextModel{}
}

func (m editTextModel) Init() tea.Cmd {
	return nil
}

func (m editTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editTextModel) View() string {
	return ""
}
