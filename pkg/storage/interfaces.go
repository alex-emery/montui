package storage

import "github.com/google/uuid"

type Storage interface {
	Transactions() TransactionGetter
	Categories() CategoryGetter
	Rules() RuleGetter
	Requisitions() RequsitionGetter
	Accounts() AccountGetter
}

// use pointers so defaulting will update the model

type RequsitionGetter interface {
	Insert(...*Requisition) error
	List() ([]*Requisition, error)
	Delete(...uuid.UUID) error
	Update(...*Requisition) error
}

type AccountGetter interface {
	Insert(...*Account) error
	List() ([]*Account, error)
	Delete(...uuid.UUID) error
}

type TransactionGetter interface {
	Insert(...*Transaction) error
	Get(*Transaction) error
	List() ([]*Transaction, error)
	Update(*Transaction) error
	SetCategory(transactionID uint, categoryID uint) error
}

type CategoryGetter interface {
	Get(*Category) error
	List() ([]*Category, error)
	Update(category *Category) error
}

type RuleGetter interface {
	Insert(...*Rule) error
	Get(*Rule) error
	List() ([]*Rule, error)
	Update(*Rule) error
}
