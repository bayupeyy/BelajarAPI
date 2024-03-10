package routes

import (
	"21-api/config"
	user "21-api/fearures/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController) {
	userRoute(c, ctl)
	todoRoute(c, tc)

}

func userRoute(c *echo.Echo, ctl user.UserController) {
	c.POST("/users", ctl.Register()) // register -> umum (boleh diakses semua orang)
	c.POST("/login", ctl.Login())
	c.GET("/users/:hp", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get profile -> butuh penanda khusus

	c.GET("/users", ctl.ListUser(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get all user -> butuh penanda khusus
	c.PUT("/users/:hp", ctl.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // update user -> butuh penanda khusus

}

func todoRoute(c *echo.Echo, tc activity.ActivityController) {
	//Menambah Kegiatan
	c.POST("/kegiatan", ctl.AddActivity(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))

	//Update Kegiatan
	c.PUT("/kegiatan/:id", ctl.UpdateActivity(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))

	//Melihat list kegiatan
	c.GET("/kegiatan", ctl.GetAllActivities(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
