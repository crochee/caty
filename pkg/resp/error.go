package resp

import (
	"errors"
	"github.com/gin-gonic/gin"

	"cca/pkg/ex"
	"cca/pkg/logx"
)

type ResponseCode struct {
	// 返回码
	// Required: true
	// Example: 20000000
	Code int `json:"code"`
	// 返回信息描述
	// Required: true
	// Example: success
	Msg string `json:"message"`
}

type ResponseError struct {
	ResponseCode
	// 具体描述信息
	// Required: true
	Result string `json:"result"`
}

// Error gin response with Code
func Error(ctx *gin.Context, code ex.Code) {
	ctx.JSON(code.StatusCode(), ResponseError{
		ResponseCode: ResponseCode{
			Code: code.Code(),
			Msg:  code.Error(),
		},
	})
}

// ErrorWith gin response with Code and message
func ErrorWith(ctx *gin.Context, code ex.Code, message string) {
	ctx.JSON(code.StatusCode(), ResponseError{
		ResponseCode: ResponseCode{
			Code: code.Code(),
			Msg:  code.Error(),
		},
		Result: message,
	})
}

// Errors gin Response with error
func Errors(ctx *gin.Context, err error) {
	logx.FromContext(ctx.Request.Context()).Errorf("has error %+v", err)
	var errResult error
	if errResult = ex.Unwrap(err); errResult == nil {
		ErrorWith(ctx, ex.ErrInternalServerError, err.Error())
		return
	}
	var errorCode ex.Code
	if errors.As(errResult, &errorCode) {
		Error(ctx, errorCode)
		return
	}
	ErrorWith(ctx, ex.ErrInternalServerError, err.Error())
}

func ErrorParam(ctx *gin.Context, err error) {
	logx.FromContext(ctx.Request.Context()).Errorf("parse param failed.Error:%+v", err)
	ctx.AbortWithStatusJSON(ex.ErrInvalidParam.StatusCode(), ResponseError{
		ResponseCode: ResponseCode{
			Code: ex.ErrInvalidParam.Code(),
			Msg:  ex.ErrInvalidParam.Error(),
		},
		Result: err.Error(),
	})
}
