package controllers

import (
	RequestUser "goravel/app/http/requests/user"
	"goravel/app/models"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}
func (r *UserController) Login(ctx http.Context) http.Response {
	var req RequestUser.Login

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var user models.User
	err := facades.Orm().Query().Where("username", req.Username).First(&user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
			"error": err.Error(),
		}).Info("Failed to query user")
		return ErrorSystem(ctx, "Failed to query user")
	}

	if user.ID == 0 || !facades.Hash().Check(req.Password, user.Password) {
		// return ctx.Response().Redirect(http.StatusFound, "/login")
		return Error(ctx, http.StatusUnauthorized, "Wrong username or password")
	}

	token_access, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
			"error": err.Error(),
		}).Info("Login failed")
		return ErrorSystem(ctx, loginErr.Error())
	}
	access, _ := facades.Auth().Parse(ctx, token_access)

	token := map[string]any{
		"access_token": token_access,
		"expires_in":   access.ExpireAt,
	}

	// if user.Type != "" {
	// 	userData.Role = append(userData.Role, "karyawan")
	// }

	// Assuming Cache().Put and Cache().Get work correctly
	cachedData := map[string]interface{}{
		"token": token,
		"name":  user.Name,
		"email": user.Email,
		"nik":   user.Nik,
		// "role":  []string{"admin"}, // Assuming this is correct
		"role": user.Type,
	}
	facades.Cache().Put("user_data", cachedData, 2*time.Minute)

	// Log the cached data for debugging
	// facades.Log().Info("Cached user data:", cachedData)

	// Redirect to the index page
	// return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	return ctx.Response().Success().Json(http.Json{
		"data": cachedData,
	})
}
func (r *UserController) Register(ctx http.Context) http.Response {
	var req RequestUser.Register

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var user models.User

	facades.Orm().Query().Where("username=? OR nik=?", req.Username, req.Nik).First(&user)
	if user.ID != 0 {
		return Error(ctx, http.StatusForbidden, "Username or Nik already exists")
	}

	if req.Type != "" {
		user.Type = req.Type
	}
	user.Username = req.Username
	user.Password, _ = facades.Hash().Make(req.Password)
	user.Email = req.Email
	user.Nik = req.Nik
	user.Name = req.Name
	user.Type = req.Type

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ErrorSystem(ctx, "Data Gagal Ditambahkan")
	}
	return Success(ctx, http.Json{
		"Success": "Data Berhasil Ditambahkan",
	})
}

func (r *UserController) Info(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
			"error": err.Error(),
		}).Info("Failed to obtain user information")
		return ErrorSystem(ctx, "Failed to obtain user information")
	}

	return Success(ctx, http.Json{
		"id":       user.ID,
		"role":     []string{"admin"},
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"nik":      user.Nik,
	})
}
