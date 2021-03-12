// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package tokenx

type Token struct {
	Domain    string            `json:"domain"`
	User      string            `json:"user"`
	ActionMap map[string]Action `json:"action_map"`
}

// Bucket permissions
type Action uint8

const (
	Read   Action = 0
	Write  Action = 1
	Delete Action = 2
	Admin  Action = 3
)
