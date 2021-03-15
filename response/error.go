// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package response

import (
	"bytes"
	"strconv"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error according to the given message structure returns an error
func Error(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// Errors according to the given message and error structure returns an error
func ErrorAll(code int, err error, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: err.Error() + "#" + message,
	}
}

func Errors(code int, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}
}

func (e *ErrorResponse) Error() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(e.Code))
	buf.WriteByte('#')
	buf.WriteString(e.Message)
	return buf.String()
}
