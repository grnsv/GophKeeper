package screens

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grnsv/GophKeeper/internal/client/app/commands"
	"github.com/grnsv/GophKeeper/internal/client/app/styles"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type listModel struct {
	svc        interfaces.Service
	records    []*models.Record
	cursor     int
	bodyHeight int
}

func NewList(svc interfaces.Service) tea.Model {
	return &listModel{svc: svc}
}

func (m listModel) Init() tea.Cmd {
	return tea.Batch(commands.Show(m.svc), tea.WindowSize())
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.RecordsMsg:
		if msg.Err != nil {
			return m, commands.Error(msg.Err)
		}
		m.records = msg.Records
		m.cursor = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, commands.BackToMenu
		case tea.KeyUp, tea.KeyShiftTab:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown, tea.KeyTab:
			if m.cursor < len(m.records)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			return m, commands.SelectRecord(m.records[m.cursor])
		case tea.KeyDelete:
			record := m.records[m.cursor]
			m.records = append(m.records[:m.cursor], m.records[m.cursor+1:]...)
			return m, commands.DeleteRecord(m.svc, record)
		}
	case tea.WindowSizeMsg:
		m.bodyHeight = styles.CalcBodyHeight(msg.Height)
		msg.Height = m.bodyHeight
	}

	return m, nil
}

func (m listModel) View() string {
	var b strings.Builder
	if len(m.records) == 0 {
		b.WriteString("Loading records or no records found...")
	} else {
		b.WriteString("Your saved records:\n\n")
		for i, record := range m.records {
			cursor := " "
			if m.cursor == i {
				cursor = styles.CursorStyle.Render(">")
			}

			var status string
			switch record.Status {
			case models.RecordStatusSynced:
				status = styles.StatusSyncedStyle.Render(string(models.RecordStatusSynced))
			case models.RecordStatusPending:
				status = styles.StatusPendingStyle.Render(string(models.RecordStatusPending))
			case models.RecordStatusConflict:
				status = styles.StatusErrorStyle.Render(string(models.RecordStatusConflict))
			case models.RecordStatusDeleted:
				status = styles.StatusDeletedStyle.Render(string(models.RecordStatusDeleted))
			default:
				status = styles.StatusErrorStyle.Render("unknown")
			}

			fmt.Fprintf(&b, "%s %s %s %s %s\n",
				cursor,
				hex.EncodeToString(record.ID[:4]),
				styles.TypeStyle.Render(string(record.Type)),
				status,
				truncateString(string(record.Data), 20),
			)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Height(m.bodyHeight).Render(b.String()),
		styles.FooterStyle.Render("Press Del to delete record, Esc to return to the menu."),
	)
}

func truncateString(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "..."
}
