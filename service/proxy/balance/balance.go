// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/5

package balance

// Node consider  registry.ServiceInstance
type Node struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints string            `json:"endpoint"`
	Weight    float64           `json:"weight"`
}

// Selector strategy algorithm
type Balancer interface {
	Next() (*Node, error)
}
