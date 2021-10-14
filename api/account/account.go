// Date: 2021/10/8

// Package account
package account

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"cca/pkg/resp"
	"cca/pkg/service/account"
)

// Register godoc
// swagger:route  POST /v1/account 账户 SAccountRegisterRequest
// 注册账户
//
// register account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 200: SAccountRegisterResponseResult
//       default: SResponseError
func Register(ctx *gin.Context) {
	var registerRequest account.CreateRequest
	if err := ctx.ShouldBindBodyWith(&registerRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	response, err := account.Create(ctx.Request.Context(), &registerRequest)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	resp.Success(ctx, response)
}

// Modify godoc
// swagger:route  PATCH /v1/account 账户 SwaggerRegisterUserRequest
// 注册账户
//
// register account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 204: SwaggerNoneResponse
//       default: SwaggerResponseError
func Modify(ctx *gin.Context) {
	var modifyRequest account.UpdateRequest
	if err := ctx.ShouldBindBodyWith(&modifyRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	if err := account.Update(ctx.Request.Context(), &modifyRequest); err != nil {
		resp.Errors(ctx, err)
		return
	}
	resp.SuccessNotContent(ctx)
}

// Retrieve godoc
// swagger:route  GET /v1/account 账户 SwaggerRegisterUserRequest
// 查询账户
//
// register account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 200: SwaggerRegisterUserResponse
//       default: SwaggerResponseError
func Retrieve(ctx *gin.Context) {
	var retrieveRequest account.RetrieveRequest
	if err := ctx.ShouldBindBodyWith(&retrieveRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	response, err := account.Retrieve(ctx.Request.Context(), &retrieveRequest)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	resp.Success(ctx, response)
}
