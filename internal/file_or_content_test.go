// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/4/28

// Package internal
package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileOrContent_Read(t *testing.T) {
	list := []struct {
		name string
		f    FileOrContent
		want string
	}{
		{
			name: "content",
			f:    "ileOrContent1.go",
			want: "ileOrContent1.go",
		},
		{
			name: "path",
			f:    "../test/test1.txt",
			want: "test",
		},
	}
	for _, tt := range list {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Read()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}
