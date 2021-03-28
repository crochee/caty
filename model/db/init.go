// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"obs/config"
	"obs/cron"
)

var db *gorm.DB

// Setup init mysql db
func Setup(ctx context.Context) error {
	var err error
	if db, err = createPool(config.Cfg.ServiceConfig.List.Mysql,
		generateGormConfig(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
			config.Cfg.ServiceConfig.List.Mysql.Debug)); err != nil {
		return err
	}
	db = db.WithContext(ctx)
	// 自动建表或者自适应表字段
	bucket := new(Bucket)
	bucketFile := new(BucketFile)
	domain := new(Domain)
	user := new(User)
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(
		bucket,
		bucketFile,
		domain,
		user,
	); err != nil {
		return err
	}
	if _, err = cron.New().AddFunc("*/20 */30 * * * *", bucket.Delete); err != nil {
		return err
	}
	if _, err = cron.New().AddFunc("*/20 */30 * * * *", bucketFile.Delete); err != nil {
		return err
	}
	if _, err = cron.New().AddFunc("*/20 */30 * * * *", domain.Delete); err != nil {
		return err
	}
	if _, err = cron.New().AddFunc("*/20 */30 * * * *", user.Delete); err != nil {
		return err
	}
	// 开启池化之后不能close  否则连接池没有作用
	// 设置数据库连接池最大连接数maxOpenConns
	// 设置数据库连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于maxIdleConns，超过的连接会被连接池关闭
	var sqlDb *sql.DB
	if sqlDb, err = db.DB(); err != nil {
		return err
	}
	if err = sqlDb.Ping(); err != nil {
		_ = sqlDb.Close()
		return err
	}
	sqlDb.SetMaxIdleConns(10)                   // 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(30)                   // 设置打开数据库连接的最大数量
	sqlDb.SetConnMaxLifetime(time.Second * 300) // 设置了连接可复用的最大时间
	return nil
}

// NewDB get gorm.DB
func NewDB() *gorm.DB {
	return db
}

// Close close db pool
func Close() error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func generateGormConfig(path string, debug bool) *gorm.Config {
	var writer io.Writer = os.Stdout
	if path != "" {
		writer = &lumberjack.Logger{
			Filename:   path,
			MaxSize:    50,    //单个日志文件最大MaxSize*M大小 // megabytes
			MaxAge:     30,    //days
			MaxBackups: 10,    //备份数量
			Compress:   false, //不压缩
			LocalTime:  true,  //备份名采用本地时间
		}
	}
	l := logger.Warn
	if debug { // 开启debug 日志模式 conn = conn.Debug()
		l = logger.Info
	}
	return &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: logger.New(
			log.New(writer, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      l,
				Colorful:      path == "",
			}),
	}
}

func createPool(cf *config.SqlConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
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
		if conn, err = gorm.Open(mysql.Open(dsn), gormConfig); err != nil {
			return nil, err
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
			cf.Host,
			cf.Port,
			cf.User,
			cf.Database,
			cf.Password)
		if conn, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsuport type %s", cf.Type)
	}
	return conn, nil
}

// AnyTime mock time
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// NewMock new a mock todo mock解除测试对数据库等中间件的依赖
func NewMock() (sqlmock.Sqlmock, error) {
	// 创建sqlmock
	slqDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	// 结合gorm、sqlmock
	if db, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      slqDb,
	}), generateGormConfig("", true)); err != nil {
		return nil, err
	}
	return mock, err
}
