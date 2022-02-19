// Date: 2021/10/24

// Package client
package client

import (
	"context"
	"net/http"

	"github.com/crochee/lirity/client"
	"github.com/crochee/lirity/e"
	"github.com/json-iterator/go"

	"caty/pkg/service/auth"
)

type Auth interface {
	Sign(ctx context.Context, request *auth.TokenClaims) (*auth.APIToken, error)
	Parse(ctx context.Context, request *auth.APIToken) (*auth.TokenClaims, error)
}

func NewAuth() Auth {
	return &AuthClient{
		Client:     client.NewStandardClient(),
		API:        jsoniter.ConfigCompatibleWithStandardLibrary,
		URLHandler: NewURLHandler(),
	}
}

type AuthClient struct {
	client.Client
	jsoniter.API
	URLHandler
}

func (a *AuthClient) Sign(ctx context.Context, request *auth.TokenClaims) (*auth.APIToken, error) {
	body, err := a.Marshal(request)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	if req, err = client.NewRequest(ctx, http.MethodPost, a.URL(ctx, "/v1/auth/sign"),
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
	var result auth.APIToken
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AuthClient) Parse(ctx context.Context, request *auth.APIToken) (*auth.TokenClaims, error) {
	body, err := a.Marshal(request)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	if req, err = client.NewRequest(ctx, http.MethodPost, a.URL(ctx, "/v1/auth/parse"),
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
	var result auth.TokenClaims
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
