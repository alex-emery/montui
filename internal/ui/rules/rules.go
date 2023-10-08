package rules

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	table  table.Model
	height int
	width  int
	edit   bool
	editModel
}

func New() *Model {
	t := table.New([]table.Column{
		table.NewFlexColumn("pattern", "Pattern", 2),
		table.NewFlexColumn("category", "Category", 1),
	}).Focused(true).
		SortByAsc("name").
		HighlightStyle(styles.RowHighlight)

	return &Model{table: t}
}
func (m Model) Init() tea.Cmd {
	return app.GetRules()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.table = m.table.WithTargetWidth(m.width).WithPageSize(m.height - 6)
	case app.GetRulesMsg:
		m.edit = false
	case app.NewRulesMsg:
		rules := msg.Rules
		rows := make([]table.Row, 0, len(rules))
		for _, rule := range rules {
			newRow := table.NewRow(table.RowData{
				"id":       rule.ID,
				"pattern":  rule.Pattern,
				"category": rule.Category.Name,
			}).WithStyle(lipgloss.NewStyle().Background(lipgloss.Color(rule.Category.Color)))
			rows = append(rows, newRow)

		}
		m.table = m.table.WithRows(rows)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.edit {
				row := m.table.HighlightedRow().Data
				m.edit = true
				m.editModel = newEditModel(
					row["id"].(uint),
					row["pattern"].(string),
					row["category"].(string),
				)
				return m, m.editModel.Init()
			}
		}
	}

	if m.edit {
		m.editModel, cmd = m.editModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := m.table.View()
	if m.edit {
		body = m.editModel.View()
	}

	return styles.RulesPage.Height(m.height).Width(m.width).Render("Rules\n" + body)
}
