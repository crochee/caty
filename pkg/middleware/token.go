// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"obs/pkg/log"
	"obs/pkg/service/business/tokenx"
	"obs/pkg/v"
)

// Token add trace_id
func Token(ctx *gin.Context) {
	xAuthToken, err := queryToken(ctx)
	if err != nil { // 缺少token 禁止访问
		log.FromContext(ctx.Request.Context()).Error(err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var claims *tokenx.TokenClaims
	if claims, err = tokenx.ParseToken(xAuthToken); err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("parse token failed.Error:%v", err)
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Set("token", claims.Token)
	ctx.Next()
}

func queryToken(ctx *gin.Context) (string, error) {
	sign := ctx.Query("sign")
	if sign != "" {
		signImpl, err := tokenx.ParseSign(sign)
		if err != nil {
			return "", err
		}
		return signImpl.Sign, nil
	}
	xAuthToken := ctx.Request.Header.Get(v.XAuthToken)
	if xAuthToken == "" {
		return "", errors.New("missing token")
	}
	return xAuthToken, nil
}
