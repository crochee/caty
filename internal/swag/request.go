// Date: 2021/9/19

// Package swag
package swag

import (
	"cca/pkg/service/account"
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

// swagger:parameters SAccountUpdateRequest
type SAccountUpdateRequest struct {
	// in: body
	Body struct {
		account.UpdateRequest
	}
}

// swagger:parameters SAccountRetrieveRequest
type SAccountRetrieveRequest struct {
	// in: body
	Body struct {
		account.RetrieveRequest
	}
}

// swagger:parameters SAccountRetrieveSingleRequest
type SAccountRetrieveSingleRequest struct {
	// in: body
	Body struct {
		account.RetrieveSingleRequest
	}
}

// swagger:parameters SAccountDeleteRequest
type SAccountDeleteRequest struct {
	// in: body
	Body struct {
		account.DeleteRequest
	}
}
