package data

import "21-api/features/activity/data"

type User struct {
	Nama       string
	Hp         string          `gorm:"type:varchar(13);uniqueIndex;primaryKey" json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Password   string          `json:"password" form:"password" validate:"required"`
	Activities []data.Activity `gorm:"foreignKey:UserID;references:Hp"`
}
