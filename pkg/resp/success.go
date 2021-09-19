// Author: crochee
// Date: 2021/9/19

// Package resp
package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success response data
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}
