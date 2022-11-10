package models

import (
	"gorm.io/gorm"
	"zgin/global"
)

type Users struct {
	BaseModel
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Status   int    `gorm:"column:status"`
	Phone    string `gorm:"column:phone"`
}

func (this *Users) TableName() string {
	return "users"
}

func (this *Users) Table() *gorm.DB {
	if this.Tx != nil {
		return this.Tx.Table(this.TableName())
	}
	return global.MysqlSub.Table(this.TableName())
}
