// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"net/http"
	"obs/e"

	"github.com/crochee/uid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go"

	"obs/logger"
	"obs/model/db"
	"obs/response"
	"obs/service/tokenx"
	"obs/service/userx"
	"obs/util"
)

// Register godoc
// @Summary register
// @Description register user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body Domain true "register request's content"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/user/register [post]
func Register(ctx *gin.Context) {
	var domainInfo Domain
	if err := ctx.ShouldBindBodyWith(&domainInfo, binding.JSON); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Check your payload"))
		return
	}
	// 检测邮箱的合法性
	if !util.VerifyEmail(domainInfo.Email) {
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Invalid email"))
		return
	}
	permission, err := jsoniter.ConfigFastest.MarshalToString(map[string]tokenx.Action{
		tokenx.AllService: tokenx.Admin,
	})
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("marshal permission failed.Error:%v", err)
		response.ErrorWithCode(ctx, http.StatusInternalServerError)
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
		logger.FromContext(ctx.Request.Context()).Errorf("insert domain failed.Error:%v", err)
		e.ErrorWith(ctx, response.Errors(http.StatusInternalServerError, err))
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/user/login [post]
func Login(ctx *gin.Context) {
	var loginInfo LoginInfo
	if err := ctx.ShouldBindBodyWith(&loginInfo, binding.JSON); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Check your payload"))
		return
	}
	// 检测邮箱的合法性
	if !util.VerifyEmail(loginInfo.Email) {
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Invalid email"))
		return
	}
	token, err := userx.UserLogin(ctx.Request.Context(), loginInfo.Email, loginInfo.PassWord)
	if err != nil {
		e.ErrorWith(ctx, err)
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/user/modify [put]
func Modify(ctx *gin.Context) {
	var modifyInfo ModifyInfo
	if err := ctx.ShouldBindBodyWith(&modifyInfo, binding.JSON); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Check your payload"))
		return
	}
	// 检测邮箱的合法性
	if !util.VerifyEmail(modifyInfo.Email) {
		e.ErrorWith(ctx, e.Error(http.StatusBadRequest, "Invalid email"))
		return
	}
	if modifyInfo.OldPassWord == "" {
		ctx.Status(http.StatusNotModified)
		return
	}

	if err := userx.ModifyUser(ctx.Request.Context(), modifyInfo.Email, modifyInfo.NewPassWord,
		modifyInfo.OldPassWord, modifyInfo.Nick); err != nil {
		e.ErrorWith(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}
