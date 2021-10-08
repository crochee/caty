package resp

import (
	"net/http"

	"github.com/crochee/lib/e"
	"github.com/gin-gonic/gin"
)

type Response struct {
	ResponseCode
	Result interface{} `json:"result"`
}

// Success response data
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(e.ErrSuccess.StatusCode(), Response{
		ResponseCode: ResponseCode{
			Code: e.ErrSuccess.Code(),
			Msg:  e.ErrSuccess.Error(),
		},
		Result: data,
	})
}

// SuccessNone response none
func SuccessNone(ctx *gin.Context) {
	ctx.JSON(e.ErrSuccess.StatusCode(), Response{
		ResponseCode: ResponseCode{
			Code: e.ErrSuccess.Code(),
			Msg:  e.ErrSuccess.Error(),
		},
	})
}

// SuccessNotContent response 204 nothing
func SuccessNotContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
