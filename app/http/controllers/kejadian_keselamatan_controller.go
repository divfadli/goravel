package controllers

import (
	"encoding/json"
	"fmt"
	kejadianKeselamatan "goravel/app/http/requests/kejadian_keselamatan"
	"goravel/app/models"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/filesystem"
)

type KejadianKeselamatanController struct {
	//Dependent services
}

func NewKejadianKeselamatanController() *KejadianKeselamatanController {
	return &KejadianKeselamatanController{
		//Inject services
	}
}

// func (r *KejadianKeselamatanController) Index(ctx http.Context) http.Response {
// 	return nil
// }

func (r *KejadianKeselamatanController) StoreKejadianKeselamatan(ctx http.Context) http.Response {
	var req kejadianKeselamatan.PostKeselamatan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	form := ctx.Request().Origin().MultipartForm
	files := form.File["files"]

	var data_keselamatan models.KejadianKeselamatan
	var list_korban models.ListKorban
	var kejadian models.JenisKejadian
	var pesan string

	// Update
	if req.IdKejadianKeselamatan != 0 {
		if err := facades.Orm().Query().
			Where("id_kejadian_keselamatan = ? AND jenis_kejadian_id=?", req.IdKejadianKeselamatan, req.JenisKejadianId).
			First(&data_keselamatan); err != nil || data_keselamatan.IdKejadianKeselamatan == 0 {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keselamatan Tidak Ditemukan!!")
		}

		if data_keselamatan.IsLocked {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keselamatan Terkunci!!")
		} else {
			data_keselamatan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
			fmt.Println(data_keselamatan.Tanggal)
			// data_keselamatan.JenisKejadianId = req.JenisKejadianId
			data_keselamatan.NamaKapal = req.NamaKapal
			data_keselamatan.SumberBerita = req.SumberBerita
			data_keselamatan.LinkBerita = req.LinkBerita
			data_keselamatan.LokasiKejadian = req.LokasiKejadian

			list_korban.KorbanHilang = req.KorbanHilang
			list_korban.KorbanSelamat = req.KorbanSelamat
			list_korban.KorbanTewas = req.KorbanTewas

			jsonKorban, _ := json.Marshal(list_korban)
			data_keselamatan.Korban = json.RawMessage(jsonKorban)

			data_keselamatan.Latitude = req.Latitude
			data_keselamatan.Longitude = req.Longitude
			data_keselamatan.Penyebab = req.Penyebab
			data_keselamatan.TipeSumberKejadian = req.KategoriSumber

			if req.PelabuhanAsal != "" {
				data_keselamatan.PelabuhanAsal = &req.PelabuhanAsal
			}

			if req.PelabuhanTujuan != "" {
				data_keselamatan.PelabuhanTujuan = &req.PelabuhanTujuan
			}

			data_keselamatan.TindakLanjut = req.TindakLanjut
			data_keselamatan.Keterangan = req.Keterangan
			data_keselamatan.Zona = req.Zona
			data_keselamatan.CreatedBy = req.Nik

			if err := facades.Orm().Query().Save(&data_keselamatan); err != nil {
				return ErrorSystem(ctx, "Data Gagal Diubah")
			}

			var fileImageKeselamatan []models.ImageKeselamatan
			facades.Orm().Query().Where("kejadian_keselamatan_id = ?", data_keselamatan.IdKejadianKeselamatan).Find(&fileImageKeselamatan)

			for _, file := range fileImageKeselamatan {
				var fileImage models.FileImage
				facades.Orm().Query().Where("id_file_image=?", file.FileImageID).First(&fileImage)
				facades.Storage().Delete(fileImage.Url)
				facades.Orm().Query().Delete(&file)
				facades.Orm().Query().Delete(&fileImage)
			}

			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					return Error(ctx, http.StatusInternalServerError, err.Error())
				}
				defer file.Close()

				if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
					return Error(ctx, http.StatusInternalServerError, "File Tidak Valid!!")
				}

				// You can directly pass the multipart.File to the storage function
				newfileIdentificator := buildFileIdentificator(fileHeader.Filename)
				newFile, _ := filesystem.NewFileFromRequest(fileHeader)

				waktu := time.Now()
				folder, err := facades.Storage().PutFileAs(strconv.Itoa(waktu.Year())+"/Photos/Keamanan/"+waktu.Month().String(), newFile, newfileIdentificator)

				if err != nil {
					return Error(ctx, http.StatusInternalServerError, err.Error())
				}

				var fileImage models.FileImage
				fileImage.Filename = newfileIdentificator
				fileImage.Extension = filepath.Ext(newfileIdentificator)
				fileImage.Url = folder

				facades.Orm().Query().Create(&fileImage)

				var fileImageKeselamatan models.ImageKeamanan
				fileImageKeselamatan.FileImageID = fileImage.IdFileImage
				fileImageKeselamatan.KejadianKeamananID = data_keselamatan.IdKejadianKeselamatan

				facades.Orm().Query().Create(&fileImageKeselamatan)
			}

			pesan = "Data Berhasil Diubah"
		}
	} else {
		data_keselamatan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
		fmt.Println(data_keselamatan.Tanggal)
		data_keselamatan.JenisKejadianId = req.JenisKejadianId
		data_keselamatan.NamaKapal = req.NamaKapal
		data_keselamatan.SumberBerita = req.SumberBerita
		data_keselamatan.LinkBerita = req.LinkBerita
		data_keselamatan.LokasiKejadian = req.LokasiKejadian

		list_korban.KorbanHilang = req.KorbanHilang
		list_korban.KorbanSelamat = req.KorbanSelamat
		list_korban.KorbanTewas = req.KorbanTewas

		jsonKorban, _ := json.Marshal(list_korban)
		data_keselamatan.Korban = json.RawMessage(jsonKorban)

		data_keselamatan.Latitude = req.Latitude
		data_keselamatan.Longitude = req.Longitude
		data_keselamatan.Penyebab = req.Penyebab
		data_keselamatan.TipeSumberKejadian = req.KategoriSumber

		if req.PelabuhanAsal != "" {
			data_keselamatan.PelabuhanAsal = &req.PelabuhanAsal
		}

		if req.PelabuhanTujuan != "" {
			data_keselamatan.PelabuhanTujuan = &req.PelabuhanTujuan
		}

		data_keselamatan.TindakLanjut = req.TindakLanjut
		data_keselamatan.Keterangan = req.Keterangan
		data_keselamatan.Zona = req.Zona
		data_keselamatan.CreatedBy = req.Nik

		if err := facades.Orm().Query().
			Where("klasifikasi_name = ? AND id_jenis_kejadian = ?", "Keselamatan Laut", req.JenisKejadianId).
			First(&kejadian); err != nil || kejadian.IDJenisKejadian == "" {
			return Error(ctx, http.StatusNotFound, "Data Jenis Kejadian Tidak Ditemukan!!")
		}

		if err := facades.Orm().Query().Create(&data_keselamatan); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return Error(ctx, http.StatusInternalServerError, err.Error())
			}
			defer file.Close()

			if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
				return Error(ctx, http.StatusInternalServerError, "File Tidak Valid!!")
			}

			// You can directly pass the multipart.File to the storage function
			newfileIdentificator := buildFileIdentificator(fileHeader.Filename)
			newFile, _ := filesystem.NewFileFromRequest(fileHeader)

			waktu := time.Now()
			folder, err := facades.Storage().PutFileAs(strconv.Itoa(waktu.Year())+"/Photos/Keselamatan/"+waktu.Month().String(), newFile, newfileIdentificator)

			if err != nil {
				return Error(ctx, http.StatusInternalServerError, err.Error())
			}

			var fileImage models.FileImage
			fileImage.Filename = newfileIdentificator
			fileImage.Extension = filepath.Ext(newfileIdentificator)
			fileImage.Url = folder

			facades.Orm().Query().Create(&fileImage)

			var fileImageKeselamatan models.ImageKeselamatan
			fileImageKeselamatan.FileImageID = fileImage.IdFileImage
			fileImageKeselamatan.KejadianKeselamatanID = data_keselamatan.IdKejadianKeselamatan

			facades.Orm().Query().Create(&fileImageKeselamatan)
		}

		pesan = "Data Berhasil Ditambahkan"
	}
	return Success(ctx, http.Json{
		"Success": pesan,
	})
}

