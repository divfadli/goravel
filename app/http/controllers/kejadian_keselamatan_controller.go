package controllers

import (
	"encoding/json"
	"fmt"
	kejadianKeselamatan "goravel/app/http/requests/kejadian_keselamatan"
	"goravel/app/models"
	template "html/template"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/filesystem"
	"github.com/xuri/excelize/v2"
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
		fmt.Println("UPDATE")
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
				data_keselamatan.PelabuhanAsal = req.PelabuhanAsal
			}

			if req.PelabuhanTujuan != "" {
				data_keselamatan.PelabuhanTujuan = req.PelabuhanTujuan
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
				folder, err := facades.Storage().PutFileAs(strconv.Itoa(waktu.Year())+"/Photos/keselamatan/"+waktu.Month().String(), newFile, newfileIdentificator)

				if err != nil {
					return Error(ctx, http.StatusInternalServerError, err.Error())
				}

				fileImage = append(fileImage, models.FileImage{
					Filename:  newfileIdentificator,
					Extension: filepath.Ext(newfileIdentificator),
					Url:       facades.Storage().Url(folder),
				})
			}

			for _, file := range fileImageKeselamatan {
				var fileImage models.FileImage
				facades.Orm().Query().Where("id_file_image=?", file.FileImageID).First(&fileImage)
				facades.Storage().Delete(fileImage.Url)
				facades.Orm().Query().Delete(&file)
				facades.Orm().Query().Delete(&fileImage)
			}

			facades.Orm().Query().Create(&fileImage)
			for _, file := range fileImage {
				var fileImageKeselamatan models.ImageKeselamatan
				fileImageKeselamatan.FileImageID = file.IdFileImage
				fileImageKeselamatan.KejadianKeselamatanID = data_keselamatan.IdKejadianKeselamatan

				facades.Orm().Query().Create(&fileImageKeselamatan)
			}

			pesan = "Data Berhasil Diubah"
		}
	} else {
		fmt.Println("CREATE")

		data_keselamatan.Tanggal = carbon.Parse(req.Tanggal).ToDateStruct()
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
			data_keselamatan.PelabuhanAsal = req.PelabuhanAsal
		}

		if req.PelabuhanTujuan != "" {
			data_keselamatan.PelabuhanTujuan = req.PelabuhanTujuan
		}

		data_keselamatan.TindakLanjut = req.TindakLanjut
		data_keselamatan.Keterangan = req.Keterangan
		data_keselamatan.Zona = req.Zona
		data_keselamatan.CreatedBy = req.Nik

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
			folder, err := facades.Storage().PutFileAs(strconv.Itoa(waktu.Year())+"/Photos/Keselamatan/"+waktu.Month().String(), newFile, newfileIdentificator)

			if err != nil {
				return Error(ctx, http.StatusInternalServerError, err.Error())
			}

			fileImage = append(fileImage, models.FileImage{
				Filename:  newfileIdentificator,
				Extension: filepath.Ext(newfileIdentificator),
				Url:       facades.Storage().Url(folder),
			})
		}

		fmt.Println(fileImage)

		if err := facades.Orm().Query().
			Where("klasifikasi_name = ? AND id_jenis_kejadian = ?", "Keselamatan Laut", req.JenisKejadianId).
			First(&kejadian); err != nil || kejadian.IDJenisKejadian == "" {
			return Error(ctx, http.StatusNotFound, "Data Jenis Kejadian Tidak Ditemukan!!")
		}

		if err := facades.Orm().Query().Create(&data_keselamatan); err != nil {
			return ErrorSystem(ctx, "Data Gagal Ditambahkan")
		}

		facades.Orm().Query().Create(&fileImage)

		for _, file := range fileImage {
			fmt.Println(file.IdFileImage)
			var fileImageKeselamatan models.ImageKeselamatan
			fileImageKeselamatan.FileImageID = file.IdFileImage
			fileImageKeselamatan.KejadianKeselamatanID = data_keselamatan.IdKejadianKeselamatan

			facades.Orm().Query().Create(&fileImageKeselamatan)
		}
		pesan = "Data Berhasil Ditambahkan"
	}
	return Success(ctx, http.Json{
		"message": pesan,
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

	query.Order("tanggal asc").Find(&data_keselamatan)

	var results []models.KejadianKeselamatanImage
	for _, data := range data_keselamatan {
		var data_keselamatan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
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
	facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
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
	if err := facades.Orm().Query().
		Where("id_kejadian_keselamatan=? AND is_locked is false", req.IdKejadianKeselamatan).
		First(&data_keselamatan); err != nil || data_keselamatan.IdKejadianKeselamatan == 0 {
		return ErrorSystem(ctx, "Data Tidak Ada")
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

	if x, err := facades.Orm().Query().Delete(&data_keselamatan); err != nil || x.RowsAffected == 0 {
		return ErrorSystem(ctx, "Gagal Menghapus Data!!")
	}

	return Success(ctx, "Data Berhasil Dihapus")
}

func (r *KejadianKeselamatanController) ExportExcel(ctx http.Context) http.Response {
	var req kejadianKeselamatan.ListKeselamatan
	if err := ctx.Request().Bind(&req); err != nil {
		return ErrorSystem(ctx, "Invalid request parameters")
	}

	// Query data with validation
	tanggalAwal, err := time.Parse(time.DateOnly, req.TanggalAwal)
	if err != nil {
		return ErrorSystem(ctx, "Invalid start date format")
	}
	tanggalAkhir, err := time.Parse(time.DateOnly, req.TanggalAkhir)
	if err != nil {
		return ErrorSystem(ctx, "Invalid end date format")
	}

	// Query with eager loading and conditions
	var kejadianKeselamatan []models.KejadianKeselamatan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")

	if req.Zona != "" && req.Zona != "null" && req.Zona != "undefined" {
		query = query.Where("lower(zona::text) like lower(?)", "%"+req.Zona+"%")
	}
	query = query.Where("tanggal BETWEEN (?) AND (?)", tanggalAwal, tanggalAkhir)

	if err := query.Order("tanggal asc").Find(&kejadianKeselamatan); err != nil {
		return ErrorSystem(ctx, "Failed to fetch data")
	}

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			facades.Log().Error(fmt.Sprintf("Failed to close Excel file: %v", err))
		}
	}()

	sheet := "Sheet1"

	// Enhanced header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0EBF5"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	// Data cell style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 10},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})

	// Set headers with improved column names
	headers := []string{
		"Tanggal",
		"Nama Kapal",
		"Jenis Pelanggaran",
		"Lokasi Kejadian",
		"Muatan",
		"Zona",
		"Dokumentasi",
	}

	// Write headers
	for i, header := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Optimize column widths
	columnWidths := map[string]float64{
		"A": 15, // Tanggal
		"B": 25, // Nama Kapal
		"C": 25, // Jenis Kejadian
		"D": 30, // Lokasi
		"E": 20, // Penyebab
		"F": 15, // Zona
		"G": 50, // Images
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheet, col, col, width)
	}

	// Write data rows with optimized image handling
	for i, item := range kejadianKeselamatan {
		row := i + 2

		// Start with base row height
		baseRowHeight := 30.0 // minimum height

		// Calculate content height based on text length and wrapping
		contentHeight := calculateContentHeightKecelakaan(item)
		rowHeight := baseRowHeight

		if contentHeight > baseRowHeight {
			rowHeight = contentHeight
		}

		// Set initial row height based on content
		f.SetRowHeight(sheet, row, rowHeight)

		// Set cell values
		cells := []struct {
			cell  string
			value interface{}
		}{
			{fmt.Sprintf("A%d", row), item.Tanggal},
			{fmt.Sprintf("B%d", row), item.NamaKapal},
			{fmt.Sprintf("C%d", row), item.JenisKejadian.NamaKejadian},
			{fmt.Sprintf("D%d", row), item.LokasiKejadian},
			{fmt.Sprintf("E%d", row), item.Penyebab},
			{fmt.Sprintf("F%d", row), item.Zona},
		}

		for _, cell := range cells {
			f.SetCellValue(sheet, cell.cell, cell.value)
			f.SetCellStyle(sheet, cell.cell, cell.cell, dataStyle)
		}

		// Handle images with improved layout and consistent sizing
		imgCell := fmt.Sprintf("G%d", row)
		f.SetCellStyle(sheet, imgCell, imgCell, dataStyle)

		var images []models.FileImage
		if err := facades.Orm().Query().
			Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keselamatan_id=?", item.IdKejadianKeselamatan).
			Find(&images); err == nil && len(images) > 0 {

			for _, img := range images {
				physicalPath := facades.Storage().Path(strings.TrimPrefix(img.Url, "/storage/app/"))

				// Get image dimensions
				imgFile, err := os.Open(physicalPath)
				if err != nil {
					facades.Log().Error(fmt.Sprintf("Failed to open image %s: %v", img.Url, err))
					continue
				}
				defer imgFile.Close()

				imgConfig, _, err := image.DecodeConfig(imgFile)
				if err != nil {
					facades.Log().Error(fmt.Sprintf("Failed to decode image %s: %v", img.Url, err))
					continue
				}

				// Calculate real height in Excel units
				realHeight := float64(imgConfig.Height)
				if realHeight > rowHeight {
					rowHeight = realHeight
				}

				err = f.AddPicture(sheet, imgCell, physicalPath, &excelize.GraphicOptions{
					Positioning:         "oneCell",
					AutoFit:             true,
					AutoFitIgnoreAspect: true,
				})
				if err != nil {
					facades.Log().Error(fmt.Sprintf("Failed to add image %s: %v", img.Url, err))
					continue
				}
			}

			f.SetRowHeight(sheet, row, rowHeight)
		}

	}

	// Generate file with timestamp
	// filename := fmt.Sprintf("kejadian_keselamatan_%s.xlsx", time.Now().Format("20060102_150405"))
	filename := "kejadian_keselamatan_temp.xlsx"
	if err := f.SaveAs(filename); err != nil {
		return ErrorSystem(ctx, "Failed to generate Excel file")
	}

	return ctx.Response().Download(filename, filename)
}

