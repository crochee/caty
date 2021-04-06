// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"obs/e"
)

// NoRoute 404
func NoRoute(ctx *gin.Context) {
	resp := &e.ErrorResponse{
		Code:    e.NotFound.String(),
		Message: e.NotFound.English(),
		Extra:   "not route",
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = e.NotFound.Chinese()
	}
	ctx.JSON(e.NotFound.Status(), resp)
}
