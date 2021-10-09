// Date: 2021/10/8

// Package account
package account

import (
	"time"

	"github.com/crochee/lib/e"
	"github.com/crochee/lib/id"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go"

	"cca/pkg/db"
	"cca/pkg/model"
	"cca/pkg/resp"
	"cca/pkg/service/business/tokenx"
	"cca/pkg/validator"
)

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

type RegisterRequest struct {
	// 账户
	// Required: true
	Account string `json:"account" binding:"required"`
	// 邮箱
	Email string `json:"email"`
	// 密码
	// Required: true
	Password string `json:"password" binding:"required"`
	// 描述信息
	Desc string `json:"desc" binding:"json"`
}

type RegisterResponseResult struct {
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

// Register godoc
// swagger:route  POST /v1/account 账户 SwaggerRegisterUserRequest
// 注册账户
//
// register account
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//     Responses:
//		 204: SwaggerRegisterUserResponse
//       default: SwaggerResponseError
func Register(ctx *gin.Context) {
	var userRequest RegisterRequest
	if err := ctx.ShouldBindBodyWith(&userRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	if userRequest.Email != "" {
		// 检测邮箱的合法性
		if err := validator.Var(userRequest.Email, "email"); err != nil {
			resp.ErrorParam(ctx, err)
			return
		}
	}

	permission, err := jsoniter.ConfigFastest.MarshalToString(map[string]tokenx.Action{
		tokenx.AllService: tokenx.Admin,
	})
	if err != nil {
		resp.ErrorWith(ctx, e.ErrInternalServerError, err.Error())
		return
	}
	var idString string
	if idString, err = id.NextIDString(); err != nil {
		resp.ErrorWith(ctx, e.ErrInternalServerError, err.Error())
		return
	}
	accountModel := &model.Account{
		AccountID:  idString,
		Account:    userRequest.Account,
		UserID:     idString,
		Password:   userRequest.Password,
		Email:      userRequest.Email,
		Permission: permission,
		Desc:       userRequest.Desc,
	}

	if err = db.With(ctx.Request.Context()).Model(accountModel).Create(accountModel).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.Success(ctx, &RegisterResponseResult{
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

type ModifyRequest struct {
	// 账户ID
	// Required: true
	AccountID string `json:"account_id" binding:"required"`
	// 用户
	// Required: true
	UserID string `json:"user_id" binding:"required"`
	// 旧密码
	// Required: true
	OldPassword string `json:"old_password" binding:"required"`
	// 账户
	Account string `json:"account"`
	// 邮箱
	Email string `json:"email"`
	// 新密码
	Password string `json:"password"`
	// 权限
	Permission string `json:"permission"`
	// 描述信息
	Desc string `json:"desc"`
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
	var modifyRequest ModifyRequest
	if err := ctx.ShouldBindBodyWith(&modifyRequest, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	updates := make(map[string]interface{})
	if modifyRequest.Email != "" {
		// 检测邮箱的合法性
		if err := validator.Var(modifyRequest.Email, "email"); err != nil {
			resp.ErrorParam(ctx, err)
			return
		}
	}
	if modifyRequest.Desc != "" {
		// 检测邮箱的合法性
		if err := validator.Var(modifyRequest.Desc, "json"); err != nil {
			resp.ErrorParam(ctx, err)
			return
		}
	}

	query := db.With(ctx.Request.Context()).Model(&model.Account{}).Updates(updates)
	if err := query.Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	if query.RowsAffected == 0 {
		resp.Error(ctx, e.ErrOperateDB)
		return
	}
	resp.SuccessNotContent(ctx)
}
