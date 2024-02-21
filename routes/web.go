package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Get("index", func(ctx http.Context) http.Response {
		// Retrieve cached user data
		userInfo := facades.Cache().Get("user_data")

		// Log the retrieved data for debugging
		// facades.Log().Info("Retrieved user data from cache:", user_info)

		// Check if data is available in cache
		if userInfo != nil {
			// Cast the retrieved data to the correct type
			// userData, ok := cachedData.(map[string]interface{})
			// if !ok {
			// 	// Handle incorrect data type
			// 	return ctx.Response().Json(http.StatusInternalServerError, "Internal Server Error")
			// }

			// Pass the retrieved data to the view
			return ctx.Response().View().Make("welcome.tmpl", map[string]interface{}{
				"version": support.Version,
				"data":    userInfo,
			})
		}

		// Handle case when no data is found in the cache
		// For instance, you might redirect the user to the login page
		return ctx.Response().Redirect(http.StatusFound, "/login")
	})
	facades.Route().Get("login", func(ctx http.Context) http.Response {
		loginURL := "/api/user/login"
		return ctx.Response().View().Make("login.php", map[string]any{
			"loginURL": loginURL,
			"version":  support.Version,
		})
	})
}
