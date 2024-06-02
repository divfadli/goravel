package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Api() {
	facades.Route().Prefix("api/").Group(func(r route.Router) {
		r.Prefix("user").Group(func(user route.Router) {
			userController := controllers.NewUserController()

			user.Get("users", userController.Show)
			user.Post("login", userController.Login)
			// user.Post("register", userController.Register)
			// user.Middleware(middleware.Jwt(), middleware.Cors()).Get("info", userController.Info)
		})

		r.Prefix("kejadian").Group(func(kejadian route.Router) {
			jenisKejadianController := controllers.NewJenisKejadianController()

			kejadian.Post("storeKejadian", jenisKejadianController.PostKejadian)
			kejadian.Post("listKejadian", jenisKejadianController.ListKejadian)
			kejadian.Get("showDetailKejadian", jenisKejadianController.ShowDetailKejadian)
			kejadian.Delete("deleteKejadian", jenisKejadianController.DeleteKejadian)

			kejadian.Prefix("keamanan").Group(func(keamanan route.Router) {
				KejadianKeamananController := controllers.NewKejadianKeamananController()

				keamanan.Post("storeKejadianKeamanan", KejadianKeamananController.StoreKejadianKeamanan)
				keamanan.Post("listKejadianKeamanan", KejadianKeamananController.ListKejadianKeamanan)
				// keamanan.Get("showDetailKejadianKeamanan", KejadianKeamananController.ShowDetailKejadianKeamanan)
				keamanan.Delete("deleteKejadianKeamanan", KejadianKeamananController.DeleteKejadianKeamanan)
			})
		})
	})

}
