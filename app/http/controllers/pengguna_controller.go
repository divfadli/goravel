package controllers

import (
	"goravel/app/models"
	"strconv"

	"github.com/golang-module/carbon/v2"
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
		idParsed := 0

		var dataKaryawan models.Karyawan

		if parseInt, err := strconv.Atoi(id); err == nil {
			idParsed = parseInt
		}

		if id != "new" && idParsed > 0 {
			facades.Orm().Query().With("Jabatan").With("User.Role").Where("user_id=?", id).Get(&dataKaryawan)
		}

		idAtasan := ""

		if dataKaryawan.IDAtasan != nil {
			idAtasan = *dataKaryawan.IDAtasan
		}

		var dataJabatan []models.Jabatan
		facades.Orm().Query().Get(&dataJabatan)

		var dataRole []models.Role
		facades.Orm().Query().Get(&dataRole)

		var dataAtasan []models.Karyawan
		facades.Orm().Query().With("Jabatan").With("User.Role").Get(&dataAtasan)

		return ctx.Response().View().Make("pengguna_detail.tmpl", map[string]interface{}{
			"title":       "Pengguna",
			"pageheading": "Pengguna",
			"version":     support.Version,
			"data":        userInfo,
			"pengguna":    dataKaryawan,
			"jabatan":     dataJabatan,
			"role":        dataRole,
			"isCreate":    id == "new",
			"listAtasan":  dataAtasan,
			"id_atasan":   idAtasan,
		})
	}

	facades.Auth().Logout(ctx)
	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}

func (r *PenggunaController) Store(ctx http.Context) http.Response {
	nama := ctx.Request().Input("nama")
	email := ctx.Request().Input("email")
	nik := ctx.Request().Input("nik")
	jabatan_id, _ := strconv.Atoi(ctx.Request().Input("jabatan"))
	role_id, _ := strconv.Atoi(ctx.Request().Input("role"))
	gender := ctx.Request().Input("gender")
	agama := ctx.Request().Input("agama")
	tanggal_lahir := ctx.Request().Input("tanggal_lahir")
	id_atasan := ctx.Request().Input("atasan")

	var karyawan models.Karyawan
	karyawan.Name = nama
	karyawan.EmpNo = nik
	karyawan.Agama = agama
	karyawan.Gender = gender
	karyawan.TanggalLahir = carbon.Parse(tanggal_lahir).ToDateStruct()
	if id_atasan != "" {
		karyawan.IDAtasan = &id_atasan
	} else {
		karyawan.IDAtasan = nil
	}
	karyawan.JabatanID = jabatan_id

	var user models.Akun
	user.Email = email
	// TODO: password default
	user.Password, _ = facades.Hash().Make("12345678")
	user.RoleID = role_id
	facades.Orm().Query().Save(&user)

	karyawan.UserID = int(user.IDUser)
	karyawan.User = user
	facades.Orm().Query().Save(&karyawan)

	return ctx.Response().Json(200, map[string]interface{}{
		"status":  "success",
		"message": "Data berhasil disimpan",
		"data":    karyawan,
	})
}

func (r *PenggunaController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	nama := ctx.Request().Input("nama")
	email := ctx.Request().Input("email")
	nik := ctx.Request().Input("nik")
	jabatan_id, _ := strconv.Atoi(ctx.Request().Input("jabatan"))
	role_id, _ := strconv.Atoi(ctx.Request().Input("role"))
	gender := ctx.Request().Input("gender")
	agama := ctx.Request().Input("agama")
	tanggal_lahir := ctx.Request().Input("tanggal_lahir")
	id_atasan := ctx.Request().Input("atasan")

	var karyawan models.Karyawan
	facades.Orm().Query().With("Jabatan").With("User.Role").Where("user_id=?", id).First(&karyawan)

	deleteCurrent := karyawan.EmpNo != nik
	oldNIK := karyawan.EmpNo

	karyawan.Name = nama
	karyawan.EmpNo = nik
	karyawan.Agama = agama
	karyawan.Gender = gender
	karyawan.TanggalLahir = carbon.Parse(tanggal_lahir).ToDateStruct()
	if id_atasan != "" {
		karyawan.IDAtasan = &id_atasan
	} else {
		karyawan.IDAtasan = nil
	}
	karyawan.JabatanID = jabatan_id
	karyawan.User.Email = email
	karyawan.User.RoleID = role_id

	facades.Orm().Query().Save(&karyawan)
	facades.Orm().Query().Save(&karyawan.User)

	if deleteCurrent {
		facades.Orm().Query().Delete(models.Karyawan{EmpNo: oldNIK})
	}

	return ctx.Response().Json(200, map[string]interface{}{
		"status":  "success",
		"message": "Data berhasil diupdate",
		"data":    karyawan,
	})
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
