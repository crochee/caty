// Date: 2021/9/19

// Package swag
package swag

import "obs/pkg/resp"

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
