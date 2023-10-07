package nordigen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mc *mockClient) Do(req *http.Request) (*http.Response, error) {
	return mc.DoFunc(req)
}
func TestBasic(t *testing.T) {

	mc := &mockClient{}

	file, err := os.ReadFile("./testdata/transactions")
	if err != nil {
		t.Fatal(err)
	}

	mc.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(file)),
		}, nil
	}
	client := ClientWithResponses{
		ClientInterface: &Client{
			Client: mc,
		},
	}

	accountId := os.Getenv("ACCOUNT_ID")
	nordigen := Nordigen{client: &client}
	transactions, err := nordigen.GetTransactions(context.Background(), accountId, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions.Booked) == 0 {
		t.Fatal(fmt.Errorf("no transactions found"))
	}
}
