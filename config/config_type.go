// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/12

package config

import "net"

// Config 配置参数
type Config struct {
	IP            net.IP
	Pid           int
	ServiceConfig *ServiceConfig
}

type ServiceConfig struct {
	ServiceInfo ServiceInformation `json:"service_info" yaml:"service_info"`
	List        Connection         `json:"list" yaml:"list"`
}

type ServiceInformation struct {
	Mode        string `json:"mode" yaml:"mode"`
	LogPath     string `json:"log_path" yaml:"log_path"`
	LogLevel    string `json:"log_level" yaml:"log_level"`
	StoragePath string `json:"storage_path" yaml:"storage_path"`
}

type Connection struct {
	Mysql *SqlConfig   `json:"mysql,omitempty" yml:"mysql,omitempty"`
	Mongo *MongoConfig `json:"mongo,omitempty" yml:"mongo,omitempty"`
}

type SqlConfig struct {
	Type     string `json:"type" yaml:"type"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
	Charset  string `json:"charset" yaml:"charset"`
	Debug    bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
}

type MongoConfig struct {
	User     string `yaml:"user" yaml:"user"`
	Password string `yaml:"password" yaml:"password"`
	Host     string `yaml:"host" yaml:"host"`
	Port     int    `yaml:"port" yaml:"port"`
	Database string `yaml:"database" yaml:"database"`
}
