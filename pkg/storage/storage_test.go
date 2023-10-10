package storage_test

import (
	"fmt"
	"testing"

	"github.com/alex-emery/montui/pkg/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestTransactions(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared", zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	// create uncategorised transaction
	transaction := storage.Transaction{
		Description: "AMAZON PRIME",
		Amount:      "123",
	}

	err = store.Transactions().Insert(&transaction)
	if err != nil {
		t.Fatal(err)
	}

	insertedTransaction := storage.Transaction{
		Model: gorm.Model{
			ID: transaction.ID,
		},
	}

	err = store.Transactions().Get(&insertedTransaction)
	if err != nil {
		t.Fatal(err)
	}

	if insertedTransaction.Description != "AMAZON PRIME" {
		t.Fatal(fmt.Errorf("model is empty %v", insertedTransaction))
	}

	category := storage.Category{
		Name: "Shopping",
	}

	err = store.Categories().Get(&category)
	if err != nil {
		t.Fatal(err)
	}

	// update transaction with category
	insertedTransaction.CategoryID = &category.ID

	err = store.Transactions().Update(&insertedTransaction)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Transactions().Get(&insertedTransaction)
	if err != nil {
		t.Fatal(err)
	}

	if *insertedTransaction.CategoryID != category.ID {
		t.Fatal(fmt.Errorf("failed to update transaction with category"))
	}
}

func TestAccountsAndReqs(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared", zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	accounts := []*storage.Account{
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
	}

	err = store.Accounts().Insert(accounts...)
	if err != nil {
		t.Fatal(err)
	}

	found, _ := store.Accounts().List()
	if len(found) != 3 {
		t.Fatalf("expected 3 accounts but found %d", len(found))
	}

	err = store.Accounts().Delete(accounts[0].ID)
	if err != nil {
		t.Fatalf("unexpected error whilst deleting account. %v", err)
	}

	found, _ = store.Accounts().List()

	for _, account := range found {
		if account.ID == accounts[0].ID {
			t.Fatal("deleted account still present")
		}
	}

	newReq := storage.Requisition{
		ID: uuid.New(),
		Accounts: []storage.Account{
			{
				ID: uuid.New(),
			},
		},
	}

	err = store.Requisitions().Insert(&newReq)
	if err != nil {
		t.Fatal(err)
	}

	reqs, _ := store.Requisitions().List()
	if reqs[0].ID.String() == reqs[0].Accounts[0].ID.String() {
		t.Fatal("IDS ARE THE SAME")
	}
}

func TestRequisitions(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared", zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*storage.Requisition{
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
	}

	err = store.Requisitions().Insert(reqs...)
	if err != nil {
		t.Fatal(err)
	}

	found, _ := store.Requisitions().List()
	if len(found) != 3 {
		t.Fatalf("expected 3 requisitions but found %d", len(found))
	}

	err = store.Requisitions().Delete(reqs[0].ID)
	if err != nil {
		t.Fatalf("unexpected error whislt deleting requisition. %v", err)
	}

	found, _ = store.Requisitions().List()

	for _, req := range found {
		if req.ID == reqs[0].ID {
			t.Fatal("deleted requisition still present")
		}
	}
}
