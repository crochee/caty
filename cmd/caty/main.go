package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/crochee/lirity/logger"
	"github.com/crochee/lirity/routine"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"

	"caty/config"
	"caty/pkg/code"
	"caty/pkg/dbx"
	"caty/pkg/message"
	"caty/pkg/transport/httpx"
	"caty/pkg/v"
	"caty/pkg/validator"
)

var configFile = flag.String("f", "./conf/caty.yaml", "the config file")

func main() {
	flag.Parse()
	// 初始化配置
	if err := config.LoadConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := code.Loading(); err != nil {
		log.Fatal(err)
	}
	if mode := strings.ToLower(viper.GetString("GIN_MODE")); mode != "" {
		gin.SetMode(mode)
	}
	// 初始化系统日志
	zap.ReplaceGlobals(logger.New(
		logger.WithFields(zap.String("service", v.ServiceName)),
		logger.WithLevel(viper.GetString("level")),
		logger.WithWriter(logger.SetWriter(viper.GetString("path")))))

	if err := run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	g := routine.NewGroup(ctx)
	srv, err := httpx.NewServer(ctx)
	if err != nil {
		return err
	}
	zap.S().Debugf("listen on %s", srv.Server.Addr)
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
	if err := validator.Init(); err != nil {
		return err
	}
	zap.S().Infof("run on %s", gin.Mode())
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
	zap.L().Info("shutting down server...")
	return srv.Stop(ctx)
}
