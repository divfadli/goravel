package controllers

import (
	"goravel/app/http/requests/kejadian"
	"goravel/app/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type KejadianController struct {
	//Dependent services
}

type DetailKejadian struct {
	IDTypeKejadian   string
	JenisPelanggaran string
	KlasifikasiId    int
	NamaKlasifikasi  string
}

func NewKejadianController() *KejadianController {
	return &KejadianController{
		//Inject services
	}
}

func (r *KejadianController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *KejadianController) PostKejadian(ctx http.Context) http.Response {
	var req kejadian.PostKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}
	// if sanitize := SanitizePost(ctx, &req); sanitize != nil {
	// 	return sanitize
	// }

	var kejadian models.Kejadian
	var klasifikasi *models.KlasifikasiKejadian

	// Update data
	if req.IdTypeKejadian != "" {
		if err := facades.Orm().Query().Where("id_type_kejadian", req.IdTypeKejadian).First(&kejadian); err != nil || kejadian.IDTypeKejadian == "" {
			return ErrorSystem(ctx, "Data Tidak Ada")
		}

		kejadian.JenisPelanggaran = req.JenisPelanggaran

		if err := facades.Orm().Query().Save(&kejadian); err != nil {
			return ErrorSystem(ctx, "Data Gagal Diubah")
		}
		return Success(ctx, http.Json{
			"Success": "Data Berhasil Diubah",
		})
	} else {
		// Create data

		if err := facades.Orm().Query().Where("id_klasifikasi", req.KlasifikasiId).First(&klasifikasi); err != nil || klasifikasi.IDKlasifikasi == 0 {
			return ErrorSystem(ctx, "Data Tidak Ada")
		}

		// Seed the random number generator
		rand.Seed(time.Now().UnixNano())

		kejadian.IDTypeKejadian = "TYP-000" + strconv.Itoa(rand.Intn(1000))
		kejadian.JenisPelanggaran = req.JenisPelanggaran
		kejadian.KlasifikasiID = uint(req.KlasifikasiId)

		if err := facades.Orm().Query().Create(&kejadian); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}
		return Success(ctx, http.Json{
			"Success": "Data Berhasil Ditambahkan",
		})
	}
}

func (r *KejadianController) ListKejadian(ctx http.Context) http.Response {
	var req kejadian.ListKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}
	// if sanitize := SanitizePost(ctx, &req); sanitize != nil {
	// 	return sanitize
	// }
	var data []struct {
		IDTypeKejadian   string
		JenisPelanggaran string
		KlasifikasiId    int
		NamaKlasifikasi  string
	}

	if err := facades.Orm().Query().Table("rekapitulasi.kejadian").
		Join("join rekapitulasi.klasifikasi_kejadian ON rekapitulasi.kejadian.klasifikasi_id = rekapitulasi.klasifikasi_kejadian.id_klasifikasi").
		Where("klasifikasi_id", req.KlasifikasiId).
		Scan(&data); err != nil {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_kejadian": data,
	})
}

func (r *KejadianController) ShowDetailKejadian(ctx http.Context) http.Response {
	var req kejadian.GetKejadian

	// if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
	// 	return SanitizeGet(ctx, chekRequestErr)
	// }
	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var data *DetailKejadian

	if err := facades.Orm().Query().Table("rekapitulasi.kejadian").
		Join("join rekapitulasi.klasifikasi_kejadian ON rekapitulasi.kejadian.klasifikasi_id = rekapitulasi.klasifikasi_kejadian.id_klasifikasi").
		Where("id_type_kejadian", req.IdTypeKejadian).
		Scan(&data); err != nil {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_kejadian": data,
	})
}

func (r *KejadianController) DeleteKejadian(ctx http.Context) http.Response {
	var req kejadian.GetKejadian

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}
	// if sanitize := SanitizePost(ctx, &req); sanitize != nil {
	// 	return sanitize
	// }

	var kejadian []models.Kejadian
	if data, err := facades.Orm().Query().Where("id_type_kejadian", req.IdTypeKejadian).Delete(&kejadian); err != nil || data.RowsAffected == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
