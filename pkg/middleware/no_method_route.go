// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"github.com/crochee/lirity/e"
	"github.com/gin-gonic/gin"
)

// NoRoute 404
func NoRoute(ctx *gin.Context) {
	e.GinErrorCode(ctx, e.ErrNotFound)
}

// NoMethod 405
func NoMethod(ctx *gin.Context) {
	e.GinErrorCode(ctx, e.ErrNotAllowMethod)
}
