// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"obs/logger"
)

// TraceId add trace_id
func TraceId(ctx *gin.Context) {
	tracedId := ctx.Request.Header.Get("trace_id")
	if tracedId != "" {
		if log, ok := logger.WithContext(ctx.Request.Context()).(*logger.Logger); ok {
			tempLog := log.Logger.With(zap.String("trace_id", tracedId))
			ctx.Request = ctx.Request.Clone(logger.NewContext(ctx.Request.Context(), &logger.Logger{
				Logger:      tempLog,
				LoggerSugar: tempLog.Sugar(),
			}))
		}
	}
	ctx.Next()
}
