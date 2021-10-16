// Date: 2021/10/16

// Package account
package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/crochee/lib/db"
	"github.com/crochee/lib/e"

	"cca/pkg/code"
	"cca/pkg/dbx"
	"cca/pkg/model"
	"cca/pkg/service/auth"
)

type LoginRequest struct {
	// 用户账号
	// Required: true
	UserID string `json:"user_id" binding:"required,numeric"`
	// 密码
	// Required: true
	Password string `json:"password" binding:"required,alphanum"`
}

// Login 用户登录
func Login(ctx context.Context, request *LoginRequest) (*auth.APIToken, error) {
	user := &model.User{}
	if err := dbx.With(ctx).Model(user).Where("id =?",
		request.UserID).First(user).Error; err != nil {
		if errors.Is(err, db.NotFound) {
			return nil, code.ErrNoAccount
		}
		return nil, fmt.Errorf("first user failed;%v;%w", err, code.ErrLoginAccount)
	}
	if user.Password != request.Password {
		return nil, code.ErrWrongPasswordAccount
	}
	token := &auth.TokenClaims{
		Now: time.Now().Unix(),
		Token: &auth.Token{
			AccountID:  FormatUint(user.AccountID),
			UserID:     FormatUint(user.ID),
			Permission: make(map[string]uint8),
		},
	}
	if err := json.Unmarshal([]byte(user.Permission), &token.Token.Permission); err != nil {
		return nil, fmt.Errorf("unmarshal failed;%v;%w", err, e.ErrInternalServerError)
	}
	return auth.Create(ctx, token)
}
