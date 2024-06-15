package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Static("/js", "./public/js")
	facades.Route().Static("/vendor", "./public/vendor")
	facades.Route().Static("/img", "./public/img")
	facades.Route().Static("/css", "./public/css")
	facades.Route().Static("/api/files/", "./storage/app")
	facades.Route().Static("/plugin", "./plugins/custom")

	facades.Route().Get("/map", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("map.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Get("/approval", func(ctx http.Context) http.Response {
		// Retrieve cached user data
		userInfo := facades.Cache().Get("user_data")

		// Check if data is available in cache
		if userInfo != nil {
			return ctx.Response().View().Make("approval.tmpl", map[string]interface{}{
				"title":       "Approval | Rekapitulasi",
				"pageheading": "Approval Laporan",
				"version":     support.Version,
				"data":        userInfo,
			})
		}

		// For instance, you might redirect the user to the login page
		return ctx.Response().Redirect(http.StatusFound, "/login")

	})

	// Login
	facades.Route().Get("login", func(ctx http.Context) http.Response {
		loginURL := "/api/user/login"
		return ctx.Response().View().Make("login.tmpl", map[string]any{
			"loginURL": loginURL,
			"version":  support.Version,
		})
	})
	// Logout
	facades.Route().Get("logout", func(ctx http.Context) http.Response {
		facades.Cache().Flush()
		facades.Auth().Logout(ctx)

		return ctx.Response().Redirect(http.StatusFound, "/login")
	})

	// Register
	facades.Route().Get("/register", func(ctx http.Context) http.Response {
		registerURL := "/api/user/register"
		return ctx.Response().View().Make("register.tmpl", map[string]any{
			"title":       "Approval | Rekapitulasi",
			"pageheading": "Approval Laporan",
			"registerURL": registerURL,
			"version":     support.Version,
		})
	})

	// Dashboard
	facades.Route().Get("dashboard", func(ctx http.Context) http.Response {
		// Retrieve cached user data
		userInfo := facades.Cache().Get("user_data")
		dataKeamananURL := "/api/kejadian/keamanan/listKejadianKeamanan"
		dataKeselamatanURL := "/api/kejadian/keselamatan/listKejadianKeselamatan"
		// fmt.Println(userInfo)

		// // Check if data is available in cache
		if userInfo != nil {
			return ctx.Response().View().Make("index.tmpl", map[string]interface{}{
				"title":           "Dashboard | Rekapitulasi",
				"pageheading":     "Dashboard",
				"version":         support.Version,
				"dataKeamananURL": dataKeamananURL,
				"dataKeselamatanURL": dataKeselamatanURL,
				"data":            userInfo,
			})
		}

		facades.Auth().Logout(ctx)
		// For instance, you might redirect the user to the login page
		return ctx.Response().Redirect(http.StatusFound, "/login")
	})
	
	facades.Route().Prefix("role_user").Group(func(router route.Router) {
		// Retrieve cached user data
		router.Get("", func(ctx http.Context) http.Response {
			userInfo := facades.Cache().Get("user_data")
			dataRoleURL := "/api/role/listRole"
			deleteRoleURL := "/api/role/deleteRole"
			if userInfo != nil {
				return ctx.Response().View().Make("role_user.tmpl", map[string]interface{}{
					"title":         "Role User | Rekapitulasi",
					"pageheading":   "Role",
					"version":       support.Version,
					"dataRoleURL":   dataRoleURL,
					"deleteRoleURL": deleteRoleURL,
					"data":          userInfo,
				})
			}

			facades.Auth().Logout(ctx)
			// For instance, you might redirect the user to the login page
			return ctx.Response().Redirect(http.StatusFound, "/login")
		})
		router.Get("form_role_user", func(ctx http.Context) http.Response {
			userInfo := facades.Cache().Get("user_data")
			idUser := ctx.Request().Query("id_user")

			dataPegawaiURL := "/api/user/dataPegawai"
			dataRoleURL := "/api/user/role/getRole"
			storeRolePegawai := "/api/user/storeRolePegawai"

			if userInfo != nil {
				return ctx.Response().View().Make("form_role_user.tmpl", map[string]interface{}{
					"title":            "Role User | Rekapitulasi",
					"pageheading":      "Role",
					"version":          support.Version,
					"dataPegawaiURL":   dataPegawaiURL,
					"getRoleURL":       dataRoleURL,
					"storeRolePegawai": storeRolePegawai,
					"idUser":           idUser,
					"data":             userInfo,
				})
			}

			facades.Auth().Logout(ctx)
			// For instance, you might redirect the user to the login page
			return ctx.Response().Redirect(http.StatusFound, "/login")
		})
	})

	facades.Route().Prefix("jenis_kejadian").Group(func(r route.Router) {
		// Retrieve cached user data
		
		r.Get("", func(ctx http.Context) http.Response {
			userInfo := facades.Cache().Get("user_data")
			dataJenisKejadianURL := "/api/kejadian/listKejadian"
			deleteJenisKejadianURL := "/api/kejadian/deleteKejadian"

			if userInfo != nil {
				return ctx.Response().View().Make("jenis_kejadian.tmpl", map[string]interface{}{
					"title":                  "Jenis Kejadian | Rekapitulasi",
					"pageheading":            "Kejadian",
					"version":                support.Version,
					"dataJenisKejadianURL":   dataJenisKejadianURL,
					"deleteJenisKejadianURL": deleteJenisKejadianURL,
					"data":                   userInfo,
				})
			}
			facades.Auth().Logout(ctx)
			// For instance, you might redirect the user to the login page
			return ctx.Response().Redirect(http.StatusFound, "/login")
		})
		r.Get("form_jenis_kejadian", func(ctx http.Context) http.Response {
			userInfo := facades.Cache().Get("user_data")
			idJenisKejadian := ctx.Request().Query("id_jenis_kejadian")

			storeJenisKejadian := "/api/kejadian/storeKejadian"
			getJenisKejadian := "/api/kejadian/showDetailKejadian"

			if userInfo != nil {
				return ctx.Response().View().Make("form_jenis_kejadian.tmpl", map[string]interface{}{
					"title":                 "Form Jenis Kejadian | Rekapitulasi",
					"pageheading":           "Form Jenis Kejadian",
					"storeJenisKejadianURL": storeJenisKejadian,
					"getJenisKejadianURL":   getJenisKejadian,
					"idJenisKejadian":       idJenisKejadian,
					"version":               support.Version,
					"data":                  userInfo,
				})
			}
			facades.Auth().Logout(ctx)
			// For instance, you might redirect the user to the login page
			return ctx.Response().Redirect(http.StatusFound, "/login")
		})
	})

	// Rekap Data Kejadian
	facades.Route().Prefix("kejadian").Group(func(r route.Router) {
		r.Prefix("keamanan").Group(func(pelanggaran route.Router) {
			// Pelanggaran
			pelanggaran.Get("", func(ctx http.Context) http.Response {
				userInfo := facades.Cache().Get("user_data")
				dataKeamananURL := "/api/kejadian/keamanan/listKejadianKeamanan"
				deleteKejadianKeamananURL := "/api/kejadian/keamanan/deleteKejadianKeamanan"
				// Retrieve cached user data
				
				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("kejadian_keamanan.tmpl", map[string]interface{}{
						"title":                     "Kejadian Keamanan",
						"pageheading":               "Pelanggaran",
						"version":                   support.Version,
						"dataKeamananURL":           dataKeamananURL,
						"deleteKejadianKeamananURL": deleteKejadianKeamananURL,
						"data":                      userInfo,
					})
				}

				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})

			pelanggaran.Get("form_kejadian_keamanan", func(ctx http.Context) http.Response {
				userInfo := facades.Cache().Get("user_data")
				idKejadianKeamanan := ctx.Request().Query("id_kejadian_keamanan")
				
				storeKejadianKeamanan := "/api/kejadian/keamanan/storeKejadianKeamanan"
				getKejadianKeamanan := "/api/kejadian/keamanan/showDetailKejadianKeamanan"

				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("form_kejadian_keamanan.tmpl", map[string]interface{}{
						"title":                  "Form Kejadian Keamanan",
						"pageheading":            "Form Pelanggaran",
						"kejadianKeamananURL":    storeKejadianKeamanan,
						"getKejadianKeamananURL": getKejadianKeamanan,
						"version":                support.Version,
						"idKejadianKeamanan":     idKejadianKeamanan,
						"data":                   userInfo,
					})
				}

				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})
		})

		r.Prefix("keselamatan").Group(func(kecelakaan route.Router) {
			// Pelanggaran
			kecelakaan.Get("", func(ctx http.Context) http.Response {
				userInfo := facades.Cache().Get("user_data")
				dataKeselamatanURL := "/api/kejadian/keselamatan/listKejadianKeselamatan"
				deleteKejadianKeselamatanURL := "/api/kejadian/keselamatan/deleteKejadianKeselamatan"
				
				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("kejadian_keselamatan.tmpl", map[string]interface{}{
						"title":                        "Kejadian Keselamatan",
						"pageheading":                  "Kecelakaan",
						"version":                      support.Version,
						"dataKeselamatanURL":           dataKeselamatanURL,
						"deleteKejadianKeselamatanURL": deleteKejadianKeselamatanURL,
						"data":                         userInfo,
					})
				}

				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})

			kecelakaan.Get("form_kejadian_keselamatan", func(ctx http.Context) http.Response {
				userInfo := facades.Cache().Get("user_data")
				idKejadianKeselamatan := ctx.Request().Query("id_kejadian_keselamatan")
				
				storeKejadianKeselamatan := "/api/kejadian/keselamatan/storeKejadianKeselamatan"
				getKejadianKeselamatan := "/api/kejadian/keselamatan/showDetailKejadianKeselamatan"

				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("form_kejadian_keselamatan.tmpl", map[string]interface{}{
						"title":                     "Form Kejadian Keselamatan",
						"pageheading":               "Form Kecelakaan",
						"kejadianKeselamatanURL":    storeKejadianKeselamatan,
						"getKejadianKeselamatanURL": getKejadianKeselamatan,
						"version":                   support.Version,
						"idKejadianKeselamatan":     idKejadianKeselamatan,
						"data":                      userInfo,
					})

				}
				
				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})
		})
		// facades.Route().Prefix().Static()
	})

}
