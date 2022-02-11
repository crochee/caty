// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/crochee/lirity/id"
	"github.com/gin-gonic/gin"

	"caty/pkg/v"
)

// TraceID add trace_id
func TraceID(ctx *gin.Context) {
	tracedID := ctx.Request.Header.Get(v.XTraceID)
	if tracedID == "" {
		tracedID = id.UV4()
	}
	ctx.Request.Header.Set(v.XTraceID, tracedID)  // 请求头
	ctx.Writer.Header().Set(v.XTraceID, tracedID) // 响应头

	ctx.Request = ctx.Request.WithContext(v.SetTraceID(ctx.Request.Context(), tracedID))

	ctx.Next()
}
