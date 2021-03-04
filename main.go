// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"obs/config"
	"obs/logger"
	"obs/router"
)

func main() {
	config.InitConfig()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg.YamlConfig.ServiceInformation.Port),
		Handler: router.GinRun(),
	}
	go func() {
		logger.Exit("obs running...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown:%v", err)
	}
	logger.Exit("obs server exit!")
}
