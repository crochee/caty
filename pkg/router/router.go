// Package router cca
//
// The purpose of this service is to provide an application
// that is using object
//
//     title: cca
//     Schemes: http,https
//     Host: localhost:8120
//     Version: 0.0.1
//
// swagger:meta
package router

import (
	"github.com/crochee/lib/log"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "cca/docs"
	"cca/pkg/middleware"
)

// @title cca Swagger API
// @version 1.0
// @description This is a cca server.

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

	router.Use(middleware.TraceId, middleware.RequestLogger(log.NewLogger()), middleware.Log, middleware.Recovery)

	v1Router := router.Group("/v1")

	registerAccount(v1Router)

	//userRouter := v1Router.Group("/user")
	//{
	//	userRouter.POST("/v1/account", user.Register)
	//	userRouter.POST("/login", user.Login)
	//	userRouter.PUT("/modify", user.Modify)
	//}
	//v1Router.Use(middleware.Token)
	//{
	//	// bucket
	//	v1Router.POST("/bucket", file.CreateBucket)
	//	v1Router.GET("/bucket/:bucket_name", file.GetBucket)
	//	v1Router.DELETE("/bucket/:bucket_name", file.DeleteBucket)
	//
	//	// file
	//	fileRouter := v1Router.Group("/bucket/:bucket_name")
	//	{
	//		fileRouter.POST("/file", file.UploadFile)
	//		fileRouter.DELETE("/file/:file_name", file.DeleteFile)
	//		fileRouter.GET("/file/:file_name/sign", file.SignFile)
	//		fileRouter.GET("/file/:file_name", file.DownloadFile)
	//		//fileRouter.GET("/files", file.FileList)
	//	}
	//}

	return router
}
