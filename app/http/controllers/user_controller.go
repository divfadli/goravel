package controllers

import (
	"fmt"
	RequestUser "goravel/app/http/requests/user"
	"goravel/app/models"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}
func getFileNameFromUrl(url string) string {
	if url == "" {
		return ""
	}
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func (r *UserController) Index(ctx http.Context) http.Response {
	userInfo := facades.Cache().Get("user_data")

	if userInfo != nil {
		userData := userInfo.(map[string]interface{})

		var dataKaryawan models.Karyawan
		facades.Orm().Query().With("Jabatan").With("User.Role").Where("emp_no=?", userData["nik"]).Get(&dataKaryawan)
		fmt.Println(dataKaryawan)

		if dataKaryawan.Ttd != nil {
			fileName := getFileNameFromUrl(*dataKaryawan.Ttd)
			dataKaryawan.NameFileTtd = &fileName
		}

		return ctx.Response().View().Make("profile_edit.tmpl", map[string]interface{}{
			"title":       "Edit Profile",
			"pageheading": "Edit profile",
			"version":     support.Version,
			"pengguna":    dataKaryawan,
			"data":        userInfo,
		})
	}

	facades.Auth().Logout(ctx)
	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}

// Show method for controller
func (r *UserController) Show(ctx http.Context) http.Response {
	return ctx.Response().View().Make("user.tmpl", map[string]interface{}{
		"title":       "User",
		"pageheading": "User",
		"version":     support.Version,
	})
}

// Store method for controller
func (r *UserController) Store(ctx http.Context) http.Response {
	return ctx.Response().View().Make("user.tmpl", map[string]interface{}{
		"title":       "User",
		"pageheading": "User",
		// "version":     support.Version,
	})
}

// Update method for controller
func (r *UserController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	nama := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	nik := ctx.Request().Input("nik")
	gender := ctx.Request().Input("gender")
	agama := ctx.Request().Input("agama")
	tanggal_lahir := ctx.Request().Input("tanggal_lahir")
	password := ctx.Request().Input("password")

	fmt.Println("masuk")
	fmt.Println(id, nama, email, nik, gender, agama, tanggal_lahir, password)
	fmt.Println("----------------------")

	var karyawan models.Karyawan
	facades.Orm().Query().With("Jabatan").With("User.Role").Where("user_id=?", id).First(&karyawan)

	// Update karyawan data
	karyawan.Name = nama
	karyawan.Gender = gender
	karyawan.Agama = agama
	karyawan.TanggalLahir = carbon.Parse(tanggal_lahir).ToDateStruct()

	// Handle signature file upload if present
	if file, err := ctx.Request().File("signature"); err == nil {
		if karyawan.Ttd != nil {
			oldFileName := getFileNameFromUrl(*karyawan.Ttd)
			fmt.Println(oldFileName)
			facades.Storage().Delete("Signatures/" + nik + "/" + oldFileName)
		}
		// newfileIdentificator := buildFileIdentificator(file.GetClientOriginalName())
		folder, err := facades.Storage().PutFileAs("Signatures/"+nik+"/", file, file.GetClientOriginalName())
		if err != nil {
			return Error(ctx, http.StatusInternalServerError, err.Error())
		}

		// Store the URL in a variable first
		fileUrl := facades.Storage().Url(folder)
		karyawan.Ttd = &fileUrl

		fmt.Println("File uploaded successfully")
	}

	// Update user data
	karyawan.User.Email = email
	if password != "" {
		hashedPassword, _ := facades.Hash().Make(password)
		karyawan.User.Password = hashedPassword
	}

	// // Save changes
	facades.Orm().Query().Save(&karyawan)
	facades.Orm().Query().Save(&karyawan.User)

	return ctx.Response().Json(200, map[string]interface{}{
		"status":  "success",
		"message": "Profile berhasil diupdate",
		"data":    karyawan,
	})
}

// Destroy method for controller
func (r *UserController) Destroy(ctx http.Context) http.Response {
	return ctx.Response().View().Make("user.tmpl", map[string]interface{}{
		"title":       "User",
		"pageheading": "User",
		// "version":     support.Version,
	})
}

func (r *UserController) GetRole(ctx http.Context) http.Response {
	var role []models.Role

	if err := facades.Orm().Query().Find(&role); err != nil || role == nil {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_role": role,
	})
}

func (r *UserController) StoreRole(ctx http.Context) http.Response {
	var req RequestUser.PostRole

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var role models.Role
	role.Name = req.Name

	if err := facades.Orm().Query().Create(&role); err != nil {
		return ErrorSystem(ctx, "Data Gagal Save!!")
	}

	return Success(ctx, http.Json{
		"message": "Data Berhasil Disimpan!!",
	})
}

func (r *UserController) StoreRolePegawai(ctx http.Context) http.Response {
	var req RequestUser.PostRolePegawai

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var akun models.Akun
	var pesan string

	if req.IDUser != 0 {
		facades.Orm().Query().Where("id_user=?", req.IDUser).First(&akun)

		akun.RoleID = req.RoleID

		if err := facades.Orm().Query().Save(&akun); err != nil {
			return ErrorSystem(ctx, "Data Gagal Save!!")
		}
		pesan = "Data Berhasil Disimpan!!"
	} else {
		facades.Orm().Query().
			Join("inner join public.karyawan kry on id_user = kry.user_id").
			Where("kry.emp_no=?", req.EmpNo).First(&akun)

		akun.RoleID = req.RoleID

		if err := facades.Orm().Query().Save(&akun); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan!!")
		}
		pesan = "Data Berhasil Ditambah!!"
	}

	return Success(ctx, http.Json{
		"message": pesan,
	})
}

func (r *UserController) FindPegawai(ctx http.Context) http.Response {
	var req RequestUser.GetPegawai

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var pegawai []models.Karyawan

	facades.Orm().Query().
		Join("inner join public.jabatan jb on jb.id_jabatan = jabatan_id").
		Join("inner join public.akun acc on acc.id_user = user_id").
		With("Jabatan").With("User").Where("(acc.role_id IS NULL) OR (acc.role_id = 0)").Find(&pegawai)

	return Success(ctx, http.Json{
		"data_karyawan": pegawai,
	})
}

func (r *UserController) Login(ctx http.Context) http.Response {
	var req RequestUser.Login

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var user models.Akun
	err := facades.Orm().Query().Where("email", req.Email).First(&user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("goravel", "User").With(map[string]any{
			"error": err.Error(),
		}).Info("Failed to query user")
		return ErrorSystem(ctx, "Failed to query user")
	}
	a, _ := facades.Crypt().DecryptString(user.Password)
	fmt.Println(a)

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

	var role models.MyRole
	facades.Orm().Query().Table("public.role").Where("id_role=?", user.RoleID).Scan(&role)
	// facades.Orm().Query().Table("public.role").Where("emp_no=?", dataKaryawan.EmpNo).Scan(&roleKaryawan)
	// defaultRole := models.MyRole{
	// 	Name: "karyawan",
	// }
	// if role == nil {
	// 	role = append(role, defaultRole)
	// }

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
		"is_staff":    role.Name == "Staff Input",
		"is_admin":    role.Name == "Admin",
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
