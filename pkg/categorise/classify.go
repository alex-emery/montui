package categorise

import (
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"go.uber.org/zap"
)

type Categorise struct {
	storage storage.Storage
	logger  *zap.Logger
}

func New(storage storage.Storage, logger *zap.Logger) *Categorise {
	return &Categorise{
		storage: storage,
		logger:  logger,
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
				c.logger.Debug("found category match", zap.String("transaction", transactions[i].Description), zap.String("category", rule.Category.Name))
				categoryID := rule.CategoryID
				transactions[i].CategoryID = &categoryID

				break
			}
		}
	}

	return nil
}
