// Date: 2021/9/22

// Package router
package router

import (
	"github.com/gin-gonic/gin"

	"caty/api/v1/account"
)

func registerAccount(v1Router *gin.RouterGroup) {
	v1Router.POST("/accounts", account.Register)
	v1Router.GET("/accounts", account.List)
	v1Router.PATCH("/accounts/:id", account.Update)
	v1Router.GET("/accounts/:id", account.Retrieve)
	v1Router.DELETE("/accounts/:id", account.Delete)
	v1Router.POST("/accounts/login", account.Login)
}
