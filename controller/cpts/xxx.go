// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package cpts
package cpts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/logger"
)

type Person struct {
	Name  string  `json:"name" form:"name"`
	Age   int     `json:"age" form:"age"`
	Score float64 `json:"score" form:"score"`
}

// FindPerson godoc
// @Summary 上传人物信息
// @Description 上传人物信息
// @Tags CPTS
// @Accept application/x-www-form-urlencoded
// @Produce  application/json
// @Param request body Person true "任务信息"
// @Success 200 {object} Person
// @Failure 400
// @Failure 500
// @Router /test/person [post]
func FindPerson(ctx *gin.Context) {
	var person Person
	if err := ctx.ShouldBindWith(&person, binding.FormPost); err != nil {
		logger.Errorf("FindPerson bind failed.Error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, person)
}
