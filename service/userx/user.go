// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package userx

import (
	"context"
	"errors"
	"time"

	"github.com/json-iterator/go"
	"gorm.io/gorm"

	"obs/e"
	"obs/logger"
	"obs/model/db"
	"obs/service/tokenx"
)

// UserLogin 登录生成token信息
func UserLogin(ctx context.Context, email, passWord string) (string, error) {
	domain := &db.Domain{}
	if err := db.NewDB().Model(domain).Where("email =?", email).Find(domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", e.New(e.NotFound, "not found record")
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", e.New(e.OperateDbFail, err.Error())
	}
	if domain.PassWord != passWord {
		return "", e.New(e.Forbidden, "wrong password")
	}
	var permission map[string]tokenx.Action
	if err := jsoniter.ConfigFastest.UnmarshalFromString(domain.Permission, &permission); err != nil {
		logger.FromContext(ctx).Errorf("Unmarshal permission failed.Error:%v", err)
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
		logger.FromContext(ctx).Errorf("Create token failed.Error:%v", err)
		return "", e.New(e.GenerateTokenFail, err.Error())
	}
	return tokenStr, nil
}

// ModifyUser 修改用户信息
func ModifyUser(ctx context.Context, email, newPassWord, oldPassWord, nick string) error {
	tx := db.NewDB().Begin()
	defer tx.Rollback()
	domain := &db.Domain{}
	if err := tx.Model(domain).Where("email =?", email).Find(domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	if domain.PassWord != oldPassWord {
		return e.New(e.Forbidden, "wrong password")
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
		err = tx.Model(domain).Select(columnList[0]).UpdateColumns(domain).Error
	default:
		err = tx.Model(domain).Select(columnList[0], columnList[1:]...).UpdateColumns(domain).Error
	}
	if err != nil {
		logger.FromContext(ctx).Errorf("update db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	tx.Commit()
	return nil
}
