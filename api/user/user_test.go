// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"

	"cca/config"
	"obs/internal"
	"obs/pkg/db"
)

func TestRegister(t *testing.T) {
	if err := config.LoadConfig("E:\\project\\cca\\conf\\cca.yml"); err != nil {
		t.Fatal(err)
	}
	if err := db.Init(context.Background()); err != nil {
		t.Fatal(err)
	}
	body := new(bytes.Buffer)
	r := &RegisterUserRequest{
		Account:  "crochee",
		Email:    "crochee@139.com",
		PassWord: "123456",
		Desc:     `{"detail":"some unknown"}`,
	}
	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/v1/account", Register)
	w := internal.PerformRequest(router, http.MethodPost, "/v1/account", body, nil)
	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
}

//func TestLogin(t *testing.T) {
//	config.LoadConfig("../../conf/cca.yml")
//	if err := db.Setup(context.Background()); err != nil {
//		t.Fatal(err)
//	}
//	body := new(bytes.Buffer)
//	r := &LoginInfo{
//		Email:    "13522570308@139.com",
//		Password: "123",
//	}
//	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
//		t.Fatal(err)
//	}
//	router := gin.With()
//	router.POST("/v1/user/login", Login)
//	w := internal.PerformRequest(router, http.MethodPost, "/v1/user/login", body, nil)
//	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
//}
//
//func TestModify(t *testing.T) {
//	config.LoadConfig("../../conf/cca.yml")
//	if err := db.Setup(context.Background()); err != nil {
//		t.Fatal(err)
//	}
//	body := new(bytes.Buffer)
//	r := &ModifyInfo{
//		Email:       "13522570308@139.com",
//		Nick:        "555",
//		NewPassWord: "123",
//		OldPassWord: "123",
//	}
//	if err := jsoniter.ConfigFastest.NewEncoder(body).Encode(r); err != nil {
//		t.Fatal(err)
//	}
//	router := gin.With()
//	router.POST("/v1/user/modify", Modify)
//	w := internal.PerformRequest(router, http.MethodPost, "/v1/user/modify", body, nil)
//	t.Logf("%+v modify:%+v body:%s", w.Result(), w.Header(), w.Body.String())
//}
