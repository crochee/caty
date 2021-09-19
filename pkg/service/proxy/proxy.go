// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/5

package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"obs/pkg/log"
	"os"
	"time"

	"github.com/pkg/errors"

	"obs/internal"
	"obs/internal/status"
)

func NewProxyBuilder() (http.Handler, error) {
	caPem, err := os.ReadFile("ca.pem")
	if err != nil {
		return nil, err
	}
	var certPem []byte
	if certPem, err = os.ReadFile("client.pem"); err != nil {
		return nil, err
	}
	var keyPem []byte
	if keyPem, err = os.ReadFile("client-key.pem"); err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPem)
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPem, keyPem); err != nil {
		return nil, err
	}
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.RequestURI = "" // Outgoing request should not have RequestURI

			if _, ok := request.Header["User-Agent"]; !ok {
				request.Header.Set("User-Agent", "OBS")
			}

			// Even if the websocket RFC says that headers should be case-insensitive,
			// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
			// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
			// https://tools.ietf.org/html/rfc6455#page-20
			request.Header["Sec-WebSocket-Key"] = request.Header["Sec-Websocket-Key"]
			request.Header["Sec-WebSocket-Extensions"] = request.Header["Sec-Websocket-Extensions"]
			request.Header["Sec-WebSocket-Accept"] = request.Header["Sec-Websocket-Accept"]
			request.Header["Sec-WebSocket-Protocol"] = request.Header["Sec-Websocket-Protocol"]
			request.Header["Sec-WebSocket-Version"] = request.Header["Sec-Websocket-Version"]
			delete(request.Header, "Sec-Websocket-Key")
			delete(request.Header, "Sec-Websocket-Extensions")
			delete(request.Header, "Sec-Websocket-Accept")
			delete(request.Header, "Sec-Websocket-Protocol")
			delete(request.Header, "Sec-Websocket-Version")
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				Certificates:           []tls.Certificate{certificate},
				RootCAs:                caPool,
				CipherSuites:           []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
				SessionTicketsDisabled: true,
				MinVersion:             tls.VersionTLS12,
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		BufferPool: internal.BufPool,
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			statusCode := http.StatusInternalServerError
			switch {
			case errors.Is(err, io.EOF):
				statusCode = http.StatusBadGateway
			case errors.Is(err, context.Canceled):
				statusCode = status.ClientClosedRequest
			default:
				var netErr net.Error
				if errors.As(err, &netErr) {
					if netErr.Timeout() {
						statusCode = http.StatusGatewayTimeout
					} else {
						statusCode = http.StatusBadGateway
					}
				}
			}
			log := log.FromContext(request.Context())
			text := status.Text(statusCode)
			log.Errorf("%+v '%d %s' caused by: %v", request.URL, statusCode, text, err)
			writer.WriteHeader(statusCode)
			if _, err = writer.Write([]byte(text)); err != nil {
				log.Errorf("Error %v while writing status code", err)
			}
		},
	}, nil
}
