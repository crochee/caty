// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package logger

import "context"

type loggerKey struct{}

// NewContext Adds fields.
func NewContext(ctx context.Context, log Builder) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// WithContext Gets the logger from context.
func WithContext(ctx context.Context) Builder {
	l, ok := ctx.Value(loggerKey{}).(Builder)
	if !ok {
		return NoLogger{}
	}
	return l
}
