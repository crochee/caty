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
		log := logger.FromContext(ctx.Request.Context())
		if log != nil {
			log.Logger = log.Logger.With(zap.String("trace_id", tracedId))
			log.LoggerSugar = log.LoggerSugar.With("trace_id", tracedId)
			logger.With(ctx.Request.Context(), log)
		}
	}
	ctx.Next()
}
