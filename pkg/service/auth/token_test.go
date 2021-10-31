// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package auth

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	var tokenImpl = &TokenClaims{
		Now: time.Now().Unix(),
		Token: &Token{
			AccountID: "123",
			UserID:    "test123",
			Permission: map[string]uint8{
				"caty": Admin,
			},
		},
	}
	value, err := tokenImpl.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
}
