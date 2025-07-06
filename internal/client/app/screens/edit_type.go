package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type editTypeModel struct {
	choices []models.RecordType
	cursor  int
}

func NewEditType() tea.Model {
	return editTypeModel{choices: []models.RecordType{
		models.RecordTypeCredentials,
		models.RecordTypeText,
		models.RecordTypeBinary,
		models.RecordTypeCard,
	}}
}

func (m editTypeModel) Init() tea.Cmd {
	return nil
}

func (m editTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "tab":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, commands.SelectType(m.choices[m.cursor])
		}
	}
	return m, nil
}

func (m editTypeModel) View() string {
	var b strings.Builder
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = styles.CursorStyle.Render(">")
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}
