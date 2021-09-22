// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"

	"obs/pkg/id"
	"obs/pkg/v"
)

// TraceId add trace_id
func TraceId(ctx *gin.Context) {
	tracedId := ctx.Request.Header.Get(v.XTraceID)
	if tracedId != "" {
		tracedId = id.Uuid()
	}
	ctx.Request.Header.Set(v.XTraceID, tracedId)  // 请求头
	ctx.Writer.Header().Set(v.XTraceID, tracedId) // 响应头

	ctx.Request = ctx.Request.WithContext(v.SetTraceID(ctx.Request.Context(), tracedId))

	ctx.Next()
}