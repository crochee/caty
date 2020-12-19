// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package uuid

import "testing"

func TestTimeUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			t.Log(New().String())
		}()
	}
}
