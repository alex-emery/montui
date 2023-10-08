package nordigen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
)

type Nordigen struct {
	client ClientWithResponsesInterface
	tokens *TokenStore
}

type TokenStore struct {
	accessToken   string
	accessExpiry  int
	refreshToken  string
	refreshExpiry int
}

const baseURL = "https://bankaccountdata.gocardless.com"

func newTokenStore(secretID, secretKey string) (*TokenStore, error) {
	client, err := NewClientWithResponses(baseURL)
	if err != nil {
		return nil, err
	}

	req := JWTObtainPairRequest{
		SecretId:  secretID, //nolint:revive //this is a generated field
		SecretKey: secretKey,
	}

	tokenRes, err := client.ObtainNewAccessrefreshTokenPairWithResponse(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if tokenRes.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", tokenRes.HTTPResponse.StatusCode)
	}

	tokens := tokenRes.JSON200

	return &TokenStore{
		accessToken:   *tokens.Access,
		accessExpiry:  *tokens.AccessExpires,
		refreshToken:  *tokens.Refresh,
		refreshExpiry: *tokens.RefreshExpires,
	}, nil
}

func (tokens *TokenStore) authProvider(_ context.Context, req *http.Request) error {
	// we need to calculate if access token has expired and auth here
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.accessToken))
	return nil
}

func New(secretID, secretKey string) (*Nordigen, error) {
	tokens, err := newTokenStore(secretID, secretKey)
	if err != nil {
		return nil, err
	}

	client, err := NewClientWithResponses(baseURL, WithRequestEditorFn(tokens.authProvider))
	if err != nil {
		return nil, err
	}

	nordigen := &Nordigen{
		client: client,
		tokens: tokens,
	}

	return nordigen, nil
}

func (c *Nordigen) ListBanks() ([]Integration, error) {
	gb := "GB"
	institutionResp, err := c.client.RetrieveAllSupportedInstitutionsInAGivenCountryWithResponse(context.TODO(), &RetrieveAllSupportedInstitutionsInAGivenCountryParams{
		Country: &gb,
	})

	if err != nil {
		return nil, err
	}

	if institutionResp.JSON200 == nil {
		return nil, fmt.Errorf("invalid status code %d", institutionResp.StatusCode())
	}

	return *institutionResp.JSON200, nil
}

func (c *Nordigen) CreateAgreement(ctx context.Context, institutionID string) (*uuid.UUID, error) {
	scopes := []interface{}{"transactions"}
	historicalDays := 90
	accessValidForDays := 30

	acceptEUARes, err := c.client.CreateEUAWithResponse(ctx, EndUserAgreementRequest{
		InstitutionId:      institutionID, //nolint:revive // generated field
		AccessScope:        &scopes,
		MaxHistoricalDays:  &historicalDays,
		AccessValidForDays: &accessValidForDays,
	})
	if err != nil {
		return nil, err
	}

	if acceptEUARes.JSON201 == nil {
		return nil, fmt.Errorf("invalid status code %d", acceptEUARes.StatusCode())
	}

	return acceptEUARes.JSON201.Id, nil
}

func (c *Nordigen) InitiateRequsition(ctx context.Context, ref string, agreementID *uuid.UUID, instituteID string) (*SpectacularRequisition, error) {
	hostRedirect := "http://localhost:3000"
	lang := "EN"
	redirect := true

	reqResp, err := c.client.CreateRequisitionWithFormdataBodyWithResponse(ctx, CreateRequisitionFormdataRequestBody{
		Redirect:          &hostRedirect,
		InstitutionId:     instituteID, //nolint:revive // generated field
		Agreement:         agreementID,
		Reference:         &ref,
		UserLanguage:      &lang,
		RedirectImmediate: &redirect,
	})

	if err != nil {
		return nil, err
	}

	if reqResp.JSON201 == nil {
		return nil, fmt.Errorf("unexpected status code")
	}

	return reqResp.JSON201, nil
}

func (c *Nordigen) GetAccounts(ctx context.Context, reqID uuid.UUID) ([]uuid.UUID, error) {
	reqIDResp, err := c.client.RequisitionByIdWithResponse(ctx, reqID) //nolint:revive // generated field
	if err != nil {
		return nil, err
	}

	if reqIDResp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code")
	}

	return *reqIDResp.JSON200.Accounts, nil
}

func WaitForRedirect(_ uuid.UUID) error {
	shutdown := make(chan interface{}, 1)
	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		// incoming := r.URL.Query().Get("ref")
		// TODO: validate the incoming ref and error if it's wrong
		shutdown <- struct{}{}
	}

	server := &http.Server{Addr: ":3000", Handler: http.HandlerFunc(handleRequest)}
	go server.ListenAndServe() //nolint: errcheck

	<-shutdown
	server.Shutdown(context.Background()) //nolint: errcheck

	return nil
}

func (c *Nordigen) GetTransactions(ctx context.Context, accountID string, dateFrom, dateTo *time.Time) (*BankTransactionStatusSchemaTransactions, error) {
	options := RetrieveAccountTransactionsParams{}
	if dateFrom != nil {
		options.DateFrom = &types.Date{Time: *dateFrom}
	}

	if dateTo != nil {
		options.DateTo = &types.Date{Time: *dateTo}
	}

	res, err := c.client.RetrieveAccountTransactionsWithResponse(ctx, uuid.MustParse(accountID), &options)
	if err != nil {
		return nil, err
	}

	return &res.JSON200.Transactions, nil
}

func FakeData(path string) (*BankTransactionStatusSchema, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction file %w", err)
	}

	responseStruct := struct {
		Transactions *BankTransactionStatusSchema `json:"transactions"`
	}{}

	err = json.Unmarshal(f, &responseStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions %w", err)
	}

	return responseStruct.Transactions, nil
}
