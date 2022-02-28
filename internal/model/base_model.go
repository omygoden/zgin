package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int64        `gorm:"column:id"`
	CreatedAt sql.NullTime `gorm:"column:created_at;<-:create"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
	Tx        *gorm.DB     `gorm:"-"`
}
