// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package response

import (
	"strconv"
	"sync"
)

type ErrorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Error according to the given message structure returns an error
func Error(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    int64(code),
		Message: message,
	}
}

// Errors according to the given message and error structure returns an error
func ErrorAll(code int, err error, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    int64(code),
		Message: err.Error() + "#" + message,
	}
}

func Errors(code int, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:    int64(code),
		Message: err.Error(),
	}
}

func (e *ErrorResponse) Error() string {
	buf := acquireBuf()
	defer releaseBuf(buf)
	strconv.AppendInt(buf, e.Code, 64)
	buf = append(buf, '.')
	buf = append(buf, '#')
	buf = append(buf, e.Message...)
	return string(buf)
}

var bufPool sync.Pool

const _size = 1024 // by default, create 1 KiB buffers

func acquireBuf() []byte {
	v := bufPool.Get()
	if v == nil {
		return make([]byte, 0, _size)
	}
	return v.([]byte)
}

func releaseBuf(buf []byte) {
	buf = buf[:0]
	bufPool.Put(buf)
}
