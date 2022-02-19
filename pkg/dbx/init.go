package dbx

import (
	"context"

	"github.com/crochee/lirity"
	"github.com/crochee/lirity/db"
	"github.com/spf13/viper"

	"time"
)

var dbClient *db.DB

// Init init database
func Init(ctx context.Context) (err error) {
	dbClient, err = db.NewClient(ctx, func(opt *db.Option) {
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
	lirity.Close(dbClient)
}

func With(ctx context.Context) *db.DB {
	return dbClient.With(ctx)
}

func New() *db.DB {
	return dbClient
}
