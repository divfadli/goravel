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
			kejadianController := controllers.NewKejadianController()

			kejadian.Post("storeKejadian", kejadianController.PostKejadian)
			kejadian.Post("listKejadian", kejadianController.ListKejadian)
			kejadian.Get("showDetailKejadian", kejadianController.ShowDetailKejadian)
			kejadian.Delete("deleteKejadian", kejadianController.DeleteKejadian)
		})

		r.Prefix("rekap").Group(func(rekap route.Router) {
			rekap.Prefix("keamanan").Group(func(keamanan route.Router) {
				// rekapKeamananController := controllers.NewRekapKejadianKeamananController()

				// keamanan.Post("storeRekapKeamanan", rekapKeamananController.StoreRekapKeamanan)
				// keamanan.Post("listRekapKeamanan", rekapKeamananController.ListRekapKeamanan)
				// keamanan.Get("showDetailRekapKeamanan", rekapKeamananController.ShowDetailRekapKeamanan)
				// keamanan.Delete("deleteRekapKeamanan", rekapKeamananController.DeleteRekapKeamanan)
			})
		})
	})

}
