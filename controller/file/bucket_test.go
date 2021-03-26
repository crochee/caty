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
	r := &Name{BucketName: "test"}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.Token)
	router.POST("/v1/bucket", CreateBucket)
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket", body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestHeadBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.GET("/v1/bucket/:bucket_name", GetBucket)
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := util.PerformRequest(router, http.MethodGet, "/v1/bucket/test", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.DELETE("/v1/bucket/:bucket_name", DeleteBucket)
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := util.PerformRequest(router, http.MethodDelete, "/v1/bucket/test", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
