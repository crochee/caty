// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package auth

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	var tokenImpl = &TokenClaims{
		Now: time.Now().Unix(),
		Token: Token{
			AccountID: "123",
			UserID:    "test123",
			Permission: map[string]uint8{
				"cca": Admin,
			},
		},
	}
	value, err := CreateToken(tokenImpl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
}

func TestParseToken(t *testing.T) {
	//input := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJ0b2tlbiI6eyJkb21haW4iOiJ0ZXN0IiwidXNlciI6IjEyMyIsImFjdGlvbl9tYXAiOnsiT0JTIjowfX19.1bjRLUcynyaQTy5yz5lMnGCFbshDWG_9Nb1XLSe2hd4"
	//input := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxMTo0Nzo1Ny44OTY3NjU2KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MH19fQ.rm5UeTu6x6ANTGP8dKpnCAZiXO6X75ZoDsBf3x0ejHo"
	inputList := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxMjowNTowMC44NTcyNzU1KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MH19fQ.REDVGlCev7LrCAutU4SHqszq8xQCVf1_YV9HvQfHS4I",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxMzoyNzozNS44NzU3NTY5KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MH19fQ.Z8i4-jzpVJUWSw1_i67sqpQn9H4tmRnC98IUTvvFpqo",
		"eyahbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzX2F0IjoiMjAyMS0wMy0xM1QxMzoyNzozNS44NzU3NTY5KzA4OjAwIiwidG9rZW4iOnsiZG9tYWluIjoidGVzdCIsInVzZXIiOiIxMjMiLCJhY3Rpb25fbWFwIjp7Ik9CUyI6MH19fQ.Z8i4-jzpVJUWSw1_i67sqpQn9H4tmRnC98IUTvvFpqo",
	}
	for _, input := range inputList {
		token, err := ParseToken(input)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", token)
	}
}
