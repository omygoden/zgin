package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"zgin/global"
	"zgin/pkg/sflogger"
	"zgin/pkg/util"
)

type Users struct {
	BaseModel
	Phone string `gorm:"column:phone"`
}

func (this *Users) TableName() string {
	return "users"
}

func (this *Users) GetOne() {
	err := global.Mysql.Model(this).Where("phone = ?", this.Phone).Take(this).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
	if err != nil {
		p, _ := json.Marshal(*this)
		sflogger.MysqlErrorLog(this.TableName(), err.Error(), util.GetMyFuncName(), string(p))
	}
}
