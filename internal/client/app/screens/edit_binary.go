package screens

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

type editBinaryModel struct {
	data           types.Binary
	focusIndex     int
	filepicker     filepicker.Model
	selectedFile   string
	metadataScreen tea.Model
}

func NewEditBinary(data []byte) (tea.Model, error) {
	m := editBinaryModel{
		filepicker: filepicker.New(),
	}

	var err error
	if m.data, err = m.decodeData(data); err != nil {
		return nil, err
	}

	m.filepicker.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m editBinaryModel) decodeData(bytes []byte) (data types.Binary, err error) {
	data.Metadata = make(types.Metadata)
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &data)
	}
	return
}

func (m editBinaryModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m editBinaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			m.focusIndex = (m.focusIndex + 1) % 2
		case tea.KeyEnter:
			if m.focusIndex == 1 {
				didSelect, path := m.filepicker.DidSelectFile(msg)
				if !didSelect {
					return m, commands.Error(errors.New("file not selected"))
				}

				data, err := os.ReadFile(path)
				if err != nil {
					return m, commands.Error(err)
				}

				m.data.Metadata["filename"] = filepath.Base(path)
				m.data.Binary = data
				m.metadataScreen = NewEditMetadata(m.data.Metadata)
				return m, m.metadataScreen.Init()
			}
		}
	}

	if m.focusIndex == 0 {
		m.filepicker, cmd = m.filepicker.Update(msg)
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
	}

	return m, cmd
}

func (m editBinaryModel) View() string {
	if m.metadataScreen != nil {
		return m.metadataScreen.View()
	}

	var b strings.Builder
	if m.selectedFile == "" {
		b.WriteString("Pick a file:")
	} else {
		b.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}

	b.WriteString("\n\n" + m.filepicker.View())

	button := "Continue"
	if m.focusIndex == 1 {
		button = styles.FocusedButtonStyle.Render(button)
	} else {
		button = styles.ButtonStyle.Render(button)
	}
	fmt.Fprintf(&b, "\n%s\n", button)

	return b.String()
}
