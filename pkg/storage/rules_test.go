package storage_test

import (
	"testing"

	"github.com/alex-emery/montui/pkg/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestRulesCRUD(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared", zap.NewNop())

	if err != nil {
		t.Fatal(err)
	}

	ruleStore := store.Rules()

	err = ruleStore.Insert(&storage.Rule{
		Pattern:    "cows go moo",
		CategoryID: 0,
	})

	if err != nil {
		t.Fatal("failed to create rule", err)
	}

	queryRule := storage.Rule{Model: gorm.Model{
		ID: 0,
	}}

	err = ruleStore.Get(&queryRule)
	if err != nil {
		t.Fatal("failed to create rule", err)
	}
}
