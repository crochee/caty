// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"obs/logger"
)

const RequestTraceId = "trace_id"

// TraceId add trace_id
func TraceId(ctx *gin.Context) {
	tracedId := ctx.Request.Header.Get(RequestTraceId)
	if tracedId != "" {
		log := logger.FromContext(ctx.Request.Context())
		if log != nil {
			log.Logger = log.Logger.With(zap.String(RequestTraceId, tracedId))
			log.LoggerSugar = log.LoggerSugar.With(RequestTraceId, tracedId)
			logger.With(ctx.Request.Context(), log)
		}
	}
	ctx.Next()
}
