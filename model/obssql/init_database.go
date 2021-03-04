// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package obssql

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"obs/config"
	"obs/logger"
)

func RegisterModel(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		if db.HasTable(table) {
			db.AutoMigrate(table)
		} else {
			db.Set("gorm:table_options",
				"ENGINE=InnoDB DEFAULT CHARSET=utf8").
				CreateTable(table) //不存在就创建新表
		}
	}
}

func CreatePool(cf *config.SqlConfig, maxOpen, maxIdle int) (*gorm.DB, error) {
	var source string
	switch cf.Type {
	case "mysql":
		source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			cf.User,
			cf.Password,
			cf.Host,
			cf.Port,
			cf.Database,
			cf.Charset)
	case "postgres":
		source = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			cf.Host,
			cf.Port,
			cf.User,
			cf.Database,
			cf.Password)
	default:
		source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			cf.User,
			cf.Password,
			cf.Host,
			cf.Port,
			cf.Database,
			cf.Charset)
	}
	conn, err := gorm.Open(cf.Type, source)
	if err != nil {
		return nil, err
	}
	// 禁止表名复数形式
	conn.SingularTable(true)
	// 开启池化之后不能close  否则连接池没有作用
	// 设置数据库连接池最大连接数maxOpenConns
	conn.DB().SetMaxOpenConns(maxOpen)
	// 设置数据库连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于maxIdleConns，超过的连接会被连接池关闭
	conn.DB().SetMaxIdleConns(maxIdle)
	// 是否开启日志模式
	conn.LogMode(cf.Debug).SetLogger(MysqlLogger{})
	return conn, nil
}

type MysqlLogger struct {
}

func (m MysqlLogger) Print(v ...interface{}) {
	values := gorm.LogFormatter(v...)
	logger.Infof(strings.Repeat("%v ", len(values)), values...)
}
