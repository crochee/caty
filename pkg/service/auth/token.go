package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"caty/pkg/code"
)

type Authentication interface {
	Create() (string, error)
	Parse(data string) error
	Secret(token *jwt.Token) (interface{}, error)
}

var (
	ExpiresTime = 30 * time.Minute

	AllService = "*"
)

type Token struct {
	// 主账号id
	// Required: true
	AccountID string `json:"account_id" binding:"required,numeric"`
	// 账户id
	// Required: true
	UserID string `json:"user_id" binding:"required,numeric"`
	// 权限列表
	// Required: true
	Permission map[string]uint8 `json:"permission" binding:"required"`
}

// TokenClaims jwt.Claims的 Token 实现
type TokenClaims struct {
	// 生成token的时间戳
	Now int64 `json:"now"`
	// token信息
	Token *Token `json:"token" binding:"required,dive"`
}

func (t *TokenClaims) Valid() error {
	if t.Now != 0 && time.Now().Add(-ExpiresTime).Unix() > t.Now {
		return code.ErrExpireAuth
	}
	return nil
}

func (t *TokenClaims) Create() (string, error) {
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	secretKey, err := t.Secret(nil)
	if err != nil {
		return "", fmt.Errorf("create secret failed;%v;%w", err, code.ErrCreateAuth)
	}
	var token string
	if token, err = tokenImpl.SignedString(secretKey); err != nil {
		return "", fmt.Errorf("signedString failed;%v;%w", err, code.ErrCreateAuth)
	}
	return token, nil
}

func (t *TokenClaims) Parse(data string) error {
	tokenImpl, err := jwt.ParseWithClaims(data, t, t.Secret)
	if err != nil {
		return fmt.Errorf("parse token failed;%v;%w", err, code.ErrParseAuth)
	}
	claims, ok := tokenImpl.Claims.(*TokenClaims)
	if !ok {
		return fmt.Errorf("cannot convert token claim;%w", code.ErrInvalidAuth)
	}
	// 验证token，如果token被修改过则为false
	if !tokenImpl.Valid {
		return code.ErrInvalidAuth
	}
	*t = *claims
	return nil
}

func (t *TokenClaims) Secret(_ *jwt.Token) (interface{}, error) {
	h := sha256.New()
	_, err := h.Write([]byte(t.Token.UserID))
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
