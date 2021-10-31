// Date: 2021/9/19

// Package swag
package swag

import (
	"caty/api"
	"caty/pkg/resp"
	"caty/pkg/service/account"
	"caty/pkg/service/auth"
)

// swagger:response SNullResponse
type SNullResponse struct {
}

// swagger:response SResponseCode
type SResponseCode struct {
	// in: body
	Body struct {
		resp.ResponseCode
	}
}

// swagger:response SAPIVersionResponse
type SAPIVersionResponse struct {
	// in: body
	Body struct {
		api.VersionResponse
	}
}

// swagger:response SAccountRegisterResponseResult
type SAccountRegisterResponseResult struct {
	// in: body
	Body struct {
		account.CreateResponseResult
	}
}

// swagger:response SAccountRetrieveResponses
type SAccountRetrieveResponses struct {
	// in: body
	Body struct {
		account.RetrieveResponses
	}
}

// swagger:response SAccountRetrieveResponse
type SAccountRetrieveResponse struct {
	// in: body
	Body struct {
		account.RetrieveResponse
	}
}

// swagger:response SAuthSignResponse
type SAuthSignResponse struct {
	// in: body
	Body struct {
		auth.APIToken
	}
}

// swagger:response SAuthParseResponse
type SAuthParseResponse struct {
	// in: body
	Body struct {
		auth.TokenClaims
	}
}
