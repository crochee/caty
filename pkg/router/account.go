// Date: 2021/9/22

// Package router
package router

import (
	"github.com/gin-gonic/gin"

	"cca/api/account"
)

func registerAccount(v1Router *gin.RouterGroup) {
	v1Router.POST("/account", account.Register)
	v1Router.PATCH("/account", account.Update)
	v1Router.GET("/account", account.Retrieve)
	v1Router.GET("/account/:id", account.RetrieveSingle)
	v1Router.DELETE("/account/:id", account.Delete)
}
