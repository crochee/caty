package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/crochee/lib"
	"github.com/crochee/lib/log"
)

var (
	NotFound           = gorm.ErrRecordNotFound
	ErrNotRowsAffected = errors.New("0 rows affected")
)

type Option struct {
	Debug bool

	MaxOpenConn int
	MaxIdleConn int

	User     string
	Password string
	IP       string
	Port     string
	Database string
	Charset  string

	Timeout         time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ConnMaxLifetime time.Duration
}

// NewClient with context.Context returns DB
func NewClient(ctx context.Context, opts ...func(*Option)) (*DB, error) {
	var c Option
	for _, opt := range opts {
		opt(&c)
	}
	client, err := gorm.Open(mysql.Open(dsn(&c)),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 不考虑表名单复数变化
			},
			DisableForeignKeyConstraintWhenMigrating: true,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
	if err != nil {
		return nil, err
	}
	session := &gorm.Session{Context: ctx}
	if c.Debug { // 是否显示sql语句
		session.Logger = client.Logger.LogMode(logger.Info)
	}
	client = client.Session(session)

	var sqlDB *sql.DB
	if sqlDB, err = client.DB(); err != nil {
		return nil, err
	}
	// 连接池配置
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)        // 默认值0，无限制
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)        // 默认值2
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime) // 默认值0，永不过期
	return &DB{DB: client, debug: c.Debug}, nil
}

type DB struct {
	*gorm.DB
	debug bool
}

// WithContext gets logger to set orm logger
func (d *DB) WithContext(ctx context.Context) *DB {
	fromContextLog := log.FromContext(ctx)
	return &DB{DB: d.Session(&gorm.Session{
		Context: ctx,
		Logger: newLog(fromContextLog, logger.Config{
			SlowThreshold: 10 * time.Second,
			LogLevel:      getLevel(fromContextLog.Opt().Level, d.debug),
		}),
	}),
		debug: d.debug,
	}
}

// Close impl io.Closer to defer close db pool
func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	lib.Close(sqlDB)
	return nil
}

func getLevel(l string, debug bool) logger.LogLevel {
	if debug {
		return logger.Info
	}
	switch l {
	case log.INFO, log.DEBUG:
		return logger.Info
	case log.WARN:
		return logger.Warn
	case log.ERROR, log.DPanic, log.PANIC, log.FATAL:
		return logger.Error
	default:
		return logger.Silent
	}
}

func dsn(opt *Option) string {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		opt.User, opt.Password, opt.IP, opt.Port, opt.Database, opt.Charset, true, "UTC")
	if opt.Timeout != 0 {
		uri += fmt.Sprintf("&timeout=%s", opt.Timeout)
	}
	if opt.ReadTimeout != 0 {
		uri += fmt.Sprintf("&readTimeout=%s", opt.ReadTimeout)
	}
	if opt.ReadTimeout != 0 {
		uri += fmt.Sprintf("&writeTimeout=%s", opt.ReadTimeout)
	}
	return uri
}
