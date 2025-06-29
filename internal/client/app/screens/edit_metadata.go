package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editMetadataModel struct {
}

func NewEditMetadata() tea.Model {
	return &editMetadataModel{}
}

func (m editMetadataModel) Init() tea.Cmd {
	return nil
}

func (m editMetadataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editMetadataModel) View() string {
	return ""
}
