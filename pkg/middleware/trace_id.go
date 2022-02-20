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
func TraceID(c *gin.Context) {
	tracedID := c.Request.Header.Get(v.XTraceID)
	if tracedID == "" {
		tracedID = id.UV4()
	}
	c.Request.Header.Set(v.XTraceID, tracedID)  // 请求头
	c.Writer.Header().Set(v.XTraceID, tracedID) // 响应头

	c.Request = c.Request.WithContext(v.SetTraceID(c.Request.Context(), tracedID))

	c.Next()
}
