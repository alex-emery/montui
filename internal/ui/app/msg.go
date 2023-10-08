package app

import (
	"github.com/alex-emery/montui/pkg/nordigen"
	"github.com/alex-emery/montui/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type GetTransactionsMsg struct{} //TODO: add in filter options
type FetchTransactionsMsg struct{}

type GetCategoriesMsg struct{}

type NewTransactionsMsg struct {
	Transactions []*storage.Transaction
}

type NewCategoriesMsg struct {
	Categories []*storage.Category
}

type SetCategoryMsg struct {
	TransactionID uint
	CategoryID    uint
}

type UpdateCategoryMsg struct {
	Category *storage.Category
}

type GetAccountsMsg struct{}

type NewAccountsMsg struct {
	Accounts []*storage.Account
}

type GetBanksMsg struct{}

type NewBanksMsg struct {
	Banks []nordigen.Integration
}
type ErrorMsg struct {
	Err error
}

type SelectBankMsg struct {
	Bank nordigen.Integration
}

type CreateLinkMsg struct {
	BankID string
}

type LinkReadyMsg struct {
}

type GetRulesMsg struct{}

type NewRulesMsg struct {
	Rules []*storage.Rule
}

type UpdateRuleMsg struct {
	Rule *storage.Rule
}

type CreateRuleMsg struct {
	Rule *storage.Rule
}

type DeleteRuleMsg struct {
	ID uint
}

func NewErrorMsg(err error) tea.Msg {
	return ErrorMsg{Err: err}
}
