// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package middleware

import (
	"github.com/gin-gonic/gin"
)

// Verify 鉴权 (每一个请求还需要横向鉴权)
func Verify(ctx *gin.Context) {
	ctx.Next()
}
