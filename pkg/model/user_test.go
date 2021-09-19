// Date: 2021/9/20

// Package model
package model

import (
	"context"
	"testing"

	"obs/config"
	"obs/pkg/db"
)

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	err := config.LoadConfig("E:\\project\\obs\\conf\\obs.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	t.Log(db.Client().Debug().Set("gorm:table_options", "ENGINE=InnoDB COMMENT='账户信息表' COLLATE='utf8mb4_bin'").AutoMigrate(&User{}))
}
