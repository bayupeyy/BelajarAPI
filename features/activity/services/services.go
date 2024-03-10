package services

import (
	"21-api/features/activity"
	"21-api/helper"
	"21-api/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m activity.ActivityModel
	v *validator.Validate
}

func NewActivityService(model activity.ActivityModel) activity.ActivityService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) AddActivity(pemilik *jwt.Token, kegiatanBaru activity.Activity) (activity.Activity, error) {
	hp := middlewares.DecodeToken(pemilik)
	if hp == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return activity.Activity{}, errors.New("data tidak valid")
	}

	err := s.v.Struct(&kegiatanBaru)
	if err != nil {
		log.Println("error validasi", err.Error())
		return activity.Activity{}, err
	}

	result, err := s.m.InsertActivity(hp, kegiatanBaru)
	if err != nil {
		return activity.Activity{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}
