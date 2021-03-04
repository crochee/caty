// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorWith gin response with format err
func ErrorWith(ctx *gin.Context, err error) {
	switch value := err.(type) {
	case *ErrorResponse:
		ctx.JSON(int(value.Code), value)
	default:
		ctx.JSON(http.StatusInternalServerError,
			Error(http.StatusInternalServerError, value.Error()))
	}
}

// ErrorWithMessage gin response with message
func ErrorWithMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Error(http.StatusInternalServerError, message))
}
