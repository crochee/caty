// Date: 2021/10/8

// Package account
package account

import (
	"time"

	"github.com/crochee/lib/e"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"cca/pkg/db"
	"cca/pkg/model"
	"cca/pkg/resp"
	"cca/pkg/service/account"
	"cca/pkg/validator"
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
	var modifyRequest account.ModifyRequest
	if err := ctx.ShouldBindBodyWith(&modifyRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	if err := account.Modify(ctx.Request.Context(), &modifyRequest); err != nil {
		resp.Errors(ctx, err)
		return
	}
	resp.SuccessNotContent(ctx)
}

type QueryRequest struct {
	// 账户ID
	// Required: true
	AccountID string `json:"account_id" binding:"required"`
	// 用户ID
	// Required: true
	UserID string `json:"user_id" binding:"required"`
}

type QueryResponseResult struct {
	// 账户ID
	AccountID string `json:"account_id"`
	// 账户
	Account string `json:"account"`
	// 用户
	UserID string `json:"user_id"`
	// 邮箱
	Email string `json:"email"`
	// 权限
	Permission string `json:"permission"`
	// 是否认证
	Verify uint8 `json:"verify"`
	// 描述
	Desc string `json:"desc"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// Query godoc
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
func Query(ctx *gin.Context) {
	var queryRequest QueryRequest
	if err := ctx.ShouldBindBodyWith(&queryRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	accountModel := &model.Account{}
	if err := db.With(ctx.Request.Context()).Model(accountModel).Where("user_id =? AND account_id =?",
		queryRequest.UserID, queryRequest.AccountID).First(accountModel).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.Success(ctx, &QueryResponseResult{
		AccountID:  accountModel.AccountID,
		Account:    accountModel.Account,
		UserID:     accountModel.UserID,
		Email:      accountModel.Email,
		Permission: accountModel.Permission,
		Verify:     accountModel.Verify,
		Desc:       accountModel.Desc,
		CreatedAt:  accountModel.CreatedAt,
		UpdatedAt:  accountModel.UpdatedAt,
	})
}

type DeleteRequest struct {
	// 账户ID
	// Required: true
	AccountID string `json:"account_id" binding:"required"`
	// 用户ID
	// Required: true
	UserID string `json:"user_id" binding:"required"`
}

// Delete godoc
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
func Delete(ctx *gin.Context) {
	var deleteRequest DeleteRequest
	if err := ctx.ShouldBindBodyWith(&deleteRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	accountModel := &model.Account{}
	if err := db.With(ctx.Request.Context()).Model(accountModel).Where("user_id =? AND account_id =?",
		deleteRequest.UserID, deleteRequest.AccountID).Delete(accountModel).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.SuccessNotContent(ctx)
}
