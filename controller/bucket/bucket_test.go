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

	"github.com/gin-gonic/gin"
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
		Action:     3,
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/bucket", CreateBucket)
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	r := &model.BucketName{
		BucketName: "cptsBuild",
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.HEAD("/v1/bucket", HeadBucket)
	w := util.PerformRequest(router, http.MethodHead, "/v1/bucket", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteBucket(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	r := &model.SimpleBucket{
		BucketName: model.BucketName{BucketName: "cptsbuild"},
		AkSk: model.AkSk{
			Ak: "2p9lq5pqg2r3g08uqpj1oa0qkphb1amjhng0",
			Sk: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJzZWNyZXQiOiJNbkE1YkhFMWNIRm5Nbkl6WnpBNGRYRndhakZ2WVRCeGEzQm9ZakZoYldwb2JtY3ciLCJidWNrZXQiOiJjcHRzYnVpbGQiLCJhY3Rpb24iOnsiMyI6e319fQ.r6JdbKTcFOZJ2uHNvGXbh8q0fJG9gvOBOQfJm5n6nHU",
		},
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.DELETE("/v1/bucket", DeleteBucket)
	w := util.PerformRequest(router, http.MethodDelete, "/v1/bucket", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
