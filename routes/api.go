package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	pdf "goravel/app/http/controllers/generator"
	"goravel/app/http/middleware"
)

func Api() {
	// laporan := controllers.NewLaporan()
	// facades.Route().Get("/listLaporan", laporan.ListLaporan)

	// facades.Route().Post("/storeLaporan", laporan.Create)
	generates := pdf.NewPdf("")
	facades.Route().Static("storage", "./storage")

	// Generate Laporan
	facades.Route().Get("generate-mingguan", generates.GenerateMingguan)
	facades.Route().Get("generate-bulanan", generates.GenerateBulanan)
	facades.Route().Get("generate-triwulan", generates.GenerateTriwulan)

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

			kejadian.Middleware(middleware.Jwt(), middleware.Cors()).Post("storeKejadian", jenisKejadianController.PostKejadian)
			kejadian.Middleware(middleware.Jwt(), middleware.Cors()).Post("listKejadian", jenisKejadianController.ListKejadian)
			kejadian.Middleware(middleware.Jwt(), middleware.Cors()).Get("showDetailKejadian", jenisKejadianController.ShowDetailKejadian)
			kejadian.Middleware(middleware.Jwt(), middleware.Cors()).Delete("deleteKejadian", jenisKejadianController.DeleteKejadian)

			kejadian.Prefix("keamanan").Group(func(keamanan route.Router) {
				kejadianKeamananController := controllers.NewKejadianKeamananController()

				keamanan.Middleware(middleware.Jwt(), middleware.Cors()).Post("export-excel", kejadianKeamananController.ExportExcel)

				keamanan.Middleware(middleware.Jwt(), middleware.Cors()).Post("storeKejadianKeamanan", kejadianKeamananController.StoreKejadianKeamanan)
				keamanan.Middleware(middleware.Jwt(), middleware.Cors()).Post("listKejadianKeamanan", kejadianKeamananController.ListKejadianKeamanan)
				keamanan.Middleware(middleware.Jwt(), middleware.Cors()).Get("showDetailKejadianKeamanan", kejadianKeamananController.ShowDetailKejadianKeamanan)
				keamanan.Middleware(middleware.Jwt(), middleware.Cors()).Delete("deleteKejadianKeamanan", kejadianKeamananController.DeleteKejadianKeamanan)
			})

			kejadian.Prefix("keselamatan").Group(func(keselamatan route.Router) {
				kejadianKeselamatanController := controllers.NewKejadianKeselamatanController()

				keselamatan.Middleware(middleware.Jwt(), middleware.Cors()).Post("export-excel", kejadianKeselamatanController.ExportExcel)

				keselamatan.Middleware(middleware.Jwt(), middleware.Cors()).Post("storeKejadianKeselamatan", kejadianKeselamatanController.StoreKejadianKeselamatan)
				keselamatan.Middleware(middleware.Jwt(), middleware.Cors()).Post("listKejadianKeselamatan", kejadianKeselamatanController.ListKejadianKeselamatan)
				keselamatan.Middleware(middleware.Jwt(), middleware.Cors()).Get("showDetailKejadianKeselamatan", kejadianKeselamatanController.ShowDetailKejadianKeselamatan)
				keselamatan.Middleware(middleware.Jwt(), middleware.Cors()).Delete("deleteKejadianKeselamatan", kejadianKeselamatanController.DeleteKejadianKeselamatan)
			})
		})

		r.Prefix("approval").Group(func(router route.Router) {
			approvalController := controllers.NewApproval()

			router.Middleware(middleware.Jwt(), middleware.Cors()).Post("/storeApproval", approvalController.StoreApproval)
			router.Middleware(middleware.Jwt(), middleware.Cors()).Get("/listApproval", approvalController.ListApproval)
		})

		r.Prefix("laporan").Group(func(router route.Router) {
			laporanController := controllers.NewLaporan()

			router.Middleware(middleware.Jwt(), middleware.Cors()).Get("/listLaporan", laporanController.ListLaporan)
		})
	})
}
