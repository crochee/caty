// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package uuid

import (
	"encoding/base32"
	"encoding/binary"
	"os"
	"sync/atomic"
	"time"
)

type UUID [19]byte

var (
	version  byte = 0x01 // set version to 1
	pid           = os.Getpid()
	clockSeq      = randUint64()
)

func New() UUID {
	return FromTime(time.Now())
}

// 8 bytes of time (ns) + 1 bytes of version + 2 byes of pid + 8 random bytes
func FromTime(aTime time.Time) UUID {
	var id UUID

	utcTime := aTime.In(time.UTC)
	// Timestamp ns, 8 bytes, big endian
	binary.BigEndian.PutUint64(id[:], uint64(utcTime.UnixNano()))
	// version, 1 bytes
	id[8] = version
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	id[9] = byte(pid >> 8)
	id[10] = byte(pid)
	// random, 8 bytes big endian
	clock := atomic.AddUint64(&clockSeq, 1)
	id[11] = byte(clock >> 56)
	id[12] = byte(clock >> 48)
	id[13] = byte(clock >> 40)
	id[14] = byte(clock >> 32)
	id[15] = byte(clock >> 24)
	id[16] = byte(clock >> 16)
	id[17] = byte(clock >> 8)
	id[18] = byte(clock)
	return id
}

var encode = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding)

func (u UUID) String() string {
	return encode.EncodeToString(u[:])
}
