// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/5

package balance

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrNoneAvailable = errors.New("none available")

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Random struct {
	mux  sync.RWMutex
	list []*Node
}

func (r *Random) Next() (*Node, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	i := rand.Int() % length
	return r.list[i], nil
}

type RoundRobin struct {
	randIndex int
	mux       sync.Mutex
	list      []*Node
}

func (r *RoundRobin) Next() (*Node, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	r.randIndex %= length
	node := r.list[r.randIndex]
	r.randIndex++
	return node, nil
}

type WeightNode struct {
	*Node
	currentWeight int //当前权重
}

type WeightRoundRobin struct {
	list []*WeightNode
}

func (w *WeightRoundRobin) Next() (*Node, error) {
	var best *WeightNode
	var total int
	for _, node := range w.list {
		node.currentWeight += node.Weight // 将当前权重与有效权重相加
		total += node.Weight              //累加总权重
		if best == nil || node.currentWeight > best.currentWeight {
			best = node
		}
	}
	if best == nil {
		return nil, ErrNoneAvailable
	}
	best.currentWeight -= total //将当前权重改为当前权重-总权重
	return best.Node, nil
}
