package accounts

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/pkg/nordigen"
	tea "github.com/charmbracelet/bubbletea"
)

type LinkModel struct {
	bank nordigen.Integration
}

func (m LinkModel) Init() tea.Cmd {
	return app.CreateLink(m.bank.ID)
}

func (m LinkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LinkModel) View() string {
	return "Waiting for callback"
}
