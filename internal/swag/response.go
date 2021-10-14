// Date: 2021/9/19

// Package swag
package swag

import (
	"cca/pkg/resp"
	"cca/pkg/service/account"
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
