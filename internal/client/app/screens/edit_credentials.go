package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editCredentialsModel struct {
}

func NewEditCredentials(data []byte) tea.Model {
	return &editCredentialsModel{}
}

func (m editCredentialsModel) Init() tea.Cmd {
	return nil
}

func (m editCredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editCredentialsModel) View() string {
	return ""
}
