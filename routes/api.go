package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	pdf "goravel/app/http/controllers/generator"
)

func Api() {
	laporan := controllers.NewLaporan()
	facades.Route().Post("/storeLaporan", laporan.Create)
	generates := pdf.NewPdf("")
	facades.Route().Post("generate", generates.GenerateKeamanan)
	facades.Route().Prefix("api").Group(func(r route.Router) {
		r.Prefix("user").Group(func(user route.Router) {
			userController := controllers.NewUserController()

			user.Get("users", userController.Show)
			user.Post("login", userController.Login)
			user.Post("register", userController.Register)
			user.Get("dataPegawai", userController.FindPegawai)

			user.Prefix("role").Group(func(role route.Router) {
				role.Get("getRole", userController.GetRole)
				role.Post("storeRole", userController.StoreRole)
			})
			// user.Middleware(middleware.Jwt(), middleware.Cors()).Get("info", userController.Info)
		})

		r.Prefix("kejadian").Group(func(kejadian route.Router) {
			jenisKejadianController := controllers.NewJenisKejadianController()

			kejadian.Post("storeKejadian", jenisKejadianController.PostKejadian)
			kejadian.Post("listKejadian", jenisKejadianController.ListKejadian)
			kejadian.Get("showDetailKejadian", jenisKejadianController.ShowDetailKejadian)
			kejadian.Delete("deleteKejadian", jenisKejadianController.DeleteKejadian)

			kejadian.Prefix("keamanan").Group(func(keamanan route.Router) {
				kejadianKeamananController := controllers.NewKejadianKeamananController()

				keamanan.Post("storeKejadianKeamanan", kejadianKeamananController.StoreKejadianKeamanan)
				keamanan.Post("listKejadianKeamanan", kejadianKeamananController.ListKejadianKeamanan)
				keamanan.Get("showDetailKejadianKeamanan", kejadianKeamananController.ShowDetailKejadianKeamanan)
				keamanan.Delete("deleteKejadianKeamanan", kejadianKeamananController.DeleteKejadianKeamanan)
			})

			kejadian.Prefix("keselamatan").Group(func(keselamatan route.Router) {
				kejadianKeselamatanController := controllers.NewKejadianKeselamatanController()

				keselamatan.Post("storeKejadianKeselamatan", kejadianKeselamatanController.StoreKejadianKeselamatan)
				keselamatan.Post("listKejadianKeselamatan", kejadianKeselamatanController.ListKejadianKeselamatan)
				keselamatan.Get("showDetailKejadianKeselamatan", kejadianKeselamatanController.ShowDetailKejadianKeselamatan)
				keselamatan.Delete("deleteKejadianKeselamatan", kejadianKeselamatanController.DeleteKejadianKeselamatan)
			})
		})

		r.Prefix("approval").Group(func(router route.Router) {
			approvalController := controllers.NewApproval()

			router.Post("/storeApproval", approvalController.StoreApproval)
			router.Get("/listApproval", approvalController.ListApproval)
		})
	})
}
