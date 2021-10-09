// Date: 2021/9/19

// Package swag
package swag

import (
	"cca/api/account"
	"cca/pkg/resp"
)

// swagger:response SwaggerNoneResponse
type SwaggerNoneResponse struct {
}

// swagger:response SwaggerResponseError
type SwaggerResponseError struct {
	// in: body
	Body struct {
		resp.ResponseError
	}
}

// swagger:response SwaggerRegisterUserResponse
type SwaggerRegisterUserResponse struct {
	// in: body
	Body struct {
		resp.ResponseCode
		Result *account.RegisterResponseResult
	}
}
