package generator

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/jung-kurt/gofpdf"
	"github.com/lib/pq"
)

func (r *Pdf) LaporanMingguan() {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)

	//html template path
	templateKeamananPath := "templates/keamanan.html"
	newTemplateKeamananPath := "keamanan.html"
	templateKeamananHeadPath := "templates/keamanan-head.html"
	newTemplateKeamananHeadPath := "keamanan-head.html"
	templateKeselamatanHeadPath := "templates/keselamatan-head.html"
	newTemplateKeselamatanHeadPath := "keselamatan-head.html"
	templateKeselamatanPath := "templates/keselamatan.html"
	newTemplateKeselamatanPath := "keselamatan.html"

	now := time.Now()
	bulan := monthNameIndonesia(now.Month())
	year := strconv.Itoa(now.Year())

	dayperweek := 7

	jumlahHari := daysInMonth(now)

	var minggu []time.Time

	// Calculate the first day of the month
	var firstDayMonth, startOfMonth time.Time
	firstDayMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfMonth = firstDayMonth
	dayOfWeek := firstDayMonth.Weekday()

	daysToAdd := 0
	switch dayOfWeek {
	case time.Monday:
		daysToAdd = 6
	case time.Tuesday:
		daysToAdd = 5
	case time.Wednesday:
		daysToAdd = 4
	case time.Thursday:
		daysToAdd = 3
	case time.Friday:
		daysToAdd = 2
	case time.Saturday:
		daysToAdd = 8
	case time.Sunday:
		daysToAdd = 7
	}

	// minggu = append(minggu, firstDay)
	nextWeek := firstDayMonth.AddDate(0, 0, daysToAdd)
	minggu = append(minggu, nextWeek)

	jumlahHari -= daysToAdd + 1
	// fmt.Println(jumlahHari)
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// // // Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		startDate := nextWeek.AddDate(0, 0, i*dayperweek)
		minggu = append(minggu, startDate)
	}

	// // // Add the start date for the remaining days if any
	if remainingDays > 0 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var startOfWeek time.Time
	var endOfWeek time.Time
	var mingguKe int
	for i, weekEnd := range minggu {
		if weekEnd.After(now) {
			// minggu ke
			mingguKe = i
			startOfWeek = minggu[i-2]
			endOfWeek = minggu[i-1]
			break
		}
	}
	fmt.Println(startOfMonth, startOfWeek, endOfWeek)

	var weekName []string
	for i, weekEnd := range minggu {
		var setName, fullName string
		if weekEnd.Day() < 10 || (i > 0 && minggu[i-1].AddDate(0, 0, 1).Day() < 10) {
			setName = "0"
		}
		if i == 0 {
			fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
		} else {
			fullName = setName + strconv.Itoa(minggu[i-1].AddDate(0, 0, 1).Day()) + "-" +
				strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
		}
		weekName = append(weekName, fullName)
	}
	fmt.Println(weekName)

	var weeklyDataKeamanans, weeklyDataKeselamatans, allWeeklyDataKeamanans, allWeeklyDataKeselamatans weeklyData

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&allWeeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&weeklyDataKeselamatans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&allWeeklyDataKeselamatans)

	allWeeklyDataKeamanan := make(map[string][]models.KejadianKeamanan)
	allWeeklyDataKeselamatan := make(map[string][]models.KejadianKeselamatan)

	var data_keamanan []models.KejadianKeamanan
	var data_keselamatan []models.KejadianKeselamatan
	for _, week := range allWeeklyDataKeamanans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keamanan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

		for i, weekEnd := range minggu {
			var setName, fullName string
			if weekEnd.Day() < 10 || (i > 0 && minggu[i-1].AddDate(0, 0, 1).Day() < 10) {
				setName = "0"
			}
			if weekEnd.After(week.WeekStart) {
				if i == 0 {
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
				}

				allWeeklyDataKeamanan[fullName] = append(allWeeklyDataKeamanan[fullName], data_keamanan...)
				break
			}
		}
	}
	for _, week := range allWeeklyDataKeselamatans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keselamatan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)
		for i, weekEnd := range minggu {
			var setName, fullName string
			if weekEnd.Day() < 10 || (i > 0 && minggu[i-1].AddDate(0, 0, 1).Day() < 10) {
				setName = "0"
			}
			if weekEnd.After(week.WeekStart) {
				if i == 0 {
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + monthNameEnglishMap[weekEnd.Month()]
				}

				allWeeklyDataKeselamatan[fullName] = append(allWeeklyDataKeselamatan[fullName], data_keselamatan...)
				break
			}
		}
	}

	var result_keamanan []models.KejadianKeamananImage
	var result_keselamatan []models.KejadianKeselamatanImage
	// Now you can loop through the grouped data
	for _, week := range weeklyDataKeamanans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keamanan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

		for _, data := range data_keamanan {
			var data_keamanan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

			result_keamanan = append(result_keamanan, models.KejadianKeamananImage{
				KejadianKeamanan: data,
				FileImage:        data_keamanan_image,
			})
		}
	}

	for _, week := range weeklyDataKeselamatans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keselamatan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)
		for _, data := range data_keselamatan {
			var data_keselamatan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

			result_keselamatan = append(result_keselamatan, models.KejadianKeselamatanImage{
				KejadianKeselamatan: data,
				FileImage:           data_keselamatan_image,
			})
		}
	}

	allWeeklyDataKeamananSorted := make([]string, 0, len(allWeeklyDataKeamanan))
	for name := range allWeeklyDataKeamanan {
		allWeeklyDataKeamananSorted = append(allWeeklyDataKeamananSorted, name)
	}
	allWeeklyDataKeselamatanSorted := make([]string, 0, len(allWeeklyDataKeselamatan))
	for name := range allWeeklyDataKeselamatan {
		allWeeklyDataKeselamatanSorted = append(allWeeklyDataKeselamatanSorted, name)
	}

	// Sort the slice of names
	sort.Strings(allWeeklyDataKeamananSorted)
	sort.Strings(allWeeklyDataKeselamatanSorted)

	var jenisKejadianKeamanan, jenisKejadianKeselamatan []models.JenisKejadian
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeamanan)
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeselamatan)

	kejadianKeamananWeek := make(map[string]map[string]int)
	kejadianKeselamatanWeek := make(map[string]map[string]int)
	for _, jenisKejadian := range jenisKejadianKeamanan {
		for _, name := range weekName {
			if kejadianKeamananWeek[jenisKejadian.NamaKejadian] == nil {
				kejadianKeamananWeek[jenisKejadian.NamaKejadian] = make(map[string]int)
			}

			if _, exists := kejadianKeamananWeek[jenisKejadian.NamaKejadian][name]; !exists {
				kejadianKeamananWeek[jenisKejadian.NamaKejadian][name] = 0
			}
		}
	}
	for _, jenisKejadian := range jenisKejadianKeselamatan {
		for _, name := range weekName {
			if kejadianKeselamatanWeek[jenisKejadian.NamaKejadian] == nil {
				kejadianKeselamatanWeek[jenisKejadian.NamaKejadian] = make(map[string]int)
			}

			if _, exists := kejadianKeselamatanWeek[jenisKejadian.NamaKejadian][name]; !exists {
				kejadianKeselamatanWeek[jenisKejadian.NamaKejadian][name] = 0
			}
		}
	}

	for _, key := range allWeeklyDataKeamananSorted {
		value := allWeeklyDataKeamanan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeamananWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}
	for _, key := range allWeeklyDataKeselamatanSorted {
		value := allWeeklyDataKeselamatan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeselamatanWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}

	countOfWeek := make(map[string]int)
	for i, weeks := range weekName {
		total := 0
		for _, kejadian := range kejadianKeamananWeek {
			total += kejadian[weeks]
		}
		name := "Minggu " + strconv.Itoa(i+1)
		countOfWeek[name] = total
	}
	// kejadianKeamananWeek, jenisKejadianKeamanan, weekName
	var images []string
	outputPath := fmt.Sprintf("storage/temp/pelanggaran%s.png", "default")

	templateDataKeamanan := struct {
		BaseURL                  string
		Bulan                    string
		BulanCapital             string
		Tahun                    string
		MingguKe                 int
		KejadianKeamananWeek     map[string]map[string]int
		WeekName                 []string
		CountOfWeek              map[string]int
		DataKejadianKeamananWeek []models.KejadianKeamananImage
	}{
		BaseURL:                  baseURL,
		Bulan:                    bulan,
		BulanCapital:             strings.ToUpper(bulan),
		Tahun:                    year,
		MingguKe:                 mingguKe,
		KejadianKeamananWeek:     kejadianKeamananWeek,
		WeekName:                 weekName,
		CountOfWeek:              countOfWeek,
		DataKejadianKeamananWeek: result_keamanan,
	}

	if err := r.ParseTemplate(templateKeamananHeadPath, newTemplateKeamananHeadPath, templateDataKeamanan); err == nil {
		// Generate Image
		success, _ := r.GenerateSlide(outputPath)
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// weeklyDataKeamananSorted := make([]string, 0, len(weeklyDataKeamanan))
	// for name := range weeklyDataKeamanan {
	// 	weeklyDataKeamananSorted = append(weeklyDataKeamananSorted, name)
	// }
	// weeklyDataKeselamatanSorted := make([]string, 0, len(weeklyDataKeselamatan))
	// for name := range weeklyDataKeselamatan {
	// 	weeklyDataKeselamatanSorted = append(weeklyDataKeselamatanSorted, name)
	// }

	// // Sort the slice of names
	// sort.Strings(weeklyDataKeamananSorted)
	// sort.Strings(weeklyDataKeselamatanSorted)

	for _, data := range result_keamanan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/pelanggaran%d.png", data.IdKejadianKeamanan)

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

		// html template data
		templateData := struct {
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
			SumberBerita     string
			Latitude         float64
			Longitude        float64
			Images           []models.FileImage
		}{
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
			SumberBerita:     data.LinkBerita,
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Images:           data.FileImage,
		}

		if err := r.ParseTemplate(templateKeamananPath, newTemplateKeamananPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	countOfWeek = make(map[string]int)
	for i, weeks := range weekName {
		total := 0
		for _, kejadian := range kejadianKeselamatanWeek {
			total += kejadian[weeks]
		}
		name := "Minggu " + strconv.Itoa(i+1)
		countOfWeek[name] = total
	}
	outputPath = fmt.Sprintf("storage/temp/kecelakaan%s.png", "default")

	templateDataKeselamatan := struct {
		BaseURL                     string
		Bulan                       string
		BulanCapital                string
		Tahun                       string
		MingguKe                    int
		KejadianKeselamatanWeek     map[string]map[string]int
		WeekName                    []string
		CountOfWeek                 map[string]int
		DataKejadianKeselamatanWeek []models.KejadianKeselamatanImage
	}{
		BaseURL:                     baseURL,
		Bulan:                       bulan,
		BulanCapital:                strings.ToUpper(bulan),
		Tahun:                       year,
		MingguKe:                    mingguKe,
		KejadianKeselamatanWeek:     kejadianKeselamatanWeek,
		WeekName:                    weekName,
		CountOfWeek:                 countOfWeek,
		DataKejadianKeselamatanWeek: result_keselamatan,
	}

	if err := r.ParseTemplate(templateKeselamatanHeadPath, newTemplateKeselamatanHeadPath, templateDataKeselamatan); err == nil {
		// Generate Image
		success, _ := r.GenerateSlide(outputPath)
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, data := range result_keselamatan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/kecelakaan%d.png", data.IdKejadianKeselamatan)

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
			return
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

		// html template data
		templateData := struct {
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
			SumberBerita     string
			Latitude         float64
			Longitude        float64
			Images           []models.FileImage
		}{
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
			SumberBerita:     data.LinkBerita,
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Images:           data.FileImage,
		}

		if err := r.ParseTemplate(templateKeselamatanPath, newTemplateKeselamatanPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	// Create a new PDF document
	pdf := gofpdf.New("L", "mm", "A4", "")

	for _, image := range images {
		// Add a new page to the PDF
		pdf.AddPage()

		// Get the image dimensions
		options := gofpdf.ImageOptions{
			ReadDpi: true,
		}
		info := pdf.RegisterImageOptions(image, options)
		width, height := info.Extent()

		// Calculate the position to center the image on the page
		pageWidth, pageHeight := pdf.GetPageSize()
		x := (pageWidth - width) / 2
		y := (pageHeight - height) / 2

		// Add the image to the PDF
		pdf.ImageOptions(image, x, y, width, height, false, options, 0, "")
	}

	// Save the PDF to a file
	// path := strconv.Itoa(req.TahunKe) + "/" + req.JenisLaporan + "/Bulan " + monthName(req.BulanKe)
	err := pdf.OutputFileAndClose("storage/laporan-keamanan-mingguan.pdf")
	if err != nil {
		fmt.Printf("Error saving PDF: %s", err)
	}

	fmt.Println("PDF created successfully!")
}
