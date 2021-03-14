// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package userx

import (
	"context"
	"net/http"
	"time"

	"github.com/json-iterator/go"

	"obs/logger"
	"obs/model/db"
	"obs/response"
	"obs/service/tokenx"
)

// UserLogin 登录生成token信息
func UserLogin(ctx context.Context, email, passWord string) (string, error) {
	domain := &db.Domain{}
	if err := db.NewDB().Model(domain).Where("email =?", email).Find(domain).Error; err != nil {
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	if domain.PassWord != passWord {
		return "", response.Error(http.StatusForbidden, "wrong password")
	}
	var permission map[string]tokenx.Action
	if err := jsoniter.ConfigFastest.UnmarshalFromString(domain.Permission, &permission); err != nil {
		logger.FromContext(ctx).Errorf("Unmarshal permission failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
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
		logger.FromContext(ctx).Errorf("Create token failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	return tokenStr, nil
}

// ModifyUser 修改用户信息
func ModifyUser(ctx context.Context, email, newPassWord, oldPassWord, nick string) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	domain := &db.Domain{}
	if err := tx.Model(domain).Where("email =?", email).Find(domain).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	if domain.PassWord != oldPassWord {
		tx.Rollback()
		return response.Error(http.StatusForbidden, "wrong password")
	}
	var columnList = make([]interface{}, 0, 2)
	if newPassWord != "" && domain.PassWord != newPassWord {
		domain.PassWord = newPassWord
		columnList = append(columnList, "pass_word")
	}
	if nick != "" && domain.Nick != nick {
		domain.Nick = nick
		columnList = append(columnList, "nick")
	}
	var err error
	switch len(columnList) {
	case 0:
		return nil
	case 1:
		err = tx.Model(domain).Select(columnList[0]).Update(domain).Error
	default:
		err = tx.Model(domain).Select(columnList[0], columnList[1:]...).Update(domain).Error
	}
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("update db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	return nil
}
