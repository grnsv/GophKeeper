package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

type Row struct {
	KeyInput   textinput.Model
	ValueInput textinput.Model
}

func newRow() Row {
	key := textinput.New()
	key.Placeholder = "Key"
	key.Width = 28
	key.CharLimit = 32

	val := textinput.New()
	val.Placeholder = "Value"
	key.Width = 28
	key.CharLimit = 32

	return Row{KeyInput: key, ValueInput: val}
}

type editMetadataModel struct {
	metadata   types.Metadata
	rows       []Row
	focusIndex int
}

func NewEditMetadata(metadata types.Metadata) tea.Model {
	m := editMetadataModel{
		metadata: metadata,
		rows:     make([]Row, 0, len(metadata)),
	}

	for k, v := range metadata {
		row := newRow()
		row.KeyInput.SetValue(k)
		row.ValueInput.SetValue(v)
		m.rows = append(m.rows, row)
	}

	if len(m.rows) > 0 {
		m.rows[0].KeyInput.Focus()
	}

	return m
}

func (m editMetadataModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editMetadataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd = m.handleEnter()
			cmds = append(cmds, cmd)
		case "tab", "shift+tab", "up", "down":
			cmd = m.handleNavigation(msg)
			cmds = append(cmds, cmd)
		}
	}

	cmd = m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *editMetadataModel) handleNavigation(msg tea.KeyMsg) tea.Cmd {
	m.blurCurrentElement()
	totalFocusPoints := len(m.rows)*3 + 2

	switch msg.String() {
	case "up", "shift+tab":
		m.focusIndex--
	case "down", "tab":
		m.focusIndex++
	}

	if m.focusIndex < 0 {
		m.focusIndex = totalFocusPoints - 1
	}
	if m.focusIndex >= totalFocusPoints {
		m.focusIndex = 0
	}

	return m.focusCurrentElement()
}

func (m *editMetadataModel) handleEnter() tea.Cmd {
	addBtnIndex := len(m.rows) * 3
	submitBtnIndex := addBtnIndex + 1

	if m.focusIndex == submitBtnIndex {
		newMetadata := make(types.Metadata)
		for _, row := range m.rows {
			key := strings.TrimSpace(row.KeyInput.Value())
			if key != "" {
				newMetadata[key] = row.ValueInput.Value()
			}
		}
		m.metadata = newMetadata
		return commands.SubmitMetadata(m.metadata)
	}

	if m.focusIndex == addBtnIndex {
		m.blurCurrentElement()
		m.rows = append(m.rows, newRow())
		m.focusIndex = (len(m.rows) - 1) * 3
		return m.focusCurrentElement()
	}

	if m.focusIndex%3 == 2 {
		rowIndex := m.focusIndex / 3
		if len(m.rows) > 0 {
			m.rows = append(m.rows[:rowIndex], m.rows[rowIndex+1:]...)
			totalFocusPoints := len(m.rows)*3 + 2
			if m.focusIndex >= totalFocusPoints {
				m.focusIndex = totalFocusPoints - 1
			}
			return m.focusCurrentElement()
		}
	} else {
		return m.handleNavigation(tea.KeyMsg{Type: tea.KeyTab})
	}

	return nil
}

func (m *editMetadataModel) updateInputs(msg tea.Msg) tea.Cmd {
	if !m.isInputFocused() {
		return nil
	}

	rowIndex := m.focusIndex / 3
	elementInRow := m.focusIndex % 3

	var cmd tea.Cmd
	if elementInRow == 0 {
		m.rows[rowIndex].KeyInput, cmd = m.rows[rowIndex].KeyInput.Update(msg)
	} else {
		m.rows[rowIndex].ValueInput, cmd = m.rows[rowIndex].ValueInput.Update(msg)
	}
	return cmd
}

func (m *editMetadataModel) isInputFocused() bool {
	if m.focusIndex >= len(m.rows)*3 {
		return false
	}
	return m.focusIndex%3 == 0 || m.focusIndex%3 == 1
}

func (m *editMetadataModel) blurCurrentElement() {
	if m.isInputFocused() {
		rowIndex := m.focusIndex / 3
		elementInRow := m.focusIndex % 3

		if elementInRow == 0 {
			m.rows[rowIndex].KeyInput.Blur()
		} else {
			m.rows[rowIndex].ValueInput.Blur()
		}
	}
}

func (m *editMetadataModel) focusCurrentElement() tea.Cmd {
	if m.isInputFocused() {
		rowIndex := m.focusIndex / 3
		elementInRow := m.focusIndex % 3

		if elementInRow == 0 {
			return m.rows[rowIndex].KeyInput.Focus()
		}
		return m.rows[rowIndex].ValueInput.Focus()
	}
	return nil
}

func (m editMetadataModel) View() string {
	var b strings.Builder
	for i, row := range m.rows {
		keyStyle, valStyle := styles.BlurredBorderStyle, styles.BlurredBorderStyle
		if m.focusIndex == i*3 {
			keyStyle = styles.FocusedBorderStyle
		}
		if m.focusIndex == i*3+1 {
			valStyle = styles.FocusedBorderStyle
		}

		deleteBtn := "[ - ]"
		if m.focusIndex == i*3+2 {
			deleteBtn = styles.FocusedButtonStyle.Render(" - ")
		}

		rowString := lipgloss.JoinHorizontal(lipgloss.Center,
			keyStyle.Render(row.KeyInput.View()),
			"  ",
			valStyle.Render(row.ValueInput.View()),
			"  ",
			deleteBtn,
		)
		b.WriteString(rowString)
		b.WriteString("\n")
	}

	addBtnIndex := len(m.rows) * 3
	addBtn := "[ + ]"
	if m.focusIndex == addBtnIndex {
		addBtn = styles.FocusedButtonStyle.Render(" + ")
	}
	fmt.Fprintf(&b, "\n%s\n", addBtn)

	submitBtnIndex := addBtnIndex + 1
	submitBtn := "Submit"
	if m.focusIndex == submitBtnIndex {
		submitBtn = styles.FocusedButtonStyle.Render(submitBtn)
	} else {
		submitBtn = styles.ButtonStyle.Render(submitBtn)
	}
	fmt.Fprintf(&b, "\n%s\n", submitBtn)

	return lipgloss.JoinVertical(lipgloss.Left,
		"You can add additional metadata.\n",
		b.String(),
	)
}
