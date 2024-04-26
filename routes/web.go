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

	facades.Route().Get("/map", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("map.tmpl", map[string]any{
			"version": support.Version,
		})
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
	facades.Route().Get("register", func(ctx http.Context) http.Response {
		registerURL := "/api/user/register"
		return ctx.Response().View().Make("register.tmpl", map[string]any{
			"registerURL": registerURL,
			"version":     support.Version,
		})
	})

	// Dashboard
	facades.Route().Get("dashboard", func(ctx http.Context) http.Response {
		// Retrieve cached user data
		userInfo := facades.Cache().Get("user_data")
		dataKeamananURL := "/api/rekap/keamanan/listRekapKeamanan"
		// fmt.Println(userInfo)

		// // Check if data is available in cache
		if userInfo != nil {
			return ctx.Response().View().Make("index.tmpl", map[string]interface{}{
				"title":           "Dashboard | Rekapitulasi",
				"pageheading":     "Dashboard",
				"version":         support.Version,
				"dataKeamananURL": dataKeamananURL,
				"data":            userInfo,
			})
		}

		facades.Auth().Logout(ctx)
		// For instance, you might redirect the user to the login page
		return ctx.Response().Redirect(http.StatusFound, "/login")
	})

	// Rekap Data Kejadian
	facades.Route().Prefix("rekap_data_kejadian").Group(func(r route.Router) {
		r.Prefix("pelanggaran").Group(func(pelanggaran route.Router) {
			// Pelanggaran
			pelanggaran.Get("", func(ctx http.Context) http.Response {
				// Retrieve cached user data
				userInfo := facades.Cache().Get("user_data")

				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("index.tmpl", map[string]interface{}{
						"title":       "Pelanggaran | Rekapitulasi",
						"pageheading": "Pelanggaran",
						"version":     support.Version,
						"data":        userInfo,
					})
				}

				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})

			pelanggaran.Get("/form_rekap_pel", func(ctx http.Context) http.Response {
				return ctx.Response().View().Make("form_rekap_pel.tmpl", map[string]interface{}{
					"title":       "Form Rekap | Pelanggaran",
					"pageheading": "Form Pelanggaran",
					"version":     support.Version,
				})
			})
		})

		r.Prefix("rekap_kecelakaan").Group(func(kecelakaan route.Router) {
			// Pelanggaran
			kecelakaan.Get("", func(ctx http.Context) http.Response {
				// Retrieve cached user data
				userInfo := facades.Cache().Get("user_data")

				// Check if data is available in cache
				if userInfo != nil {
					return ctx.Response().View().Make("index.tmpl", map[string]interface{}{
						"title":       "Kecelakaan | Rekapitulasi",
						"pageheading": "Kecelakaan",
						"version":     support.Version,
						"data":        userInfo,
					})
				}

				// For instance, you might redirect the user to the login page
				return ctx.Response().Redirect(http.StatusFound, "/login")
			})
		})
		// facades.Route().Prefix().Static()
	})

}
