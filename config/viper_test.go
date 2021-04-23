// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/4/23

// Package config
package config

import "testing"

func TestViperInfo_Decode(t *testing.T) {
	v := ViperInfo{path: "../conf/config.yml"}
	info, err := v.Decode()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)

}
