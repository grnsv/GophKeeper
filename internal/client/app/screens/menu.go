package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
)

type MenuMode int

const (
	MenuGuest MenuMode = iota
	MenuAuth
)

type menuModel struct {
	svc     interfaces.Service
	choices []string
	cursor  int
}

func NewMenu(svc interfaces.Service, mode MenuMode) tea.Model {
	m := &menuModel{
		svc: svc,
	}
	switch mode {
	case MenuGuest:
		m.choices = []string{
			"Login",
			"Register",
			"About",
		}
	case MenuAuth:
		m.choices = []string{
			"Show",
			"Add",
			"Sync",
			"About",
		}
	}

	return m
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, commands.Select(m.choices[m.cursor])
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	var b strings.Builder
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	b.WriteString("\nPress q to quit.\n")
	return b.String()
}
