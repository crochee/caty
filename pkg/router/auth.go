// Date: 2021/10/16

// Package router
package router

import (
	"github.com/gin-gonic/gin"

	"cca/api/v1/auth"
)

func registerAuth(v1Router *gin.RouterGroup) {
	v1Router.POST("/auth/sign", auth.Sign)
	v1Router.POST("/auth/parse", auth.Parse)
}
