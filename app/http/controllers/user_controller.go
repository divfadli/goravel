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

	var user models.Akun
	err := facades.Orm().Query().Where("email", req.Email).First(&user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
			"error": err.Error(),
		}).Info("Failed to query user")
		return ErrorSystem(ctx, "Failed to query user")
	}

	if user.IDUser == 0 || !facades.Hash().Check(req.Password, user.Password) {
		// return ctx.Response().Redirect(http.StatusFound, "/login")
		return Error(ctx, http.StatusUnauthorized, "Username atau Password salah")
	}

	var dataKaryawan models.Karyawan
	facades.Orm().Query().Where("user_id", user.IDUser).First(&dataKaryawan)

	var count int64
	// facades.Orm().Query().Table("public.karyawan").Where("sup_pos_id=?", dataKaryawan.PosID).Count(&count)
	facades.Orm().Query().Table("public.karyawan").Where("id_atasan", dataKaryawan.EmpNo).Count(&count)

	is_superior := true
	if count == 0 {
		is_superior = false
	}

	var role []models.MyRole
	facades.Orm().Query().Table("public.role").Where("id_role=?", user.RoleID).Scan(&role)
	// facades.Orm().Query().Table("public.role").Where("emp_no=?", dataKaryawan.EmpNo).Scan(&roleKaryawan)
	defaultRole := models.MyRole{
		Name: "karyawan",
	}
	if role == nil {
		role = append(role, defaultRole)
	}

	token_access, loginErr := facades.Auth().LoginUsingID(ctx, user.IDUser)
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

	// Assuming Cache().Put and Cache().Get work correctly
	cachedData := map[string]interface{}{
		"token":       token,
		"name":        dataKaryawan.Name,
		"email":       user.Email,
		"nik":         dataKaryawan.EmpNo,
		"is_superior": is_superior,
		"role":        role,
	}
	facades.Cache().Put("user_data", cachedData, 2*time.Hour)

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

	var user models.Akun

	facades.Orm().Query().Where("email=?", req.Email).First(&user)
	if user.IDUser != 0 {
		return Error(ctx, http.StatusForbidden, "Email already exists")
	}

	user.Password, _ = facades.Hash().Make(req.Password)
	user.Email = req.Email

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ErrorSystem(ctx, "Data Gagal Ditambahkan")
	}
	return Success(ctx, http.Json{
		"Success": "Data Berhasil Ditambahkan",
	})
}

// func (r *UserController) Info(ctx http.Context) http.Response {
// 	var user models.Akun
// 	err := facades.Auth().User(ctx, &user)
// 	if err != nil {
// 		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
// 			"error": err.Error(),
// 		}).Info("Failed to obtain user information")
// 		return ErrorSystem(ctx, "Failed to obtain user information")
// 	}

// 	return Success(ctx, http.Json{
// 		"id":       user.IDUser,
// 		"role":     []string{"admin"},
// 		"username": user.Username,
// 		"name":     user.Name,
// 		"email":    user.Email,
// 		"nik":      user.Nik,
// 	})
// }
