// Date: 2021/10/24

// Package client
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crochee/lib/client"
	"github.com/json-iterator/go"

	"cca/pkg/resp"
	"cca/pkg/service/account"
)

type Account interface {
	Register(ctx context.Context, request *account.CreateRequest) (*account.CreateResponseResult, error)
	Retrieves(ctx context.Context, request *account.RetrievesRequest) (*account.RetrieveResponses, error)
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
	if req, err = client.NewRequest(ctx, http.MethodPost, a.Url(ctx, "/v1/account"),
		body, a.Header(ctx)); err != nil {
		return nil, err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var result resp.ResponseCode
		if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("http code %d,but not 200,%w", response.StatusCode, err)
		}
		return nil, &result
	}
	var result account.CreateResponseResult
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AccountClient) Retrieves(ctx context.Context,
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

	req, err := client.NewRequest(ctx, http.MethodGet, a.UrlWithQuery(ctx, "/v1/account", params),
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
		var result resp.ResponseCode
		if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("http code %d,but not 200,%w", response.StatusCode, err)
		}
		return nil, &result
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
	if req, err = client.NewRequest(ctx, http.MethodPatch, a.Url(ctx, "/v1/account"+user.ID),
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
	var result resp.ResponseCode
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	return &result
}

func (a *AccountClient) Retrieve(ctx context.Context, user *account.User) (*account.RetrieveResponse, error) {
	req, err := client.NewRequest(ctx, http.MethodGet, a.Url(ctx, "/v1/account"+user.ID), nil, a.Header(ctx))
	if err != nil {
		return nil, err
	}
	var response *http.Response
	if response, err = a.Do(req); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var result resp.ResponseCode
		if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("http code %d,but not 200,%w", response.StatusCode, err)
		}
		return nil, &result
	}
	var result account.RetrieveResponse
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AccountClient) Delete(ctx context.Context, user *account.User) error {
	req, err := client.NewRequest(ctx, http.MethodDelete, a.Url(ctx, "/v1/account"+user.ID), nil, a.Header(ctx))
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
	var result resp.ResponseCode
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	return &result
}
