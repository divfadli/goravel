package controllers

import (
	"fmt"
	kejadianKeamanan "goravel/app/http/requests/kejadian_keamanan"
	"goravel/app/models"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type KejadianKeamananController struct {
	//Dependent services
}

func NewKejadianKeamananController() *KejadianKeamananController {
	return &KejadianKeamananController{
		//Inject services
	}
}

// func (r *KejadianKeamananController) Index(ctx http.Context) http.Response {
// 	return nil
// }

func (r *KejadianKeamananController) StoreKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.PostKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	// upload_file, handler, err := ctx.Request().Origin().FormFile("file")

	var data_keamanan models.KejadianKeamanan
	var kejadian models.JenisKejadian
	var pesan string

	// Update
	if req.IdKejadianKeamanan != 0 {
		if err := facades.Orm().Query().Table("public.kejadian_keamanan").
			Join(`inner join public.kejadian k ON k.id_jenis_kejadian = public.kejadian_keamanan.jenis_kejadian_id`).
			Where("k.klasifikasi_name = ? AND id_kejadian_keamanan = ? AND k.id_jenis_kejadian=?", "Keamanan Laut", req.IdKejadianKeamanan, req.JenisKejadianId).
			First(&data_keamanan); err != nil || data_keamanan.IdKejadianKeamanan == 0 {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keamanan Tidak Ditemukan!!")
		}

		if data_keamanan.IsLocked {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keamanan Terkunci!!")
		} else {
			data_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
			fmt.Println(data_keamanan.Tanggal)
			data_keamanan.JenisKejadianId = req.JenisKejadianId
			data_keamanan.NamaKapal = req.NamaKapal
			data_keamanan.SumberBerita = req.SumberBerita
			data_keamanan.LinkBerita = req.LinkBerita
			data_keamanan.LokasiKejadian = req.LokasiKejadian
			data_keamanan.Muatan = req.Muatan
			if req.Asal != "" {
				data_keamanan.Asal = &req.Asal
			}
			if req.Bendera != "" {
				data_keamanan.Bendera = &req.Bendera
			}
			if req.Tujuan != "" {
				data_keamanan.Tujuan = &req.Tujuan
			}
			data_keamanan.Latitude = req.Latitude
			data_keamanan.Longitude = req.Longitude
			data_keamanan.KategoriSumber = req.KategoriSumber
			data_keamanan.TindakLanjut = req.TindakLanjut
			if req.IMOKapal != "" {
				data_keamanan.IMOKapal = &req.IMOKapal
			}
			if req.MMSIKapal != "" {
				data_keamanan.MMSIKapal = &req.MMSIKapal
			}
			data_keamanan.InformasiKategori = req.InformasiKategori
			data_keamanan.Zona = req.Zona
			data_keamanan.CreatedBy = req.Nik

			if err := facades.Orm().Query().Save(&data_keamanan); err != nil {
				return ErrorSystem(ctx, "Data Gagal Diubah")
			}

			pesan = "Data Berhasil Diubah"
		}
	} else {
		data_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
		fmt.Println(data_keamanan.Tanggal)
		data_keamanan.JenisKejadianId = req.JenisKejadianId
		data_keamanan.NamaKapal = req.NamaKapal
		data_keamanan.SumberBerita = req.SumberBerita
		data_keamanan.LinkBerita = req.LinkBerita
		data_keamanan.LokasiKejadian = req.LokasiKejadian
		data_keamanan.Muatan = req.Muatan
		if req.Asal != "" {
			data_keamanan.Asal = &req.Asal
		}
		if req.Bendera != "" {
			data_keamanan.Bendera = &req.Bendera
		}
		if req.Tujuan != "" {
			data_keamanan.Tujuan = &req.Tujuan
		}
		data_keamanan.Latitude = req.Latitude
		data_keamanan.Longitude = req.Longitude
		data_keamanan.KategoriSumber = req.KategoriSumber
		data_keamanan.TindakLanjut = req.TindakLanjut
		if req.IMOKapal != "" {
			data_keamanan.IMOKapal = &req.IMOKapal
		}
		if req.MMSIKapal != "" {
			data_keamanan.MMSIKapal = &req.MMSIKapal
		}
		data_keamanan.InformasiKategori = req.InformasiKategori
		data_keamanan.Zona = req.Zona
		data_keamanan.CreatedBy = req.Nik

		if err := facades.Orm().Query().
			Where("klasifikasi_name = ? AND id_jenis_kejadian = ?", "Keamanan Laut", req.JenisKejadianId).
			First(&kejadian); err != nil || kejadian.IDJenisKejadian == "" {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keamanan Tidak Ditemukan!!")
		}

		if err := facades.Orm().Query().Create(&data_keamanan); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}

		pesan = "Data Berhasil Ditambahkan"
	}
	return Success(ctx, http.Json{
		"Success": pesan,
	})
}

