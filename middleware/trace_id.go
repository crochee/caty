// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"

	"obs/logger"
)

// TraceId add trace_id
func TraceId(ctx *gin.Context) {
	tracedId := ctx.Request.Header.Get("trace_id")
	if tracedId != "" {
		log := logger.FromContext(ctx.Request.Context()).With("trace_id", tracedId)
		ctx.Request = ctx.Request.WithContext(logger.WithContext(ctx.Request.Context(), log))
	}
	ctx.Next()
}
