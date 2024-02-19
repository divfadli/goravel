package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Api() {
	// userController := controllers.NewUserController()
	// facades.Route().Middleware(middleware.Cors()).Get("/users/{id}", userController.Show)
	// facades.Route().Middleware(user_middleware.Cors()).Get("users", userController.Show)
	facades.Route().Prefix("api/").Group(func(r route.Router) {
		r.Prefix("user").Group(func(user route.Router) {
			userController := controllers.NewUserController()

			user.Get("users", userController.Show)
			// user.Post("login" useruserController.Login),
		})
	})
}
