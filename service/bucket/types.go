// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import "time"

type Info struct {
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
	Name         string    `json:"name"`
	Count        int       `json:"count"`
}
