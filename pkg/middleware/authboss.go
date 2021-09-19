// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package middleware

import (
	"github.com/gin-gonic/gin"
	authboss "github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
)

func Auth(ctx *gin.Context) {
	ab := authboss.New()

	//ab.Config.Storage.Server = myDatabaseImplementation
	//ab.Config.Storage.SessionState = mySessionImplementation
	//ab.Config.Storage.CookieState = myCookieImplementation

	ab.Config.Paths.Mount = "/authboss"
	ab.Config.Paths.RootURL = "https://www.example.com/"

	// This is using the renderer from: github.com/volatiletech/authboss
	//ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	// Probably want a MailRenderer here too.

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, false, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}

	// Mount the router to a path (this should be the same as the Mount path above)
	// mux in this example is a chi router, but it could be anything that can route to
	// the Core.Router.
	//mux.Mount("/authboss", http.StripPrefix("/authboss", ab.Config.Core.Router))
}
