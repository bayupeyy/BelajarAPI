package data

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Kegiatan  string `json:"kegiatan" form:"kegiatan" validate:"required"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required"`
}
