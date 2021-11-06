package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig init Config
func LoadConfig(path string) error {
	viper.AddConfigPath("/") // 设置配置文件路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	viper.SetConfigName(absPath)                // 设置配置文件名称
	viper.SetConfigType(filepath.Ext(path)[1:]) // 设置配置文件类型

	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed.%w", err)
	}
	return nil
}
