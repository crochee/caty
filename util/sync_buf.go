// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/24

package util

import "sync"

var bufPool sync.Pool

const DefaultSize = 512

func AcquireBuf() []byte {
	buf := bufPool.Get()
	if buf == nil {
		return make([]byte, DefaultSize)
	}
	return buf.([]byte)
}

func ReleaseBuf(buf []byte) {
	buf = buf[:0]
	bufPool.Put(buf)
}
