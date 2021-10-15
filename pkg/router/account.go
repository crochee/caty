// Date: 2021/9/22

// Package router
package router

import (
	"github.com/gin-gonic/gin"

	"cca/api/v1/account"
)

func registerAccount(v1Router *gin.RouterGroup) {
	v1Router.POST("/account", account.Register)
	v1Router.GET("/account", account.Retrieves)
	v1Router.PATCH("/account/:id", account.Update)
	v1Router.GET("/account/:id", account.Retrieve)
	v1Router.DELETE("/account/:id", account.Delete)
}
