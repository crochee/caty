// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"cca/pkg/db"
	"net/http"

	"github.com/crochee/lib/e"
	"github.com/crochee/lib/id"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go"

	"cca/pkg/model"
	"cca/pkg/resp"
	"cca/pkg/service/business/tokenx"
	"cca/pkg/service/business/userx"
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
	user := &model.Account{
		AccountID:  idString,
		Account:    userRequest.Account,
		UserID:     idString,
		PassWord:   userRequest.PassWord,
		Email:      userRequest.Email,
		Permission: permission,
		Desc:       userRequest.Desc,
	}

	if err = db.With(ctx.Request.Context()).Model(user).Create(user).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.SuccessNotContent(ctx)
}

type ModifyUserRequest struct {
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

// Modify godoc
// @Summary Modify
// @Description user modify
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body ModifyInfo true "request's content"
// @Success 200
// @Success 304
// @Failure 400 {object} ex.Response
// @Failure 403 {object} ex.Response
// @Failure 500 {object} ex.Response
// @Router /v1/user/modify [put]
func Modify(ctx *gin.Context) {
	var modifyInfo ModifyInfo
	if err := ctx.ShouldBindBodyWith(&modifyInfo, binding.JSON); err != nil {
		logx.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		ex.ErrorWith(ctx, ex.ParsePayloadFailed, err.Error())
		return
	}
	// 检测邮箱的合法性
	if !internal.VerifyEmail(modifyInfo.Email) {
		ex.Error(ctx, ex.InvalidEmail)
		return
	}
	if modifyInfo.OldPassWord == "" {
		ctx.Status(http.StatusNotModified)
		return
	}

	if err := userx.ModifyUser(ctx.Request.Context(), modifyInfo.Email, modifyInfo.NewPassWord,
		modifyInfo.OldPassWord, modifyInfo.Nick); err != nil {
		ex.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

// Login godoc
// @Summary Login
// @Description user login
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body LoginInfo true "login request's content"
// @Success 200 {string} string
// @Failure 400 {object} ex.Response
// @Failure 500 {object} ex.Response
// @Router /v1/user/login [post]
func Login(ctx *gin.Context) {
	var loginInfo LoginInfo
	if err := ctx.ShouldBindBodyWith(&loginInfo, binding.JSON); err != nil {
		resp.ErrorParam(ctx, err)
		return
	}
	// 检测邮箱的合法性
	if !internal.VerifyEmail(loginInfo.Email) {
		resp.Error(ctx, e.InvalidEmail)
		return
	}
	token, err := userx.UserLogin(ctx.Request.Context(), loginInfo.Email, loginInfo.PassWord)
	if err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}
