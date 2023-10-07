package accounts

import (
	"github.com/alex-emery/montui/internal/ui/app"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthWorkflow struct {
	steps tea.Model
}

func newAuthWorkflow() *AuthWorkflow {
	return &AuthWorkflow{
		steps: newInstitutionPicker(),
	}
}

func (m AuthWorkflow) Init() tea.Cmd {
	return m.steps.Init()
}

func (m AuthWorkflow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	model, cmd := m.steps.Update(msg)
	m.steps = model
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case app.SelectBankMsg:
		m.steps = LinkModel{bank: msg.Bank}
		return m, m.steps.Init()
	}
	return m, tea.Batch(cmds...)
}

func (m AuthWorkflow) View() string {
	return m.steps.View()
}
