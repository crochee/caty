// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"obs/config"
	"obs/logger"
	"os"
)

// Register godoc
// @Summary register
// @Description register user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param request body RegisterRequest true "register request's content"
// @Success 200
// @Failure 400
// @Failure 500 {object} util.HttpErr
// @Router /user/register [post]
func Register(ctx *gin.Context) {
	var registerRequest RegisterRequest
	err := ctx.ShouldBindBodyWith(&registerRequest, binding.JSON)
	if err != nil {
		logger.Errorf("register bind body fail,%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	// 检测邮箱的合法性
	if !util.VerifyEmail(registerRequest.Email) {
		util.Error(ctx, util.RegisterErr)
		return
	}
	// 邮件确认
	//if err = email.SendEmail(email.FromQQEmail{}, &email.NeedParameterfForEmail{
	//	From:        "console",
	//	To:          []string{registerRequest.Email},
	//	Subject:     "[Console] Remind!",
	//	Files:       nil,
	//	ContentType: "",
	//	Content:     Remind,
	//}); err != nil {
	//	logger.Errorf("send email failed.Error:%v", err)
	//	util.Error(ctx, util.RegisterErr)
	//	return
	//}
	user := &cmysql.UserInfo{
		Email:    registerRequest.Email,
		Nick:     registerRequest.Nick,
		PassWord: registerRequest.PassWord,
	}
	// 开启事务
	dbTx := user.Conn().Begin()
	if err = dbTx.Table(user.TableName()).Create(user).Error; err != nil {
		dbTx.Rollback() // 事务回滚
		logger.Errorf("register user fail,%v", err)
		util.Error(ctx, util.RegisterErr)
		return
	}
	// 创建用户空间
	if err = os.MkdirAll(config.Cfg.SavePath+user.Email, os.ModePerm); err != nil {
		dbTx.Rollback() // 事务回滚
		logger.Errorf("register mkdir fail,%v", err)
		util.Error(ctx, util.RegisterErr)
		return
	}
	dbTx.Commit() // 事务提交
	ctx.Status(http.StatusOK)
}
