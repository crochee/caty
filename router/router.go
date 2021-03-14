// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"obs/controller/file"
	_ "obs/docs"
	"obs/middleware"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Token

// GinRun gin router
func GinRun() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain, middleware.TraceId, middleware.Recovery, middleware.Log)
	if gin.Mode() != gin.ReleaseMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}

	v1Router := router.Group("/v1")
	v1Router.Use(middleware.Token)
	{
		// bucket
		v1Router.POST("/bucket", file.CreateBucket)
		v1Router.HEAD("/bucket/:bucket_name", file.HeadBucket)
		v1Router.DELETE("/bucket/:bucket_name", file.DeleteBucket)

		// file
		fileRouter := v1Router.Group("/bucket/:bucket_name")
		{
			fileRouter.POST("/file", file.UploadFile)
			fileRouter.DELETE("/file/:file_name", file.DeleteFile)
			fileRouter.HEAD("/file/:file_name", file.SignFile)
			fileRouter.GET("/file/:file_name", file.DownloadFile)
			//fileRouter.GET("/files", file.FileList)
		}
	}

	return router
}
