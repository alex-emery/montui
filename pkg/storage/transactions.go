package storage

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type transactionStore struct {
	db *gorm.DB
}

var _ TransactionGetter = &transactionStore{}

func (s *transactionStore) Insert(transactions ...*Transaction) error {
	result := s.db.Create(transactions)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%w\n%w", ErrUniqueConstraintFailed, result.Error)
		}
	}

	return result.Error
}

func (s *transactionStore) Get(transaction *Transaction) error {

	result := s.db.Preload("Category").First(transaction)
	return result.Error
}

func (s *transactionStore) List() ([]*Transaction, error) {
	transactions := []*Transaction{}
	result := s.db.Preload("Category").Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (s *transactionStore) Update(transaction *Transaction) error {
	result := s.db.Model(&transaction).Updates(transaction)

	return result.Error
}

func (s *transactionStore) SetCategory(transactionID uint, categoryID uint) error {
	transaction := Transaction{
		Model: gorm.Model{
			ID: transactionID,
		},
		CategoryID: &categoryID,
	}

	res := s.db.Model(&transaction).Update("category_id", categoryID)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
