// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"bytes"
	"context"
	"net/http"
	"obs/pkg/model/db"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"

	"obs/config"
	"obs/internal"
)

func TestRegister(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	body := new(bytes.Buffer)
	r := &User{
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
	w := internal.PerformRequest(router, http.MethodPost, "/v1/user/register", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestLogin(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
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
	w := internal.PerformRequest(router, http.MethodPost, "/v1/user/login", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

func TestModify(t *testing.T) {
	config.LoadConfig("../../conf/obs.yml")
	if err := db.Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
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
	w := internal.PerformRequest(router, http.MethodPost, "/v1/user/modify", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}
