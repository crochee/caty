// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/5

package balance

import "testing"

func TestWeightRoundRobin_Next(t *testing.T) {
	wr := &WeightRoundRobin{
		list: []*WeightNode{
			{
				Node: &Node{
					Endpoints: "1",
					Weight:    3,
				},
			},
			{
				Node: &Node{
					Endpoints: "2",
					Weight:    1,
				},
			},
			{
				Node: &Node{
					Endpoints: "3",
					Weight:    1,
				},
			},
		},
	}
	for i := 0; i < 12; i++ {
		node, err := wr.Next()
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", node)
	}
}

func BenchmarkWeightRoundRobin_Next(b *testing.B) {
	wr := &WeightRoundRobin{
		list: []*WeightNode{
			{
				Node: &Node{
					Endpoints: "1",
					Weight:    3,
				},
			},
			{
				Node: &Node{
					Endpoints: "2",
					Weight:    1,
				},
			},
			{
				Node: &Node{
					Endpoints: "3",
					Weight:    1,
				},
			},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wr.Next()
	}
}
