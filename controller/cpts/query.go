// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package cpts
package cpts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"obs/logger"
)

type QueryRequest struct {
	Name string `json:"name" form:"name"`
	Age  int    `json:"age" form:"age"`
}

// QueryEncode godoc
// @Summary 请求参数处理
// @Description 请求参数编码
// @Tags CPTS
// @Accept application/json
// @Produce  application/json
// @Param name query string false "名字"
// @Param age query int false "年龄"
// @Success 200 {object} QueryRequest
// @Failure 400
// @Failure 500
// @Router /test/query [get]
func QueryEncode(ctx *gin.Context) {
	var req QueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Errorf("QueryEncode bind failed.Error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, req)
}
