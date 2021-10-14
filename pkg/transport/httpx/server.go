// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/2

package httpx

import (
	"context"
	"net/http"
	"time"

	"cca/pkg/registry"
)

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
