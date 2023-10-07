package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountStore struct {
	db *gorm.DB
}

var _ AccountGetter = &accountStore{}

func (s *accountStore) Insert(accounts ...Account) error {
	res := s.db.Save(accounts)

	return res.Error
}

func (s *accountStore) List() ([]Account, error) {
	accounts := []Account{}
	result := s.db.Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (s *accountStore) Delete(ids ...uuid.UUID) error {
	accounts := make([]Account, 0, len(ids))
	for _, id := range ids {
		accounts = append(accounts, Account{
			ID: id,
		})
	}

	res := s.db.Delete(accounts)

	return res.Error
}
