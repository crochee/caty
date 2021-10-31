// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"github.com/crochee/lirity/e"
	"github.com/gin-gonic/gin"

	"cca/pkg/resp"
)

// NoRoute 404
func NoRoute(ctx *gin.Context) {
	resp.Error(ctx, e.ErrNotFound)
}

// NoMethod 405
func NoMethod(ctx *gin.Context) {
	resp.Error(ctx, e.ErrMethodNotAllow)
}