func (r *KejadianKeamananController) ListKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.ListKeamanan
	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var rekap_keamanan []models.KejadianKeamanan
	facades.Orm().Query().With("JenisKejadian").Where("nama_kejadian=?", "Pelanggaran Wilayah").Find(&rekap_keamanan)
	// var rekap_keamanan []models.KejadianKeamanan
	// query := facades.Orm().Query().Table("public.kejadian_keamanan").
	// 	Join(`inner join public.jenis_kejadian k ON k.id_jenis_kejadian = public.kejadian_keamanan.jenis_kejadian_id`)

	// if req.Key != "" {
	// 	query = query.Where(`(
	// 							lower(k.nama_kejadian) like lower(?) OR
	// 						 	lower(public.kejadian_keamanan.nama_kapal) like lower(?) OR
	// 						 	lower(public.kejadian_keamanan.lokasi_kejadian) like lower(?)
	// 						)`, "%"+req.Key+"%", "%"+req.Key+"%", "%"+req.Key+"%")
	// }

	// if req.Zona != "" {
	// 	query = query.Where("lower(public.kejadian_keamanan.zona::text) like lower(?)", "%"+req.Zona+"%")
	// }

	// tanggal_awal, _ := time.Parse(time.DateOnly, req.TanggalAwal)
	// tanggal_akhir, _ := time.Parse(time.DateOnly, req.TanggalAkhir)

	// query = query.Where("public.kejadian_keamanan.tanggal between ? AND ?",
	// 	tanggal_awal, tanggal_akhir)

	// if err := query.Order("public.kejadian_keamanan.tanggal asc").Find(&)Scan(&rekap_keamanan); err != nil || len(rekap_keamanan) == 0 {
	// 	return ErrorSystem(ctx, "Data Tidak Ada")
	// }

	// var kejadian_keamanan models.KejadianKeamanan
	// facades.Orm().Query().Association("JenisKejadian").Find(&kejadian_keamanan)
	// var jenis_kejadians []models.JenisKejadian

	// facades.Orm().Query().Where("public.kejadian_keamanan.jenis_kejadian_id =? ", "TYP-000024").
	// 	Find(&jenis_kejadians)

	// if err := facades.Orm().Query().Join(`inner join public.jenis_kejadian k ON k.id_jenis_kejadian = public.kejadian_keamanan.jenis_kejadian_id`).
	// 	Where("public.kejadian_keamanan.jenis_kejadian_id =? ", "TYP-000024").
	// 	With("JenisKejadian").Find(&rekap_keamanan).
	// 	Error; err != nil || len(rekap_keamanan) == 0 {
	// 	return ErrorSystem(ctx, "Data Tidak Ada")
	// }
	// if err := facades.Orm().Query().Table("public.kejadian_keamanan").
	// 	Join(`inner join public.jenis_kejadian k ON k.id_jenis_kejadian = public.kejadian_keamanan.jenis_kejadian_id`).
	// 	Where("public.kejadian_keamanan.jenis_kejadian_id=?", "TYP-000024").
	// 	Association("JenisKejadianId").Find(&rekap_keamanan).Error; err != nil || len(rekap_keamanan) == 0 {
	// 	return ErrorSystem(ctx, "Data Tidak Ada")
	// }

	return Success(ctx, http.Json{
		"data": rekap_keamanan,
	})
}

// func (r *KejadianKeamananController) ShowDetailKejadianKeamanan(ctx http.Context) http.Response {
// 	var req kejadianKeamanan.GetKeamanan

// 	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
// 		return sanitize
// 	}

// 	var data_keamanan models.DetailKejadianKeamanan
// 	query := facades.Orm().Query().Table("public.kejadian_keamanan").
// 		Join(`inner join public.kejadian k ON k.id_jenis_kejadian = public.kejadian_keamanan.jenis_kejadian_id`)

// 	query = query.Where("public.kejadian_keamanan.id_kejadian_keamanan", req.IdKejadianKeamanan)

// 	if err := query.First(&data_keamanan); err != nil || data_keamanan.IdKejadianKeamanan == 0 {
// 		return ErrorSystem(ctx, "Data Tidak Ada")
// 	}

// 	return Success(ctx, http.Json{
// 		"data": data_keamanan,
// 	})
// }

func (r *KejadianKeamananController) DeleteKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.GetKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data_keamanan []models.KejadianKeamanan
	if x, err := facades.Orm().Query().Where("id_kejadian_keamanan = ? AND is_locked is false", req.IdKejadianKeamanan).Delete(&data_keamanan); err != nil || x.RowsAffected == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada / Data Tidak Dapat Dihapus")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
