package routes

import (
	"21-api/config"
	activity "21-api/features/activity"
	user "21-api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController, ac activity.ActivityController) {
	userRoute(c, ctl)
	activityRoute(c, ac)

}

func userRoute(c *echo.Echo, ctl user.UserController) {
	c.POST("/login", ctl.Login())
	c.GET("/users/:hp", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get profile -> butuh penanda khusus

	// c.GET("/users", ctl.ListUser(), echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(config.JWTSECRET),
	// })) // get all user -> butuh penanda khusus
	// c.PUT("/users/:hp", ctl.Update(), echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(config.JWTSECRET),
	// })) // update user -> butuh penanda khusus

}

func activityRoute(c *echo.Echo, ac activity.ActivityController) {
	//Menambah Kegiatan
	c.POST("/kegiatan", ac.Add(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
