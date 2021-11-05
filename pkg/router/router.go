// Package router caty
//
// The purpose of this service is to provide an application
// that is using auth
//
// title: caty
// Schemes: http, https
// Host: localhost:8120
// Version: 1.0.1
// security:
// - token:
// 		type: apiKey
//		name: X-Auth-Token
//		in: header
//  - ak:
//      type: apiKey
//      name: ak
//      in: query
// swagger:meta
package router

import (
	"github.com/crochee/lirity/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"caty/api"
	"caty/pkg/middleware"
	"caty/pkg/v"
)

// New gin router
func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain())
	router.NoRoute(middleware.NoRoute)
	router.NoMethod(middleware.NoMethod)

	router.Use(middleware.TraceID,
		middleware.RequestLogger(log.NewLogger(func(option *log.Option) {
			option.Path = viper.GetString("path")
			option.Level = log.JudgeLevel(viper.GetString("level"), gin.Mode())
		})),
		middleware.Log,
		middleware.Recovery,
	)

	router.GET("/version", api.Version)
	v1Router := router.Group("/" + v.V1API)

	registerAccount(v1Router)
	registerAuth(v1Router)

	return router
}
