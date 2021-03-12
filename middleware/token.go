// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"obs/logger"
)

const XAuthToken = "X-Auth-Token"

// TraceId add trace_id
func Token(ctx *gin.Context) {
	xAuthToken := ctx.Request.Header.Get(XAuthToken)
	if xAuthToken == "" { // 缺少token 禁止访问
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, err := ctoken.ParseTokenThenRefresh(tokenString, ctx.Request.RemoteAddr, ok)
	if err != nil {
		logger.Errorf("parse token failed.Error:%v", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	defer claims.Put()
	ctx.Set("email", claims.Email)
	ctx.Set("nick", claims.Nick)
	ctx.Set("permissions", claims.Permissions)
	ctx.Set("token", tokenString)
	ctx.Next()
	claims, err := ctoken.ParseTokenThenRefresh(tokenString, ctx.Request.RemoteAddr, ok)
	if err != nil {
		logger.Errorf("parse token failed.Error:%v", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if tracdId != "" {
		log := logger.FromContext(ctx.Request.Context())
		log.Logger = log.Logger.With(zap.String(RequestTraceId, tracdId))
		log.LoggerSugar = log.LoggerSugar.With(RequestTraceId, tracdId)
		logger.With(ctx.Request.Context(), log)
	}
	ctx.Next()
}
