package app

import (
	"context"

	"github.com/alex-emery/montui/pkg/montui"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	montui *montui.Montui
	ctx    context.Context
}

func New(montui *montui.Montui) *App {
	return &App{
		montui: montui,
		ctx:    context.Background(),
	}
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case FetchTransactionsMsg:
		return func() tea.Msg {
			transactions, err := c.montui.FetchTransactions(c.ctx, "", nil, nil)
			if err != nil {
				return NewErrorMsg(err)
			}
			return NewTransactionsMsg{Transactions: transactions}
		}
	case GetTransactionsMsg:
		return func() tea.Msg {
			transactions, err := c.montui.GetTransactions("", nil, nil)
			if err != nil {
				return NewErrorMsg(err)
			}
			return NewTransactionsMsg{Transactions: transactions}
		}

	case GetCategoriesMsg:
		return func() tea.Msg {
			categories, err := c.montui.GetCategories()
			if err != nil {
				return NewCategoriesMsg{Categories: categories}
			} else {
				return NewCategoriesMsg{Categories: categories}
			}
		}

	case SetCategoryMsg:
		return func() tea.Msg {
			categories, err := c.montui.GetCategories()
			if err != nil {
				return NewCategoriesMsg{Categories: categories}
			} else {
				return NewCategoriesMsg{Categories: categories}
			}
		}
	case UpdateCategoryMsg:
		return func() tea.Msg {
			err := c.montui.UpdateCategory(msg.Category)
			if err != nil {
				return NewErrorMsg(err)
			}
			return GetCategoriesMsg{}
		}
	case GetAccountsMsg:
		return func() tea.Msg {
			accounts, err := c.montui.GetAccounts()
			if err != nil {
				return NewErrorMsg(err)
			}
			return NewAccountsMsg{Accounts: accounts}
		}
	case GetBanksMsg:
		return func() tea.Msg {
			banks, err := c.montui.ListBanks()
			if err != nil {
				return NewErrorMsg(err)
			}
			return NewBanksMsg{Banks: banks}
		}
	case CreateLinkMsg:
		return func() tea.Msg {
			c.montui.Link(c.ctx, msg.BankID)
			return LinkReadyMsg{}
		}
	case GetRulesMsg:
		return func() tea.Msg {
			rules, err := c.montui.ListRules()
			if err != nil {
				return NewErrorMsg(err)
			}
			return NewRulesMsg{Rules: rules}
		}
	case UpdateRuleMsg:
		return func() tea.Msg {
			err := c.montui.UpdateRule(msg.Rule)
			if err != nil {
				return NewErrorMsg(err)
			} else {
				return GetRulesMsg{}
			}
		}
	case CreateRuleMsg:
		return func() tea.Msg {
			err := c.montui.CreateRule(msg.Rule)
			if err != nil {
				return NewErrorMsg(err)
			}
			return GetRulesMsg{}
		}
	case DeleteRuleMsg:
		return func() tea.Msg {
			err := c.montui.DeleteRule(msg.ID)
			if err != nil {
				return NewErrorMsg(err)
			}
			return GetRulesMsg{}
		}
	case CategoriseMsg:
		return func() tea.Msg {
			err := c.montui.CategoriseTransactions()
			if err != nil {
				return NewErrorMsg(err)
			}
			return GetTransactionsMsg{}
		}
	}

	return nil
}
