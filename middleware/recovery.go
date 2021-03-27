// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package middleware
package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"obs/e"
	"obs/logger"
)

// Recovery panic log
func Recovery(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			var brokenPipe bool
			if ne, ok := r.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne.Err, &se) {
					brokenPipe = strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer")
				}
			}
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
			logger.FromContext(ctx.Request.Context()).Errorf("[Recovery] %s\n%v\n%s", httpRequest, r, debug.Stack())
			if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
				if brokenPipe {
					ctx.AbortWithStatusJSON(e.Recovery.Status(), &e.ErrorResponse{
						Code:    e.Recovery.String(),
						Message: e.Recovery.Chinese(),
						Extra:   fmt.Sprintf("broken pipe or connection reset by peer;%v", r),
					})
					return
				}
				ctx.AbortWithStatusJSON(e.Recovery.Status(), &e.ErrorResponse{
					Code:    e.Recovery.String(),
					Message: e.Recovery.Chinese(),
					Extra:   fmt.Sprint(r),
				})
				return
			}
			if brokenPipe {
				ctx.AbortWithStatusJSON(e.Recovery.Status(), &e.ErrorResponse{
					Code:    e.Recovery.String(),
					Message: e.Recovery.English(),
					Extra:   fmt.Sprintf("broken pipe or connection reset by peer;%v", r),
				})
				return
			}
			ctx.AbortWithStatusJSON(e.Recovery.Status(), &e.ErrorResponse{
				Code:    e.Recovery.String(),
				Message: e.Recovery.English(),
				Extra:   fmt.Sprint(r),
			})
		}
	}()
	ctx.Next()
}
