// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/4

package logger

type Builder interface {
	Debugf(format string, v ...interface{})
	Debug(message string)
	Infof(format string, v ...interface{})
	Info(message string)
	Errorf(format string, v ...interface{})
	Error(message string)
}

type NopLogger struct {
}

func (n NopLogger) Debugf(format string, v ...interface{}) {
}

func (n NopLogger) Debug(message string) {
}

func (n NopLogger) Infof(format string, v ...interface{}) {
}

func (n NopLogger) Info(message string) {
}

func (n NopLogger) Errorf(format string, v ...interface{}) {
}

func (n NopLogger) Error(message string) {
}
