package middleware

import (
	"github.com/crochee/lirity/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"caty/pkg/v"
)

// RequestLogger 设置请求日志
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var fieldList []zap.Field
		if traceID := v.GetTraceID(ctx); traceID != "" {
			fieldList = append(fieldList, zap.String("trace_id", traceID))
		}
		fieldList = append(fieldList, zap.String("client_ip", c.ClientIP()))
		c.Request = c.Request.WithContext(logger.With(ctx, log.With(fieldList...)))
		c.Next()
	}
}
