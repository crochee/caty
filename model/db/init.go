// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package db

import (
	"database/sql"
	"fmt"
	"io"
	"log"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"obs/config"
	loggers "obs/logger"
)

var db *gorm.DB

// Setup init mysql db
func Setup() {
	var err error
	if db, err = createPool(config.Cfg.ServiceConfig.List.Mysql,
		loggers.SetLoggerWriter(config.Cfg.ServiceConfig.ServiceInfo.LogPath)); err != nil {
		loggers.Fatal(err.Error())
		return
	}
	// 开启池化之后不能close  否则连接池没有作用
	// 设置数据库连接池最大连接数maxOpenConns
	// 设置数据库连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于maxIdleConns，超过的连接会被连接池关闭
	var sqlDb *sql.DB
	if sqlDb, err = db.DB(); err != nil {
		loggers.Fatal(err.Error())
		return
	}
	if err = sqlDb.Ping(); err != nil {
		_ = sqlDb.Close()
		loggers.Fatal(err.Error())
		return
	}
	sqlDb.SetMaxIdleConns(10)                   // 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(30)                   // 设置打开数据库连接的最大数量
	sqlDb.SetConnMaxLifetime(time.Second * 300) // 设置了连接可复用的最大时间

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
		AutoMigrate(new(Bucket), new(BucketFile), new(Domain), new(User)); err != nil {
		loggers.Fatal(err.Error())
	}
}

// NewDB get gorm.DB
func NewDB() *gorm.DB {
	return db
}

func createPool(cf *config.SqlConfig, writer io.Writer) (*gorm.DB, error) {
	var (
		conn *gorm.DB
		err  error
	)
	switch cf.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			cf.User,
			cf.Password,
			cf.Host,
			cf.Port,
			cf.Database,
			cf.Charset)
		if conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
			return nil, err
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
			cf.Host,
			cf.Port,
			cf.User,
			cf.Database,
			cf.Password)
		if conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsuport type %s", cf.Type)
	}
	if cf.Debug {
		conn = conn.Session(&gorm.Session{Logger: logger.New(
			log.New(writer, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      true,
			})}) // 开启debug 日志模式 conn = conn.Debug()
	}
	return conn, nil
}

var Mock sqlmock.Sqlmock

func SetUpMock() {
	// 创建sqlmock
	var (
		slqDb *sql.DB
		err   error
	)
	if slqDb, Mock, err = sqlmock.New(); err != nil {
		loggers.Fatal(err.Error())
		return
	}
	// 结合gorm、sqlmock
	if db, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      slqDb,
	}), &gorm.Config{}); err != nil {
		loggers.Fatal(err.Error())
	}
}
