// Date: 2021/10/16

// Package account
package account

import (
	"net/http"

	"github.com/crochee/lirity/e"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"caty/pkg/service/account"
)

// Login godoc
// swagger:operation POST /v1/account/login 账户 SAccountLoginRequest
// ---
// summary: 用户登录
// description: 用户登录获取token信息
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
func Login(ctx *gin.Context) {
	var request account.LoginRequest
	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		e.GinErrorCode(ctx, e.ErrInvalidParam.WithResult(err))
		return
	}
	if err := account.ValidPassword(request.Password); err != nil {
		e.GinErrorCode(ctx, e.ErrInvalidParam.WithResult(err))
		return
	}
	response, err := account.Login(ctx.Request.Context(), &request)
	if err != nil {
		e.GinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
