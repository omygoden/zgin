package model

type Users struct {
	BaseModel
	Phone string `gorm:"column:phone"`
}

func (this *Users) TableName() string {
	return "users"
}
