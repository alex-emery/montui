package categorise_test

import (
	"testing"

	"github.com/alex-emery/montui/pkg/categorise"
	"github.com/alex-emery/montui/pkg/storage"
	"go.uber.org/zap"
)

func TestCategoryise(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared", zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	c := categorise.New(store, zap.NewNop())
	transactions := []*storage.Transaction{
		{
			Description: "AMAZON PRIME",
		},
	}

	err = c.Categorise(transactions...)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	if transactions[0].CategoryID == nil {
		t.Fatalf("failed to classify %v", transactions[0])
	}
}
