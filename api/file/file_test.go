// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/25

package file

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"obs/pkg/middleware"
	"obs/pkg/model/db"
	"testing"

	"github.com/gin-gonic/gin"

	"obs/config"
	"obs/internal"
)

func TestUploadFile(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	f, _ := mw.CreateFormFile("file", "lcf.txt")
	_, _ = f.Write([]byte(`hello world`))
	_ = mw.Close()
	router := gin.New()
	router.Use(middleware.Token)
	header := make(http.Header)
	header.Add("Content-Type", mw.FormDataContentType())
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	router.POST("/v1/bucket/:bucket_name/file", UploadFile)
	w := internal.PerformRequest(router, http.MethodPost, "/v1/bucket/test/file", body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteFile(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.Token)
	router.DELETE("/v1/bucket/:bucket_name/file/:file_name", DeleteFile)
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := internal.PerformRequest(router, http.MethodDelete, "/v1/bucket/test/file/lcf.txt", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestSignFile(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.Token)
	router.GET("/v1/bucket/:bucket_name/file/:file_name", SignFile)
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := internal.PerformRequest(router, http.MethodGet, "/v1/bucket/test/file/lcf.txt", nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDownloadFile(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.Token)
	router.GET("/v1/bucket/:bucket_name/file/:file_name", DownloadFile)
	uri := "/v1/bucket/test/file/lcf.txt"
	header := make(http.Header)
	header.Add("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xNFQxNToxMToyMy42NDA3MDk1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoiMTIzIiwidXNlciI6InRlc3QxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6M319fQ.ZOX-KpOVeDhOV9qN4SWw5DWPDsl4LY1NrrXHv1yqNSU")
	w := internal.PerformRequest(router, http.MethodGet, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
