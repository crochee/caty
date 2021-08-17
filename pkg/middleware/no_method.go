// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"github.com/gin-gonic/gin"

	"obs/e"
)

// NoMethod 405
func NoMethod(ctx *gin.Context) {
	e.Error(ctx, e.MethodNotAllow)
}
