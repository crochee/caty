// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorWith gin response with format err
func ErrorWith(ctx *gin.Context, err error) {
	var errorResponse *ErrorResponse
	if errors.As(err, &errorResponse) {
		ctx.JSON(errorResponse.Code, errorResponse)
		return
	}
	ctx.JSON(http.StatusInternalServerError,
		Error(http.StatusInternalServerError, err.Error()))
}

// ErrorWithMessage gin response with message
func ErrorWithMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Error(http.StatusInternalServerError, message))
}

// ErrorWithCode gin response with http code
func ErrorWithCode(ctx *gin.Context, code int) {
	ctx.JSON(http.StatusInternalServerError, Error(code, http.StatusText(code)))
}
