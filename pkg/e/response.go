// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package e

import (
	"obs/pkg/logx"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Extra   string `json:"extra"`
}

// ErrorWith gin response with with Code and message
func ErrorWith(ctx *gin.Context, code Code, message string) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logx.FromContext(ctx.Request.Context()).Errorf("[ERROR] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.JSON(code.Status(), resp)
}

// Error gin response with Code
func Error(ctx *gin.Context, code Code) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logx.FromContext(ctx.Request.Context()).Errorf("[ERROR] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.JSON(code.Status(), resp)
}

// Errors gin response with error
func Errors(ctx *gin.Context, err error) {
	var errResult error
	if errResult = errors.Cause(err); errResult == nil {
		ErrorWith(ctx, Success, err.Error())
		return
	}
	var errorCode *ErrorCode
	if errors.As(errResult, &errorCode) {
		ErrorWith(ctx, errorCode.code, errorCode.message)
		return
	}
	ErrorWith(ctx, Unknown, err.Error())
}

// AbortWith gin response with with Code and message
func AbortWith(ctx *gin.Context, code Code, message string) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logx.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}

// Abort gin response with with Code and message
func Abort(ctx *gin.Context, code Code) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logx.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}
