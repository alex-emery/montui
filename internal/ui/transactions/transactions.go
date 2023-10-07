package transactions

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/internal/ui/helpers"
	"github.com/alex-emery/montui/internal/ui/styles"
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type TransactionModel struct {
	height         int
	width          int
	table          table.Model
	edit           bool
	categoryPicker categoryPicker
	spinner        spinner.Model
	loading        bool
}

func (m TransactionModel) Init() tea.Cmd {
	return tea.Batch(app.GetTransactions(), m.categoryPicker.Init())
}

func (m TransactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.edit {
		// basically just waiting for the "select message"
		// probably add some selector to the message so
		// check if the message is for us.
		if msg, ok := msg.(CategorySelect); ok {
			transaction := m.table.HighlightedRow().Data
			id := transaction["id"].(uint)
			categoryID := msg.ID

			m.edit = false
			return m, app.SetCategory(id, categoryID)
		}

		newModel, cmd := m.categoryPicker.Update(msg)
		m.categoryPicker = newModel
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.table = m.table.WithTargetWidth(m.width).WithPageSize(m.height - 7) // 6 = 3 row header + 3 footer + 1 title
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.edit = true
			return m, app.GetCategories()
		case "r":
			m.loading = true
			cmds = append(cmds, func() tea.Msg {
				return app.FetchTransactionsMsg{}
			}, m.spinner.Tick)
		}
	case app.NewTransactionsMsg:
		rows := TransactionsToRows(msg.Transactions)
		m.table = m.table.WithRows(rows)
		m.loading = false
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
func (m TransactionModel) View() string {

	title := "Transactions"
	if m.loading {
		title += " fetching " + m.spinner.View()
	}

	title += "\n"
	tableView := styles.TransactionPage.Height(m.height).Render(title + m.table.View())
	if m.edit {
		overlayText := styles.Overlay.Render(m.categoryPicker.View())
		return helpers.PlaceOverlay(m.width/2, m.height/2, overlayText, tableView, false)
	}

	return tableView
}

func TransactionsToRows(transactions []*storage.Transaction) []table.Row {

	rows := make([]table.Row, 0, len(transactions))

	for _, transaction := range transactions {
		var categoryID *uint
		var color = ""
		if transaction.Category != nil {
			categoryID = &transaction.Category.ID
			color = transaction.Category.Color
		}

		style := lipgloss.NewStyle().Background(lipgloss.Color(color))

		rows = append(rows,
			table.NewRow(table.RowData{
				"id":          transaction.ID,
				"date":        transaction.Date,
				"name":        transaction.Name,
				"description": transaction.Description,
				"amount":      transaction.Amount,
				"category":    safe(transaction.Category),
				"category_id": categoryID,
			}).WithStyle(style))
	}

	return rows
}

func safe(category *storage.Category) string { //todo: lol
	if category == nil {
		return ""
	}

	return category.Name
}

func New() (*TransactionModel, error) {

	t := table.New([]table.Column{
		table.NewColumn("date", "Date", 12),
		table.NewFlexColumn("name", "Name", 1),
		table.NewFlexColumn("description", "Description", 3),
		table.NewColumn("amount", "Amount", 8),
		table.NewColumn("category", "Category", 15),
	}).
		Focused(true).
		SortByDesc("date").
		HighlightStyle(styles.RowHighlight)

	m := &TransactionModel{table: t, categoryPicker: categoryPicker{}, spinner: spinner.New()}

	return m, nil
}
