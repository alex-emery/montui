package ui

import (
	"fmt"

	"github.com/alex-emery/montui/internal/ui/accounts"
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/internal/ui/categories"
	"github.com/alex-emery/montui/internal/ui/rules"
	"github.com/alex-emery/montui/internal/ui/styles"
	"github.com/alex-emery/montui/internal/ui/transactions"
	"github.com/alex-emery/montui/pkg/montui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BorderStyle(lipgloss.NormalBorder()).
// BorderForeground(lipgloss.Color("240"))

type VisiblePage int

const (
	AccountPage = iota
	TransactionPage
	CategoriesPage
	RulesPage
)

type model struct {
	height   int
	width    int
	pages    []tea.Model
	selected VisiblePage
	app      *app.App
	errMsg   string
}

func (m model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.pages))
	for _, page := range m.pages {
		cmds = append(cmds, page.Init())
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case app.ErrorMsg:
		m.errMsg = msg.Err.Error()
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		msg.Height -= 2
		msg.Width -= 2

		// now we update the children with our new
		for index := range m.pages {
			m.pages[index], cmd = m.pages[index].Update(msg)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)

	case tea.KeyMsg: // keep this minimal, if they exist here they can't be used in nested models.
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.selected = VisiblePage(int(m.selected+1) % len(m.pages))
			return m, nil
		}

		m.pages[m.selected], cmd = m.pages[m.selected].Update(msg)
		cmds = append(cmds, cmd)
	}

	if _, ok := msg.(tea.KeyMsg); !ok {
		for index := range m.pages {
			m.pages[index], cmd = m.pages[index].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	cmds = append(cmds, m.app.Update(msg)) // this handles all interactions between UI and montui

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return styles.MainPage.
		Height(m.height).
		Width(m.width).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, m.pages[m.selected].View(),
				styles.ErrorBar.Render(m.errMsg),
			),
		)
}

func New(montui *montui.Montui) (*model, error) {
	transactionModel, err := transactions.New()
	if err != nil {
		return nil, err
	}

	categoriesModel, err := categories.New()
	if err != nil {
		return nil, err
	}

	accountModel := accounts.New()

	rulesModel := rules.New()

	pages := []tea.Model{
		accountModel,
		transactionModel,
		categoriesModel,
		rulesModel,
	}

	m := &model{app: app.New(montui), pages: pages, selected: 0}

	return m, nil
}

func Run(m *model) error {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}

	return nil

}
