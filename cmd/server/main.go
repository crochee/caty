// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"obs/pkg/etcdx"
	"obs/pkg/routine"
	"obs/pkg/v"
	"obs/pkg/validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crochee/uid"
	kratos "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"obs/config"
	"obs/pkg/cron"
	"obs/pkg/log"
	"obs/pkg/message"
	"obs/pkg/model/db"
	"obs/pkg/router"
	"obs/pkg/service/transport/httpx"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()
	// 初始化配置
	config.LoadConfig(*configFile)

	// 初始化系统日志
	log.InitSystemLogger(func(option *log.Option) {
		option.Path = viper.GetString("system-log-path")
		option.Level = viper.GetString("system-log-level")
	})

	if err := server(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err.Error())
	}

	// 初始化请求日志
	requestLog := log.NewLogger(pathFunc, levelFunc)
	//httpSrv:=http.NewServer(
	//	http.Address(":8150"),
	//	http.Logger(requestLog),
	//	)
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
				cron.Cron().Stop() // 关闭定时器
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
		grpc.Logger(requestLog),
		grpc.Timeout(30*time.Second),
	)
	if err := message.Setup(ctx); err != nil {
		log.Fatal(err.Error())
	}
	etcd, err := etcdx.NewEtcdRegistry(func(option *etcdx.Option) {
		option.Context = ctx
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	//grpcClient,err:=grpc.Dial(
	//	ctx,
	//	grpc.WithDiscovery(etcd),
	//)
	app := kratos.New(
		kratos.ID(uid.New().String()),
		kratos.Name(v.ServiceName),
		kratos.Version(v.Version),
		kratos.Server(httpSrv, grpcSrv),
		kratos.Signal(syscall.SIGINT, syscall.SIGTERM),
		kratos.Registrar(etcd),
		kratos.Logger(requestLog),
	)
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func tlsConfig() *tls.Config {
	caPem, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	var certPem []byte
	if certPem, err = os.ReadFile("server.pem"); err != nil {
		log.Error(err.Error())
		return nil
	}
	var keyPem []byte
	if keyPem, err = os.ReadFile("server-key.pem"); err != nil {
		log.Error(err.Error())
		return nil
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPem)

	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPem, keyPem); err != nil {
		log.Error(err.Error())
		return nil
	}

	return &tls.Config{
		Certificates:           []tls.Certificate{certificate},
		ClientAuth:             tls.RequireAndVerifyClientCert,
		ClientCAs:              caPool,
		CipherSuites:           []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		SessionTicketsDisabled: true,
		MinVersion:             tls.VersionTLS12,
	}
}

func server() error {
	ctx := context.Background()
	g := routine.NewGroup(ctx)

	srv := &http.Server{
		Addr:      "0.0.0.0:30085",
		Handler:   nil,
		TLSConfig: tlsConfig(),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	// 服务启动流程
	g.Go(func(ctx context.Context) error {
		return startServer(ctx, srv)
	})
	// 服务关闭流程
	g.Go(func(ctx context.Context) error {
		return shutdownServer(ctx, srv)
	})
	g.Go(message.Setup)
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func startServer(ctx context.Context, srv *http.Server) error {
	// 初始化数据库
	if err := db.Init(ctx); err != nil {
		return err
	}
	defer db.Close()
	if err := validator.Init(); err != nil {
		return err
	}
	log.FromContext(ctx).Infof("%s run on %s", v.ServiceName, gin.Mode())
	return srv.ListenAndServe()
}

const DefaultStopTime = 10 * time.Second

func shutdownServer(ctx context.Context, srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	}
	log.FromContext(ctx).Info("shutting down server...")

	newCtx, cancel := context.WithTimeout(ctx, DefaultStopTime)
	defer cancel()
	message.Close()
	err := srv.Shutdown(newCtx)
	log.FromContext(ctx).Info("Server exiting")
	return err
}
