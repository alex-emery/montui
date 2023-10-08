package app

import (
	"github.com/alex-emery/montui/pkg/nordigen"
	"github.com/alex-emery/montui/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func GetTransactions() tea.Cmd {
	return func() tea.Msg {
		return GetTransactionsMsg{}
	}
}

func GetCategories() tea.Cmd {
	return func() tea.Msg {
		return GetCategoriesMsg{}
	}
}

func SetCategory(id, categoryID uint) tea.Cmd {
	return func() tea.Msg {
		return SetCategoryMsg{
			TransactionID: id,
			CategoryID:    categoryID,
		}
	}
}

func GetAccounts() tea.Cmd {
	return func() tea.Msg {
		return GetAccountsMsg{}
	}
}

func GetBanks() tea.Cmd {
	return func() tea.Msg {
		return GetBanksMsg{}
	}
}

func SelectBank(bank nordigen.Integration) tea.Cmd {
	return func() tea.Msg {
		return SelectBankMsg{
			Bank: bank,
		}
	}
}

func CreateLink(bankID string) tea.Cmd {
	return func() tea.Msg {
		return CreateLinkMsg{
			BankID: bankID,
		}
	}
}

func UpdateCategory(category *storage.Category) tea.Cmd {
	return func() tea.Msg {
		return UpdateCategoryMsg{Category: category}
	}
}

func GetRules() tea.Cmd {
	return func() tea.Msg {
		return GetRulesMsg{}
	}
}

func UpdateRule(rule *storage.Rule) tea.Cmd {
	return func() tea.Msg {
		return UpdateRuleMsg{
			Rule: rule,
		}
	}
}

func CreateRule(rule *storage.Rule) tea.Cmd {
	return func() tea.Msg {
		return CreateRuleMsg{
			Rule: rule,
		}
	}
}

func DeleteRule(id uint) tea.Cmd {
	return func() tea.Msg {
		return DeleteRuleMsg{
			ID: id,
		}
	}
}
func SendError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{
			Err: err,
		}
	}
}
