package categories

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	table    table.Model
	height   int
	width    int
	edit     bool
	category categoryEditModel
}

func (m Model) Init() tea.Cmd {
	return app.GetCategories()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.table = m.table.WithTargetWidth(m.width).WithPageSize(m.height - 6)
	case app.GetCategoriesMsg:
		// close modal
		m.edit = false
	case app.NewCategoriesMsg:
		categories := msg.Categories
		rows := make([]table.Row, 0, len(categories))
		for _, category := range categories {
			newRow := table.NewRow(table.RowData{
				"id":    category.ID,
				"name":  category.Name,
				"color": category.Color,
			}).WithStyle(lipgloss.NewStyle().Background(lipgloss.Color(category.Color)))

			rows = append(rows, newRow)
		}

		m.table = m.table.WithRows(rows)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.edit {

				row := m.table.HighlightedRow().Data
				m.edit = true
				m.category = newCategoryEdit(row["id"].(uint), row["name"].(string), row["color"].(string))
				return m, m.category.Init()
			}
		}
	}

	if m.edit {
		m.category, cmd = m.category.Update(msg)
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
		body = m.category.View()
	}
	return styles.CategoriesPage.Height(m.height).Width(m.width).Render("Categories" + "\n" + body)
}

func New() (*Model, error) {
	t := table.New([]table.Column{
		table.NewFlexColumn("name", "Name", 1),
		table.NewFlexColumn("color", "Color", 1),
	}).Focused(true).
		SortByAsc("name").
		HighlightStyle(styles.RowHighlight)

	m := &Model{table: t}

	return m, nil
}
