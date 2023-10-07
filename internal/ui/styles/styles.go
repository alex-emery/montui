package styles

import "github.com/charmbracelet/lipgloss"

var (
	MainPage = lipgloss.NewStyle()

	AccountPage = lipgloss.NewStyle()

	TransactionPage = lipgloss.NewStyle()

	CategoriesPage = lipgloss.NewStyle()

	ErrorBar = lipgloss.NewStyle().Height(2)

	Table = lipgloss.NewStyle()

	TableRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("f8fe6c"))
	TableRowAlt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f092ff"))

	Overlay = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	RowHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#f88"))
)
