// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"obs/pkg/db"
	"obs/pkg/e"
	"obs/pkg/id"
	"obs/pkg/model"
	"obs/pkg/resp"
	"obs/pkg/validator"
)

//// 账户ID
//// Required: true
//AccountID string `json:"account_id" binding:"required,numeric"`
//// 用户ID
//// Required: true
//UserID string `json:"user_id" binding:"required,numeric"`

type RegisterUserRequest struct {
	// 账户
	// Required: true
	Account string `json:"account" binding:"required"`
	// 昵称
	Nick string `json:"nick"`
	// 邮箱
	Email string `json:"email"`
	// 密码
	PassWord string `json:"pass_word" binding:"required"`
	// 权限信息
	Permission *json.RawMessage `json:"permission" binding:"json"`
	// 描述信息
	Desc *json.RawMessage `json:"desc" binding:"json"`
}

// Register godoc
// @Summary register
// @Description register user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body User true "register request's content"
// @Success 200
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/user/register [post]

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
		if err := validator.Default.Validate.Var(userRequest.Email, "email"); err != nil {
			resp.ErrorParam(ctx, validator.Default.Translate(err))
			return
		}
	}

	//permission, err := jsoniter.ConfigFastest.MarshalToString(map[string]tokenx.Action{
	//	tokenx.AllService: tokenx.Admin,
	//})
	//if err != nil {
	//	log.FromContext(ctx.Request.Context()).Errorf("marshal permission failed.Error:%v", err)
	//	e.Error(ctx, e.MarshalFail)
	//	return
	//}
	//domain := &model.Domain{
	//	Domain:     uid.New().String(),
	//	Email:      userRequest.Email,
	//	Nick:       userRequest.Nick,
	//	PassWord:   userRequest.PassWord,
	//	Permission: permission,
	//}
	idString, err := id.NextIDString()
	if err != nil {
		resp.ErrorWith(ctx, e.ErrInternalServerError, err.Error())
		return
	}
	user := &model.User{
		AccountID:  idString,
		UserID:     idString,
		Nick:       "",
		PassWord:   "",
		Email:      "",
		Permission: nil,
		Verify:     0,
		Desc:       "",
	}

	if err = db.New(ctx.Request.Context()).Model(user).Create(user).Error; err != nil {
		resp.ErrorWith(ctx, e.ErrOperateDB, err.Error())
		return
	}
	resp.SuccessNotContent(ctx)
}

//// Login godoc
//// @Summary Login
//// @Description user login
//// @Tags user
//// @Accept application/json
//// @Produce application/json
//// @Param request body LoginInfo true "login request's content"
//// @Success 200 {string} string
//// @Failure 400 {object} e.Response
//// @Failure 500 {object} e.Response
//// @Router /v1/user/login [post]
//func Login(ctx *gin.Context) {
//	var loginInfo LoginInfo
//	if err := ctx.ShouldBindBodyWith(&loginInfo, binding.JSON); err != nil {
//		log.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
//		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
//		return
//	}
//	// 检测邮箱的合法性
//	if !internal.VerifyEmail(loginInfo.Email) {
//		e.Error(ctx, e.InvalidEmail)
//		return
//	}
//	token, err := userx.UserLogin(ctx.Request.Context(), loginInfo.Email, loginInfo.PassWord)
//	if err != nil {
//		e.Errors(ctx, err)
//		return
//	}
//	ctx.JSON(http.StatusOK, token)
//}
//
//// Modify godoc
//// @Summary Modify
//// @Description user modify
//// @Tags user
//// @Accept application/json
//// @Produce application/json
//// @Param request body ModifyInfo true "request's content"
//// @Success 200
//// @Success 304
//// @Failure 400 {object} e.Response
//// @Failure 403 {object} e.Response
//// @Failure 500 {object} e.Response
//// @Router /v1/user/modify [put]
//func Modify(ctx *gin.Context) {
//	var modifyInfo ModifyInfo
//	if err := ctx.ShouldBindBodyWith(&modifyInfo, binding.JSON); err != nil {
//		log.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
//		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
//		return
//	}
//	// 检测邮箱的合法性
//	if !internal.VerifyEmail(modifyInfo.Email) {
//		e.Error(ctx, e.InvalidEmail)
//		return
//	}
//	if modifyInfo.OldPassWord == "" {
//		ctx.Status(http.StatusNotModified)
//		return
//	}
//
//	if err := userx.ModifyUser(ctx.Request.Context(), modifyInfo.Email, modifyInfo.NewPassWord,
//		modifyInfo.OldPassWord, modifyInfo.Nick); err != nil {
//		e.Errors(ctx, err)
//		return
//	}
//	ctx.Status(http.StatusOK)
//}
