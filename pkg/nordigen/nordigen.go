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

func newTokenStore(secretId, secretKey string) (*TokenStore, error) {
	client, err := NewClient(baseURL)
	if err != nil {
		return nil, err
	}

	body := JWTObtainPairRequest{
		SecretId:  secretId,
		SecretKey: secretKey,
	}

	res, err := client.ObtainNewAccessrefreshTokenPair(context.Background(), body)
	if err != nil {
		return nil, err
	}

	tokenRes, err := ParseObtainNewAccessrefreshTokenPairResponse(res)
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

func (tokens *TokenStore) authProvider(ctx context.Context, req *http.Request) error {
	// we need to calculate if access token has expired and auth here
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.accessToken))
	return nil
}

func New(secretId, secretKey string) (*Nordigen, error) {

	tokens, err := newTokenStore(secretId, secretKey)
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
		InstitutionId:      institutionID,
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
		InstitutionId:     instituteID,
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

func (c *Nordigen) GetAccounts(ctx context.Context, reqId uuid.UUID) ([]uuid.UUID, error) {
	reqIdResp, err := c.client.RequisitionByIdWithResponse(ctx, reqId)
	if err != nil {
		return nil, err
	}

	if reqIdResp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code")
	}

	return *reqIdResp.JSON200.Accounts, nil

}
func WaitForRedirect(ref uuid.UUID) error {
	shutdown := make(chan interface{}, 1)
	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		// incoming := r.URL.Query().Get("ref")
		// TODO: validate the incoming ref and error if it's wrong
		shutdown <- struct{}{}
	}

	server := &http.Server{Addr: ":3000", Handler: http.HandlerFunc(handleRequest)}
	go server.ListenAndServe()

	<-shutdown
	server.Shutdown(context.Background())
	return nil
}

func (c *Nordigen) GetTransactions(ctx context.Context, accountId string, dateFrom, dateTo *time.Time) (*BankTransactionStatusSchema, error) {

	options := RetrieveAccountTransactionsParams{}
	if dateFrom != nil {
		options.DateFrom = &types.Date{Time: *dateFrom}
	}

	res, err := c.client.RetrieveAccountTransactionsWithResponse(ctx, uuid.MustParse(accountId), &options)
	if err != nil {
		return nil, err
	}

	return res.JSON200, nil

}

func FakeData(path string) (*BankTransactionStatusSchema, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction file %v", err)
	}

	responseStruct := struct {
		Transactions *BankTransactionStatusSchema `json:"transactions"`
	}{}

	err = json.Unmarshal(f, &responseStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions %v", err)
	}

	return responseStruct.Transactions, nil
}
