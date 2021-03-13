// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package util

import (
	"os"
	"path/filepath"
)

// DirSize 统计文件夹的大小
func DirSize(path string) (int64, int, error) {
	var (
		size  int64
		count int
	)
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return err
	})
	return size, count, err
}
