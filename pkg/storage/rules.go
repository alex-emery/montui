package storage

import (
	"gorm.io/gorm"
)

//nolint:gochecknoglobals
var rulesMap = map[string]string{
	"ADOBESYSTEM":                  "Other",
	"AMAZON PRIME":                 "Other",
	"AMAZON.CO.UK":                 "Shopping",
	"AMZNMKTPLACE":                 "Shopping",
	"APPLE.COM":                    "Shopping",
	"ARGOSDIRECT":                  "Shopping",
	"AUTO TRADER":                  "Other",
	"AVIVA ZERO":                   "Other",
	"BEEFY BOYS":                   "Dining",
	"BRAID BARBERS":                "Other",
	"CLOUDFLARE":                   "Other",
	"DELIVEROO":                    "Dining",
	"EBAY":                         "Shopping",
	"JUST EAT.CO.UK":               "Dining",
	"M&S":                          "Groceries",
	"MCDONALDS":                    "Dining",
	"MFG":                          "Fuel",
	"MY NCP PARKING":               "Transport",
	"NO38 PARK HOTEL TRUST":        "Other",
	"PAYMENT RECEIVED - THANK YOU": "Income",
	"OCUKLIMITED":                  "Shopping",
	"SCANCOMPUTE":                  "Shopping",
	"SAINSBURY'S":                  "Groceries",
	"STARBUCKS":                    "Dining",
	"STEAM GAMES":                  "Other",
	"TFL TRAVEL CHARGE":            "Transport",
	"W H SMITH":                    "Shopping",
	"LULULEMON ATHLETICA":          "Shopping",
	"APPLE STORE":                  "Shopping",
	"TESCO STORES":                 "Groceries",
	"TESCO PETROL":                 "Fuel",
	"SQUARE MILE COFFE":            "Groceries",
	"CURRYS":                       "Shopping",
	"WELCOME BREAK":                "Food",
}

type ruleStore struct {
	db *gorm.DB
}

var _ RuleGetter = &ruleStore{}

func (s *ruleStore) Get(rule *Rule) error {
	result := s.db.First(&rule)
	return result.Error
}

func (s *ruleStore) List() ([]*Rule, error) {
	var rules []*Rule
	res := s.db.Preload("Category").Find(&rules)

	if res.Error != nil {
		return nil, res.Error
	}

	return rules, nil
}

func (s *ruleStore) Insert(rules ...*Rule) error {
	result := s.db.Create(rules)
	return result.Error
}

func (s *ruleStore) Update(rule *Rule) error {
	result := s.db.Model(&rule).Updates(rule)

	return result.Error
}
