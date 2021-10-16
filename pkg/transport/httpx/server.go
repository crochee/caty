// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/2

package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/crochee/lib/log"
	"github.com/crochee/uid"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"cca/internal/host"
	"cca/pkg/etcdx"
	"cca/pkg/registry"
	"cca/pkg/router"
	"cca/pkg/tlsx"
	"cca/pkg/v"
)

func NewServer(ctx context.Context) (*HTTPServer, error) {
	r, err := etcdx.NewEtcdRegistry(func(option *etcdx.Option) {
		option.AddrList = viper.GetStringSlice("etcd.url")
	})
	if err != nil {
		return nil, err
	}
	var ip string
	if gin.Mode() == gin.ReleaseMode {
		if ip, err = createHost("eth0"); err != nil {
			ip = "0.0.0.0"
		}
	}
	srv := &HTTPServer{
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
		Registrar: r,
	}
	var (
		cfg *tls.Config
		uri = &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:8120", ip),
		}
	)
	if cfg, err = tlsx.TLSConfig(tls.RequireAndVerifyClientCert, tlsx.Config{
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

type HTTPServer struct {
	Server    *http.Server
	Instance  *registry.ServiceInstance
	Registrar registry.Registrar
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if s.Registrar != nil {
		if err := s.Registrar.Register(ctx, s.Instance); err != nil {
			return err
		}
	}
	return s.Server.ListenAndServe()
}

var DefaultStopTime = 10 * time.Second

func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.Registrar != nil {
		if err := s.Registrar.Deregister(ctx, s.Instance); err != nil {
			return err
		}
	}
	newCtx, cancel := context.WithTimeout(ctx, DefaultStopTime)
	defer cancel()
	return s.Server.Shutdown(newCtx)
}
