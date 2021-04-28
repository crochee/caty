// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
	"obs/service/transport/httpx"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(*configFile)
	// 初始化系统日志
	pathFunc := func(option *logger.Option) {
		option.Path = config.Cfg.ServiceConfig.ServiceInfo.LogPath
	}
	levelFunc := func(option *logger.Option) {
		option.Path = config.Cfg.ServiceConfig.ServiceInfo.LogLevel
	}
	logger.InitSystemLogger(pathFunc, levelFunc)
	// 初始化请求日志
	requestLog := logger.NewLogger(pathFunc, levelFunc)
	httpSrv := httpx.New(":8150",
		httpx.WithContext(ctx),
		httpx.WithTls(tlsConfig()),
		httpx.WithLog(requestLog),
		httpx.WithHandler(router.GinRun()),
		httpx.WithBeforeStart(
			func(ctx context.Context) error {
				cron.Setup() // cron 初始化
				return nil
			},
			func(ctx context.Context) error {
				return db.Setup(ctx)
			}),
		httpx.WithAfterStop(
			func(ctx context.Context) error {
				cron.New().Stop() // 关闭定时器
				return nil
			},
			func(ctx context.Context) error {
				// 数据库关闭连接池
				return db.Close()
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
	etcd, err := etcdx.NewEtcdRegistry(func(option *etcdx.Option) {
		option.Context = ctx
	})
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
		kratos.Signal(syscall.SIGINT, syscall.SIGTERM),
		kratos.Registrar(etcd),
		kratos.Logger(zapgrpc.NewLogger(requestLog.Logger)),
	)
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

func tlsConfig() *tls.Config {
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPEMBlock)

	return &tls.Config{
		Rand:                        nil,
		Time:                        nil,
		Certificates:                nil,
		NameToCertificate:           nil,
		GetCertificate:              nil,
		GetClientCertificate:        nil,
		GetConfigForClient:          nil,
		VerifyPeerCertificate:       nil,
		VerifyConnection:            nil,
		RootCAs:                     nil,
		NextProtos:                  nil,
		ServerName:                  "",
		ClientAuth:                  tls.RequireAndVerifyClientCert,
		ClientCAs:                   nil,
		InsecureSkipVerify:          false,
		CipherSuites:                nil,
		PreferServerCipherSuites:    false,
		SessionTicketsDisabled:      false,
		SessionTicketKey:            [32]byte{},
		ClientSessionCache:          nil,
		MinVersion:                  0,
		MaxVersion:                  0,
		CurvePreferences:            nil,
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}
}
