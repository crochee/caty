// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/22

// Package bucket
package bucket

import (
	"bytes"
	"fmt"
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
	r := &model.BucketAction{Action: 3}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/bucket/:bucket_name", CreateBucket)
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket/cptsbuild", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig()
	router := gin.New()
	router.HEAD("/v1/bucket/:bucket_name", HeadBucket)
	w := util.PerformRequest(router, http.MethodHead, "/v1/bucket/cptsBuild", nil, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteBucket(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/bucket/%s", "cptsbuild")
	router := gin.New()
	router.DELETE("/v1/bucket/:bucket_name", DeleteBucket)
	w := util.PerformRequest(router, http.MethodDelete, uri, nil, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
