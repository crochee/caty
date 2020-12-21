// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/21

// Package verify
package verify

import (
	"testing"
)

func TestNewToken(t *testing.T) {
	tmp := NewToken("cpts")
	tmp.AddAction(Admin)
	ak, sk, err := tmp.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ak:%s sk:%s", ak, sk)
	if err = tmp.Verify(sk); err != nil {
		t.Fatal(err)
	}
}
