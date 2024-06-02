package controllers

import (
	jenisKejadian "goravel/app/http/requests/jenis_kejadian"
	"goravel/app/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type JenisKejadianController struct {
	//Dependent services
}

func NewJenisKejadianController() *JenisKejadianController {
	return &JenisKejadianController{
		//Inject services
	}
}

func (r *JenisKejadianController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *JenisKejadianController) PostKejadian(ctx http.Context) http.Response {
	var req jenisKejadian.PostKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}
	// if sanitize := SanitizePost(ctx, &req); sanitize != nil {
	// 	return sanitize
	// }

	var data models.JenisKejadian
	var pesan string

	// Update data
	if req.IdJenisKejadian != "" {
		if err := facades.Orm().Query().Where("id_jenis_kejadian", req.IdJenisKejadian).First(&data); err != nil || data.IDJenisKejadian == "" {
			return ErrorSystem(ctx, "Data Tidak Ada")
		}

		data.NamaKejadian = req.NamaKejadian

		if err := facades.Orm().Query().Save(&data); err != nil {
			return ErrorSystem(ctx, "Data Gagal Diubah")
		}

		pesan = "Data Berhasil Diubah"
	} else {
		// Create data

		// Seed the random number generator
		rand.Seed(time.Now().UnixNano())

		data.IDJenisKejadian = "TYP-000" + strconv.Itoa(rand.Intn(1000))
		data.NamaKejadian = req.NamaKejadian
		data.KlasifikasiName = req.KlasifikasiName

		if err := facades.Orm().Query().Create(&data); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}

		pesan = "Data Berhasil Ditambahkan"
	}

	return Success(ctx, http.Json{
		"Success": pesan,
	})
}

func (r *JenisKejadianController) ListKejadian(ctx http.Context) http.Response {
	var req jenisKejadian.ListKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}
	// if sanitize := SanitizePost(ctx, &req); sanitize != nil {
	// 	return sanitize
	// }
	var data []models.JenisKejadian

	if err := facades.Orm().Query().Where("klasifikasi_name", req.KlasifikasiName).Find(&data); err != nil || len(data) == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_kejadian": data,
	})
}

func (r *JenisKejadianController) ShowDetailKejadian(ctx http.Context) http.Response {
	var req jenisKejadian.GetKejadian

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var data models.JenisKejadian

	if err := facades.Orm().Query().Where("id_jenis_kejadian=?", req.IdJenisKejadian).First(&data); err != nil || data.IDJenisKejadian == "" {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_kejadian": data,
	})
}

func (r *JenisKejadianController) DeleteKejadian(ctx http.Context) http.Response {
	var req jenisKejadian.GetKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data []models.JenisKejadian
	if x, err := facades.Orm().Query().Where("id_jenis_kejadian", req.IdJenisKejadian).Delete(&data); err != nil || x.RowsAffected == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
