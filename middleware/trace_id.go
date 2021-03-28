// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http/httptrace"

	"obs/logger"
)

// TraceId add trace_id
func TraceId(ctx *gin.Context) {
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
		},
		ConnectStart: func(network, addr string) {
		},
		ConnectDone: func(network, addr string, err error) {

		},
		WroteHeaders: func() {

		},
	}
	ctx.Request = ctx.Request.WithContext(httptrace.WithClientTrace(ctx.Request.Context(), trace))

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
