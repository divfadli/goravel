package controllers

import (
	"fmt"
	kejadianKeamanan "goravel/app/http/requests/kejadian_keamanan"
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
//

func (r *KejadianKeamananController) StoreKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.PostKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	form := ctx.Request().Origin().MultipartForm
	files := form.File["files"]

	var data_keamanan models.KejadianKeamanan
	var kejadian models.JenisKejadian
	var pesan string

	// Update
	if req.IdKejadianKeamanan != 0 {
		if err := facades.Orm().Query().
			Where("id_kejadian_keamanan = ? AND jenis_kejadian_id=?", req.IdKejadianKeamanan, req.JenisKejadianId).
			First(&data_keamanan); err != nil || data_keamanan.IdKejadianKeamanan == 0 {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keamanan Tidak Ditemukan!!")
		}

		if data_keamanan.IsLocked {
			return Error(ctx, http.StatusNotFound, "Data Kejadian Keamanan Terkunci!!")
		} else {
			data_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
			fmt.Println(data_keamanan.Tanggal)
			// data_keamanan.JenisKejadianId = req.JenisKejadianId
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

			var fileImageKeamanan []models.ImageKeamanan
			facades.Orm().Query().Where("kejadian_keamanan_id = ?", data_keamanan.IdKejadianKeamanan).Find(&fileImageKeamanan)

			var fileImage []models.FileImage
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

				fileImage = append(fileImage, models.FileImage{
					Filename:  newfileIdentificator,
					Extension: filepath.Ext(newfileIdentificator),
					Url:       folder,
				})
			}

			if err := facades.Orm().Query().Save(&data_keamanan); err != nil {
				return ErrorSystem(ctx, "Data Gagal Diubah")
			}

			for _, file := range fileImageKeamanan {
				var fileImage models.FileImage
				facades.Orm().Query().Where("id_file_image=?", file.FileImageID).First(&fileImage)
				facades.Storage().Delete(fileImage.Url)
				facades.Orm().Query().Delete(&file)
				facades.Orm().Query().Delete(&fileImage)
			}

			facades.Orm().Query().Create(&fileImage)
			for _, file := range fileImage {
				var fileImageKeamanan models.ImageKeamanan
				fileImageKeamanan.FileImageID = file.IdFileImage
				fileImageKeamanan.KejadianKeamananID = data_keamanan.IdKejadianKeamanan

				facades.Orm().Query().Create(&fileImageKeamanan)
			}

			pesan = "Data Berhasil Diubah"
		}
	} else {
		data_keamanan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
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

		var fileImage []models.FileImage
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

			fileImage = append(fileImage, models.FileImage{
				Filename:  newfileIdentificator,
				Extension: filepath.Ext(newfileIdentificator),
				Url:       folder,
			})
		}

		if err := facades.Orm().Query().
			Where("klasifikasi_name = ? AND id_jenis_kejadian = ?", "Keamanan Laut", req.JenisKejadianId).
			First(&kejadian); err != nil || kejadian.IDJenisKejadian == "" {
			return Error(ctx, http.StatusNotFound, "Data Jenis Kejadian Tidak Ditemukan!!")
		}

		if err := facades.Orm().Query().Create(&data_keamanan); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}

		facades.Orm().Query().Create(&fileImage)

		for _, file := range fileImage {
			var fileImageKeamanan models.ImageKeamanan
			fileImageKeamanan.FileImageID = file.IdFileImage
			fileImageKeamanan.KejadianKeamananID = data_keamanan.IdKejadianKeamanan

			facades.Orm().Query().Create(&fileImageKeamanan)
		}

		pesan = "Data Berhasil Ditambahkan"
	}
	return Success(ctx, http.Json{
		"message": pesan,
	})
}

func (r *KejadianKeamananController) ListKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.ListKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data_keamanan []models.KejadianKeamanan

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

	query.Order("tanggal asc").Find(&data_keamanan)

	var results []models.KejadianKeamananImage
	for _, data := range data_keamanan {
		var data_keamanan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

		results = append(results, models.KejadianKeamananImage{
			KejadianKeamanan: data,
			FileImage:        data_keamanan_image,
		})
	}

	return Success(ctx, http.Json{
		"data_kejadian_keamanan": results,
	})
}

func (r *KejadianKeamananController) ShowDetailKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.GetKeamanan

	if sanitize := SanitizePost(ctx, &req); sanitize != nil {
		return sanitize
	}

	var data_keamanan models.KejadianKeamanan

	if err := facades.Orm().Query().With("JenisKejadian").Where("id_kejadian_keamanan=?", req.IdKejadianKeamanan).
		First(&data_keamanan); err != nil || data_keamanan.IdKejadianKeamanan == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	var data_keamanan_image []models.FileImage

	facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
		Where("imk.kejadian_keamanan_id=?", data_keamanan.IdKejadianKeamanan).Find(&data_keamanan_image)

	results := models.KejadianKeamananImage{
		KejadianKeamanan: data_keamanan,
		FileImage:        data_keamanan_image,
	}

	return Success(ctx, http.Json{
		"data_kejadian_keamanan": results,
	})
}

func (r *KejadianKeamananController) DeleteKejadianKeamanan(ctx http.Context) http.Response {
	var req kejadianKeamanan.GetKeamanan

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var data_keamanan models.KejadianKeamanan
	if err := facades.Orm().Query().
		Where("id_kejadian_keamanan=? AND is_locked is false", req.IdKejadianKeamanan).
		First(&data_keamanan); err != nil || data_keamanan.IdKejadianKeamanan == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
	}

	var fileImageKeamanan []models.ImageKeamanan
	facades.Orm().Query().Where("kejadian_keamanan_id = ?", data_keamanan.IdKejadianKeamanan).Find(&fileImageKeamanan)

	for _, file := range fileImageKeamanan {
		var fileImage models.FileImage
		facades.Orm().Query().Where("id_file_image=?", file.FileImageID).First(&fileImage)
		facades.Storage().Delete(fileImage.Url)
		facades.Orm().Query().Delete(&file)
		facades.Orm().Query().Delete(&fileImage)
	}

	if x, err := facades.Orm().Query().Delete(&data_keamanan); err != nil || x.RowsAffected == 0 {
		return ErrorSystem(ctx, "Gagal Menghapus Data!!")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}
