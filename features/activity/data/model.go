package data

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Kegiatan string `json:"kegiatan" form:"kegiatan" validate:"required"`
	Pemilik  string `gorm:"type:varchar(13);"`
}
