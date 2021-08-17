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

	"os"
	"time"

	"github.com/pkg/errors"

	"obs/internal/host"
	"obs/pkg/logx"
)

type (
	Option struct {
		tls         *tls.Config
		ctx         context.Context
		logger      logx.Builder
		handler     http.Handler
		beforeStart []func(ctx context.Context) error
		afterStop   []func(ctx context.Context) error

		shutdownDelayTimeout time.Duration
		readTimeout          time.Duration
		readHeaderTimeout    time.Duration
		writeTimeout         time.Duration
		idleTimeout          time.Duration
		maxHeaderBytes       int
	}

	Server struct {
		*http.Server
		ctx                  context.Context
		logger               logx.Builder
		sigList              []os.Signal
		beforeStart          []func(ctx context.Context) error
		afterStop            []func(ctx context.Context) error
		shutdownDelayTimeout time.Duration
	}
)

func WithTls(tls *tls.Config) func(*Option) {
	return func(opt *Option) {
		opt.tls = tls
	}
}

func WithContext(ctx context.Context) func(*Option) {
	return func(opt *Option) {
		opt.ctx = ctx
	}
}

func WithHandler(handler http.Handler) func(*Option) {
	return func(opt *Option) {
		opt.handler = handler
	}
}

func WithLog(log logx.Builder) func(*Option) {
	return func(opt *Option) {
		opt.logger = log
	}
}

func WithShutdownDelayTimeout(t time.Duration) func(*Option) {
	return func(opt *Option) {
		opt.shutdownDelayTimeout = t
	}
}

func WithBeforeStart(beforeStart ...func(ctx context.Context) error) func(*Option) {
	return func(opt *Option) {
		opt.beforeStart = beforeStart
	}
}

func WithAfterStop(afterStop ...func(ctx context.Context) error) func(*Option) {
	return func(opt *Option) {
		opt.afterStop = afterStop
	}
}

func New(addr string, opts ...func(*Option)) *Server {
	option := &Option{
		ctx:                  context.Background(),
		handler:              http.DefaultServeMux,
		shutdownDelayTimeout: 15 * time.Second,
	}
	for _, opt := range opts {
		opt(option)
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           option.handler,
		TLSConfig:         option.tls,
		ReadTimeout:       option.readTimeout,
		ReadHeaderTimeout: option.readHeaderTimeout,
		WriteTimeout:      option.writeTimeout,
		IdleTimeout:       option.idleTimeout,
		MaxHeaderBytes:    option.maxHeaderBytes,
		BaseContext: func(listener net.Listener) context.Context {
			return option.ctx
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			if option.logger == nil {
				return option.ctx
			}
			return logx.WithContext(ctx, option.logger)
		},
	}
	return &Server{
		Server:               srv,
		ctx:                  option.ctx,
		logger:               option.logger,
		beforeStart:          option.beforeStart,
		afterStop:            option.afterStop,
		shutdownDelayTimeout: option.shutdownDelayTimeout,
	}
}

func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.Addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

func (s *Server) Start() error {
	for _, f := range s.beforeStart {
		if err := f(s.ctx); err != nil {
			return err
		}
	}
	go func() {
		s.logger.Info("http server running...")
		if s.Server.TLSConfig == nil {
			if err := s.Server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				s.logger.Errorf(err.Error())
			}
			return
		}
		if err := s.Server.ListenAndServeTLS("", ""); err != nil &&
			errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf(err.Error())
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if s.shutdownDelayTimeout != 0 {
		ctx, cancel = context.WithTimeout(s.ctx, s.shutdownDelayTimeout)
	} else {
		ctx, cancel = context.WithCancel(s.ctx)
	}
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := s.Server.Shutdown(ctx); err != nil {
		s.logger.Errorf(err.Error())
	}

	for _, f := range s.afterStop {
		if err := f(ctx); err != nil {
			s.logger.Errorf(err.Error())
		}
	}
	return nil
}
