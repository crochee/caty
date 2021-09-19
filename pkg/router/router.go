// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"obs/api/file"
	"obs/api/user"
	_ "obs/docs"
	"obs/pkg/middleware"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Token

// New gin router
func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain)
	router.NoRoute(middleware.NoRoute)
	router.NoMethod(middleware.NoMethod)
	if gin.Mode() == gin.DebugMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}

	router.Use(middleware.TraceId, middleware.RequestLogger(nil), middleware.Log, middleware.Recovery)

	v1Router := router.Group("/v1")
	userRouter := v1Router.Group("/user")
	{
		userRouter.POST("/register", user.Register)
		userRouter.POST("/login", user.Login)
		userRouter.PUT("/modify", user.Modify)
	}
	v1Router.Use(middleware.Token)
	{
		// bucket
		v1Router.POST("/bucket", file.CreateBucket)
		v1Router.GET("/bucket/:bucket_name", file.GetBucket)
		v1Router.DELETE("/bucket/:bucket_name", file.DeleteBucket)

		// file
		fileRouter := v1Router.Group("/bucket/:bucket_name")
		{
			fileRouter.POST("/file", file.UploadFile)
			fileRouter.DELETE("/file/:file_name", file.DeleteFile)
			fileRouter.GET("/file/:file_name/sign", file.SignFile)
			fileRouter.GET("/file/:file_name", file.DownloadFile)
			//fileRouter.GET("/files", file.FileList)
		}
	}

	return router
}
