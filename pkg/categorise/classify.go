package categorise

import (
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Categorise struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Categorise {
	return &Categorise{
		storage: storage,
	}
}

// Goes through every "rule" and tries to substring match it
// to the transaction.
// TODO: set the category to other by default.
func (c *Categorise) Categorise(transactions ...*storage.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	rules, err := c.storage.Rules().List()
	if err != nil {
		return err
	}

	for i := range transactions {
		for _, rule := range rules {
			if fuzzy.MatchFold(rule.Pattern, transactions[i].Description) {
				categoryID := rule.CategoryID
				transactions[i].CategoryID = &categoryID

				break
			}
		}
	}

	return nil
}
