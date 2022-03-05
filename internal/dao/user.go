package model

import (
	"zgin/global"
	"zgin/internal/model"
)

type User model.Users

func (this *User) GetOne() {
	global.Mysql.Model(this).Where("phone = ?", this.Phone).Take(this)
}
