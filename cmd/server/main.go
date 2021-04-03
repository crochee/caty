// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"flag"
	"syscall"
	"time"

	"github.com/crochee/uid"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.uber.org/zap/zapgrpc"

	"obs/cmd"
	"obs/config"
	"obs/cron"
	"obs/logger"
	"obs/message"
	"obs/model/db"
	"obs/model/etcdx"
	"obs/router"
	"obs/transport/httpx"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(*configFile)
	// 初始化系统日志
	logger.InitSystemLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
		config.Cfg.ServiceConfig.ServiceInfo.LogLevel)
	// 初始化请求日志
	requestLog := logger.NewLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
		config.Cfg.ServiceConfig.ServiceInfo.LogLevel)
	httpSrv := httpx.New(":8150",
		httpx.WithContext(ctx),
		httpx.WithLog(requestLog),
		httpx.WithHandler(router.GinRun()),
		httpx.WithBeforeStart(
			func(ctx context.Context) {
				cron.Setup() // cron 初始化
			},
			func(ctx context.Context) {
				if err := db.Setup(ctx); err != nil {
					logger.Fatalf(err.Error())
				}
			}),
		httpx.WithAfterStop(
			func(ctx context.Context) {
				cron.New().Stop() // 关闭定时器
			},
			func(ctx context.Context) {
				// 数据库关闭连接池
				if err := db.Close(); err != nil {
					logger.Errorf("db forced to shutdown:%v", err)
				}
			},
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Logger(zapgrpc.NewLogger(requestLog.Logger)),
		grpc.Timeout(30*time.Second),
	)
	if err := message.Setup(ctx); err != nil {
		logger.Fatal(err.Error())
	}
	etcd, err := etcdx.NewEtcdRegistry()
	if err != nil {
		logger.Fatal(err.Error())
	}
	//grpcClient,err:=grpc.Dial(
	//	ctx,
	//	grpc.WithDiscovery(etcd),
	//)
	app := kratos.New(
		kratos.ID(uid.New().String()),
		kratos.Name(cmd.ServiceName),
		kratos.Version(cmd.Version),
		kratos.Server(httpSrv, grpcSrv),
		kratos.Signal(syscall.SIGINT),
		kratos.Registrar(etcd),
		kratos.Logger(zapgrpc.NewLogger(requestLog.Logger)),
	)
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}
