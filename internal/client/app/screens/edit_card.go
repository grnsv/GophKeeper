package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editCardModel struct {
}

func NewEditCard() tea.Model {
	return &editCardModel{}
}

func (m editCardModel) Init() tea.Cmd {
	return nil
}

func (m editCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m editCardModel) View() string {
	return ""
}
