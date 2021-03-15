// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"obs/config"
	"obs/logger"
	"obs/model/db"
	"obs/router"
)

func main() {
	// 初始化
	setup()
	// 初始化请求日志
	requestLog := logger.NewLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
		config.Cfg.ServiceConfig.ServiceInfo.LogLevel)
	srv := &http.Server{
		Addr:    ":8150",
		Handler: router.GinRun(),
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return logger.With(ctx, requestLog)
		},
	}
	go func() {
		logger.Info("obs running...")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown:%v", err)
	}
	logger.Exit("obs server exit!")
}

func setup() {
	// 初始化配置
	config.InitConfig(os.Getenv("config"))
	// 初始化系统日志
	logger.InitSystemLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
		config.Cfg.ServiceConfig.ServiceInfo.LogLevel)
	// 初始化数据库
	db.Setup()
}
