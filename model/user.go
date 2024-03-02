package model

import "gorm.io/gorm"

type User struct {
	Nama     string `json:"nama" form:"nama" validate:"required"`
	Hp       string `json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Password string `json:"password" form:"password" validate:"required"`
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
