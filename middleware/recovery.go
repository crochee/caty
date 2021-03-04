// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package middleware
package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"obs/logger"
	"obs/response"
)

// Recovery panic log
func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					brokenPipe = strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer")
				}
			}
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
			logger.Errorf("[Recovery] %v\n%v\n%v", string(httpRequest), err, string(debug.Stack()))
			if brokenPipe {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					response.Error(http.StatusInternalServerError,
						fmt.Sprintf("broken pipe or connection reset by peer;%v", err)))
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, fmt.Sprint(err)))
		}
	}()
	ctx.Next()
}
