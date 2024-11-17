package controllers

import (
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

type PenggunaController struct {
	//Dependent services
}

func NewPenggunaController() *PenggunaController {
	return &PenggunaController{
		//Inject services
	}
}

func (r *PenggunaController) Index(ctx http.Context) http.Response {
	userInfo := facades.Cache().Get("user_data")

	var dataKaryawan []models.Karyawan
	facades.Orm().Query().With("Jabatan").With("User.Role").Get(&dataKaryawan)

	if userInfo != nil {
		return ctx.Response().View().Make("pengguna.tmpl", map[string]interface{}{
			"title":       "Pengguna",
			"pageheading": "Pengguna",
			"version":     support.Version,
			"data":        userInfo,
			"pengguna":    dataKaryawan,
		})
	}

	facades.Auth().Logout(ctx)
	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}

func (r *PenggunaController) Show(ctx http.Context) http.Response {
	userInfo := facades.Cache().Get("user_data")

	if userInfo != nil {
		id := ctx.Request().Route("id")
		var dataKaryawan models.Karyawan
		facades.Orm().Query().With("Jabatan").With("User.Role").Where("user_id=?", id).Get(&dataKaryawan)
		var dataJabatan []models.Jabatan
		facades.Orm().Query().Get(&dataJabatan)
		var dataRole []models.Role
		facades.Orm().Query().Get(&dataRole)
		return ctx.Response().View().Make("pengguna_detail.tmpl", map[string]interface{}{
			"title":       "Pengguna",
			"pageheading": "Pengguna",
			"version":     support.Version,
			"data":        userInfo,
			"pengguna":    dataKaryawan,
			"jabatan":     dataJabatan,
			"role":        dataRole,
		})
	}

	facades.Auth().Logout(ctx)
	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}

func (r *PenggunaController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *PenggunaController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *PenggunaController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var karyawan models.Karyawan
	facades.Orm().Query().With("User").Where("user_id=?", id).First(&karyawan)
	facades.Orm().Query().Delete(&karyawan)
	facades.Orm().Query().Delete(&karyawan.User)
	return ctx.Response().Json(200, map[string]interface{}{
		"status":  "success",
		"message": "Data berhasil dihapus",
		"data":    karyawan,
	})
}
