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

const (
	resource = iota
	login
	password
)

type editCredentialsModel struct {
	data           types.Credentials
	focusIndex     int
	inputs         []textinput.Model
	metadataScreen tea.Model
}

func NewEditCredentials(data []byte) (tea.Model, error) {
	m := editCredentialsModel{
		inputs: make([]textinput.Model, 3),
	}

	var err error
	if m.data, err = m.decodeData(data); err != nil {
		return nil, err
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32
		t.Width = 32
		t.Prompt = ""

		switch i {
		case resource:
			t.Placeholder = "https://passport.yandex.ru/"
			t.CharLimit = 64
			t.Width = 64
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			t.SetValue(m.data.Resource)
		case login:
			t.Placeholder = "user@yandex.ru"
			t.SetValue(m.data.Login)
		case password:
			t.Placeholder = "Pa$$w0rd"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.SetValue(m.data.Password)
		}

		m.inputs[i] = t
	}

	return m, nil
}

func (m editCredentialsModel) decodeData(bytes []byte) (data types.Credentials, err error) {
	data.Metadata = make(types.Metadata)
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &data)
	}
	return
}

func (m editCredentialsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editCredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.MetadataMsg:
		m.data.Metadata = msg.Metadata
		return m, commands.SubmitData(m.data)
	}

	var cmd tea.Cmd
	if m.metadataScreen != nil {
		m.metadataScreen, cmd = m.metadataScreen.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.data.Resource = m.inputs[resource].Value()
				m.data.Login = m.inputs[login].Value()
				m.data.Password = m.inputs[password].Value()
				m.metadataScreen = NewEditMetadata(m.data.Metadata)
				return m, m.metadataScreen.Init()
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

	cmd = m.updateInputs(msg)

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
	if m.metadataScreen != nil {
		return m.metadataScreen.View()
	}

	var b strings.Builder
	fmt.Fprintf(&b, `Enter the resource name, login and password

 %s
 %s

 %s  %s
 %s  %s
`,
		styles.InputTextStyle.Width(64).Render("Resource"),
		m.inputs[resource].View(),
		styles.InputTextStyle.Width(33).Render("Login"), styles.InputTextStyle.Width(33).Render("Password"),
		m.inputs[login].View(), m.inputs[password].View(),
	)

	button := "Continue"
	if m.focusIndex == len(m.inputs) {
		button = styles.FocusedButtonStyle.Render(button)
	} else {
		button = styles.ButtonStyle.Render(button)
	}
	fmt.Fprintf(&b, "\n%s\n", button)

	return b.String()
}
