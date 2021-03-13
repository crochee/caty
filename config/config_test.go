// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package config

import "testing"

func TestInitConfig(t *testing.T) {
	cf := &ServiceConfig{
		ServiceInfo: ServiceInformation{
			Mode:        "debug",
			LogPath:     "./log/obs.log",
			LogLevel:    "debug",
			StoragePath: "/obs/",
		},
		List: Connection{
			Mysql: &SqlConfig{
				Type:     "mysql",
				User:     "root",
				Password: "1234567",
				Host:     "192.168.31.62",
				Port:     3307,
				Database: "obs",
				Charset:  "utf8",
				Debug:    true,
			},
			Mongo: &MongoConfig{
				User:     "root",
				Password: "123456",
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
