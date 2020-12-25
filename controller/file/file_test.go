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
	"obs/util"
)

func TestUploadFile(t *testing.T) {
	config.InitConfig()
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	_ = mw.WriteField("path", "data")
	f, _ := mw.CreateFormFile("file", "self.txt")
	_, _ = f.Write([]byte(`hello world`))
	_ = mw.Close()
	uri := fmt.Sprintf("/v1/file/bucket/%s", "cptsbuild")
	header := make(http.Header)
	header.Add("Content-Type", mw.FormDataContentType())
	router := gin.New()
	router.POST("/v1/file/bucket/:bucket_name", UploadFile)
	w := util.PerformRequest(router, http.MethodPost, uri, body, header)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestDeleteFile(t *testing.T) {
	config.InitConfig()
	uri := fmt.Sprintf("/v1/file/bucket/%s?path=%s",
		"cptsbuild",
		"data/self.txt",
	)
	router := gin.New()
	router.DELETE("/v1/file/bucket/:bucket_name", DeleteFile)
	w := util.PerformRequest(router, http.MethodDelete, uri, nil, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