func (r *KejadianKeselamatanController) ListKejadianKeselamatan(ctx http.Context) http.Response {
	var req kejadianKeselamatan.ListKeselamatan
	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data_keselamatan []models.KejadianKeselamatan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")

	if req.Key != "" {
		query = query.Where(`(
								lower(k.nama_kejadian) like lower(?) OR
							 	lower(nama_kapal) like lower(?) OR
							 	lower(lokasi_kejadian) like lower(?)
							)`, "%"+req.Key+"%", "%"+req.Key+"%", "%"+req.Key+"%")
	}

	if req.Zona != "" {
		query = query.Where("lower(zona::text) like lower(?)", "%"+req.Zona+"%")
	}

	tanggal_awal, _ := time.Parse(time.DateOnly, req.TanggalAwal)
	tanggal_akhir, _ := time.Parse(time.DateOnly, req.TanggalAkhir)

	query = query.Where("tanggal between ? AND ?",
		tanggal_awal, tanggal_akhir)

	if err := query.Order("tanggal asc").Find(&data_keselamatan); err != nil || len(data_keselamatan) == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	var results []models.KejadianKeselamatanImage
	for _, data := range data_keselamatan {
		var data_keselamatan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

		results = append(results, models.KejadianKeselamatanImage{
			KejadianKeselamatan: data,
			FileImage:           data_keselamatan_image,
		})
	}

	return Success(ctx, http.Json{
		"data_kejadian_keselamatan": results,
	})
}

