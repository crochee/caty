// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package config

import "testing"

func TestInitConfig(t *testing.T) {
	cf := &ServiceConfig{
		ServiceInfo: ServiceInformation{
			Port:        8150,
			Mode:        "debug",
			LogPath:     "./log/obs.log",
			LogLevel:    "debug",
			StoragePath: "/obs/",
		},
		List: Connection{
			Mysql: &SqlConfig{
				User:     "root",
				Password: "123456",
				Host:     "192.168.31.62",
				Port:     3306,
				Database: "obs",
				Charset:  "utf8",
			},
			Mongo: &MongoConfig{
				User:     "root",
				PassWord: "123456",
				Host:     "192.168.31.62",
				Port:     27017,
				Database: "log",
			},
		},
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
