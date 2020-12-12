// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/12

package config

import "net"

// Config 配置参数
type Config struct {
	IP         net.IP
	Pid        int
	YamlConfig *YamlConfig
}

type YamlConfig struct {
	ServiceInformation ServiceInformation `yaml:"service_information"`
	ConnectionList     Connection         `yaml:"connection_list"`
}

type ServiceInformation struct {
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	Version      string `yaml:"version"`
	LogPath      string `yaml:"log_path"`
	LogLevel     string `yaml:"log_level"`
	SaveRootPath string `yaml:"save_root_path"`
}

type Connection struct {
	LogConfig MongoConfig `yaml:"log_config"`
}

type MongoConfig struct {
	User     string `yaml:"user"`
	PassWord string `yaml:"pass_word"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}
