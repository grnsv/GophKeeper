package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

type editTextModel struct {
	data           types.Text
	focusIndex     int
	textarea       textarea.Model
	metadataScreen tea.Model
}

func NewEditText(data []byte) (tea.Model, error) {
	m := editTextModel{
		textarea: textarea.New(),
	}

	var err error
	if m.data, err = m.decodeData(data); err != nil {
		return nil, err
	}

	m.textarea.Placeholder = "Once upon a time..."
	m.textarea.Focus()
	m.textarea.SetValue(m.data.Text)

	return m, nil
}

func (m editTextModel) decodeData(bytes []byte) (data types.Text, err error) {
	data.Metadata = make(types.Metadata)
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &data)
	}
	return
}

func (m editTextModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m editTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			m.focusIndex = (m.focusIndex + 1) % 2
			if m.focusIndex == 0 {
				cmds = append(cmds, m.textarea.Focus())
			} else {
				m.textarea.Blur()
			}
		case tea.KeyEnter:
			if m.focusIndex == 1 {
				m.data.Text = m.textarea.Value()
				m.metadataScreen = NewEditMetadata(m.data.Metadata)
				return m, m.metadataScreen.Init()
			}
		}
	}

	if m.focusIndex == 0 {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m editTextModel) View() string {
	if m.metadataScreen != nil {
		return m.metadataScreen.View()
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Tell me a story.\n\n%s\n\n", m.textarea.View())

	button := "Continue"
	if m.focusIndex == 1 {
		button = styles.FocusedButtonStyle.Render(button)
	} else {
		button = styles.ButtonStyle.Render(button)
	}
	fmt.Fprintf(&b, "\n%s\n", button)

	return b.String()
}
