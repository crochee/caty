// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/2

package config

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"obs/util"
)

// Config 配置参数
type Config struct {
	ConfigPath     string
	IP             net.IP
	Pid            int
	ServiceInfo    ServiceInformation
	ConnectionInfo Connection
}

type ServiceInformation struct {
	Port         int
	Mode         string
	Log          string
	Version      string
	SaveRootPath string
}

type Connection struct {
	Log MongoConfig
}

type MongoConfig struct {
	User     string
	PassWord string
	Host     string
	Port     int
	Database string
}

// Cfg 全局配置参数
var Cfg Config

func InitConfig() {
	var ok bool
	if Cfg.ConfigPath, ok = os.LookupEnv("config_path"); !ok {
		Cfg.ConfigPath = "conf/config.yml"
	}
	Cfg.Pid = os.Getpid() // pid获取
	var err error
	if Cfg.IP, err = util.ExternalIP(); err != nil { // ip获取
		log.Fatal(err)
	}
	// 配置文件加载
	if err = loadConfigFile(); err != nil {
		log.Fatal(err)
	}
	Cfg.ServiceInfo = ServiceInformation{
		Port:         viper.GetInt("service_information.port"),
		Mode:         viper.GetString("service_information.mode"),
		Log:          viper.GetString("service_information.log"),
		Version:      viper.GetString("service_information.version"),
		SaveRootPath: viper.GetString("service_information.save_root_path"),
	}
	Cfg.ConnectionInfo = Connection{Log: MongoConfig{
		User:     viper.GetString("connection.mongo.log.user"),
		PassWord: viper.GetString("connection.mongo.log.pass"),
		Host:     viper.GetString("connection.mongo.log.host"),
		Port:     viper.GetInt("connection.mongo.log.port"),
		Database: viper.GetString("connection.mongo.log.database"),
	}}
}

// LoadConfig 文件配置加载
func loadConfigFile() error {
	configPath, err := filepath.Abs(Cfg.ConfigPath)
	if err != nil {
		return err
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
