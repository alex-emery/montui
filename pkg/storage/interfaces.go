package storage

import "github.com/google/uuid"

type Storage interface {
	Transactions() TransactionGetter
	Categories() CategoryGetter
	Rules() RuleGetter
	Requisitions() RequsitionGetter
	Accounts() AccountGetter
}

type RequsitionGetter interface {
	Insert(...Requisition) error
	List() ([]Requisition, error)
	Delete(...uuid.UUID) error
	Update(...Requisition) error
}

type AccountGetter interface {
	Insert(...Account) error
	List() ([]Account, error)
	Delete(...uuid.UUID) error
}

type TransactionGetter interface {
	Insert(...*Transaction) error
	Get(id uint) (*Transaction, error)
	List() ([]*Transaction, error)
	Update(*Transaction) error
	SetCategory(transactionID uint, categoryID uint) error
}

type CategoryGetter interface {
	Get(name string) (*Category, error)
	List() ([]Category, error)
	Where(query interface{}, args ...interface{}) ([]Category, error) // not used.
	Update(category Category) error
}

type RuleGetter interface {
	Get(name string) (*Rule, error)
	List() ([]Rule, error)
}
