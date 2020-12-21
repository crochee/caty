// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/22

// Package util
package util

import (
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func PerformRequest(f gin.HandlerFunc, method, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(method, path, body)
	f(ctx)
	return w
}
