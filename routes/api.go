package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
)

func Api() {
	facades.Route().Prefix("api/").Group(func(r route.Router) {
		r.Prefix("user").Group(func(user route.Router) {
			userController := controllers.NewUserController()

			user.Get("users", userController.Show)
			user.Post("login", userController.Login)
			user.Post("register", userController.Register)
			user.Middleware(middleware.Jwt()).Get("info", userController.Info)
		})
	})
}
