// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/20

package uuid

import (
	"hash/crc32"
	"io/ioutil"
	"os"
)

func Pid() int {
	pid := os.Getpid()
	b, err := ioutil.ReadFile("/proc/self/cpuset")
	if err == nil && len(b) > 1 {
		pid ^= int(crc32.ChecksumIEEE(b))
	}
	return pid
}
