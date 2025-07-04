package screens

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
)

const (
	ccn = iota
	exp
	cvv
	cch
)

type editCardModel struct {
	data           types.Card
	focusIndex     int
	inputs         []textinput.Model
	metadataScreen tea.Model
}

func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("card number is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("card number is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("card number must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func NewEditCard(data []byte) (tea.Model, error) {
	m := editCardModel{
		inputs: make([]textinput.Model, 4),
	}

	var err error
	if m.data, err = m.decodeData(data); err != nil {
		return nil, err
	}

	m.inputs[ccn] = textinput.New()
	m.inputs[ccn].Placeholder = "4505 **** **** 1234"
	m.inputs[ccn].Focus()
	m.inputs[ccn].CharLimit = 20
	m.inputs[ccn].Width = 30
	m.inputs[ccn].Prompt = ""
	m.inputs[ccn].Validate = ccnValidator
	m.inputs[ccn].SetValue(m.data.CardNumber)

	m.inputs[exp] = textinput.New()
	m.inputs[exp].Placeholder = "MM/YY "
	m.inputs[exp].CharLimit = 5
	m.inputs[exp].Width = 5
	m.inputs[exp].Prompt = ""
	m.inputs[exp].Validate = expValidator
	m.inputs[exp].SetValue(m.data.EXP)

	m.inputs[cvv] = textinput.New()
	m.inputs[cvv].Placeholder = "XXX"
	m.inputs[cvv].CharLimit = 3
	m.inputs[cvv].Width = 5
	m.inputs[cvv].Prompt = ""
	m.inputs[cvv].Validate = cvvValidator
	m.inputs[cvv].SetValue(m.data.CVV)

	m.inputs[cch] = textinput.New()
	m.inputs[cch].Placeholder = "John Galt"
	m.inputs[cch].CharLimit = 64
	m.inputs[cch].Width = 64
	m.inputs[cch].Prompt = ""
	m.inputs[cch].SetValue(m.data.CardHolder)

	return m, nil
}

func (m editCardModel) decodeData(bytes []byte) (data types.Card, err error) {
	data.Metadata = make(types.Metadata)
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &data)
	}
	return
}

func (m editCardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				for _, input := range m.inputs {
					if input.Err != nil {
						return m, commands.Error(input.Err)
					}
				}
				m.data.CardNumber = m.inputs[ccn].Value()
				m.data.EXP = m.inputs[exp].Value()
				m.data.CVV = m.inputs[cvv].Value()
				m.data.CardHolder = m.inputs[cch].Value()
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

func (m *editCardModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m editCardModel) View() string {
	if m.metadataScreen != nil {
		return m.metadataScreen.View()
	}

	var b strings.Builder
	fmt.Fprintf(&b, `Enter your bank card details

 %s
 %s

 %s  %s
 %s  %s

 %s
 %s
`,
		styles.InputTextStyle.Width(30).Render("Card Number"),
		m.inputs[ccn].View(),
		styles.InputTextStyle.Width(6).Render("EXP"), styles.InputTextStyle.Width(6).Render("CVV"),
		m.inputs[exp].View(), m.inputs[cvv].View(),
		styles.InputTextStyle.Width(30).Render("Card Holder"),
		m.inputs[cch].View(),
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
