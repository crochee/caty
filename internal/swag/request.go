// Date: 2021/9/19

// Package swag
package swag

import "obs/api/user"

// swagger:parameters SwaggerNoneRequest
type SwaggerNoneRequest struct {
}

// swagger:parameters SwaggerDomain
type SwaggerDomain struct {
	// in: body
	Body struct {
		user.User
	}
}
