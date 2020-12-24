// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package verify

import (
	"errors"
	"obs/util"

	"github.com/crochee/uid"
	"github.com/dgrijalva/jwt-go"
)

type AkSk interface {
	Create() (string, string, error)
	Verify(string) error
}

// Bucket permissions
type BucketAction uint8

const (
	Read   BucketAction = 0
	Write  BucketAction = 1
	Delete BucketAction = 2
	Admin  BucketAction = 3
)

type Token struct {
	AkSecret []byte                    `json:"secret"`
	Bucket   string                    `json:"bucket"`
	Action   map[BucketAction]struct{} `json:"action"`
}

func NewToken(bucket string) *Token {
	return &Token{
		Bucket: bucket,
		Action: make(map[BucketAction]struct{}),
	}
}

func (t *Token) Valid() error {
	return nil
}

func (t *Token) AddAction(action BucketAction) {
	t.Action[action] = struct{}{}
}

func (t *Token) Create() (string, string, error) {
	secret := uid.New().String()
	t.AkSecret = util.Slice(secret)
	tokenImpl := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	skToken, err := tokenImpl.SignedString(t.AkSecret)
	if err != nil {
		return "", "", err
	}
	return secret, skToken, nil
}

func (t *Token) Verify(skToken string) error {
	tokenImpl, err := jwt.ParseWithClaims(skToken, t, func(token *jwt.Token) (interface{}, error) {
		return t.AkSecret, nil
	})
	if err != nil {
		return err
	}
	if !tokenImpl.Valid {
		return errors.New("token is invalid")
	}
	thisToken, ok := tokenImpl.Claims.(*Token)
	if !ok {
		return errors.New("claim is not Token")
	}
	if thisToken.Bucket != t.Bucket {
		return errors.New("bucket is not right")
	}
	t.Action = thisToken.Action
	return nil
}

func VerifyAuthentication(token *Token, action BucketAction) bool {
	for action := range token.Action {
		if action >= action {
			return true
		}
	}
	return false
}
