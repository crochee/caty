// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/29

// Package proxy
package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
)

func ReverseProxy(ctx *gin.Context) {
	derector := func() {

	}
	proxyRouter := &httputil.ReverseProxy{
		Director:       nil,
		Transport:      nil,
		FlushInterval:  0,
		ErrorLog:       nil,
		BufferPool:     nil,
		ModifyResponse: nil,
		ErrorHandler:   nil,
	}
	proxyRouter.ServeHTTP(ctx.Writer, ct)
}
