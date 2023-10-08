package montui

import (
	"context"
	"errors"
	"time"

	"github.com/alex-emery/montui/pkg/categorise"
	"github.com/alex-emery/montui/pkg/nordigen"
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"gorm.io/gorm"
)

type Montui struct {
	client   *nordigen.Nordigen
	store    storage.Storage
	classify *categorise.Categorise
}

func New(nordigenSecretID, nordigenSecretKey, dir string) (*Montui, error) {
	store, err := storage.New(dir)
	if err != nil {
		return nil, err
	}

	client, err := nordigen.New(nordigenSecretID, nordigenSecretKey)
	if err != nil {
		return nil, err
	}

	classify := categorise.New(store)

	return &Montui{store: store, client: client, classify: classify}, nil
}

func (s *Montui) GetAccounts() ([]*storage.Account, error) {
	return s.store.Accounts().List()
}

// Attempts to fetch transactions from nordigen and
// store them in the DB.
func (s *Montui) FetchTransactions(ctx context.Context, accountID string, dateTo, dateFrom *string) ([]*storage.Transaction, error) {
	// get all accounts
	accounts, err := s.store.Accounts().List()
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		transactions, err := s.client.GetTransactions(ctx, account.ID.String(), account.LastFetch, nil)
		if err != nil {
			return nil, err
		}

		safe := func(input *string) string {
			if input == nil {
				return ""
			}

			return *input
		}

		finalised := make([]*storage.Transaction, 0, len(transactions.Booked))

		var latest time.Time

		for _, transaction := range transactions.Booked {
			dbTransaction := &storage.Transaction{
				AccountID:             account.ID,
				InternalTransactionID: safe(transaction.TransactionId),
				Name:                  safe(transaction.CreditorName),
				Date:                  safe(transaction.BookingDate),
				Description:           safe(transaction.RemittanceInformationUnstructured),
				Amount:                transaction.TransactionAmount.Amount,
			}

			s.classify.Categorise(dbTransaction) //nolint:errcheck

			if dbTransaction.CategoryID != nil {
				s.store.Transactions().SetCategory(dbTransaction.ID, *dbTransaction.CategoryID) //nolint:errcheck
			}

			finalised = append(finalised, dbTransaction)

			// don't care if this errors
			parsedBookingDate, _ := time.Parse("2006-01-02", safe(transaction.BookingDate))
			if parsedBookingDate.After(latest) {
				latest = parsedBookingDate
			}
		}

		err = s.store.Transactions().Insert(finalised...)
		if err != nil && !errors.Is(err, storage.ErrUniqueConstraintFailed) {
			return s.GetTransactions(accountID, dateTo, dateFrom)
		}

		account.LastFetch = &latest
		err = s.store.Accounts().Insert(account)

		if err != nil {
			return nil, err
		}
	}

	return s.GetTransactions(accountID, dateTo, dateFrom)
}

// Fetches transactions from storage.
func (s *Montui) GetTransactions(accountID string, dateTo, dateFrom *string) ([]*storage.Transaction, error) { //nolint:revive //these are in TODO
	return s.store.Transactions().List()
}

// Fetches categories from storage.
func (s *Montui) GetCategories() ([]*storage.Category, error) {
	return s.store.Categories().List()
}

// Sets the category for a specific transaction.
func (s *Montui) SetCategory(transactionID uint, categoryID uint) (*storage.Transaction, error) {
	transaction := &storage.Transaction{
		Model: gorm.Model{
			ID: transactionID,
		},
		CategoryID: &categoryID,
	}

	err := s.store.Transactions().Update(transaction)
	if err != nil {
		return nil, err
	}

	// and get the update so we populate the category field
	return transaction, nil
}

func (s *Montui) UpdateCategory(category *storage.Category) error {
	return s.store.Categories().Update(category)
}

func (s *Montui) ListRules() ([]*storage.Rule, error) {
	return s.store.Rules().List()
}

func (s *Montui) findCategoryForRule(rule *storage.Rule) error {
	category := &storage.Category{
		Name: rule.Category.Name,
	}

	err := s.store.Categories().Get(category)
	if err != nil {
		return err
	}

	rule.CategoryID = category.ID
	rule.Category = *category

	return nil
}

func (s *Montui) CreateRule(rule *storage.Rule) error {
	if err := s.findCategoryForRule(rule); err != nil {
		return err
	}

	return s.store.Rules().Insert(rule)
}

func (s *Montui) UpdateRule(rule *storage.Rule) error {
	if err := s.findCategoryForRule(rule); err != nil {
		return err
	}

	return s.store.Rules().Update(rule)
}

func (s *Montui) DeleteRule(id uint) error {
	return s.store.Rules().Delete(id)
}

func (s *Montui) ListBanks() ([]nordigen.Integration, error) {
	return s.client.ListBanks()
}

func (s *Montui) Link(ctx context.Context, institutionID string) error {
	ref := uuid.New()
	req, err := s.client.InitiateRequsition(ctx, ref.String(), nil, institutionID)

	if err != nil {
		return err
	}

	_ = browser.OpenURL(*req.Link) //TODO: expose the link to the user without relying on this.

	err = nordigen.WaitForRedirect(ref)
	if err != nil {
		return err
	}

	requisition := storage.Requisition{ID: *req.Id}

	err = s.store.Requisitions().Insert(&requisition)
	if err != nil {
		return err
	}

	accounts, err := s.client.GetAccounts(ctx, *req.Id)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		requisition.Accounts = append(requisition.Accounts,
			storage.Account{
				ID:           account,
				RequsitionID: *req.Id,
			})
	}

	err = s.store.Requisitions().Update(&requisition)
	if err != nil {
		return err
	}

	return nil
}

func (s *Montui) CategoriseTransactions() error {
	transactions, err := s.store.Transactions().List()
	if err != nil {
		return err
	}

	err = s.classify.Categorise(transactions...)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		err = s.store.Transactions().Update(transaction)
		if err != nil {
			return err
		}
	}

	return nil
}
