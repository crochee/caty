// Author: crochee
// Date: 2021/9/19

// Package resp
package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cca/pkg/ex"
)

type Response struct {
	ResponseCode
	Result interface{} `json:"result"`
}

// Success response data
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(ex.ErrSuccess.StatusCode(), Response{
		ResponseCode: ResponseCode{
			Code: ex.ErrSuccess.Code(),
			Msg:  ex.ErrSuccess.Error(),
		},
		Result: data,
	})
}

// SuccessNone response none
func SuccessNone(ctx *gin.Context) {
	ctx.JSON(ex.ErrSuccess.StatusCode(), Response{
		ResponseCode: ResponseCode{
			Code: ex.ErrSuccess.Code(),
			Msg:  ex.ErrSuccess.Error(),
		},
	})
}

// SuccessNotContent response 204 nothing
func SuccessNotContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
