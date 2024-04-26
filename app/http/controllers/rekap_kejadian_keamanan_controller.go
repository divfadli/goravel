package controllers

import (
	"fmt"
	"goravel/app/http/requests/rekap"
	"goravel/app/models"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type RekapKejadianKeamananController struct {
	//Dependent services
}

func NewRekapKejadianKeamananController() *RekapKejadianKeamananController {
	return &RekapKejadianKeamananController{
		//Inject services
	}
}

func (r *RekapKejadianKeamananController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *RekapKejadianKeamananController) StoreRekapKeamanan(ctx http.Context) http.Response {
	var req rekap.PostKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	// upload_file, handler, err := ctx.Request().Origin().FormFile("file")

	var rekap_keamanan models.RekapKejadianDataKeamanan
	var kejadian models.Kejadian

	if req.IdRekapKeamanan != 0 {
		if err := facades.Orm().Query().Table("rekapitulasi.rekap_kejadian_data_keamanan").
			Join(`inner join rekapitulasi.kejadian k ON k.id_type_kejadian = rekapitulasi.rekap_kejadian_data_keamanan.type_kejadian_id
				inner join rekapitulasi.klasifikasi_kejadian kk ON k.klasifikasi_id = kk.id_klasifikasi`).
			Where("kk.id_klasifikasi = ? AND id_rekap_keamanan = ? AND k.id_type_kejadian=?", "1", req.IdRekapKeamanan, req.TypeKejadianId).
			First(&rekap_keamanan); err != nil || rekap_keamanan.IdRekapKeamanan == 0 {
			return Error(ctx, http.StatusNotFound, "Rekap Keamanan Not Found")
		}

		if rekap_keamanan.IsLocked {
			return Error(ctx, http.StatusNotFound, "Rekap Keamanan Locked")
		} else {
			rekap_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
			fmt.Println(rekap_keamanan.Tanggal)
			rekap_keamanan.TypeKejadianId = req.TypeKejadianId
			rekap_keamanan.NamaKapal = req.NamaKapal
			rekap_keamanan.SumberBerita = req.SumberBerita
			rekap_keamanan.LinkBerita = req.LinkBerita
			rekap_keamanan.LokasiKejadian = req.LokasiKejadian
			rekap_keamanan.Muatan = req.Muatan
			if req.Asal != "" {
				rekap_keamanan.Asal = &req.Asal
			}
			if req.Bendera != "" {
				rekap_keamanan.Bendera = &req.Bendera
			}
			if req.Tujuan != "" {
				rekap_keamanan.Tujuan = &req.Tujuan
			}
			rekap_keamanan.Latitude = req.Latitude
			rekap_keamanan.Longitude = req.Longitude
			rekap_keamanan.KategoriSumber = req.KategoriSumber
			rekap_keamanan.TindakLanjut = req.TindakLanjut
			if req.IMOKapal != "" {
				rekap_keamanan.IMOKapal = &req.IMOKapal
			}
			if req.MMSIKapal != "" {
				rekap_keamanan.MMSIKapal = &req.MMSIKapal
			}
			rekap_keamanan.InformasiKategori = req.InformasiKategori
			rekap_keamanan.Zona = req.Zona
			rekap_keamanan.CreatedBy = req.Nik

			if err := facades.Orm().Query().Save(&rekap_keamanan); err != nil {
				return ErrorSystem(ctx, "Data Gagal Diubah")
			}

			return Success(ctx, http.Json{
				"Success": "Data Berhasil Diubah",
			})
		}
	} else {
		rekap_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
		fmt.Println(rekap_keamanan.Tanggal)
		rekap_keamanan.TypeKejadianId = req.TypeKejadianId
		rekap_keamanan.NamaKapal = req.NamaKapal
		rekap_keamanan.SumberBerita = req.SumberBerita
		rekap_keamanan.LinkBerita = req.LinkBerita
		rekap_keamanan.LokasiKejadian = req.LokasiKejadian
		rekap_keamanan.Muatan = req.Muatan
		if req.Asal != "" {
			rekap_keamanan.Asal = &req.Asal
		}
		if req.Bendera != "" {
			rekap_keamanan.Bendera = &req.Bendera
		}
		if req.Tujuan != "" {
			rekap_keamanan.Tujuan = &req.Tujuan
		}
		rekap_keamanan.Latitude = req.Latitude
		rekap_keamanan.Longitude = req.Longitude
		rekap_keamanan.KategoriSumber = req.KategoriSumber
		rekap_keamanan.TindakLanjut = req.TindakLanjut
		if req.IMOKapal != "" {
			rekap_keamanan.IMOKapal = &req.IMOKapal
		}
		if req.MMSIKapal != "" {
			rekap_keamanan.MMSIKapal = &req.MMSIKapal
		}
		rekap_keamanan.InformasiKategori = req.InformasiKategori
		rekap_keamanan.Zona = req.Zona
		rekap_keamanan.CreatedBy = req.Nik

		if err := facades.Orm().Query().Table("rekapitulasi.kejadian").
			Join(`inner join rekapitulasi.klasifikasi_kejadian kk ON rekapitulasi.kejadian.klasifikasi_id = kk.id_klasifikasi`).
			Where("kk.id_klasifikasi = ? AND rekapitulasi.kejadian.id_type_kejadian = ?", "1", req.TypeKejadianId).
			First(&kejadian); err != nil || kejadian.IDTypeKejadian == "" {
			return Error(ctx, http.StatusNotFound, "Data Keamanan Not Found")
		}

		if err := facades.Orm().Query().Create(&rekap_keamanan); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}
		return Success(ctx, http.Json{
			"Success": "Data Berhasil Ditambahkan",
		})
	}
}

