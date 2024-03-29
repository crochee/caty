// Date: 2021/9/19

// Package swag
package swag

import (
	"caty/pkg/service/account"
	"caty/pkg/service/auth"
)

// swagger:parameters SNullRequest
type SNullRequest struct {
}

// swagger:parameters SAccountRegisterRequest
type SAccountRegisterRequest struct {
	// in: body
	Body struct {
		account.CreateRequest
	}
}

// swagger:parameters SAccountRetrievesRequest
type SAccountRetrievesRequest struct {
	account.RetrievesRequest
}

// swagger:parameters SAccountUpdateRequest
type SAccountUpdateRequest struct {
	// in: body
	Body struct {
		account.UpdateRequest
	}
	account.User
}

// swagger:parameters SAccountRetrieveRequest
type SAccountRetrieveRequest struct {
	account.User
}

// swagger:parameters SAccountDeleteRequest
type SAccountDeleteRequest struct {
	account.User
}

// swagger:parameters SAuthSignRequest
type SAuthSignRequest struct {
	// in: body
	Body struct {
		auth.TokenClaims
	}
}

// swagger:parameters SAuthParseRequest
type SAuthParseRequest struct {
	// in: body
	Body struct {
		auth.APIToken
	}
}

// swagger:parameters SAccountLoginRequest
type SAccountLoginRequest struct {
	// in: body
	Body struct {
		account.LoginRequest
	}
}
