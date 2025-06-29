package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

type editCredentialsModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewEditCredentials(data []byte) (tea.Model, error) {
	m := editCredentialsModel{
		inputs: make([]textinput.Model, 3),
	}

	credentials := types.Credentials{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &credentials); err != nil {
			return nil, err
		}
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 64
		t.Width = 64

		switch i {
		case 0:
			t.Placeholder = "Resource"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			if credentials.Resource != "" {
				t.SetValue(credentials.Resource)
			}
		case 1:
			t.Placeholder = "Login"
			if credentials.Login != "" {
				t.SetValue(credentials.Login)
			}
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			if credentials.Password != "" {
				t.SetValue(credentials.Password)
			}
		}

		m.inputs[i] = t
	}

	return m, nil
}

func (m editCredentialsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editCredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				data, err := json.Marshal(types.Credentials{
					Resource: m.inputs[0].Value(),
					Login:    m.inputs[1].Value(),
					Password: m.inputs[2].Value(),
				})
				if err != nil {
					return m, commands.Error(err)
				}
				return m, commands.SubmitData(data)
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.FocusedStyle
					m.inputs[i].TextStyle = styles.FocusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = styles.NoStyle
				m.inputs[i].TextStyle = styles.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *editCredentialsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m editCredentialsModel) View() string {
	var b strings.Builder
	b.WriteString("Enter the resource name, login and password\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := "Submit"
	if m.focusIndex == len(m.inputs) {
		button = styles.FocusedButton(button)
	} else {
		button = styles.BlurredButton(button)
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	return b.String()
}
