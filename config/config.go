// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/2

package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"obs/logger"
	"obs/util"
)

// Cfg 全局配置参数
var Cfg Config

func InitConfig() {
	configYaml, err := loadYaml()
	if err != nil {
		panic(err)
	}
	gin.SetMode(configYaml.ServiceInformation.Mode)
	logger.InitLogger(configYaml.ServiceInformation.LogPath, configYaml.ServiceInformation.LogLevel)

	logger.Debugf("config:%+v", configYaml)
	Cfg.YamlConfig = configYaml

	Cfg.Pid = os.Getpid()                            // pid获取
	if Cfg.IP, err = util.ExternalIP(); err != nil { // ip获取
		logger.Fatal(err.Error())
	}
}

func loadYaml() (*YamlConfig, error) {
	configPath, ok := os.LookupEnv("config_path")
	if !ok {
		configPath = "conf/config.yml"
	}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config YamlConfig
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
