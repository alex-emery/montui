package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	InternalTransactionID string `gorm:"uniqueIndex"`
	Date                  string
	Name                  string
	Description           string
	Amount                string
	AccountID             uuid.UUID
	CategoryID            *uint
	Category              *Category
}

type Category struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex"`
	Color string
}

type Rule struct {
	gorm.Model
	Pattern    string
	CategoryID uint
	Category   Category
}

type Account struct {
	ID           uuid.UUID `gorm:"type:uuid;primarykey;"`
	RequsitionID uuid.UUID
	LastFetch    *time.Time
}

type Requisition struct {
	ID       uuid.UUID `gorm:"type:uuid;primarykey;"`
	Accounts []Account `gorm:"foreignKey:RequsitionID"`
}
