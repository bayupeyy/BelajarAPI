package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	ID         uint       `gorm:"primaryKey"`
	Nama       string     `json:"nama" form:"nama" validate:"required"`
	Hp         string     `gorm:"uniqueIndex" json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Password   string     `json:"password" form:"password" validate:"required"`
	Activities []Activity `gorm:"foreignKey:UserID"`
}

type Activity struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user_id"`
	Kegiatan  string `json:"kegiatan" form:"kegiatan" validate:"required"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required"`
}

type ActivityResponse struct {
	ID        uint   `json:"id"`
	Kegiatan  string `json:"kegiatan"`
	Deskripsi string `json:"deskripsi"`
}

type UserModel struct {
	Connection *gorm.DB
}

func (um *UserModel) AddUser(newData User) error {
	err := um.Connection.Create(&newData).Error
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) CekUser(hp string) bool {
	var data User
	if err := um.Connection.Where("hp = ?", hp).First(&data).Error; err != nil {
		return false
	}
	return true
}

func (um *UserModel) Update(hp string, data User) error {
	if err := um.Connection.Model(&data).Where("hp = ?", hp).Update("nama", data.Nama).Update("password", data.Password).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) GetAllUser() ([]User, error) {
	var result []User

	if err := um.Connection.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (um *UserModel) GetProfile(hp string) (User, error) {
	var result User
	if err := um.Connection.Where("hp = ?", hp).First(&result).Error; err != nil {
		return User{}, err
	}
	return result, nil
}

func (um *UserModel) Login(hp string, password string) (User, error) {
	var result User
	if err := um.Connection.Where("hp = ? AND password = ?", hp, password).First(&result).Error; err != nil {
		return User{}, err
	}
	return result, nil
}

// Fungsi untuk menambah kegiatan
func (um *UserModel) AddActivity(hp string, activity Activity) error {
	//Mendapatkan pengguna berdasarkan No Hp
	var user User
	if err := um.Connection.Where("hp = ?", hp).First(&user).Error; err != nil {
		return err
	}

	//Set no Hp pengguna untuk kegiatan yang akan ditambahkan
	activity.UserID = user.ID

	//Menambahkan ke dalam DB
	if err := um.Connection.Create(&activity).Error; err != nil {
		return err
	}

	return nil
}

// UpdateActivity mengubah kegiatan yang dimiliki oleh pengguna yang sedang login
func (um *UserModel) UpdateActivity(userID uint, activityID uint, updatedActivity Activity) error {
	// Mencari kegiatan berdasarkan ID
	var activity Activity
	if err := um.Connection.Where("id = ? AND user_id = ?", activityID, userID).First(&activity).Error; err != nil {
		return err
	}

	// Memastikan kegiatan yang ingin diubah dimiliki oleh pengguna yang sedang login
	if activity.UserID != userID {
		return errors.New("Anda tidak memiliki izin untuk mengubah kegiatan ini")
	}

	// Memperbarui kegiatan yang ditemukan
	if err := um.Connection.Model(&activity).Updates(&updatedActivity).Error; err != nil {
		return err
	}

	return nil
}

// GetActivitiesByHp mengambil daftar kegiatan berdasarkan nomor HP pengguna
func (um *UserModel) GetActivitiesByHp(hp string) ([]ActivityResponse, error) {
	var activities []Activity

	// Mendapatkan pengguna berdasarkan nomor HP
	var user User
	if err := um.Connection.Where("hp = ?", hp).First(&user).Error; err != nil {
		return nil, err
	}

	// Mengambil daftar kegiatan berdasarkan ID pengguna
	if err := um.Connection.Model(&user).Association("Activities").Find(&activities); err != nil {
		return nil, err
	}

	// Membuat slice untuk menyimpan respons kegiatan
	var activityResponses []ActivityResponse

	// Mengonversi setiap kegiatan menjadi respons kegiatan yang diinginkan
	for _, activity := range activities {
		activityResponse := ActivityResponse{
			ID:        activity.ID,
			Kegiatan:  activity.Kegiatan,
			Deskripsi: activity.Deskripsi,
		}
		activityResponses = append(activityResponses, activityResponse)
	}

	return activityResponses, nil
}
