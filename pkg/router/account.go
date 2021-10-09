// Date: 2021/9/22

// Package router
package router

import (
	"cca/api/account"
	"github.com/gin-gonic/gin"
)

func registerAccount(v1Router *gin.RouterGroup) {
	v1Router.POST("/account", account.Register)
	v1Router.PATCH("/account", account.Modify)
	v1Router.GET("/account", account.Query)
	v1Router.DELETE("/account", account.Delete)
}
