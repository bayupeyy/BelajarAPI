package main

import (
	"21-api/config"
	td "21-api/features/activity/data"
	th "21-api/features/activity/handler"
	ts "21-api/features/activity/services"
	"21-api/features/user/data"
	"21-api/features/user/handler"
	"21-api/features/user/services"
	"21-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()            // inisiasi echo
	cfg := config.InitConfig() // baca seluruh system variable
	db := config.InitSQL(cfg)  // Koneksi ke DB

	userData := data.New(db)
	userService := services.NewService(userData)
	userHandler := handler.NewUserHandler(userService)

	activityData := td.New(db)
	activityService := ts.NewActivityService(activityData)
	activityHandler := th.NewHandler(activityService)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // ini aja cukup
	routes.InitRoute(e, userHandler, activityHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
