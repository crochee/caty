// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package tokenx

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type ObsClaims struct {
	Ip    string `json:"ip"`
	Token *Token `json:"token"`
}

func (o *ObsClaims) Valid() error {
	return nil
}

func (u *UserClaims) RefreshExpiresTime() error {
	client := credis.FirstRedisDB()
	_, err := client.Expire(context.Background(), tokenPrefix+u.Email, ExpiresTime).Result()
	return err
}

func (u *UserClaims) ValidThenRefreshExpiresTime(tokenString, ip string, ipEnable bool) error {
	if ipEnable && !strings.Contains(ip, u.Ip) { // ip做不做校验 只比较ip不比较端口
		return fmt.Errorf("ip is changed,create ip:%s,but ip:%s", u.Ip, ip)
	}
	// 验证redis存的值是否与之一致
	tokenStringValue, err := u.Get()
	if err != nil {
		return err
	}
	if tokenStringValue != tokenString {
		return errors.New("claim has changed")
	}
	// 刷新过期时间
	if err = u.RefreshExpiresTime(); err != nil {
		return err
	}
	return nil
}

var (
	tokenPrefix = "console:token:"
	secret      = []byte(`console secret`)
	ExpiresTime = 30 * time.Minute
)

func CreateTokenThenSave(claims *UserClaims) (string, error) {
	tokenString, err := createToken(claims)
	if err != nil {
		return "", err
	}
	// 存入redis
	if err = claims.Set(tokenString); err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseTokenThenRefresh(tokenString, ip string, ipEnable bool) (*ObsClaims, error) {
	uc, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if err = uc.ValidThenRefreshExpiresTime(tokenString, ip, ipEnable); err != nil {
		return nil, err
	}
	return uc, nil
}

func createToken(claims *ObsClaims) (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenImpl.SignedString(secret)
}

func parseToken(tokenString string) (*ObsClaims, error) {
	claims := new(ObsClaims)
	tokenImpl, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		return nil, err
	}
	var ok bool
	if claims, ok = tokenImpl.Claims.(*ObsClaims); !ok {
		return nil, errors.New("cannot convert claim")
	}
	//验证token，如果token被修改过则为false
	if !tokenImpl.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
