package branchmodel

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"zgin/global"
	"zgin/internal/models"
)

type Orders struct {
	models.BaseModel
	UserId int64 `gorm:"column:user_id"`
	Tx *gorm.DB `gorm:"-"`
}

const smsOrderSubTableBaseNum = 200

func (this *Orders) TableName() string {
	return "orders"
}

func (this *Orders)Table() *gorm.DB {
	tableName := fmt.Sprintf("%s_%d",this.TableName(),math.Ceil(float64(this.UserId)/baseNum))
	if this.Tx != nil {
		return this.Tx.Table(tableName)
	}
	return global.MysqlSub.Table(tableName)
}

func (this *Orders) IsExists() int64 {
	var id int64
	global.MysqlSub.Table(this.TableName()).Where("id = ?", this.Id).Pluck("id", &id)

	return id
}
