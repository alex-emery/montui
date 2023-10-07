package categorise_test

import (
	"context"
	"testing"

	"github.com/alex-emery/montui/pkg/categorise"
	"github.com/alex-emery/montui/pkg/storage"
)

func TestCategoryise(t *testing.T) {
	store, err := storage.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	c := categorise.New(store)
	transactions := []*storage.Transaction{
		{
			Description: "AMAZON PRIME",
		},
	}
	c.Categorise(context.Background(), transactions...)

	if transactions[0].CategoryID == nil {
		t.Fatalf("failed to classify %v", transactions[0])
	}
}
