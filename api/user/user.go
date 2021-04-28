// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"net/http"

	"github.com/crochee/uid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go"

	"obs/e"
	"obs/internal"
	"obs/logger"
	"obs/model/db"
	"obs/service/business/tokenx"
	"obs/service/business/userx"
)

// Register godoc
// @Summary register
// @Description register user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body Domain true "register request's content"
// @Success 200
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/user/register [post]
func Register(ctx *gin.Context) {
	var domainInfo Domain
	if err := ctx.ShouldBindBodyWith(&domainInfo, binding.JSON); err != nil {
		logger.WithContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	// 检测邮箱的合法性
	if !internal.VerifyEmail(domainInfo.Email) {
		e.Error(ctx, e.InvalidEmail)
		return
	}
	permission, err := jsoniter.ConfigFastest.MarshalToString(map[string]tokenx.Action{
		tokenx.AllService: tokenx.Admin,
	})
	if err != nil {
		logger.WithContext(ctx.Request.Context()).Errorf("marshal permission failed.Error:%v", err)
		e.Error(ctx, e.MarshalFail)
		return
	}
	domain := &db.Domain{
		Domain:     uid.New().String(),
		Email:      domainInfo.Email,
		Nick:       domainInfo.Nick,
		PassWord:   domainInfo.PassWord,
		Permission: permission,
	}
	if err = db.NewDB().Create(domain).Error; err != nil {
		logger.WithContext(ctx.Request.Context()).Errorf("insert domain failed.Error:%v", err)
		e.ErrorWith(ctx, e.OperateDbFail, err.Error())
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
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/user/login [post]
func Login(ctx *gin.Context) {
	var loginInfo LoginInfo
	if err := ctx.ShouldBindBodyWith(&loginInfo, binding.JSON); err != nil {
		logger.WithContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	// 检测邮箱的合法性
	if !internal.VerifyEmail(loginInfo.Email) {
		e.Error(ctx, e.InvalidEmail)
		return
	}
	token, err := userx.UserLogin(ctx.Request.Context(), loginInfo.Email, loginInfo.PassWord)
	if err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

// Modify godoc
// @Summary Modify
// @Description user modify
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body ModifyInfo true "request's content"
// @Success 200
// @Success 304
// @Failure 400 {object} e.Response
// @Failure 403 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/user/modify [put]
func Modify(ctx *gin.Context) {
	var modifyInfo ModifyInfo
	if err := ctx.ShouldBindBodyWith(&modifyInfo, binding.JSON); err != nil {
		logger.WithContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	// 检测邮箱的合法性
	if !internal.VerifyEmail(modifyInfo.Email) {
		e.Error(ctx, e.InvalidEmail)
		return
	}
	if modifyInfo.OldPassWord == "" {
		ctx.Status(http.StatusNotModified)
		return
	}

	if err := userx.ModifyUser(ctx.Request.Context(), modifyInfo.Email, modifyInfo.NewPassWord,
		modifyInfo.OldPassWord, modifyInfo.Nick); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}
