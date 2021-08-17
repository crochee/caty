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

// CrossDomain skip the cross-domain phase
func CrossDomain(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodOptions &&
		ctx.GetHeader("Access-Control-Request-Method") != "" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,HEAD,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,X-Auth-Token")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
