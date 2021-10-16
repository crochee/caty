package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/crochee/lib/log"
	"github.com/crochee/lib/routine"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"cca/config"
	"cca/pkg/code"
	"cca/pkg/dbx"
	"cca/pkg/message"
	"cca/pkg/transport/httpx"
	"cca/pkg/v"
	"cca/pkg/validator"
)

var configFile = flag.String("f", "./conf/cca.yml", "the config file")

func main() {
	flag.Parse()
	// 初始化配置
	if err := config.LoadConfig(*configFile); err != nil {
		panic(err)
	}
	if err := code.Loading(); err != nil {
		panic(err)
	}
	// 初始化系统日志
	log.InitSystemLogger(func(option *log.Option) {
		option.Path = viper.GetString("path")
		option.Level = viper.GetString("level")
	})

	if err := run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err.Error())
		log.Sync()
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	log.Sync()
}

func run() error {
	ctx := context.Background()
	g := routine.NewGroup(ctx)
	srv, err := httpx.NewServer(ctx)
	if err != nil {
		return err
	}
	log.Debugf("listen on %s", srv.Server.Addr)
	// 服务启动流程
	g.Go(func(ctx context.Context) error {
		return startAction(ctx, srv)
	})
	// 服务关闭流程
	g.Go(func(ctx context.Context) error {
		return shutdownAction(ctx, srv)
	})
	// 启动mq
	g.Go(message.Setup)
	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func startAction(ctx context.Context, srv *httpx.HTTPServer) error {
	// 初始化数据库
	if err := dbx.Init(ctx); err != nil {
		return err
	}
	defer dbx.Close()
	if err := validator.Init(); err != nil {
		return err
	}
	log.Infof("%s run on %s", v.ServiceName, gin.Mode())
	return srv.Start(ctx)
}

func shutdownAction(ctx context.Context, srv *httpx.HTTPServer) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	}
	message.Close()
	log.Info("shutting down server...")
	return srv.Stop(ctx)
}
