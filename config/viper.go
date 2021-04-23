// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/4/23

// Package config
package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ViperInfo struct {
	path     string
	fileName string
	Ext      string
}

func (v ViperInfo) Decode() (*ServiceConfig, error) {
	viper.AddConfigPath(v.path)
	index := strings.LastIndexByte(v.path, '.')
	if index == -1 {
		return nil, errors.New("path hasn't ext")
	}
	if !strings.HasSuffix(filepath.Ext(v.path), v.path[index:]) {
		return nil, errors.New("path hasn't ext")
	}

	viper.AddConfigPath("./")
	viper.SetConfigType(strings.TrimPrefix(v.path[index:], "."))
	viper.SetConfigName(v.path[:index])
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config ServiceConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&config); err != nil {
			log.Println(in.String(), err)
			return
		}
	})
	viper.WatchConfig()
	return &config, nil
}

func (v ViperInfo) Encode(config *ServiceConfig) error {
	return viper.SafeWriteConfigAs(v.path)
}