func calculateContentHeightKecelakaan(item models.KejadianKeselamatan) float64 {
	// Base height per line of text
	lineHeight := 15.0

	// Count approximate number of lines based on content length
	lines := 1.0

	// Add lines for each field that might wrap
	if len(item.NamaKapal) > 30 {
		lines += float64(len(item.NamaKapal)) / 30
	}
	if len(item.LokasiKejadian) > 40 {
		lines += float64(len(item.LokasiKejadian)) / 40
	}
	if len(item.Penyebab) > 25 {
		lines += float64(len(item.Penyebab)) / 25
	}

	return lines * lineHeight
}

func GetDetailKejadianKeselamatan(ctx http.Context) http.Response {
	userInfo := facades.Cache().Get("user_data")

	id := ctx.Request().Route("id")

	if userInfo != nil {
		baseURL := "http://" + ctx.Request().Host()

		var data models.KejadianKeselamatan
		var jenisKejadian []models.JenisKejadian

		facades.Orm().Query().Where("klasifikasi_name =? AND deleted_at IS NULL", "Keselamatan Laut").
			Order("nama_kejadian asc").Get(&jenisKejadian)

		if err := facades.Orm().Query().With("JenisKejadian").Where("id_kejadian_keselamatan", id).
			First(&data); err != nil || data.IdKejadianKeselamatan == 0 {
			return ErrorSystem(ctx, "Data Tidak Ada")
		}

		var data_keselamatan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

		var perpindahanAwal string
		if data.PelabuhanAsal != "-" && data.PelabuhanAsal != "" {
			perpindahanAwal = data.PelabuhanAsal
		}
		var perpindahanAkhir string
		if data.PelabuhanTujuan != "-" && data.PelabuhanTujuan != "" {
			perpindahanAwal = data.PelabuhanTujuan
		}

		var perpindahan string
		if perpindahanAwal != "" && perpindahanAkhir != "" {
			perpindahan = perpindahanAwal + " - " + perpindahanAkhir
		} else if perpindahanAwal != "" && perpindahanAkhir == "" {
			perpindahan = perpindahanAwal + " - "
		} else if perpindahanAwal == "" && perpindahanAkhir != "" {
			perpindahan = " - " + perpindahanAkhir
		} else {
			perpindahan = " - "
		}

		var korbanData models.ListKorban

		var korban string
		if err := json.Unmarshal(data.Korban, &korbanData); err != nil {
			fmt.Println("ERROR", err)
			return nil
		}

		if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas != 0 {
			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, selamat " +
				strconv.Itoa(korbanData.KorbanSelamat) + " orang, dan tewas " +
				strconv.Itoa(korbanData.KorbanTewas) + " orang"
		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas == 0 {
			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, dan selamat " +
				strconv.Itoa(korbanData.KorbanSelamat) + " orang"
		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas != 0 {
			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, dan tewas " +
				strconv.Itoa(korbanData.KorbanTewas) + " orang"
		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas != 0 {
			korban = "Korban selamat " + strconv.Itoa(korbanData.KorbanSelamat) + " orang, dan tewas " +
				strconv.Itoa(korbanData.KorbanTewas) + " orang"
		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas == 0 {
			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang"
		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas == 0 {
			korban = "Korban selamat " + strconv.Itoa(korbanData.KorbanSelamat) + " orang"
		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas != 0 {
			korban = "Korban tewas " + strconv.Itoa(korbanData.KorbanTewas) + " orang"
		} else {
			korban = "tidak ada korban jiwa"
		}

		linkBerita := strings.Split(data.LinkBerita, "\n")
		var sumberBerita []string

		for _, v := range linkBerita {
			if strings.Contains(v, "http") {
				trimmed := strings.TrimSpace(v)
				berita := fmt.Sprintf(`<a href="%s" target="_blank">%s</a><br>`, trimmed, trimmed)
				sumberBerita = append(sumberBerita, berita)
			} else {
				trimmed := strings.TrimSpace(v)
				berita := fmt.Sprintf(`<u>%s</u><br>`, trimmed)
				sumberBerita = append(sumberBerita, berita)
			}
		}

		// html template data
		templateData := struct {
			JenisKejadian    []models.JenisKejadian
			BaseURL          string
			Title            string
			NamaKapal        string
			Kejadian         string
			Penyebab         string
			Lokasi           string
			Korban           string
			Perpindahan      string
			Keterangan       string
			Waktu            string
			InstansiPenindak string
			SumberBerita     template.HTML
			Latitude         float64
			Longitude        float64
			Images           []models.FileImage
		}{
			JenisKejadian:    jenisKejadian,
			BaseURL:          baseURL,
			Title:            data.JenisKejadian.NamaKejadian,
			NamaKapal:        data.NamaKapal,
			Kejadian:         data.JenisKejadian.NamaKejadian,
			Penyebab:         data.Penyebab,
			Lokasi:           data.LokasiKejadian,
			Korban:           korban,
			Perpindahan:      perpindahan,
			Keterangan:       data.TindakLanjut,
			Waktu:            data.Tanggal.ToDateString(),
			InstansiPenindak: data.SumberBerita,
			SumberBerita:     template.HTML(strings.Join(sumberBerita, "")),
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Images:           data_keselamatan_image,
		}

		return ctx.Response().View().Make("detail-keselamatan.tmpl", templateData)
	}

	facades.Auth().Logout(ctx)

	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}
