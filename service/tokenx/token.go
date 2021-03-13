// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package tokenx

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// "github.com/crochee/uid"
	// uid.New().String()
	secret      = []byte(`2plsm96getq7o0bub6uifr4dov90vht5oo10`)
	ExpiresTime = 30 * time.Minute
)

type Token struct {
	Domain    string            `json:"domain"`
	User      string            `json:"user"`
	ActionMap map[string]Action `json:"action_map"` //service-action
}

// TokenClaims jwt.Claims的 Token 实现
type TokenClaims struct {
	Now   time.Time `json:"expires_at"`
	Token *Token    `json:"token"`
}

func (t *TokenClaims) Valid() error {
	if !t.Now.IsZero() && t.Now.Add(ExpiresTime).Before(time.Now()) {
		return errors.New("token is overdue")
	}
	return nil
}

// Action service's allow action
type Action uint8

const (
	Read   Action = 0
	Write  Action = 1
	Delete Action = 2
	Admin  Action = 3
)

// CreateToken 生成token
//
// @param claims jwt.Claims的token实现
// @Success string token加密信息
// @Failure error 标准错误
func CreateToken(claims *TokenClaims) (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenImpl.SignedString(secret)
}

// ParseToken 解析出token信息
//
// @param tokenString token的加密信息
// @Success TokenClaims jwt.Claims的token实现
// @Failure error 标准错误
func ParseToken(tokenString string) (*TokenClaims, error) {
	claims := new(TokenClaims)
	tokenImpl, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		return nil, err
	}
	var ok bool
	if claims, ok = tokenImpl.Claims.(*TokenClaims); !ok {
		return nil, errors.New("cannot convert token claim")
	}
	//验证token，如果token被修改过则为false
	if !tokenImpl.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

// Signature jwt.Claims的签名实现
type Signature string

func (s Signature) Valid() error {
	return nil
}

// CreateSign 生成签名加密信息
//
// @param claims jwt.Claims的签名实现
// @Success string 签名加密信息
// @Failure error 标准错误
func CreateSign(claims Signature) (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenImpl.SignedString(secret)
}

// ParseSign 解析出签名信息
//
// @param signString 签名的加密信息
// @Success Signature jwt.Claims的签名实现
// @Failure error 标准错误
func ParseSign(signString string) (Signature, error) {
	var claims Signature
	tokenImpl, err := jwt.ParseWithClaims(signString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		return "", err
	}
	var ok bool
	if claims, ok = tokenImpl.Claims.(Signature); !ok {
		return "", errors.New("cannot convert sign claim")
	}
	//验证token，如果token被修改过则为false
	if !tokenImpl.Valid {
		return "", errors.New("sign is invalid")
	}
	return claims, nil
}
