// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package verify

import (
	"github.com/crochee/uid"
	"github.com/dgrijalva/jwt-go"
)

type AkSk interface {
	Create() (string, string, error)
	Verify() error
}

type Token struct {
	Bucket string `json:"bucket"`
	jwt.StandardClaims
}

func (t Token) Create() (string, string, error) {
	akSecret := uid.New().String()
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	skToken, err := tokenImpl.SignedString(akSecret)
	if err != nil {
		return "", "", err
	}
	return akSecret, skToken, nil
}

func (t Token) Verify() error {
	panic("implement me")
}
