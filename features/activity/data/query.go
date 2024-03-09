package data

import (
	"21-api/features/activity"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) activity.activityModel {
	return &model{
		connection: db,
	}
}

func (am *model) InsertActivity(deskripsi string, kegiatanBaru activity.Activity) (activity.Activity, error) {
	var inputProcess = Activity{Kegiatan: kegiatanBaru.Kegiatan, Deskripsi: deskripsi}
	if err := am.connection.Create(&inputProcess).Error; err != nil {
		return activity.Activity{}, err
	}

	return activity.Activity{Kegiatan: inputProcess.Kegiatan}, nil
}
