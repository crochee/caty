// Date: 2021/9/19

// Package swag
package swag

import (
	"cca/pkg/resp"
	"cca/pkg/service/account"
)

// swagger:response SNoneResponse
type SNoneResponse struct {
}

// swagger:response SResponseError
type SResponseError struct {
	// in: body
	Body struct {
		resp.ResponseError
	}
}

// swagger:response SAccountRegisterResponseResult
type SAccountRegisterResponseResult struct {
	// in: body
	Body struct {
		resp.ResponseCode
		Result *account.CreateResponseResult
	}
}
