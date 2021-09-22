// Date: 2021/9/22

// Package router
package router

import (
	"github.com/gin-gonic/gin"

	"obs/api/user"
)

func registerAccount(v1Router *gin.RouterGroup) {
	v1Router.POST("account", user.Register)
}
