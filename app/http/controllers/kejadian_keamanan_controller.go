package controllers

import (
	"encoding/json"
	"fmt"
	kejadianKeamanan "goravel/app/http/requests/kejadian_keamanan"
	"goravel/app/models"
	template "html/template"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/filesystem"
	excelize "github.com/xuri/excelize/v2"
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
					Url:       facades.Storage().Url(folder),
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
				Url:       facades.Storage().Url(folder),
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

func (r *KejadianKeamananController) ExportExcel(ctx http.Context) http.Response {
	var req kejadianKeamanan.ListKeamanan
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
	var kejadianKeamanan []models.KejadianKeamanan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")

	if req.Zona != "" && req.Zona != "null" && req.Zona != "undefined" {
		query = query.Where("lower(zona::text) like lower(?)", "%"+req.Zona+"%")
	}

	query = query.Where("tanggal BETWEEN (?) AND (?)", tanggalAwal, tanggalAkhir)

	if err := query.Order("tanggal asc").Find(&kejadianKeamanan); err != nil {
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
		"E": 20, // Muatan
		"F": 15, // Zona
		"G": 50, // Images
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheet, col, col, width)
	}

	// Write data rows with optimized image handling
	for i, item := range kejadianKeamanan {
		row := i + 2

		// Start with base row height
		baseRowHeight := 30.0 // minimum height

		// Calculate content height based on text length and wrapping
		contentHeight := calculateContentHeightPelanggaran(item)
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
			{fmt.Sprintf("E%d", row), item.Muatan},
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
			Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keamanan_id=?", item.IdKejadianKeamanan).
			Find(&images); err == nil && len(images) > 0 {

			// Enhanced image placement configuration
			const (
				imageScale   = 0.25 // Reduced scale for better fit
				baseOffsetY  = 5    // Increased top padding
				imagesPerRow = 2    // 2 images per row for better organization
				imageSpacing = 20   // Increased spacing between images
				rowSpacing   = 25   // Vertical spacing between rows
			)

			for idx, img := range images {
				physicalPath := facades.Storage().Path(strings.TrimPrefix(img.Url, "/storage/app/"))

				// Calculate grid position
				rowPosition := idx / imagesPerRow
				colPosition := idx % imagesPerRow

				// Calculate precise positioning
				offsetX := float64(colPosition * imageSpacing)
				offsetY := float64(rowPosition*rowSpacing) + baseOffsetY

				err := f.AddPicture(sheet, imgCell, physicalPath, &excelize.GraphicOptions{
					ScaleX:      imageScale,
					ScaleY:      imageScale,
					Positioning: "oneCell",
					OffsetX:     int(offsetX),
					OffsetY:     int(offsetY),
				})
				if err != nil {
					facades.Log().Error(fmt.Sprintf("Failed to add image %s: %v", img.Url, err))
					continue
				}
			}

			// Dynamic row height based on number of images
			numRows := (len(images) + imagesPerRow - 1) / imagesPerRow
			rowHeight := float64(baseOffsetY + (numRows * rowSpacing))
			if rowHeight < 30 {
				rowHeight = 30 // Minimum row height
			}
			f.SetRowHeight(sheet, row, rowHeight)
		}

	}

	// Generate file with timestamp
	// filename := fmt.Sprintf("kejadian_keamanan_%s.xlsx", time.Now().Format("20060102_150405"))
	filename := "kejadian_keamanan_temp.xlsx"
	if err := f.SaveAs(filename); err != nil {
		return ErrorSystem(ctx, "Failed to generate Excel file")
	}

	return ctx.Response().Download(filename, filename)
}

func calculateContentHeightPelanggaran(item models.KejadianKeamanan) float64 {
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
	if len(item.Muatan) > 25 {
		lines += float64(len(item.Muatan)) / 25
	}

	return lines * lineHeight
}

func GetDetailKejadianKeamanan(ctx http.Context) http.Response {
	userInfo := facades.Cache().Get("user_data")

	id := ctx.Request().Route("id")

	if userInfo != nil {
		baseURL := "http://" + ctx.Request().Host()

		var data models.KejadianKeamanan
		var jenisKejadian []models.JenisKejadian

		facades.Orm().Query().Where("klasifikasi_name =? AND deleted_at IS NULL", "Keamanan Laut").
			Order("nama_kejadian asc").Get(&jenisKejadian)

		if err := facades.Orm().Query().With("JenisKejadian").Where("id_kejadian_keamanan=?", id).
			First(&data); err != nil || data.IdKejadianKeamanan == 0 {
			return ErrorSystem(ctx, "Data Tidak Ada")
		}
		jsonData, _ := json.MarshalIndent(jenisKejadian, "", "    ")
		fmt.Println(string(jsonData))

		var data_keamanan_image []models.FileImage

		facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

		var abk string
		if strings.Contains(data.Muatan, "ABK") {
			re := regexp.MustCompile(`\b\d+\s+orang\b`)
			matches := re.FindAllString(data.Muatan, -1)
			if len(matches) > 0 {
				abk = matches[0]
			} else {
				abk = " - "
			}
		} else {
			abk = " - "
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

		templateData := struct {
			JenisKejadian    []models.JenisKejadian
			BaseURL          string
			Title            string
			NamaKapal        string
			Kejadian         string
			Penyebab         string
			Lokasi           string
			ABK              string
			Muatan           string
			InstansiPenindak string
			Keterangan       string
			Waktu            string
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
			Penyebab:         "-",
			Lokasi:           data.LokasiKejadian,
			ABK:              abk,
			Muatan:           data.Muatan,
			InstansiPenindak: data.SumberBerita,
			Keterangan:       data.TindakLanjut,
			Waktu:            data.Tanggal.ToDateString(),
			SumberBerita:     template.HTML(strings.Join(sumberBerita, "")),
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Images:           data_keamanan_image,
		}

		return ctx.Response().View().Make("detail-keamanan.tmpl", templateData)
	}

	facades.Auth().Logout(ctx)

	// For instance, you might redirect the user to the login page
	return ctx.Response().Redirect(http.StatusFound, "/login")
}
