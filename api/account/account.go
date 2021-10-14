// Date: 2021/10/8

// Package account
package account

import (
	"net/http"

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
//       default: SResponseCode
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
	ctx.JSON(http.StatusOK, response)
}

// Update godoc
// swagger:route  PATCH /v1/account 账户 SAccountUpdateRequest
// 编辑账户
//
// Update account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 204: SNullResponse
//       default: SResponseCode
func Update(ctx *gin.Context) {
	var modifyRequest account.UpdateRequest
	if err := ctx.ShouldBindBodyWith(&modifyRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	if err := account.Update(ctx.Request.Context(), &modifyRequest); err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Retrieve godoc
// swagger:route  GET /v1/account 账户 SAccountRetrieveRequest
// 查询账户
//
// retrieve account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 200: SAccountRetrieveResponses
//       default: SResponseCode
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
	ctx.JSON(http.StatusOK, response)
}

// RetrieveSingle godoc
// swagger:route  GET /v1/account/{id} 账户 SAccountRetrieveSingleRequest
// 查询指定账户
//
// retrieve specified account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 200: SAccountRetrieveResponse
//       default: SResponseCode
func RetrieveSingle(ctx *gin.Context) {
	var retrieveRequest account.RetrieveSingleRequest
	if err := ctx.ShouldBindBodyWith(&retrieveRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	response, err := account.RetrieveSingle(ctx.Request.Context(), &retrieveRequest)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Delete godoc
// swagger:route  DELETE /v1/account/{id} 账户 SAccountDeleteRequest
// 删除账户
//
// delete account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 204: SNullResponse
//       default: SResponseCode
func Delete(ctx *gin.Context) {
	var retrieveRequest account.DeleteRequest
	if err := ctx.ShouldBindBodyWith(&retrieveRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	err := account.Delete(ctx.Request.Context(), &retrieveRequest)
	if err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
