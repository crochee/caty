package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/crochee/lib/log"
	"github.com/crochee/lib/routine"
	"github.com/crochee/uid"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"cca/config"
	"cca/internal/host"
	"cca/pkg/db"
	"cca/pkg/etcdx"
	"cca/pkg/ex"
	"cca/pkg/message"
	"cca/pkg/registry"
	"cca/pkg/router"
	"cca/pkg/tlsx"
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
	if err := ex.Loading(); err != nil {
		panic(err)
	}
	// 初始化系统日志
	log.InitSystemLogger(func(option *log.Option) {
		option.Path = viper.GetString("system-log-path")
		option.Level = viper.GetString("system-log-level")
	})

	if err := run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err.Error())
	}
}

func run() error {
	ctx := context.Background()
	g := routine.NewGroup(ctx)
	srv, err := NewServer(ctx)
	if err != nil {
		return err
	}
	// 服务启动流程
	g.Go(srv.Start)
	// 服务关闭流程
	g.Go(srv.Stop)
	// 启动mq
	g.Go(message.Setup)
	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func startAction(ctx context.Context) error {
	// 初始化数据库
	if err := db.Init(ctx); err != nil {
		return err
	}
	defer db.Close()
	if err := validator.Init(); err != nil {
		return err
	}
	log.FromContext(ctx).Infof("%s run on %s", v.ServiceName, gin.Mode())
	return nil
}

func shutdownAction(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	}
	message.Close()
	log.FromContext(ctx).Info("shutting down server...")
	return nil
}

func NewServer(ctx context.Context) (*httpx.HTTPServer, error) {
	r, err := etcdx.NewEtcdRegistry()
	if err != nil {
		return nil, err
	}
	var ip string
	if ip, err = createHost("WLAN"); err != nil {
		return nil, err
	}
	srv := &httpx.HTTPServer{
		Server: &http.Server{
			Handler: router.New(),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		Instance: &registry.ServiceInstance{
			ID:      uid.New().String(),
			Name:    v.ServiceName,
			Version: v.Version,
		},
		Registrar:   r,
		BeforeStart: startAction,
		BeforeStop:  shutdownAction,
	}
	var (
		cfg *tls.Config
		uri = &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:8120", ip),
		}
	)
	if cfg, err = tlsx.TlsConfig(tls.RequireAndVerifyClientCert, tlsx.Config{
		Ca:   "ca.pem",
		Cert: "server.pem",
		Key:  "server-key.pem",
	}); err != nil {
		uri.Scheme = "http"
		log.Warn(err.Error())
	}
	srv.Server.TLSConfig = cfg
	srv.Server.Addr = uri.Host
	srv.Instance.Endpoints = []string{uri.String()}
	return srv, nil
}

func createHost(name string) (string, error) {
	ip, err := host.GetIPByName(name)
	if err == nil {
		return ip.String(), nil
	}
	if ip, err = host.ExternalIP(); err != nil {
		return "", err
	}
	return ip.String(), nil
}
