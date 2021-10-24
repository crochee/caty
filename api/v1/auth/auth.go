// Date: 2021/10/16

// Package auth
package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"cca/pkg/resp"
	"cca/pkg/service/auth"
)

// Sign godoc
// swagger:operation POST /v1/auth/sign 鉴权 SAuthSignRequest
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
		resp.ErrorParam(ctx, err)
		return
	}
	token, err := auth.Create(ctx.Request.Context(), &request)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

// Parse godoc
// swagger:operation POST /v1/auth/parse 鉴权 SAuthParseRequest
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
		resp.ErrorParam(ctx, err)
		return
	}
	token, err := auth.Parse(ctx.Request.Context(), &apiToken)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}