func (r *KejadianKeselamatanController) ShowDetailKejadianKeselamatan(ctx http.Context) http.Response {
	var req kejadianKeselamatan.GetKeselamatan

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var data_keselamatan models.KejadianKeselamatan

	if err := facades.Orm().Query().With("JenisKejadian").Where("id_kejadian_keselamatan", req.IdKejadianKeselamatan).
		First(&data_keselamatan); err != nil || data_keselamatan.IdKejadianKeselamatan == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	var data_keselamatan_image []models.FileImage
	facades.Orm().Query().Join("inner join public.image_keselematan imk ON id_file_image = imk.file_image_id").
		Where("imk.kejadian_keselamatan_id=?", data_keselamatan.IdKejadianKeselamatan).Find(&data_keselamatan_image)

	results := models.KejadianKeselamatanImage{
		KejadianKeselamatan: data_keselamatan,
		FileImage:           data_keselamatan_image,
	}

	return Success(ctx, http.Json{
		"data_kejadian_keselamatan": results,
	})
}

func (r *KejadianKeselamatanController) DeleteKejadianKeselamatan(ctx http.Context) http.Response {
	var req kejadianKeselamatan.GetKeselamatan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data_keselamatan models.KejadianKeselamatan
	facades.Orm().Query().Where("id_kejadian_keselamatan=? AND is_locked is false", req.IdKejadianKeselamatan).First(&data_keselamatan)

	var fileImageKeselamatan []models.ImageKeselamatan
	facades.Orm().Query().Where("kejadian_keselamatan_id = ?", data_keselamatan.IdKejadianKeselamatan).Find(&fileImageKeselamatan)

	for _, file := range fileImageKeselamatan {
		var fileImage models.FileImage
		facades.Orm().Query().Where("id_file_image=?", file.FileImageID).First(&fileImage)
		facades.Storage().Delete(fileImage.Url)
		facades.Orm().Query().Delete(&file)
		facades.Orm().Query().Delete(&fileImage)
	}

	if x, err := facades.Orm().Query().Delete(&data_keselamatan); err != nil || x.RowsAffected == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada / Data Tidak Dapat Dihapus")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
