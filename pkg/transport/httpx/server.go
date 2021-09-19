// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/2

package httpx

import (
	"context"
	"net/http"
	"time"

	"obs/pkg/registry"
)

type HTTPServer struct {
	Server      *http.Server
	Instance    *registry.ServiceInstance
	Registrar   registry.Registrar
	BeforeStart func(ctx context.Context) error
	BeforeStop  func(ctx context.Context) error
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if s.BeforeStart != nil {
		if err := s.BeforeStart(ctx); err != nil {
			return err
		}
	}
	if s.Registrar != nil {
		if err := s.Registrar.Register(ctx, s.Instance); err != nil {
			return err
		}
	}
	return s.Server.ListenAndServe()
}

var DefaultStopTime = 10 * time.Second

func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.BeforeStart != nil {
		if err := s.BeforeStart(ctx); err != nil {
			return err
		}
	}
	if s.Registrar != nil {
		if err := s.Registrar.Deregister(ctx, s.Instance); err != nil {
			return err
		}
	}
	newCtx, cancel := context.WithTimeout(ctx, DefaultStopTime)
	defer cancel()
	return s.Server.Shutdown(newCtx)
}
