// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package uuid

import (
	"crypto/rand"
	"fmt"
)

// randInt generates a random uint64
func randUint64() uint64 {
	b := make([]byte, 8)
	if _, err := rand.Reader.Read(b); err != nil {
		panic(fmt.Sprintf("cann't generate random number: %v", err))
	}
	var id uint64
	id |= uint64(b[0])
	id <<= 8
	id |= uint64(b[1])
	id <<= 8
	id |= uint64(b[2])
	id <<= 8
	id |= uint64(b[3])
	id <<= 8
	id |= uint64(b[4])
	id <<= 8
	id |= uint64(b[5])
	id <<= 8
	id |= uint64(b[6])
	id <<= 8
	id |= uint64(b[7])
	return id
}
