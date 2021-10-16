// Date: 2021/10/16

// Package auth
package auth

import (
	"context"
	"time"
)

type APIToken struct {
	// 加密token信息
	// Required: true
	Token string `json:"token" binding:"required"`
}

// Create 生成token
func Create(_ context.Context, token *TokenClaims) (*APIToken, error) {
	if token.Now == 0 {
		token.Now = time.Now().Unix()
	}
	permission, err := token.Create()
	if err != nil {
		return nil, err
	}
	return &APIToken{Token: permission}, nil
}

// Parse 解析token
func Parse(_ context.Context, token string) (*TokenClaims, error) {
	tokenImpl := &TokenClaims{}
	if err := tokenImpl.Parse(token); err != nil {
		return nil, err
	}
	return tokenImpl, nil
}
