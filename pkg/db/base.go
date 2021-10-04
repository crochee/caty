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
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:current_timestamp();comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:current_timestamp() on update current_timestamp();comment:更新时间"`
	DeletedAt DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
}
