// Date: 2021/9/19

// Package db
package db

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"

	"obs/internal"
)

var (
	dbClient *DB

	ErrNotRowsAffected = errors.New("0 rows affected")
	ErrNotFound        = gorm.ErrRecordNotFound
)

// Init init database
func Init(ctx context.Context) (err error) {
	dbClient, err = NewClient(ctx, func(opt *Option) {
		opt.Debug = viper.GetString("mode") == "debug"
		opt.User = viper.GetString("mysql.user")
		opt.Password = viper.GetString("mysql.password")
		opt.IP = viper.GetString("mysql.ip")
		opt.Port = viper.GetString("mysql.port")
		opt.Database = viper.GetString("mysql.database")
		opt.Charset = viper.GetString("mysql.charset")
		opt.MaxOpenConn = viper.GetInt("mysql.max_open_conns")
		opt.MaxIdleConn = viper.GetInt("mysql.max_idle_conns")
		opt.ConnMaxLifetime = time.Duration(viper.GetInt("mysql.conn_max_lifetime")) * time.Second
	})
	return
}

func Close() {
	internal.Close(dbClient)
}

func New(ctx context.Context) *DB {
	return dbClient.WithContext(ctx)
}

func Client() *DB {
	return dbClient
}
