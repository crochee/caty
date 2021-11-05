// Date: 2021/10/24

// Package client
package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crochee/lirity/client"
	"github.com/json-iterator/go"

	"caty/pkg/resp"
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
		var result resp.ResponseCode
		if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("http code %d,but not 200,%w", response.StatusCode, err)
		}
		return nil, &result
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
		var result resp.ResponseCode
		if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("http code %d,but not 200,%w", response.StatusCode, err)
		}
		return nil, &result
	}
	var result auth.TokenClaims
	if err = a.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
