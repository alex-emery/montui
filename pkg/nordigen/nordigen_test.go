package nordigen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mc *mockClient) Do(req *http.Request) (*http.Response, error) {
	return mc.DoFunc(req)
}

func TestBasic(t *testing.T) {
	mc := &mockClient{}

	mockResponse := BankTransactionStatusSchema{
		Transactions: BankTransactionStatusSchemaTransactions{
			Booked: []TransactionSchema{
				{
					RemittanceInformationUnstructured: func(a string) *string { return &a }("test 123"),
				},
			},
		},
	}

	responseBytes, err := json.Marshal(mockResponse)
	if err != nil {
		t.Fatal("failed to marshal mock response")
	}

	mc.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBytes)),
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		}, nil
	}

	client := ClientWithResponses{
		ClientInterface: &Client{
			Client: mc,
		},
	}

	accountID := uuid.New().String()
	nordigen := Nordigen{client: &client}

	transactions, err := nordigen.GetTransactions(context.Background(), accountID, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions.Booked) == 0 {
		t.Fatal(fmt.Errorf("no transactions found"))
	}
}
