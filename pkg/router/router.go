// Package router cca
//
// The purpose of this service is to provide an application
// that is using object
//
// title: cca
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
	"github.com/crochee/lib/log"
	"github.com/gin-gonic/gin"

	"cca/api"
	"cca/pkg/middleware"
	"cca/pkg/v"
)

// New gin router
func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain())
	router.NoRoute(middleware.NoRoute)
	router.NoMethod(middleware.NoMethod)

	router.Use(middleware.TraceID, middleware.RequestLogger(log.NewLogger()), middleware.Log, middleware.Recovery)

	router.GET("/", api.Version)
	v1Router := router.Group("/" + v.V1API)

	registerAccount(v1Router)
	registerAuth(v1Router)

	return router
}
