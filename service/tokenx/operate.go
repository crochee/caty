// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package tokenx

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// QueryToken 查询 Token
func QueryToken(ctx *gin.Context) (*Token, error) {
	token, ok := ctx.Get("token")
	if !ok {
		return nil, errors.New("token isn't exists")
	}
	var xToken *Token
	if xToken, ok = token.(*Token); !ok {
		return nil, errors.New("token's type isn't Token")
	}
	return xToken, nil
}
