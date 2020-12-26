// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"obs/controller/file"

	"obs/controller/bucket"
	"obs/controller/cpts"
	_ "obs/docs"
	"obs/middleware"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// GinRun gin router
func GinRun() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Limit, middleware.Log, middleware.CrossDomain,
		middleware.Recovery, middleware.Verify)

	if gin.Mode() != gin.ReleaseMode {
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	testRouter := router.Group("/test")
	{
		testRouter.GET("/query", cpts.QueryEncode)
		testRouter.POST("/person", cpts.FindPerson)
		testRouter.POST("/all", cpts.PayAll)
	}

	v1Router := router.Group("/v1")

	bucketRouter := v1Router.Group("/bucket")
	{
		bucketRouter.POST("/:bucket_name", bucket.CreateBucket)
		bucketRouter.HEAD("/:bucket_name", bucket.HeadBucket)
		bucketRouter.DELETE("/:bucket_name", bucket.DeleteBucket)
	}

	fileRouter := v1Router.Group("/file")
	{
		fileRouter.POST("/bucket/:bucket_name", file.UploadFile)
		fileRouter.DELETE("/bucket/:bucket_name", file.DeleteFile)
		fileRouter.DELETE("/bucket/:bucket_name", file.SignFile)
		fileRouter.GET("/bucket/:bucket_name", file.DownloadFile)
	}
	return router
}
