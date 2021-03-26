// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package e

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Extra   string `json:"extra"`
}

// ErrorWith gin response with with Code and message
func ErrorWith(ctx *gin.Context, code Code, message string) {
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		ctx.JSON(code.Status(), &ErrorResponse{
			Code:    code.String(),
			Message: code.Chinese(),
			Extra:   message,
		})
		return
	}
	ctx.JSON(code.Status(), &ErrorResponse{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	})
}

// Error gin response with Code
func Error(ctx *gin.Context, code Code) {
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		ctx.JSON(code.Status(), &ErrorResponse{
			Code:    code.String(),
			Message: code.Chinese(),
		})
		return
	}
	ctx.JSON(code.Status(), &ErrorResponse{
		Code:    code.String(),
		Message: code.English(),
	})
}

// Errors gin response with error
func Errors(ctx *gin.Context, err error) {
	var errorCode *ErrorCode
	if errors.As(err, &errorCode) {
		ErrorWith(ctx, errorCode.code, errorCode.message)
		return
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		ctx.JSON(Unknown.Status(), &ErrorResponse{
			Code:    Unknown.String(),
			Message: Unknown.Chinese(),
		})
		return
	}
	ctx.JSON(Unknown.Status(), &ErrorResponse{
		Code:    Unknown.String(),
		Message: Unknown.English(),
	})
}
