// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/25

package file

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"obs/config"
	"obs/middleware"
	"obs/util"
)

func TestUploadFile(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	_ = mw.WriteField("path", "data")
	f, _ := mw.CreateFormFile("file", "test.txt")
	_, _ = f.Write([]byte(`hello world`))
	_ = mw.Close()
	uri := fmt.Sprintf("/v1/file/%s", "cptsbuild")
	router := gin.New()
	router.Use(middleware.Verify)
	header := make(http.Header)
	header.Add("Content-Type", mw.FormDataContentType())
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	router.POST("/v1/file/:bucket_name", UploadFile)
	w := util.PerformRequest(router, http.MethodPost, uri, body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteFile(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/file/%s?path=%s",
		"cptsbuild",
		"data/self.txt",
	)
	router := gin.New()
	router.Use(middleware.Verify)
	router.DELETE("/v1/file/:bucket_name", DeleteFile)
	header := make(http.Header)
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	w := util.PerformRequest(router, http.MethodDelete, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestSignFile(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/file/%s?path=%s",
		"cptsbuild",
		"data/test.txt",
	)
	router := gin.New()
	router.Use(middleware.Verify)
	router.HEAD("/v1/file/:bucket_name", SignFile)
	header := make(http.Header)
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	w := util.PerformRequest(router, http.MethodHead, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDownloadFile(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/file/%s?path=%s",
		"cptsbuild",
		"data/test.txt",
	)
	router := gin.New()
	router.Use(middleware.Verify)
	router.GET("/v1/file/:bucket_name", DownloadFile)
	header := make(http.Header)
	header.Add("ak", "2pa4kh0996gk008uqpj2nq3hj9vs7lablep0")
	header.Add("sk", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
		"eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJzZWNyZXQiOiJNbkJoTkd0b01EazVObWRyTURBNGRYRndhakp1Y1ROb2FqbDJjemRzWVdKc1pYQXciLCJhY3Rpb24iOnsiMyI6e319fQ.5W25ws5EGn6nle7ts62V3RgNSuMrMNIIyrP7PBzHjJM")
	w := util.PerformRequest(router, http.MethodGet, uri, nil, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
