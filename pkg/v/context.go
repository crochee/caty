// Date: 2021/9/19

// Package v
package v

import "context"

type traceIDKey struct{}

// SetTraceID Add traceID to context.Context.
func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID Get the traceID from context.Context.
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return ""
	}
	return traceID
}

type hostKey struct{}

// SetHost Add host to context.Context.
func SetHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, hostKey{}, host)
}

// GetHost Get the host from context.Context.
func GetHost(ctx context.Context) string {
	host, ok := ctx.Value(hostKey{}).(string)
	if !ok {
		return ""
	}
	return host
}
