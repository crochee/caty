// Date: 2021/10/15

// Package model
package model

import (
	"context"
	"testing"

	"cca/config"
	"cca/pkg/dbx"
)

func TestUser_TableName(t *testing.T) {
	ctx := context.Background()
	err := config.LoadConfig("E:/project/cca/conf/cca.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err = dbx.Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer dbx.Close()
	d := dbx.New().Debug()
	u := &User{}
	if !d.Migrator().HasTable(u) {
		t.Log(d.Set("gorm:table_options",
			"ENGINE=InnoDB COMMENT='账户信息表' COLLATE='utf8mb4_bin' DEFAULT CHARSET='utf8mb4'").
			AutoMigrate(u))
	}
}
