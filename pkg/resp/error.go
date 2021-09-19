package resp

import (
	"errors"

	"github.com/gin-gonic/gin"

	"obs/pkg/e"
	"obs/pkg/log"
)

type ResponseCode struct {
	Code   int    `json:"code"`
	Msg    string `json:"message"`
	Detail string `json:"detail,omitempty"`
}

// Error gin response with Code
func Error(ctx *gin.Context, code e.Code) {
	ctx.JSON(code.StatusCode(), ResponseCode{
		Code: code.Code(),
		Msg:  code.Error(),
	})
}

// ErrorWith gin response with Code and message
func ErrorWith(ctx *gin.Context, code e.Code, message string) {
	ctx.JSON(code.StatusCode(), &ResponseCode{
		Code:   code.Code(),
		Msg:    code.Error(),
		Detail: message,
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
