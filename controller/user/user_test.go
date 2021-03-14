// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"bytes"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"obs/config"
	"obs/model/db"
	"obs/util"
	"testing"
)

func TestRegister(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	body := new(bytes.Buffer)
	r := &Domain{
		Nick: "nick",
		LoginInfo: LoginInfo{
			Email:    "13522570308@139.com",
			PassWord: "123",
		},
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/user/register", Register)
	w := util.PerformRequest(router, http.MethodPost, "/v1/user/register", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestLogin(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	body := new(bytes.Buffer)
	r := &LoginInfo{
		Email:    "13522570308@139.com",
		PassWord: "123",
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/user/login", Login)
	w := util.PerformRequest(router, http.MethodPost, "/v1/user/login", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestModify(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	db.Setup()
	body := new(bytes.Buffer)
	r := &ModifyInfo{
		Email:       "13522570308@139.com",
		Nick:        "555",
		NewPassWord: "123",
		OldPassWord: "123",
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/user/modify", Modify)
	w := util.PerformRequest(router, http.MethodPost, "/v1/user/modify", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
