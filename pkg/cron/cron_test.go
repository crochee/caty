// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/28

package cron

import (
	"fmt"
	"obs/pkg/log"
	"testing"
	"time"

	"obs/config"
)

func TestCronSetup(t *testing.T) {
	config.LoadConfig("../conf/obs.yml")
	log.InitSystemLogger()
	Setup()
	// 0 0/5 * * * ?
	id, err := timeCron.AddFunc("*/20 */30 * * * *", func() {
		fmt.Println("3min run...")
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
	time.Sleep(10 * time.Minute)
}
