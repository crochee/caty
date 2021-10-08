// Date: 2021/9/20

// Package model
package model

import (
	"context"
	"testing"

	"cca/config"
	"obs/pkg/db"
)

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	err := config.LoadConfig("E:\\project\\cca\\conf\\cca.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := db.New().Debug()
	u := &User{
		AccountID:  "lcf",
		UserID:     "44",
		Account:    "90",
		PassWord:   "5",
		Email:      "6",
		Permission: `{"op":90}`,
		Verify:     0,
		Desc:       "lcf_desc",
	}
	if !d.Migrator().HasTable(u) {
		t.Log(d.Set("gorm:table_options",
			"ENGINE=InnoDB COMMENT='账户信息表' COLLATE='utf8mb4_bin'").AutoMigrate(u))
	}
	//t.Log(d.Model(u).Create(u))
	//t.Log(d.Model(u).Delete(u))
	//t.Log(d.Model(u).First(u))
}
