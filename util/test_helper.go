// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/22

// Package util
package util

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// PerformRequest unit test function httptest.ResponseRecorder
func PerformRequest(r http.Handler, method, path string, body io.Reader,
	headers http.Header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header = headers
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
