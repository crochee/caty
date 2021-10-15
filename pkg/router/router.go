// Package router cca
//
// The purpose of this service is to provide an application
// that is using object
//
//		title: cca
//		Schemes: http, https
//		Host: localhost:8120
//		Version: 0.0.1
//		security:
//			token:
//				type: apiKey
//				name: X-Auth-Token
//				in: header
// swagger:meta
package router

import (
	"github.com/crochee/lib/log"
	"github.com/gin-gonic/gin"

	"cca/pkg/middleware"
)

// New gin router
func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain())
	router.NoRoute(middleware.NoRoute)
	router.NoMethod(middleware.NoMethod)

	router.Use(middleware.TraceID, middleware.RequestLogger(log.NewLogger()), middleware.Log, middleware.Recovery)

	v1Router := router.Group("/v1")

	registerAccount(v1Router)

	return router
}
