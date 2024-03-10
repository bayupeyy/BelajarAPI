package activity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ActivityController interface {
	Add() echo.HandlerFunc
	// Update() echo.HandlerFunc
	// Delete() echo.HandlerFunc
	// ShowMyTodo() echo.HandlerFunc
}

type ActivityModel interface {
	InsertActivity(pemilik string, kegiatanBaru Activity) (Activity, error)
	UpdateActivity(pemilik string, activityID uint, data Activity) (Activity, error)
	// DeleteActivity()
	GetActivityByOwner(pemilik string) ([]Activity, error)
}

type ActivityService interface {
	AddActivity(pemilik *jwt.Token, kegiatanBaru Activity) (Activity, error)
	// UpdateTodo(pemilik *jwt.Token, todoID string, data Todo) (Todo, error)
}

type Activity struct {
	Kegiatan string
}
