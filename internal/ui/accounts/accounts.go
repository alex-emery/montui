package accounts

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/internal/ui/styles"
	"github.com/alex-emery/montui/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type Accounts struct {
	accounts []storage.Account
	workflow tea.Model
	width    int
	height   int
}

func New() *Accounts {
	return &Accounts{}
}

func (m Accounts) Init() tea.Cmd {
	return app.GetAccounts()
}

func (m Accounts) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.workflow != nil {
		if _, ok := msg.(app.LinkReadyMsg); ok {
			m.workflow = nil
			return m, app.GetAccounts()
		}

		workflow, cmd := m.workflow.Update(msg)
		m.workflow = workflow
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.workflow = newAuthWorkflow()
			return m, m.workflow.Init()
		}
	case app.NewAccountsMsg:
		m.accounts = msg.Accounts
	}

	return m, nil
}

func (m Accounts) Footer() string {
	return "Press n to link new account"
}
func (m Accounts) View() string {
	if m.workflow != nil {
		return m.workflow.View()
	}
	accountList := "Accounts"
	for _, account := range m.accounts {
		accountList += "\n" + account.ID.String()
	}

	accountList += "\n" + m.Footer()
	return styles.AccountPage.Height(m.height).Render(accountList)
}
