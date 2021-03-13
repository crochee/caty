// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/25

package file

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"obs/config"
	"obs/middleware"
	"obs/model/db"
	"obs/util"
)

func TestUploadFile(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	_ = mw.WriteField("path", "data")
	f, _ := mw.CreateFormFile("file", "test1.txt")
	_, _ = f.Write([]byte(`hello world`))
	_ = mw.Close()
	router := gin.New()
	router.Use(middleware.Token)
	header := make(http.Header)
	header.Add("Content-Type", mw.FormDataContentType())
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQwMDoyNDowMi4yNzkzOTU5KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.mmcFJ0I9vXbHUPN4zdf0eo-eIz72-FW43RwyP5SY12Y")
	router.POST("/v1/bucket/:bucket_id/file", UploadFile)
	w := util.PerformRequest(router, http.MethodPost, "/v1/bucket/9/file", body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteFile(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.DELETE("/v1/bucket/:bucket_id/file/:file_id", DeleteFile)
	header := make(http.Header)
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQwMDoyNDowMi4yNzkzOTU5KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.mmcFJ0I9vXbHUPN4zdf0eo-eIz72-FW43RwyP5SY12Y")
	w := util.PerformRequest(router, http.MethodDelete, "/v1/bucket/9/file/1", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestSignFile(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.HEAD("/v1/bucket/:bucket_id/file/:file_id", SignFile)
	header := make(http.Header)
	header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQwMDo1NDoyOS4wMDcxMDE4KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.e3BuyBDOb-Pgj1mceXTxGChYDO6M9cy34TPFbIWKdoA")
	w := util.PerformRequest(router, http.MethodHead, "/v1/bucket/9/file/3", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDownloadFile(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	router := gin.New()
	router.Use(middleware.Token)
	router.GET("/v1/bucket/:bucket_id/file/:file_id", DownloadFile)
	uri := "/v1/bucket/9/file/3?sign=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ImV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhCcGNtVnpYMkYwSWpvaU1qQXlNUzB3TXkweE5GUXdNVG94TlRveE9TNHhOak13T1RZNEt6QTRPakF3SWl3aWRHOXJaVzRpT25zaVpHOXRZV2x1SWpvaWRHVnpkQ0lzSW5WelpYSWlPaUl4TWpNaUxDSmhZM1JwYjI1ZmJXRndJanA3SWs5Q1V5STZNSDE5ZlEuYjBBSGw5cEZfQlJYSHBOakpyR0VTbFIxQU5SY3RLZDRfcXpSSWRIZ2xBYyI.ZTeUBgJyQUOhY_ragehji8j2E7I19GckU3u0bNUEAvI"
	header := make(http.Header)
	//header.Add(middleware.XAuthToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQwMDo1NDoyOS4wMDcxMDE4KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.e3BuyBDOb-Pgj1mceXTxGChYDO6M9cy34TPFbIWKdoA")
	w := util.PerformRequest(router, http.MethodGet, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
