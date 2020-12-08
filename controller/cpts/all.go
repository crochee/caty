// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package cpts
package cpts

import (
	"io"

	"github.com/gin-gonic/gin"

	"obs/logger"
)

// QueryEncode godoc
// @Summary 返回请求数据
// @Description 返回请求数据
// @Tags CPTS
// @Accept application/octet-stream
// @Produce  application/octet-stream
// @Param request body string true "数据流"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /test/all [any]
func PayAll(ctx *gin.Context) {
	n, err := io.Copy(ctx.Writer, ctx.Request.Body)
	if err != nil {
		logger.Errorf("payAll %d failed.Error:%v", n, err)
	}
}
