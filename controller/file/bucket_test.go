// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/22

// Package bucket
package file

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"

	"obs/config"
	"obs/middleware"
	"obs/model/db"
	"obs/util"
)

func TestCreateBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	body := new(bytes.Buffer)
	r := &Name{BucketName: "obs"}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.Token)
	router.POST("/v1/bucket", CreateBucket)
	header := make(http.Header)
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxNTozODoyNy4wNDQ4OTM1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MX19fQ.5gpZiBgclzoftN1k2npgPmlHE5Dukcf7MkSrdLZfRSs")
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket", body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.HEAD("/v1/bucket/:bucket_id", HeadBucket)
	header := make(http.Header)
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxNTozODoyNy4wNDQ4OTM1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MX19fQ.5gpZiBgclzoftN1k2npgPmlHE5Dukcf7MkSrdLZfRSs")
	w := util.PerformRequest(router, http.MethodHead, "/v1/bucket/9", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.DELETE("/v1/bucket/:bucket_id", DeleteBucket)
	header := make(http.Header)
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxNTozODoyNy4wNDQ4OTM1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MX19fQ.5gpZiBgclzoftN1k2npgPmlHE5Dukcf7MkSrdLZfRSs")
	w := util.PerformRequest(router, http.MethodDelete, "/v1/bucket/9", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
