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

	AllService = "*"
)

type Token struct {
	AccountID string            `json:"account_id"`
	UserID    string            `json:"user_id"`
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
	Not    Action = 0
	Read   Action = 1
	Write  Action = 2
	Delete Action = 3
	Admin  Action = 4
)

var ActionString = map[Action]string{
	Not:    "not",
	Read:   "read",
	Write:  "write",
	Delete: "delete",
	Admin:  "admin",
}

// CreateToken 生成token
func CreateToken(claims *TokenClaims) (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenImpl.SignedString(secret)
}

// ParseToken 解析出token信息 TokenClaims
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
type Signature struct {
	Sign string
}

func (s *Signature) Valid() error {
	return nil
}

// CreateSign 生成签名加密信息
func CreateSign(claims *Signature) (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenImpl.SignedString(secret)
}

// ParseSign 解析出签名信息 Signature
func ParseSign(signString string) (*Signature, error) {
	claims := new(Signature)
	tokenImpl, err := jwt.ParseWithClaims(signString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		return nil, err
	}
	var ok bool
	if claims, ok = tokenImpl.Claims.(*Signature); !ok {
		return nil, errors.New("cannot convert sign claim")
	}
	//验证token，如果token被修改过则为false
	if !tokenImpl.Valid {
		return nil, errors.New("sign is invalid")
	}
	return claims, nil
}
