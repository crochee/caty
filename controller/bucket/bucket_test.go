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
	"obs/middleware"
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
	router.Use(middleware.Verify)
	router.POST("/v1/bucket/:bucket_name", CreateBucket)
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket/cptsbuild", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig()
	router := gin.New()
	router.Use(middleware.Verify)
	router.HEAD("/v1/bucket/:bucket_name", HeadBucket)
	header := make(http.Header)
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	w := util.PerformRequest(router, http.MethodHead, "/v1/bucket/cptsBuild", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteBucket(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/bucket/%s", "cptsbuild")
	router := gin.New()
	router.Use(middleware.Verify)
	router.DELETE("/v1/bucket/:bucket_name", DeleteBucket)
	header := make(http.Header)
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	w := util.PerformRequest(router, http.MethodDelete, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
