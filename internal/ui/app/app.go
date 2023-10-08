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
		return Future(func(channel chan tea.Msg) {
			transactions, _ := c.montui.FetchTransactions(c.ctx, "", nil, nil)
			channel <- NewTransactionsMsg{Transactions: transactions}
		})

	case GetTransactionsMsg:
		return Future(func(channel chan tea.Msg) {
			transactions, _ := c.montui.GetTransactions("", nil, nil)
			channel <- NewTransactionsMsg{Transactions: transactions}
		})

	case GetCategoriesMsg:
		return Future(func(channel chan tea.Msg) {
			categories, err := c.montui.GetCategories()
			if err != nil {
				channel <- NewCategoriesMsg{Categories: categories}
			} else {
				channel <- NewCategoriesMsg{Categories: categories}
			}
		})

	case SetCategoryMsg:
		return Future(func(channel chan tea.Msg) {
			_, err := c.montui.SetCategory(msg.TransactionID, msg.CategoryID)
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- GetTransactionsMsg{}
			}
		})
	case UpdateCategoryMsg:
		return Future(func(channel chan tea.Msg) {
			err := c.montui.UpdateCategory(msg.Category)
			if err != nil {
				channel <- NewErrorMsg(err)
				return
			}
			channel <- GetCategoriesMsg{}
		})
	case GetAccountsMsg:
		return Future(func(channel chan tea.Msg) {
			accounts, _ := c.montui.GetAccounts()
			channel <- NewAccountsMsg{Accounts: accounts}
		})
	case GetBanksMsg:
		return Future(func(channel chan tea.Msg) {
			banks, _ := c.montui.ListBanks()
			channel <- NewBanksMsg{Banks: banks}
		})
	case CreateLinkMsg:
		return Future(func(channel chan tea.Msg) {
			c.montui.Link(c.ctx, msg.BankID)
			channel <- LinkReadyMsg{}
		})
	case GetRulesMsg:
		return Future(func(channel chan tea.Msg) {
			rules, err := c.montui.ListRules()
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- NewRulesMsg{Rules: rules}
			}
		})
	case UpdateRuleMsg:
		return Future(func(channel chan tea.Msg) {
			err := c.montui.UpdateRule(msg.Rule)
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- GetRulesMsg{}
			}
		})
	case CreateRuleMsg:
		return Future(func(channel chan tea.Msg) {
			err := c.montui.CreateRule(msg.Rule)
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- GetRulesMsg{}
			}
		})
	case DeleteRuleMsg:
		return Future(func(channel chan tea.Msg) {
			err := c.montui.DeleteRule(msg.ID)
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- GetRulesMsg{}
			}
		})
	case CategoriseMsg:
		return Future(func(channel chan tea.Msg) {
			err := c.montui.CategoriseTransactions()
			if err != nil {
				channel <- NewErrorMsg(err)
			} else {
				channel <- GetTransactionsMsg{}
			}

		})
	}

	return nil
}

type FutureFn func(channel chan tea.Msg)

func Future(fn FutureFn) tea.Cmd {
	sub := make(chan tea.Msg)
	go fn(sub)

	return func() tea.Msg {
		return <-sub
	}
}
