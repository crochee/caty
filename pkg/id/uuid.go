// Date: 2021/9/19

// Package id
package id

import "github.com/satori/go.uuid"

func Uuid() string {
	return uuid.NewV4().String()
}
