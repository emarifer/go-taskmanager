package UITable

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	ColumnKeyID          = "id"
	ColumnKeyName        = "name"
	ColumnKeyDescription = "description"
	ColumnKeyStatus      = "status"
)

var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "╥",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "╨",
		InnerJunction:  "╫",

		InnerDivider: "║",
	}
)

type Model struct {
	tableModel table.Model
}

func NewModel(rows []table.Row) /* Model */ {
	columns := []table.Column{
		table.NewColumn(ColumnKeyID, "ID", 5).WithStyle(
			lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.Color("#88f")).
				Align(lipgloss.Center)),
		table.NewColumn(ColumnKeyName, "Name", 10),
		table.NewColumn(ColumnKeyDescription, "Description", 30),
		table.NewColumn(ColumnKeyStatus, "Status", 8),
	}

	/* rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyID: "abc",
			// Missing name
			columnKeyDescription: "The first table entry, ever",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "123",
			columnKeyName:        "Oh no",
			columnKeyDescription: "Super bold!",
			columnKeyCount:       17,
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#af00d7")).Bold(true)),
		table.NewRow(table.RowData{
			columnKeyID: "def",
			// Apply a style to this cell
			columnKeyName:        table.NewStyledCell("Styled", lipgloss.NewStyle().Foreground(lipgloss.Color("#8ff"))),
			columnKeyDescription: "This is a really, really, really long description that will get cut off",
			columnKeyCount:       table.NewStyledCell(0, lipgloss.NewStyle().Faint(true)),
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
	} */

	// Start with the default key map and change it slightly, just for demoing
	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	model := Model{
		// Throw features in... the point is not to look good, it's just reference!
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
			SelectableRows(true).
			Focused(true).
			Border(customBorder).
			WithKeyMap(keys).
			WithStaticFooter("Footer!").
			WithPageSize(5).
			WithSelectedText(" ", "✓").
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("#a38")).
					Foreground(lipgloss.Color("#a7a")).
					Align(lipgloss.Left),
			).
			// SortByAsc(ColumnKeyID).
			WithMissingDataIndicatorStyled(table.StyledCell{
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#faa")),
				Data:  "<nil>",
			}),
	}

	model.updateFooter()

	// return model
	if _, err := tea.NewProgram(model).Run(); err != nil {
		log.Fatal(err)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updateFooter() {
	highlightedRow := m.tableModel.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Description of the current task: %s",
		m.tableModel.CurrentPage(),
		m.tableModel.MaxPages(),
		highlightedRow.Data[ColumnKeyDescription],
	)

	m.tableModel = m.tableModel.WithStaticFooter(footerText)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	// We control the footer text, so make sure to update it
	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "i":
			m.tableModel = m.tableModel.WithHeaderVisibility(!m.tableModel.GetHeaderVisibility())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	// body.WriteString("A (chaotic) table demo with all features enabled!\n")
	body.WriteString("\nPress left/right or page up/down to move pages\n")
	body.WriteString("Press 'i' to toggle the header visibility\n")
	body.WriteString("Press space/enter to select a row, q or ctrl+c to quit\n\n")

	// selectedIDs := []string{}

	// for _, row := range m.tableModel.SelectedRows() {
	// 	// Slightly dangerous type assumption but fine for demo
	// 	selectedIDs = append(selectedIDs, row.Data[columnKeyID].(string))
	// }

	// body.WriteString(fmt.Sprintf("SelectedIDs: %s\n", strings.Join(selectedIDs, ", ")))
	// selectedIDs := m.tableModel.HighlightedRow().Data[columnKeyID].(string)
	// if len(selectedIDs) != 0 {
	// 	body.WriteString(fmt.Sprintf("Selected ID: %s\n", selectedIDs))
	// }

	body.WriteString(m.tableModel.View())

	body.WriteString("\n")

	return body.String()
}
