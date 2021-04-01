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
	resp := &ErrorResponse{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	ctx.JSON(code.Status(), resp)
}

// Error gin response with Code
func Error(ctx *gin.Context, code Code) {
	resp := &ErrorResponse{
		Code:    code.String(),
		Message: code.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	ctx.JSON(code.Status(), resp)
}

// Errors gin response with error
func Errors(ctx *gin.Context, err error) {
	var errorCode *ErrorCode
	if errors.As(err, &errorCode) {
		if errorCode.message == "" {
			Error(ctx, errorCode.code)
			return
		}
		ErrorWith(ctx, errorCode.code, errorCode.message)
		return
	}
	resp := &ErrorResponse{
		Code:    Unknown.String(),
		Message: Unknown.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = Unknown.Chinese()
	}
	ctx.JSON(Unknown.Status(), resp)
}
