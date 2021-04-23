// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/4/23

// Package internal
package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Stack 根据skip输出堆栈
func Stack(skip int) []byte {
	var (
		buffer       = new(bytes.Buffer)
		lines        [][]byte
		lastFilePath string
	)

	for i := skip; ; i++ {
		pc, filePath, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buffer, "%s:%d (0x%x)\n", cleanPath(filePath), line, pc)
		if filePath != lastFilePath {
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFilePath = filePath
		}
		_, _ = fmt.Fprintf(buffer, "\t%s: %s\n", functionName(pc), fixSource(lines, line))
	}
	return buffer.Bytes()
}

func functionName(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	funcName := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(funcName, slash); lastSlash >= 0 {
		funcName = funcName[lastSlash+1:]
	}
	if period := bytes.Index(funcName, dot); period >= 0 {
		funcName = funcName[period+1:]
	}
	return bytes.Replace(funcName, centerDot, dot, -1)
}

func fixSource(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func cleanPath(filePath string) string {
	indexByte := strings.LastIndexByte(filePath, '/')
	if indexByte == -1 {
		return filePath
	}
	indexByte = strings.LastIndexByte(filePath[:indexByte], '/')
	if indexByte == -1 {
		return filePath
	}
	return filePath[indexByte+1:]
}
