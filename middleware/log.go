// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package middleware
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"obs/logger"
)

func Log(ctx *gin.Context) {
	// Start timer
	start := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	// Process request
	ctx.Next()
	// Log only when path is not being skipped

	param := gin.LogFormatterParams{
		Request: ctx.Request,
		Keys:    ctx.Keys,
	}
	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	param.ClientIP = ctx.ClientIP()
	param.Method = ctx.Request.Method
	param.StatusCode = ctx.Writer.Status()
	param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

	param.BodySize = ctx.Writer.Size()

	if raw != "" {
		path = path + "?" + raw
	}
	param.Path = path
	logger.Info(defaultLogFormatter(param))
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[INFO] %v |%s %3d %s| %13v | %15s |%s %-7s %s |%8d| %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.BodySize,
		param.Path,
		param.ErrorMessage,
	)
}
