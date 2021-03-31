// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package cmd

import "strconv"

const Version = "v1.0.1"

func getNextChunk(version string, n, p int) (int, int, error) {
	if p > n-1 {
		return 0, p, nil
	}
	pEnd := p
	for pEnd < n && version[pEnd] != '.' {
		pEnd++
	}
	var (
		i   int
		err error
	)
	if pEnd != n-1 {
		i, err = strconv.Atoi(version[p:pEnd])
		if err != nil {
			return 0, 0, err
		}
	} else {
		i, err = strconv.Atoi(version[p:n])
		if err != nil {
			return 0, 0, err
		}
	}
	p = pEnd + 1
	return i, p, nil
}

func CompareVersion(version1, version2 string) (int, error) {
	var (
		p1, p2 int
		i1, i2 int
		err    error
	)
	n1, n2 := len(version1), len(version2)
	for p1 < n1 || p2 < n2 {
		if i1, p1, err = getNextChunk(version1, n1, p1); err != nil {
			return 0, err
		}
		if i2, p2, err = getNextChunk(version2, n2, p2); err != nil {
			return 0, err
		}
		if i1 != i2 {
			if i1 > i2 {
				return 1, nil
			}
			return -1, nil
		}
	}
	return 0, nil
}
