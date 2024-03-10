package data

import (
	"21-api/features/activity"
	"errors"

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

func (am *model) InsertActivity(pemilik string, kegiatanBaru activity.Activity) (activity.Activity, error) {
	var inputProcess = Activity{Kegiatan: kegiatanBaru.Kegiatan, Pemilik: pemilik}
	if err := am.connection.Create(&inputProcess).Error; err != nil {
		return activity.Activity{}, err
	}

	return activity.Activity{Kegiatan: inputProcess.Kegiatan}, nil
}

// Fungsi untuk Update kegiatan
func (am *model) UpdateActivity(pemilik string, activityID uint, data activity.Activity) (activity.Activity, error) {
	var qry = am.connection.Where("pemilik = ? AND id = ?", pemilik, activityID).Updates(data)
	if err := qry.Error; err != nil {
		return activity.Activity{}, err
	}

	if qry.RowsAffected < 1 {
		return activity.Activity{}, errors.New("no data affected")
	}

	return data, nil
}

func (am *model) GetActivityByOwner(pemilik string) ([]activity.Activity, error) {
	var result []activity.Activity
	if err := am.connection.Where("pemilik = ?", pemilik).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
