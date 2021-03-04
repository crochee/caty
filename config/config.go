// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/2

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"obs/logger"
	"obs/util"
)

// Cfg 全局配置参数
var Cfg Config

// InitConfig init Config
func InitConfig(path string) {
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}

	gin.SetMode(config.ServiceInfo.Mode)
	logger.InitLogger(config.ServiceInfo.LogPath, config.ServiceInfo.LogLevel)

	logger.Debugf("config:%+v", config)
	Cfg.ServiceConfig = config

	Cfg.Pid = os.Getpid()                            // pid获取
	if Cfg.IP, err = util.ExternalIP(); err != nil { // ip获取
		logger.Fatal(err.Error())
	}
}

type DecodeEncode interface {
	Decode() (*ServiceConfig, error)
	Encode(config *ServiceConfig) error
}

func loadConfig(path string) (*ServiceConfig, error) {
	var lc DecodeEncode
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".json":
		lc = Json{path: path}
	case ".yml", ".yaml":
		lc = Yml{path: path}
	default:
		return nil, fmt.Errorf("unsupport config extension %s", ext)
	}
	return lc.Decode()
}
