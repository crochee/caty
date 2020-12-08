// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package middleware
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CrossDomain(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,X-Auth-Token")
	ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
	if ctx.Request.Method == http.MethodOptions {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,PATCH,OPTIONS")
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
