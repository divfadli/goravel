package middleware

import (
	"errors"

	"github.com/goravel/framework/contracts/http"

	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/facades"
)

func Jwt() http.Middleware {
	return func(ctx http.Context) {
		token := ctx.Request().Header("Authorization", ctx.Request().Header("Sec-WebSocket-Protocol"))
		if len(token) == 0 {
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    401,
				"message": "Silahkan Login!!",
			})
			return
		}

		// Authentication JWT
		if _, err := facades.Auth().Parse(ctx, token); err != nil {
			if errors.Is(err, auth.ErrorTokenExpired) {
				token, err = facades.Auth().Refresh(ctx)
				if err != nil {
					// Refresh Token Expired
					ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
						"code":    401,
						"message": "Login telah kedaluwarsa",
					})
					return
				}

				token = "Bearer " + token
			} else {
				ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
					"code":    401,
					"message": "Login telah kedaluwarsa",
				})
				return
			}
		}

		ctx.Response().Header("Authorization", token)

		ctx.Request().Next()
	}
}
