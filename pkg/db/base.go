package db

import (
	"time"

	"gorm.io/gorm"

	"obs/pkg/id"
)

type SnowID struct{}

func (SnowID) BeforeCreate(tx *gorm.DB) error {
	_, ok := tx.Statement.Schema.FieldsByDBName["id"]
	if !ok {
		return nil
	}
	snowID, err := id.NextID()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("id", snowID)
	return nil
}

type Base struct {
	SnowID
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}
