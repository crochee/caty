// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package userx

import (
	"context"
	"errors"
	"obs/pkg/logx"
	db2 "obs/pkg/model/db"
	"time"

	"github.com/json-iterator/go"
	"gorm.io/gorm"

	"obs/e"
	"obs/service/business/tokenx"
)

// UserLogin 登录生成token信息
func UserLogin(ctx context.Context, email, passWord string) (string, error) {
	domain := &db2.Domain{}
	if err := db2.NewDBWithContext(ctx).Model(domain).Where("email =?", email).First(domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", e.New(e.NotFound, "not found record")
		}
		logx.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", e.New(e.OperateDbFail, err.Error())
	}
	if domain.PassWord != passWord {
		return "", e.New(e.Forbidden, "wrong password")
	}
	var permission map[string]tokenx.Action
	if err := jsoniter.ConfigFastest.UnmarshalFromString(domain.Permission, &permission); err != nil {
		logx.FromContext(ctx).Errorf("Unmarshal permission failed.Error:%v", err)
		return "", e.New(e.UnmarshalFail, err.Error())
	}
	token := &tokenx.TokenClaims{
		Now: time.Now(),
		Token: &tokenx.Token{
			Domain:    domain.Domain,
			User:      domain.Domain,
			ActionMap: permission,
		},
	}
	tokenStr, err := tokenx.CreateToken(token)
	if err != nil {
		logx.FromContext(ctx).Errorf("Create token failed.Error:%v", err)
		return "", e.New(e.GenerateTokenFail, err.Error())
	}
	return tokenStr, nil
}

// ModifyUser 修改用户信息
func ModifyUser(ctx context.Context, email, newPassWord, oldPassWord, nick string) error {
	tx := db2.NewDBWithContext(ctx).Begin()
	defer tx.Rollback()
	domain := &db2.Domain{}
	if err := tx.Model(domain).Where("email =?", email).First(domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		logx.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	if domain.PassWord != oldPassWord {
		return e.New(e.Forbidden, "wrong password")
	}
	var columnMap = make(map[string]interface{}, 2)
	if newPassWord != "" && domain.PassWord != newPassWord {
		columnMap["pass_word"] = newPassWord
	}
	if nick != "" && domain.Nick != nick {
		columnMap["nick"] = nick
	}
	if len(columnMap) == 0 {
		return nil
	}
	if err := tx.Model(domain).UpdateColumns(columnMap).Error; err != nil {
		logx.FromContext(ctx).Errorf("update db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	tx.Commit()
	return nil
}
