// Date: 2021/9/19

// Package middleware
package middleware

import (
	"github.com/gin-gonic/gin"

	"obs/pkg/log"
	"obs/pkg/v"
)

// RequestLogger 设置请求日志
func RequestLogger(logger log.Builder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fieldList []log.Field
		ctx := c.Request.Context()
		if traceID := v.GetTraceID(ctx); traceID != "" {
			fieldList = append(fieldList, log.Field{
				Key:   "TraceId",
				Value: traceID,
			})
		}
		fieldList = append(fieldList, log.Field{
			Key:   "ClientIp",
			Value: c.ClientIP(),
		})
		c.Request = c.Request.WithContext(log.WithContext(ctx, logger.With(fieldList...)))
		c.Next()
	}
}
