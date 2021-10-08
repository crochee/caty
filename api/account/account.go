// Date: 2021/10/8

// Package account
package account

import (
	"github.com/crochee/lib/e"
	"github.com/crochee/lib/id"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go"
	"time"

	"cca/pkg/model"
	"cca/pkg/resp"
	"cca/pkg/service/business/tokenx"
	"cca/pkg/validator"
)

type RegisterUserRequest struct {
	// 账户
	// Required: true
	Account string `json:"account" binding:"required"`
	// 邮箱
	Email string `json:"email"`
	// 密码
	// Required: true
	PassWord string `json:"pass_word" binding:"required"`
	// 描述信息
	Desc string `json:"desc" binding:"json"`
}

type RegisterUserResponseResult struct {
	AccountID  string
	Account    string
	UserID     string
	Email      string
	Permission string
	Verify     uint8
	Desc       string
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
//		 204: SwaggerNoneResponse
//       default: SwaggerResponseError
func Register(ctx *gin.Context) {
	var userRequest RegisterUserRequest
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
	user := &model.User{
		AccountID:  idString,
		Account:    userRequest.Account,
		UserID:     idString,
		PassWord:   userRequest.PassWord,
		Email:      userRequest.Email,
		Permission: permission,
		Desc:       userRequest.Desc,
	}

	if err = model.With(ctx.Request.Context()).Model(user).Create(user).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.Success(ctx, &RegisterUserResponseResult{
		AccountID:  user.AccountID,
		Account:    user.Account,
		UserID:     user.UserID,
		Email:      user.Email,
		Permission: user.Permission,
		Verify:     user.Verify,
		Desc:       user.Desc,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	})
}
