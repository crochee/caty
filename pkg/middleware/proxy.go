// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/5

package middleware

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cca/internal"
)

const abortIndex int8 = math.MaxInt8 / 2 // 中间件最大数+1

type ProxyMiddle struct {
	middleList   []http.Handler
	ProxyBuilder http.Handler
	prefix       string
	index        int8
}

// Proxy example: /proxy/cca/...
func (p *ProxyMiddle) Proxy(ctx *gin.Context) {
	list := strings.SplitN(ctx.Request.URL.Path, "/", 3)
	if len(list) > 1 {
		if list[1] == p.prefix {
			prefix := "/" + p.prefix
			ctx.Request.URL.Path = internal.EnsureLeadingSlash(strings.TrimPrefix(ctx.Request.URL.Path, prefix))
			if ctx.Request.URL.RawPath != "" {
				ctx.Request.URL.RawPath = internal.EnsureLeadingSlash(strings.TrimPrefix(ctx.Request.URL.RawPath, prefix))
			}
			p.Use().ServeHTTP(ctx.Writer, ctx.Request) // 执行中间件
			p.ProxyBuilder.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}
	}
	ctx.Next()
}

// Use Chain returns a Middleware that specifies the chained handler for endpoint.
func (p *ProxyMiddle) Use(middleware ...http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		for index := 0; index < len(middleware); index++ {
			middleware[index].ServeHTTP(writer, request)
		}
	}
}
