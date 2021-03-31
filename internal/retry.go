// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/1

package internal

import "context"

func Retry(ctx context.Context, fn func(ctx context.Context) error, attempts int) error {
	if err := fn(ctx); err != nil {
		if attempts--; attempts > 0 {
			return Retry(ctx, fn, attempts)
		}
		return err
	}
	return nil
}
