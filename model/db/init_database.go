// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"obs/config"
	"obs/logger"
)

var db *gorm.DB

func Setup() {
	var err error
	if db, err = createPool(config.Cfg.ServiceConfig.List.Mysql); err != nil {
		logger.Fatal(err.Error())
		return
	}
	// 开启池化之后不能close  否则连接池没有作用
	// 设置数据库连接池最大连接数maxOpenConns
	// 设置数据库连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于maxIdleConns，超过的连接会被连接池关闭
	db.DB().SetMaxIdleConns(10)                   //最大空闲连接数
	db.DB().SetMaxOpenConns(30)                   //最大连接数
	db.DB().SetConnMaxLifetime(time.Second * 300) //设置连接空闲超时

	registerModel(db, new(Bucket), new(BucketFile), new(Domain), new(User))
}

func NewDB() *gorm.DB {
	if err := db.DB().Ping(); err != nil {
		_ = db.Close()
		if db, err = createPool(config.Cfg.ServiceConfig.List.Mysql); err != nil {
			logger.Fatal(err.Error())
		}
	}
	return db
}

func registerModel(db *gorm.DB, tables ...interface{}) {
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

func createPool(cf *config.SqlConfig) (*gorm.DB, error) {
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
	// 是否开启日志模式
	conn.LogMode(cf.Debug).SetLogger(MysqlLogger{})
	return conn, nil
}

type MysqlLogger struct {
}

func (m MysqlLogger) Print(v ...interface{}) {

	values := gorm.LogFormatter(v...)
	_, _ = os.Stdout.WriteString(fmt.Sprintf(strings.Repeat("%v ", len(values)), values...))
	_ = os.Stdout.Sync()
	//logger.Infof(strings.Repeat("%v ", len(values)), values...)
}
