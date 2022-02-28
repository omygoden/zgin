package branchmodel

import (
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
	"zgin/global"
	"zgin/internal/model"
	"zgin/pkg/sflogger"
	"zgin/pkg/util"
)

type Orders struct {
	model.BaseModel
	UserId int64 `gorm:"column:user_id"`
}

const smsOrderSubTableBaseNum = 200

func (this *Orders) TableName() string {
	sub := strconv.Itoa(int(this.UserId) / smsOrderSubTableBaseNum)
	return "orders_" + sub
}

func (this *Orders) IsExists() int64 {
	var id int64
	err := global.MysqlSub.Table(this.TableName()).Where("id = ?", this.Id).Pluck("id", &id).Error

	if err == gorm.ErrRecordNotFound {
		return id
	}

	if err != nil {
		p, _ := json.Marshal(*this)
		sflogger.MysqlErrorLog(this.TableName(), err.Error(), util.GetMyFuncName(), string(p))
	}
	return id
}