func (r *RekapKejadianKeamananController) ListRekapKeamanan(ctx http.Context) http.Response {
	var req rekap.ListKeamanan
	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var rekap_keamanan []models.RekapKejadianDataKeamanan
	query := facades.Orm().Query().Table("rekapitulasi.rekap_kejadian_data_keamanan").
		Join(`inner join rekapitulasi.kejadian k ON k.id_type_kejadian = rekapitulasi.rekap_kejadian_data_keamanan.type_kejadian_id
			  inner join rekapitulasi.klasifikasi_kejadian kk ON k.klasifikasi_id = kk.id_klasifikasi`)

	if req.Key != "" {
		query = query.Where(`(
								lower(k.jenis_pelanggaran) like lower(?) OR
							 	lower(rekapitulasi.rekap_kejadian_data_keamanan.nama_kapal) like lower(?) OR
							 	lower(rekapitulasi.rekap_kejadian_data_keamanan.lokasi_kejadian) like lower(?)
							)`, "%"+req.Key+"%", "%"+req.Key+"%", "%"+req.Key+"%")
	}

	if req.Zona != "" {
		query = query.Where("lower(rekapitulasi.rekap_kejadian_data_keamanan.zona) like lower(?)", req.Zona)
	}

	tanggal_awal, _ := time.Parse(time.DateOnly, req.TanggalAwal)
	tanggal_akhir, _ := time.Parse(time.DateOnly, req.TanggalAkhir)

	query = query.Where("rekapitulasi.rekap_kejadian_data_keamanan.tanggal between ? AND ?",
		tanggal_awal, tanggal_akhir)

	query = query.Where("kk.id_klasifikasi", "1")

	if err := query.Order("rekapitulasi.rekap_kejadian_data_keamanan.tanggal asc").Find(&rekap_keamanan); err != nil {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_rekap_keamanan": rekap_keamanan,
	})
}

func (r *RekapKejadianKeamananController) ShowDetailRekapKeamanan(ctx http.Context) http.Response {
	var req rekap.GetKeamanan

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var rekap_keamanan models.RekapKejadianDataKeamanan
	query := facades.Orm().Query().Table("rekapitulasi.rekap_kejadian_data_keamanan").
		Join(`inner join rekapitulasi.kejadian k ON k.id_type_kejadian = rekapitulasi.rekap_kejadian_data_keamanan.type_kejadian_id
			  inner join rekapitulasi.klasifikasi_kejadian kk ON k.klasifikasi_id = kk.id_klasifikasi`)

	query = query.Where("rekapitulasi.rekap_kejadian_data_keamanan.id_rekap_keamanan", req.IdRekapKeamanan)

	if err := query.First(&rekap_keamanan); err != nil || rekap_keamanan.IdRekapKeamanan == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	return Success(ctx, http.Json{
		"data_rekap_keamanan": rekap_keamanan,
	})
}

func (r *RekapKejadianKeamananController) DeleteRekapKeamanan(ctx http.Context) http.Response {
	var req rekap.GetKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var rekap_keamanan []models.RekapKejadianDataKeamanan
	if data, err := facades.Orm().Query().Where("id_rekap_keamanan = ? AND is_locked is false", req.IdRekapKeamanan).Delete(&rekap_keamanan); err != nil || data.RowsAffected == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada / Data Tidak Dapat Dihapus")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
