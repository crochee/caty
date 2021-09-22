package resp

import (
	"errors"

	"github.com/gin-gonic/gin"

	"obs/pkg/e"
	"obs/pkg/log"
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
func Error(ctx *gin.Context, code e.Code) {
	ctx.JSON(code.StatusCode(), ResponseError{
		ResponseCode: ResponseCode{
			Code: code.Code(),
			Msg:  code.Error(),
		},
	})
}

// ErrorWith gin response with Code and message
func ErrorWith(ctx *gin.Context, code e.Code, message string) {
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
	log.FromContext(ctx.Request.Context()).Errorf("has error %+v", err)
	var errResult error
	if errResult = e.Unwrap(err); errResult == nil {
		ErrorWith(ctx, e.ErrInternalServerError, err.Error())
		return
	}
	var errorCode e.Code
	if errors.As(errResult, &errorCode) {
		Error(ctx, errorCode)
		return
	}
	ErrorWith(ctx, e.ErrInternalServerError, err.Error())
}

func ErrorParam(ctx *gin.Context, err error) {
	ErrorWith(ctx, e.ErrInvalidParam, err.Error())
}