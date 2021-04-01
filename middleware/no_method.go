// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"github.com/gin-gonic/gin"
	"obs/e"
	"strings"
)

// NoMethod 405
func NoMethod(ctx *gin.Context) {
	resp := &e.ErrorResponse{
		Code:    e.MethodNotAllow.String(),
		Message: e.MethodNotAllow.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = e.MethodNotAllow.Chinese()
	}
	ctx.JSON(e.MethodNotAllow.Status(), resp)
}
