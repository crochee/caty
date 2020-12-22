// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/22

// Package bucket
package bucket

import (
	"bytes"
	"net/http"

	"testing"

	"github.com/json-iterator/go"

	"obs/config"
	"obs/model"
	"obs/util"
)

func TestCreateBucket(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	r := &model.CreateBucket{
		BucketName: model.BucketName{BucketName: "cptsbuild"},
		Action:     0,
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	w := util.PerformRequest(CreateBucket, http.MethodPost, "/v1/bucket", body)
	if w.Code != 200 {
		t.Fatalf("code got:%d want:%d", w.Code, 200)
	}
	t.Log(w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	r := &model.BucketName{
		BucketName: "cptsbuild",
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	w := util.PerformRequest(HeadBucket, http.MethodHead, "/v1/bucket", body)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
