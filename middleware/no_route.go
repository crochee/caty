// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"github.com/gin-gonic/gin"

	"obs/e"
)

// NoRoute 404
func NoRoute(ctx *gin.Context) {
	e.Error(ctx, e.NotFound)
}
