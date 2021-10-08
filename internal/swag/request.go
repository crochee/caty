// Date: 2021/9/19

// Package swag
package swag

import "cca/api/user"

// swagger:parameters SwaggerNoneRequest
type SwaggerNoneRequest struct {
}

// swagger:parameters SwaggerRegisterUserRequest
type SwaggerRegisterUserRequest struct {
	// in: body
	Body struct {
		user.RegisterUserRequest
	}
}
