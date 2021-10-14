// Date: 2021/9/19

// Package swag
package swag

import (
	"cca/pkg/service/account"
)

// swagger:parameters SNoneRequest
type SNoneRequest struct {
}

// swagger:parameters SAccountRegisterRequest
type SAccountRegisterRequest struct {
	// in: body
	Body struct {
		account.CreateRequest
	}
}
