// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"obs/controller/cpts"
	_ "obs/docs"
	"obs/middleware"
)

// @title console Swagger API
// @version 1.0
// @description This is a console server celler server.

// @host 10.78.74.37:8150

func GinRun() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Limit, middleware.Log, middleware.CrossDomain, middleware.Recovery)

	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	testRouter := router.Group("/test")
	{
		testRouter.GET("/query", cpts.QueryEncode)
		testRouter.POST("/person", cpts.FindPerson)
		testRouter.POST("/all", cpts.PayAll)
	}
	return router
}
