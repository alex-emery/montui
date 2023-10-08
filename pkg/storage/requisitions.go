package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type requisitionStore struct {
	db *gorm.DB
}

var _ RequsitionGetter = &requisitionStore{}

func (s *requisitionStore) Insert(requisitions ...*Requisition) error {
	res := s.db.Save(requisitions)

	return res.Error
}

func (s *requisitionStore) List() ([]*Requisition, error) {
	requisitions := []*Requisition{}
	result := s.db.Preload("Accounts").Find(&requisitions)

	if result.Error != nil {
		return nil, result.Error
	}

	return requisitions, nil
}

func (s *requisitionStore) Update(requsition ...*Requisition) error {
	result := s.db.Save(requsition)

	return result.Error
}

func (s *requisitionStore) Delete(ids ...uuid.UUID) error {
	reqs := make([]Requisition, 0, len(ids))
	for _, id := range ids {
		reqs = append(reqs, Requisition{
			ID: id,
		})
	}

	res := s.db.Delete(reqs)

	return res.Error
}
