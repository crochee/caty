// Date: 2021/10/16

// Package auth
package auth

import (
	"net/http"

	"github.com/crochee/lirity/e"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"caty/pkg/service/auth"
)

// Sign godoc
// swagger:operation POST /v1/auths/sign 鉴权 SAuthSignRequest
// ---
// summary: 生成token
// description: 生成token信息
// Consumes:
// - application/json
// produces:
// - application/json
// responses:
//   '200':
//     type: object
//     "$ref": "#/responses/SAuthSignResponse"
//   default:
//     type: object
//     "$ref": "#/responses/SResponseCode"
func Sign(ctx *gin.Context) {
	var request auth.TokenClaims
	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		e.Code(ctx, e.ErrInvalidParam.WithResult(err))
		return
	}
	token, err := auth.Create(ctx.Request.Context(), &request)
	if err != nil {
		e.Error(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

// Parse godoc
// swagger:operation POST /v1/auths/parse 鉴权 SAuthParseRequest
// ---
// summary: 解析token
// description: 解析token信息
// Consumes:
// - application/json
// produces:
// - application/json
// responses:
//   '200':
//     type: object
//     "$ref": "#/responses/SAuthParseResponse"
//   default:
//     type: object
//     "$ref": "#/responses/SResponseCode"
func Parse(ctx *gin.Context) {
	var apiToken auth.APIToken
	if err := ctx.ShouldBindBodyWith(&apiToken, binding.JSON); err != nil {
		e.Code(ctx, e.ErrInvalidParam.WithResult(err))
		return
	}
	token, err := auth.Parse(ctx.Request.Context(), &apiToken)
	if err != nil {
		e.Error(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}
