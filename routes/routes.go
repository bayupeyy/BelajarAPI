package routes

import (
	"21-api/config"
	"21-api/controller/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController) {
	c.POST("/users", ctl.Register()) // register -> umum (boleh diakses semua orang)
	c.POST("/login", ctl.Login())
	c.GET("/users", ctl.ListUser(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get all user -> butuh penanda khusus
	c.GET("/users/:hp", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get profile -> butuh penanda khusus
	c.PUT("/users/:hp", ctl.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // update user -> butuh penanda khusus

	//Menambah Kegiatan
	c.POST("/kegiatan", ctl.AddActivity(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}