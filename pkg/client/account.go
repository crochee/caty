// Date: 2021/10/24

// Package client
package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/crochee/lirity/client"
	"github.com/crochee/lirity/e"
	"github.com/json-iterator/go"

	"caty/pkg/service/account"
)

type Account interface {
	Register(ctx context.Context, request *account.CreateRequest) (*account.CreateResponseResult, error)
	List(ctx context.Context, request *account.RetrievesRequest) (*account.RetrieveResponses, error)
	Update(ctx context.Context, user *account.User, request *account.UpdateRequest) error
	Retrieve(ctx context.Context, user *account.User) (*account.RetrieveResponse, error)
	Delete(ctx context.Context, user *account.User) error
}

func NewAccount() Account {
	return &AccountClient{
		Client:     client.NewStandardClient(),
		API:        jsoniter.ConfigCompatibleWithStandardLibrary,
		URLHandler: NewURLHandler(),
	}
}

type AccountClient struct {
	client.Client
	jsoniter.API
	URLHandler
}

func (a *AccountClient) Register(ctx context.Context,
	request *account.CreateRequest) (*account.CreateResponseResult, error) {
	body, err := a.Marshal(request)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	if req, err = client.NewRequest(ctx, http.MethodPost, a.URL(ctx, "/v1/account"),
		body, a.Header(ctx)); err != nil {
		return nil, err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, e.From(response)
	}
	var result account.CreateResponseResult
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AccountClient) List(ctx context.Context,
	request *account.RetrievesRequest) (*account.RetrieveResponses, error) {
	params := url.Values{}
	if request.AccountID != "" {
		params.Add("account-id", request.AccountID)
	}
	if request.ID != "" {
		params.Add("id", request.ID)
	}
	if request.Account != "" {
		params.Add("account", request.Account)
	}
	if request.Email != "" {
		params.Add("email", request.Email)
	}

	req, err := client.NewRequest(ctx, http.MethodGet, a.URLWithQuery(ctx, "/v1/account", params),
		nil, a.Header(ctx))
	if err != nil {
		return nil, err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, e.From(response)
	}
	var result account.RetrieveResponses
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AccountClient) Update(ctx context.Context, user *account.User, request *account.UpdateRequest) error {
	body, err := a.Marshal(request)
	if err != nil {
		return err
	}
	var req *http.Request
	if req, err = client.NewRequest(ctx, http.MethodPatch, a.URL(ctx, "/v1/account/"+user.ID),
		body, a.Header(ctx)); err != nil {
		return err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return e.From(response)
}

func (a *AccountClient) Retrieve(ctx context.Context, user *account.User) (*account.RetrieveResponse, error) {
	req, err := client.NewRequest(ctx, http.MethodGet, a.URL(ctx, "/v1/account/"+user.ID), nil, a.Header(ctx))
	if err != nil {
		return nil, err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, e.From(response)
	}
	var result account.RetrieveResponse
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AccountClient) Delete(ctx context.Context, user *account.User) error {
	req, err := client.NewRequest(ctx, http.MethodDelete, a.URL(ctx, "/v1/account/"+user.ID), nil, a.Header(ctx))
	if err != nil {
		return err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return e.From(response)
}
