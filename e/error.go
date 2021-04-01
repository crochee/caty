// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package e

import (
	"bytes"
	"fmt"
)

type ErrorCode struct {
	code    Code
	message string
}

func (e *ErrorCode) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.code.English())
	if e.message != "" {
		buf.WriteString(".In addition,")
		buf.WriteString(e.message)
	}
	return buf.String()
}

func New(code Code, message string) *ErrorCode {
	return &ErrorCode{
		code:    code,
		message: message,
	}
}

func Errorf(code Code, format string, a ...interface{}) *ErrorCode {
	return &ErrorCode{
		code:    code,
		message: fmt.Sprintf(format, a...),
	}
}

func NewCode(code Code) *ErrorCode {
	return &ErrorCode{
		code: code,
	}
}
