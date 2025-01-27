package generator

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/lib/pq"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// LAPORAN MINGGUAN (DONE)
func (r *Pdf) LaporanMingguanLastMonth(date time.Time) {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)

	// html template path
	templateKeamananPath := "templates/keamanan.html"
	newTemplateKeamananPath := "keamanan.html"
	templateKeamananHeadPath := "templates/keamanan-head.html"
	newTemplateKeamananHeadPath := "keamanan-head.html"
	templateKeselamatanHeadPath := "templates/keselamatan-head.html"
	newTemplateKeselamatanHeadPath := "keselamatan-head.html"
	templateKeselamatanPath := "templates/keselamatan.html"
	newTemplateKeselamatanPath := "keselamatan.html"

	var year, bulan string
	var yearInt int
	var bulanInt time.Month
	if date.Month() == 1 {
		bulan = MonthNameIndonesia(12)
		bulanInt = 12
		year = strconv.Itoa(date.Year() - 1)
		yearInt = date.Year() - 1
	} else {
		bulan = MonthNameIndonesia(date.Month() - 1)
		bulanInt = date.Month() - 1
		year = strconv.Itoa(date.Year())
		yearInt = date.Year()
	}

	var firstDayMonth, startOfMonth time.Time
	firstDayMonth = time.Date(yearInt, bulanInt, 1, 0, 0, 0, 0, time.UTC)
	startOfMonth = firstDayMonth
	dayOfWeek := firstDayMonth.Weekday()

	fmt.Println(year, bulan)

	dayperweek := 7

	jumlahHari := DaysInMonth(startOfMonth)

	var minggu []time.Time

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

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var startOfWeek time.Time
	var endOfWeek time.Time
	mingguKe := len(minggu)
	startOfWeek = minggu[len(minggu)-2]
	endOfWeek = minggu[len(minggu)-1]
	fmt.Println(startOfMonth, startOfWeek, endOfWeek, mingguKe)

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

	var weeklyDataKeamanans, weeklyDataKeselamatans, allWeeklyDataKeamanans, allWeeklyDataKeselamatans WeeklyData

	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeselamatans)

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
	sort.Slice(result_keamanan, func(i, j int) bool {
		if result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian == result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian {
			return result_keamanan[i].KejadianKeamanan.Tanggal.ToStdTime().Before(result_keamanan[j].KejadianKeamanan.Tanggal.ToStdTime())
		}
		return result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian < result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian
	})

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
	sort.Slice(result_keselamatan, func(i, j int) bool {
		if result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian == result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian {
			return result_keselamatan[i].KejadianKeselamatan.Tanggal.ToStdTime().Before(result_keselamatan[j].KejadianKeselamatan.Tanggal.ToStdTime())
		}
		return result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian < result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian
	})

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
	var images []string
	outputPath := fmt.Sprintf("storage/temp/create-last-month/Mingguan/pelanggaran/pelanggaran-%s.pdf", "head-last-month")
	outputPathTempPelanggaran := fmt.Sprintf("storage/temp/create-last-month/Mingguan/pelanggaran")
	_ = os.MkdirAll(outputPathTempPelanggaran, 0755)

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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_create-last-month/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keamanan_arr []int64
	for _, data := range result_keamanan {
		id_keamanan_arr = append(id_keamanan_arr, data.IdKejadianKeamanan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/create-last-month/Mingguan/pelanggaran/pelanggaran-%d.pdf", data.IdKejadianKeamanan)

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
			KejadianKeamananWeek map[string]map[string]int
			BaseURL              string
			Title                string
			NamaKapal            string
			Kejadian             string
			Penyebab             string
			Lokasi               string
			ABK                  string
			Muatan               string
			InstansiPenindak     string
			Keterangan           string
			Waktu                string
			SumberBerita         string
			Latitude             float64
			Longitude            float64
			Images               []models.FileImage
		}{
			KejadianKeamananWeek: kejadianKeamananWeek,
			BaseURL:              baseURL,
			Title:                data.JenisKejadian.NamaKejadian,
			NamaKapal:            data.NamaKapal,
			Kejadian:             data.JenisKejadian.NamaKejadian,
			Penyebab:             "-",
			Lokasi:               data.LokasiKejadian,
			ABK:                  abk,
			Muatan:               data.Muatan,
			InstansiPenindak:     data.SumberBerita,
			Keterangan:           data.TindakLanjut,
			Waktu:                data.Tanggal.ToDateString(),
			SumberBerita:         data.LinkBerita,
			Latitude:             data.Latitude,
			Longitude:            data.Longitude,
			Images:               data.FileImage,
		}

		if err := r.ParseTemplate(templateKeamananPath, newTemplateKeamananPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_create-last-month/")
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
	outputPath = fmt.Sprintf("storage/temp/create-last-month/Mingguan/kecelakaan/kecelakaan-%s.pdf", "head-last-month")
	outputPathTempKecelakaan := fmt.Sprintf("storage/temp/create-last-month/Mingguan/kecelakaan")

	_ = os.MkdirAll(outputPathTempKecelakaan, 0755)

	templateDataKeselamatan := struct {
		KejadianKeamananLength      int
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
		KejadianKeamananLength:      len(kejadianKeamananWeek),
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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_create-last-month/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keselamatan_arr []int64
	for _, data := range result_keselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, data.IdKejadianKeselamatan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/create-last-month/Mingguan/kecelakaan/kecelakaan-%d.pdf", data.IdKejadianKeselamatan)

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
			KejadianKeamananLength  int
			KejadianKeselamatanWeek map[string]map[string]int
			BaseURL                 string
			Title                   string
			NamaKapal               string
			Kejadian                string
			Penyebab                string
			Lokasi                  string
			Korban                  string
			Perpindahan             string
			Keterangan              string
			Waktu                   string
			InstansiPenindak        string
			SumberBerita            string
			Latitude                float64
			Longitude               float64
			Images                  []models.FileImage
		}{
			KejadianKeamananLength:  len(kejadianKeamananWeek),
			KejadianKeselamatanWeek: kejadianKeselamatanWeek,
			BaseURL:                 baseURL,
			Title:                   data.JenisKejadian.NamaKejadian,
			NamaKapal:               data.NamaKapal,
			Kejadian:                data.JenisKejadian.NamaKejadian,
			Penyebab:                data.Penyebab,
			Lokasi:                  data.LokasiKejadian,
			Korban:                  korban,
			Perpindahan:             perpindahan,
			Keterangan:              data.TindakLanjut,
			Waktu:                   data.Tanggal.ToDateString(),
			InstansiPenindak:        data.SumberBerita,
			SumberBerita:            data.LinkBerita,
			Latitude:                data.Latitude,
			Longitude:               data.Longitude,
			Images:                  data.FileImage,
		}

		if err := r.ParseTemplate(templateKeselamatanPath, newTemplateKeselamatanPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_create-last-month/")
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	// ttd
	tgl := fmt.Sprintf("%d %s %d", date.Day(), MonthNameIndonesia(date.Month()), date.Year())
	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	templatePath := "templates/ttd-mingguan.html"
	newTemplatePath := "ttd-mingguan.html"
	outputTtdPath := "storage/temp/ttd/output-ttd-mingguan-acc.pdf"
	outputPathTempTtd := fmt.Sprintf("storage/temp/ttd")
	_ = os.MkdirAll(outputPathTempTtd, 0755)

	templateData := struct {
		BaseURL    string
		Tanggal    string
		Jabatan    string
		Nama       string
		IsApproved bool
		Ttd        *string
		Nik        string
	}{
		BaseURL:    baseURL,
		Tanggal:    tgl,
		Jabatan:    deputi.Jabatan.Name,
		Nama:       deputi.Name,
		IsApproved: false == true,
		Ttd:        deputi.Ttd,
		Nik:        deputi.EmpNo,
	}

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		// Generate Image
		success, _ := r.GenerateSlide(outputTtdPath, "Mingguan_create-last-month/")
		if success {
			images = append(images, outputTtdPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Save the PDF to a file
	// Create directory path
	yearStr := strconv.Itoa(yearInt)
	rootPath := "storage/app/"
	dirPath := yearStr + "/Laporan Mingguan/Bulan " + bulan + "/"
	fullPath := rootPath + dirPath
	_ = os.MkdirAll(fullPath, 0755)

	lastTwoDigits := yearStr[len(yearStr)-2:]
	nameFile := "LAP MING KE-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishMap[bulanInt] + "'" + lastTwoDigits + ".pdf"
	path := fullPath + nameFile
	// jenisFile := Laporan Mingguan

	err := api.MergeCreateFile(images, path, false, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return
	}

	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	jenisFile := "Laporan Mingguan"
	namaLaporan := "Laporan Minggu ke-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishTitleCase(bulanInt) + "'" + lastTwoDigits
	document := dirPath + nameFile
	document = strings.ReplaceAll(document, "/", "\\")

	// save to Laporan
	laporan := models.Laporan{
		NamaLaporan:  namaLaporan,
		JenisLaporan: jenisFile,
		TahunKe:      yearInt,
		BulanKe:      int(bulanInt),
		MingguKe:     mingguKe,
		Dokumen:      document,
	}

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", result_keamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	if err := facades.Orm().Query().Create(&laporan); err != nil {
		fmt.Println("Data Gagal Ditambahkan:", err)
		return
	}

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Create(&approval)

	_ = os.RemoveAll(outputPathTempPelanggaran)
	_ = os.RemoveAll(outputPathTempKecelakaan)
	_ = os.RemoveAll(outputPathTempTtd)

	fmt.Println("PDF created successfully!")
}
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
	bulan := MonthNameIndonesia(now.Month())
	year := strconv.Itoa(now.Year())

	dayperweek := 7

	jumlahHari := DaysInMonth(now)

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
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var startOfWeek time.Time
	var endOfWeek time.Time
	var mingguKe int
	for i, weekEnd := range minggu {
		if weekEnd.After(now) {
			if i == 0 {
				r.LaporanMingguanLastMonth(now)
				return
			} else if i == 1 {
				if minggu[i-1].Day() < 7 {
					r.LaporanMingguanLastMonth(now)
				}
				mingguKe = i
				startOfWeek = startOfMonth
				endOfWeek = minggu[i-1]
			} else {
				mingguKe = i
				startOfWeek = minggu[i-2]
				endOfWeek = minggu[i-1]
			}
			break
		} else if weekEnd.Equal(now) {
			mingguKe = i
			startOfWeek = minggu[i-2]
			endOfWeek = minggu[i-1]
			break
		}
	}

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

	var weeklyDataKeamanans, weeklyDataKeselamatans, allWeeklyDataKeamanans, allWeeklyDataKeselamatans WeeklyData

	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeselamatans)

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
	sort.Slice(result_keamanan, func(i, j int) bool {
		if result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian == result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian {
			return result_keamanan[i].KejadianKeamanan.Tanggal.ToStdTime().Before(result_keamanan[j].KejadianKeamanan.Tanggal.ToStdTime())
		}
		return result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian < result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian
	})

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
	sort.Slice(result_keselamatan, func(i, j int) bool {
		if result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian == result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian {
			return result_keselamatan[i].KejadianKeselamatan.Tanggal.ToStdTime().Before(result_keselamatan[j].KejadianKeselamatan.Tanggal.ToStdTime())
		}
		return result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian < result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian
	})

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
	outputPath := fmt.Sprintf("storage/temp/create/Mingguan/pelanggaran/pelanggaran-%s.pdf", "head")
	outputPathTempPelanggaran := fmt.Sprintf("storage/temp/create/Mingguan/pelanggaran")
	_ = os.MkdirAll(outputPathTempPelanggaran, 0755)

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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_create/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keamanan_arr []int64
	for _, data := range result_keamanan {
		id_keamanan_arr = append(id_keamanan_arr, data.IdKejadianKeamanan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/create/Mingguan/pelanggaran/pelanggaran-%d.pdf", data.IdKejadianKeamanan)

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
			KejadianKeamananWeek map[string]map[string]int
			BaseURL              string
			Title                string
			NamaKapal            string
			Kejadian             string
			Penyebab             string
			Lokasi               string
			ABK                  string
			Muatan               string
			InstansiPenindak     string
			Keterangan           string
			Waktu                string
			SumberBerita         string
			Latitude             float64
			Longitude            float64
			Images               []models.FileImage
		}{
			KejadianKeamananWeek: kejadianKeamananWeek,
			BaseURL:              baseURL,
			Title:                data.JenisKejadian.NamaKejadian,
			NamaKapal:            data.NamaKapal,
			Kejadian:             data.JenisKejadian.NamaKejadian,
			Penyebab:             "-",
			Lokasi:               data.LokasiKejadian,
			ABK:                  abk,
			Muatan:               data.Muatan,
			InstansiPenindak:     data.SumberBerita,
			Keterangan:           data.TindakLanjut,
			Waktu:                data.Tanggal.ToDateString(),
			SumberBerita:         data.LinkBerita,
			Latitude:             data.Latitude,
			Longitude:            data.Longitude,
			Images:               data.FileImage,
		}

		if err := r.ParseTemplate(templateKeamananPath, newTemplateKeamananPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_create/")
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
	outputPath = fmt.Sprintf("storage/temp/create/Mingguan/kecelakaan/kecelakaan-%s.pdf", "head")
	outputPathTempKecelakaan := fmt.Sprintf("storage/temp/create/Mingguan/kecelakaan")
	_ = os.MkdirAll(outputPathTempKecelakaan, 0755)

	templateDataKeselamatan := struct {
		KejadianKeamananLength      int
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
		KejadianKeamananLength:      len(kejadianKeamananWeek),
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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_create/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keselamatan_arr []int64
	for _, data := range result_keselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, data.IdKejadianKeselamatan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/create/Mingguan/kecelakaan/kecelakaan-%d.pdf", data.IdKejadianKeselamatan)

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
			KejadianKeamananLength  int
			KejadianKeselamatanWeek map[string]map[string]int
			BaseURL                 string
			Title                   string
			NamaKapal               string
			Kejadian                string
			Penyebab                string
			Lokasi                  string
			Korban                  string
			Perpindahan             string
			Keterangan              string
			Waktu                   string
			InstansiPenindak        string
			SumberBerita            string
			Latitude                float64
			Longitude               float64
			Images                  []models.FileImage
		}{
			KejadianKeamananLength:  len(kejadianKeamananWeek),
			KejadianKeselamatanWeek: kejadianKeselamatanWeek,
			BaseURL:                 baseURL,
			Title:                   data.JenisKejadian.NamaKejadian,
			NamaKapal:               data.NamaKapal,
			Kejadian:                data.JenisKejadian.NamaKejadian,
			Penyebab:                data.Penyebab,
			Lokasi:                  data.LokasiKejadian,
			Korban:                  korban,
			Perpindahan:             perpindahan,
			Keterangan:              data.TindakLanjut,
			Waktu:                   data.Tanggal.ToDateString(),
			InstansiPenindak:        data.SumberBerita,
			SumberBerita:            data.LinkBerita,
			Latitude:                data.Latitude,
			Longitude:               data.Longitude,
			Images:                  data.FileImage,
		}

		if err := r.ParseTemplate(templateKeselamatanPath, newTemplateKeselamatanPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_create/")
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	// ttd
	date := fmt.Sprintf("%d %s %d", now.Day(), bulan, now.Year())
	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	templatePath := "templates/ttd-mingguan.html"
	newTemplatePath := "ttd-mingguan.html"
	outputTtdPath := "storage/temp/ttd/output-ttd-mingguan-acc.pdf"
	outputPathTempTtd := fmt.Sprintf("storage/temp/ttd")
	_ = os.MkdirAll(outputPathTempTtd, 0755)

	templateData := struct {
		BaseURL    string
		Tanggal    string
		Jabatan    string
		Nama       string
		IsApproved bool
		Ttd        *string
		Nik        string
	}{
		BaseURL:    baseURL,
		Tanggal:    date,
		Jabatan:    deputi.Jabatan.Name,
		Nama:       deputi.Name,
		IsApproved: false == true,
		Ttd:        deputi.Ttd,
		Nik:        deputi.EmpNo,
	}

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		// Generate Image
		success, _ := r.GenerateSlide(outputTtdPath, "Mingguan_create/")
		if success {
			images = append(images, outputTtdPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Mingguan/Bulan " + bulan + "/"
	fullPath := rootPath + dirPath
	_ = os.MkdirAll(fullPath, 0755)

	lastTwoDigits := year[len(year)-2:]
	nameFile := "LAP MING KE-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishMap[now.Month()] + "'" + lastTwoDigits + ".pdf"
	path := fullPath + nameFile

	err := api.MergeCreateFile(images, path, false, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return
	}

	var kejadianKeamanan models.KejadianKeamanan

	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	jenisFile := "Laporan Mingguan"
	namaLaporan := "Laporan Minggu ke-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishTitleCase(now.Month()) + "'" + lastTwoDigits
	document := dirPath + nameFile
	document = strings.ReplaceAll(document, "/", "\\")

	// save to Laporan
	laporan := models.Laporan{
		NamaLaporan:  namaLaporan,
		JenisLaporan: jenisFile,
		TahunKe:      now.Year(),
		BulanKe:      int(now.Month()),
		MingguKe:     mingguKe,
		Dokumen:      document,
	}

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", result_keamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	if err := facades.Orm().Query().Create(&laporan); err != nil {
		fmt.Println("Data Gagal Ditambahkan:", err)
		return
	}

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Create(&approval)

	_ = os.RemoveAll(outputPathTempPelanggaran)
	_ = os.RemoveAll(outputPathTempKecelakaan)
	_ = os.RemoveAll(outputPathTempTtd)

	fmt.Println("PDF created successfully!")
}
func (r *Pdf) LaporanMingguanUpdate(id_laporan int64, weeksTo int, month int, years int) {
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

	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	bulan := MonthNameIndonesia(now.Month())
	year := strconv.Itoa(now.Year())

	dayperweek := 7

	jumlahHari := DaysInMonth(now)

	var minggu []time.Time

	// Calculate the first day of the month
	var firstDayMonth, startOfMonth time.Time
	firstDayMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfMonth = firstDayMonth
	dayOfWeek := firstDayMonth.Weekday()
	fmt.Println(dayOfWeek)

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
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var startOfWeek time.Time
	var endOfWeek time.Time
	var mingguKe int
	for i, weekEnd := range minggu {
		if i == weeksTo {
			mingguKe = i
			if i == 1 {
				startOfWeek = startOfMonth
				endOfWeek = minggu[i-1]
				fmt.Println(i, weekEnd, mingguKe)
			} else {
				startOfWeek = minggu[i-2]
				endOfWeek = minggu[i-1]
				fmt.Println(i, weekEnd, mingguKe)
			}
			break
		}
	}

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

	var weeklyDataKeamanans, weeklyDataKeselamatans, allWeeklyDataKeamanans, allWeeklyDataKeselamatans WeeklyData

	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&allWeeklyDataKeselamatans)

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
	sort.Slice(result_keamanan, func(i, j int) bool {
		if result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian == result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian {
			return result_keamanan[i].KejadianKeamanan.Tanggal.ToStdTime().Before(result_keamanan[j].KejadianKeamanan.Tanggal.ToStdTime())
		}
		return result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian < result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian
	})

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
	sort.Slice(result_keselamatan, func(i, j int) bool {
		if result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian == result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian {
			return result_keselamatan[i].KejadianKeselamatan.Tanggal.ToStdTime().Before(result_keselamatan[j].KejadianKeselamatan.Tanggal.ToStdTime())
		}
		return result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian < result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian
	})

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
	outputPath := fmt.Sprintf("storage/temp/update/Mingguan/pelanggaran/pelanggaran-%s.pdf", "head")
	outputPathTempPelanggaran := fmt.Sprintf("storage/temp/update/Mingguan/pelanggaran")
	_ = os.MkdirAll(outputPathTempPelanggaran, 0755)

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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_update/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keamanan_arr []int64
	for _, data := range result_keamanan {
		id_keamanan_arr = append(id_keamanan_arr, data.IdKejadianKeamanan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/update/Mingguan/pelanggaran/pelanggaran-%d.pdf", data.IdKejadianKeamanan)

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
			KejadianKeamananWeek map[string]map[string]int
			BaseURL              string
			Title                string
			NamaKapal            string
			Kejadian             string
			Penyebab             string
			Lokasi               string
			ABK                  string
			Muatan               string
			InstansiPenindak     string
			Keterangan           string
			Waktu                string
			SumberBerita         string
			Latitude             float64
			Longitude            float64
			Images               []models.FileImage
		}{
			KejadianKeamananWeek: kejadianKeamananWeek,
			BaseURL:              baseURL,
			Title:                data.JenisKejadian.NamaKejadian,
			NamaKapal:            data.NamaKapal,
			Kejadian:             data.JenisKejadian.NamaKejadian,
			Penyebab:             "-",
			Lokasi:               data.LokasiKejadian,
			ABK:                  abk,
			Muatan:               data.Muatan,
			InstansiPenindak:     data.SumberBerita,
			Keterangan:           data.TindakLanjut,
			Waktu:                data.Tanggal.ToDateString(),
			SumberBerita:         data.LinkBerita,
			Latitude:             data.Latitude,
			Longitude:            data.Longitude,
			Images:               data.FileImage,
		}

		if err := r.ParseTemplate(templateKeamananPath, newTemplateKeamananPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_update/")
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
	outputPath = fmt.Sprintf("storage/temp/update/Mingguan/kecelakaan/kecelakaan-%s.pdf", "head")
	outputPathTempKecelakaan := fmt.Sprintf("storage/temp/update/Mingguan/kecelakaan")
	_ = os.MkdirAll(outputPathTempKecelakaan, 0755)

	templateDataKeselamatan := struct {
		KejadianKeamananLength      int
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
		KejadianKeamananLength:      len(kejadianKeamananWeek),
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
		success, _ := r.GenerateSlide(outputPath, "Mingguan_update/")
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var id_keselamatan_arr []int64
	for _, data := range result_keselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, data.IdKejadianKeselamatan)
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/update/Mingguan/kecelakaan/kecelakaan-%d.pdf", data.IdKejadianKeselamatan)

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
			KejadianKeamananLength  int
			KejadianKeselamatanWeek map[string]map[string]int
			BaseURL                 string
			Title                   string
			NamaKapal               string
			Kejadian                string
			Penyebab                string
			Lokasi                  string
			Korban                  string
			Perpindahan             string
			Keterangan              string
			Waktu                   string
			InstansiPenindak        string
			SumberBerita            string
			Latitude                float64
			Longitude               float64
			Images                  []models.FileImage
		}{
			KejadianKeamananLength:  len(kejadianKeamananWeek),
			KejadianKeselamatanWeek: kejadianKeselamatanWeek,
			BaseURL:                 baseURL,
			Title:                   data.JenisKejadian.NamaKejadian,
			NamaKapal:               data.NamaKapal,
			Kejadian:                data.JenisKejadian.NamaKejadian,
			Penyebab:                data.Penyebab,
			Lokasi:                  data.LokasiKejadian,
			Korban:                  korban,
			Perpindahan:             perpindahan,
			Keterangan:              data.TindakLanjut,
			Waktu:                   data.Tanggal.ToDateString(),
			InstansiPenindak:        data.SumberBerita,
			SumberBerita:            data.LinkBerita,
			Latitude:                data.Latitude,
			Longitude:               data.Longitude,
			Images:                  data.FileImage,
		}

		if err := r.ParseTemplate(templateKeselamatanPath, newTemplateKeselamatanPath, templateData); err == nil {
			// Generate Image
			success, _ := r.GenerateSlide(outputPath, "Mingguan_update/")
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	// ttd
	newDate := time.Now()
	date := fmt.Sprintf("%d %s %d", newDate.Day(), MonthNameIndonesia(newDate.Month()), newDate.Year())
	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	templatePath := "templates/ttd-mingguan.html"
	newTemplatePath := "ttd-mingguan.html"
	outputTtdPath := "storage/temp/ttd/output-ttd-mingguan-acc.pdf"
	outputPathTempTtd := fmt.Sprintf("storage/temp/ttd")
	_ = os.MkdirAll(outputPathTempTtd, 0755)

	templateData := struct {
		BaseURL    string
		Tanggal    string
		Jabatan    string
		Nama       string
		IsApproved bool
		Ttd        *string
		Nik        string
	}{
		BaseURL:    baseURL,
		Tanggal:    date,
		Jabatan:    deputi.Jabatan.Name,
		Nama:       deputi.Name,
		IsApproved: false == true,
		Ttd:        deputi.Ttd,
		Nik:        deputi.EmpNo,
	}

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		// Generate Image
		success, _ := r.GenerateSlide(outputTtdPath, "Mingguan_update/")
		if success {
			fmt.Println("ttd")
			images = append(images, outputTtdPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Mingguan/Bulan " + bulan + "/"
	fullPath := rootPath + dirPath
	_ = os.MkdirAll(fullPath, 0755)

	lastTwoDigits := year[len(year)-2:]
	nameFile := "LAP MING KE-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishMap[now.Month()] + "'" + lastTwoDigits + ".pdf"
	path := fullPath + nameFile

	err := api.MergeCreateFile(images, path, false, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return
	}

	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	jenisFile := "Laporan Mingguan"
	namaLaporan := "Laporan Minggu ke-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishTitleCase(now.Month()) + "'" + lastTwoDigits
	document := dirPath + nameFile
	document = strings.ReplaceAll(document, "/", "\\")

	// save to Laporan
	var laporan models.Laporan
	facades.Orm().Query().Where("id_laporan=?", id_laporan).First(&laporan)
	laporan.NamaLaporan = namaLaporan
	laporan.JenisLaporan = jenisFile
	laporan.TahunKe = now.Year()
	laporan.BulanKe = int(now.Month())
	laporan.MingguKe = mingguKe
	laporan.Dokumen = document

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", result_keamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	facades.Orm().Query().Save(&laporan)

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}
	facades.Orm().Query().Delete(&models.Approval{}, "laporan_id = ?", laporan.IDLaporan)

	facades.Orm().Query().Create(&approval)

	_ = os.RemoveAll(outputPathTempPelanggaran)
	_ = os.RemoveAll(outputPathTempKecelakaan)
	_ = os.RemoveAll(outputPathTempTtd)

	fmt.Println("PDF created successfully!")
}

// LAPORAN BULANAN
func (r *Pdf) LaporanBulanan() {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)
	now := time.Now()
	bulan := MonthNameIndonesia(now.Month())
	year := strconv.Itoa(now.Year())
	dayperweek := 7

	jumlahHari := DaysInMonth(now)

	var minggu []time.Time

	// Calculate the first day of the month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dayOfWeek := firstDay.Weekday()

	var periodeTanggal []string
	periodeTanggal = append(periodeTanggal, "01")
	periodeTanggal = append(periodeTanggal, strconv.Itoa(jumlahHari))

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

	nextWeek := firstDay.AddDate(0, 0, daysToAdd)
	minggu = append(minggu, nextWeek)

	jumlahHari -= daysToAdd + 1
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var weekName []string
	for i, weekEnd := range minggu {
		var setName, fullName string
		if weekEnd.Day() < 10 || (i > 0 && minggu[i-1].AddDate(0, 0, 1).Day() < 10) {
			setName = "0"
		}
		if i == 0 {
			fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
		} else {
			fullName = setName + strconv.Itoa(minggu[i-1].AddDate(0, 0, 1).Day()) + "-" +
				strconv.Itoa(weekEnd.Day()) + " " + bulan
		}
		weekName = append(weekName, fullName)
	}

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var weeklyDataKeamanans, weeklyDataKeselamatans WeeklyData

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	weeklyDataKeamanan := make(map[string][]models.KejadianKeamanan)
	weeklyDataKeselamatan := make(map[string][]models.KejadianKeselamatan)

	var data_keamanan []models.KejadianKeamanan
	var data_keselamatan []models.KejadianKeselamatan
	// Now you can loop through the grouped data
	for _, week := range weeklyDataKeamanans {
		fmt.Printf("Week starting %s: IDs %v\n", week.WeekStart.Format("2006-01-02"), week.KejadianIDs)
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
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + bulan
				}

				weeklyDataKeamanan[fullName] = append(weeklyDataKeamanan[fullName], data_keamanan...)
				break
			}
		}
	}
	for _, week := range weeklyDataKeselamatans {
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
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + bulan
				}

				weeklyDataKeselamatan[fullName] = append(weeklyDataKeselamatan[fullName], data_keselamatan...)
				break
			}
		}
	}

	weeklyDataKeamananSorted := make([]string, 0, len(weeklyDataKeamanan))
	for name := range weeklyDataKeamanan {
		weeklyDataKeamananSorted = append(weeklyDataKeamananSorted, name)
	}
	weeklyDataKeselamatanSorted := make([]string, 0, len(weeklyDataKeselamatan))
	for name := range weeklyDataKeselamatan {
		weeklyDataKeselamatanSorted = append(weeklyDataKeselamatanSorted, name)
	}

	// Sort the slice of names
	sort.Strings(weeklyDataKeamananSorted)
	sort.Strings(weeklyDataKeselamatanSorted)

	var jenisKejadianKeamanan, jenisKejadianKeselamatan []models.JenisKejadian
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeamanan)
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeselamatan)

	facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

	facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)

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

	for _, key := range weeklyDataKeamananSorted {
		value := weeklyDataKeamanan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeamananWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}
	for _, key := range weeklyDataKeselamatanSorted {
		value := weeklyDataKeselamatan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeselamatanWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range data_keamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	// Create a slice of keys (NamaKejadian values)
	sortedNames := make([]string, 0, len(groupedByJenisKeamanan))
	for name := range groupedByJenisKeamanan {
		sortedNames = append(sortedNames, name)
	}

	// Sort the slice of names
	sort.Strings(sortedNames)

	var groupKeamanan []GroupingKeamanan
	var keamananBarat []GroupingKeamananBarat
	var keamananTimur []GroupingKeamananTimur
	var keamananTengah []GroupingKeamananTengah
	var groupKeamananBarat []models.KejadianKeamanan
	var groupKeamananTimur []models.KejadianKeamanan
	var groupKeamananTengah []models.KejadianKeamanan
	// Print the grouped data
	for _, jenisName := range sortedNames {
		kejadianGroup := groupedByJenisKeamanan[jenisName]
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		for _, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				jumlahTengah++
				groupKeamananTengah = append(groupKeamananTengah, index)
			}
			jumlah++
		}

		if jumlahBarat != 0 {
			keamananBarat = append(keamananBarat, GroupingKeamananBarat{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananBarat,
				Jumlah:           jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keamananTimur = append(keamananTimur, GroupingKeamananTimur{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTimur,
				Jumlah:           jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keamananTengah = append(keamananTengah, GroupingKeamananTengah{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTengah,
				Jumlah:           jumlahTengah,
			})
		}

		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
			NamaKejadian:     jenisName,
			KejadianKeamanan: kejadianGroup,
			Jumlah:           jumlah,
			JumlahZonaBarat:  jumlahBarat,
			JumlahZonaTimur:  jumlahTimur,
			JumlahZonaTengah: jumlahTengah,
		})
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
	for _, kejadian := range data_keselamatan {
		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	sortedNames = make([]string, 0, len(groupedByJenisKeselamatan))
	for name := range groupedByJenisKeselamatan {
		sortedNames = append(sortedNames, name)
	}

	// Sort the slice of names
	sort.Strings(sortedNames)

	var groupKeselamatan []GroupingKeselamatan
	var keselamatanBarat []GroupingKeselamatanBarat
	var keselamatanTimur []GroupingKeselamatanTimur
	var keselamatanTengah []GroupingKeselamatanTengah
	var groupKeselamatanBarat []models.KejadianKeselamatanKorban
	var groupKeselamatanTimur []models.KejadianKeselamatanKorban
	var groupKeselamatanTengah []models.KejadianKeselamatanKorban
	// Print the grouped data
	for _, jenisName := range sortedNames {
		kejadianGroup := groupedByJenisKeselamatan[jenisName]
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		var list_korban []models.KejadianKeselamatanKorban

		for _, data := range kejadianGroup {
			var x models.ListKorban
			err := json.Unmarshal(data.Korban, &x)
			if err != nil {
				return
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			if data.Zona == "BARAT" {
				jumlahBarat++
				groupKeselamatanBarat = append(groupKeselamatanBarat, temp)
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
				groupKeselamatanTimur = append(groupKeselamatanTimur, temp)
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
				groupKeselamatanTengah = append(groupKeselamatanTengah, temp)
			}

			list_korban = append(list_korban, temp)
			jumlah++
		}

		if jumlahBarat != 0 {
			keselamatanBarat = append(keselamatanBarat, GroupingKeselamatanBarat{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanBarat,
				Jumlah:              jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keselamatanTimur = append(keselamatanTimur, GroupingKeselamatanTimur{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTimur,
				Jumlah:              jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keselamatanTengah = append(keselamatanTengah, GroupingKeselamatanTengah{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTengah,
				Jumlah:              jumlahTengah,
			})
		}

		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
			NamaKejadian:        jenisName,
			KejadianKeselamatan: list_korban,
			Jumlah:              jumlah,
			JumlahZonaBarat:     jumlahBarat,
			JumlahZonaTimur:     jumlahTimur,
			JumlahZonaTengah:    jumlahTengah,
		})
	}

	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	firstDayNextMonth := now.AddDate(0, 1, -now.Day()+1)
	date := fmt.Sprintf("%d %s %d", firstDayNextMonth.Day(), MonthNameIndonesia(firstDayNextMonth.Month()), firstDayNextMonth.Year())
	fmt.Printf("%d %s %d\n", firstDayNextMonth.Day(), MonthNameIndonesia(firstDayNextMonth.Month()), firstDayNextMonth.Year())
	// html template data
	templateData := struct {
		BaseURL                   string
		Tanggal                   string
		Jabatan                   string
		Nama                      string
		IsApproved                bool
		Ttd                       *string
		Nik                       string
		Bulan                     string
		BulanCapital              string
		Tahun                     string
		JumlahKejadianKeamanan    int
		JumlahKejadianKeselamatan int
		KejadianKeamanan          []GroupingKeamanan
		KejadianKeamananBarat     []GroupingKeamananBarat
		KejadianKeamananTimur     []GroupingKeamananTimur
		KejadianKeamananTengah    []GroupingKeamananTengah
		KejadianKeselamatan       []GroupingKeselamatan
		KejadianKeselamatanBarat  []GroupingKeselamatanBarat
		KejadianKeselamatanTimur  []GroupingKeselamatanTimur
		KejadianKeselamatanTengah []GroupingKeselamatanTengah
		KejadianKeamananWeek      map[string]map[string]int
		KejadianKeselamatanWeek   map[string]map[string]int
		JenisKejadianKeamanan     []models.JenisKejadian
		JenisKejadianKeselamatan  []models.JenisKejadian
		WeekName                  []string
		PeriodeTanggal            []string
	}{
		BaseURL:                   baseURL,
		Tanggal:                   date,
		Jabatan:                   deputi.Jabatan.Name,
		Nama:                      deputi.Name,
		IsApproved:                false == true,
		Ttd:                       deputi.Ttd,
		Nik:                       deputi.EmpNo,
		Bulan:                     bulan,
		BulanCapital:              strings.ToUpper(bulan),
		Tahun:                     year,
		JumlahKejadianKeamanan:    len(data_keamanan),
		JumlahKejadianKeselamatan: len(data_keselamatan),
		KejadianKeamanan:          groupKeamanan,
		KejadianKeamananBarat:     keamananBarat,
		KejadianKeamananTimur:     keamananTimur,
		KejadianKeamananTengah:    keamananTengah,
		KejadianKeselamatan:       groupKeselamatan,
		KejadianKeselamatanBarat:  keselamatanBarat,
		KejadianKeselamatanTimur:  keselamatanTimur,
		KejadianKeselamatanTengah: keselamatanTengah,
		KejadianKeamananWeek:      kejadianKeamananWeek,
		KejadianKeselamatanWeek:   kejadianKeselamatanWeek,
		JenisKejadianKeamanan:     jenisKejadianKeamanan,
		JenisKejadianKeselamatan:  jenisKejadianKeselamatan,
		WeekName:                  weekName,
		PeriodeTanggal:            periodeTanggal,
	}

	//html template path
	templatePath := "templates/laporan-bulanan.html"
	newTemplatePath := "laporan-bulanan.html"

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Bulanan/"
	fullPath := rootPath + dirPath

	_ = os.MkdirAll(fullPath, 0755)

	nameFile := strconv.Itoa(int(now.Month())) + ". Laporan bulan " + bulan + " Keamanan & Keselamatan di Wilayah Perairan Indonesia"
	path := fullPath + nameFile + ".pdf"

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(path, 0, "Bulanan_create/")
	} else {
		fmt.Println(err)
	}

	jenisFile := "Laporan Bulanan"
	document := dirPath + nameFile + ".pdf"

	// save to Laporan
	laporan := models.Laporan{
		NamaLaporan:  nameFile,
		JenisLaporan: jenisFile,
		TahunKe:      now.Year(),
		BulanKe:      int(now.Month()),
		Dokumen:      document,
	}

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", data_keamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	if err := facades.Orm().Query().Create(&laporan); err != nil {
		fmt.Println("Data Gagal Ditambahkan:", err)
		return
	}

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Create(&approval)

	fmt.Println("PDF created successfully!")
}
func (r *Pdf) LaporanBulananUpdate(id_laporan int64, month int, years int) {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)
	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	bulan := MonthNameIndonesia(now.Month())
	year := strconv.Itoa(now.Year())
	dayperweek := 7

	jumlahHari := DaysInMonth(now)

	var minggu []time.Time

	// Calculate the first day of the month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dayOfWeek := firstDay.Weekday()

	var periodeTanggal []string
	periodeTanggal = append(periodeTanggal, "01")
	periodeTanggal = append(periodeTanggal, strconv.Itoa(jumlahHari))

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

	nextWeek := firstDay.AddDate(0, 0, daysToAdd)
	minggu = append(minggu, nextWeek)

	jumlahHari -= daysToAdd + 1
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var weekName []string
	for i, weekEnd := range minggu {
		var setName, fullName string
		if weekEnd.Day() < 10 || (i > 0 && minggu[i-1].AddDate(0, 0, 1).Day() < 10) {
			setName = "0"
		}
		if i == 0 {
			fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
		} else {
			fullName = setName + strconv.Itoa(minggu[i-1].AddDate(0, 0, 1).Day()) + "-" +
				strconv.Itoa(weekEnd.Day()) + " " + bulan
		}
		weekName = append(weekName, fullName)
	}

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var weeklyDataKeamanans, weeklyDataKeselamatans WeeklyData

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	weeklyDataKeamanan := make(map[string][]models.KejadianKeamanan)
	weeklyDataKeselamatan := make(map[string][]models.KejadianKeselamatan)

	var data_keamanan []models.KejadianKeamanan
	var data_keselamatan []models.KejadianKeselamatan
	// Now you can loop through the grouped data
	for _, week := range weeklyDataKeamanans {
		fmt.Printf("Week starting %s: IDs %v\n", week.WeekStart.Format("2006-01-02"), week.KejadianIDs)
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
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + bulan
				}

				weeklyDataKeamanan[fullName] = append(weeklyDataKeamanan[fullName], data_keamanan...)
				break
			}
		}
	}
	for _, week := range weeklyDataKeselamatans {
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
					fullName = "01-" + setName + strconv.Itoa(weekEnd.Day()) + " " + bulan
				} else {
					fullName = setName + strconv.Itoa(week.WeekStart.Day()) + "-" +
						strconv.Itoa(weekEnd.Day()) + " " + bulan
				}

				weeklyDataKeselamatan[fullName] = append(weeklyDataKeselamatan[fullName], data_keselamatan...)
				break
			}
		}
	}

	weeklyDataKeamananSorted := make([]string, 0, len(weeklyDataKeamanan))
	for name := range weeklyDataKeamanan {
		weeklyDataKeamananSorted = append(weeklyDataKeamananSorted, name)
	}
	weeklyDataKeselamatanSorted := make([]string, 0, len(weeklyDataKeselamatan))
	for name := range weeklyDataKeselamatan {
		weeklyDataKeselamatanSorted = append(weeklyDataKeselamatanSorted, name)
	}

	// Sort the slice of names
	sort.Strings(weeklyDataKeamananSorted)
	sort.Strings(weeklyDataKeselamatanSorted)

	var jenisKejadianKeamanan, jenisKejadianKeselamatan []models.JenisKejadian
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeamanan)
	facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
		Order("nama_kejadian asc").Find(&jenisKejadianKeselamatan)

	facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

	facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)

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

	for _, key := range weeklyDataKeamananSorted {
		value := weeklyDataKeamanan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeamananWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}
	for _, key := range weeklyDataKeselamatanSorted {
		value := weeklyDataKeselamatan[key]
		for _, data := range value {
			// fmt.Println(key, i, j)
			for _, weeks := range weekName {
				if key == weeks {
					kejadianKeselamatanWeek[data.JenisKejadian.NamaKejadian][weeks]++
				}
			}
		}
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range data_keamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	// Create a slice of keys (NamaKejadian values)
	sortedNames := make([]string, 0, len(groupedByJenisKeamanan))
	for name := range groupedByJenisKeamanan {
		sortedNames = append(sortedNames, name)
	}

	// Sort the slice of names
	sort.Strings(sortedNames)

	var groupKeamanan []GroupingKeamanan
	var keamananBarat []GroupingKeamananBarat
	var keamananTimur []GroupingKeamananTimur
	var keamananTengah []GroupingKeamananTengah
	var groupKeamananBarat []models.KejadianKeamanan
	var groupKeamananTimur []models.KejadianKeamanan
	var groupKeamananTengah []models.KejadianKeamanan
	// Print the grouped data
	for _, jenisName := range sortedNames {
		kejadianGroup := groupedByJenisKeamanan[jenisName]
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		for _, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				jumlahTengah++
				groupKeamananTengah = append(groupKeamananTengah, index)
			}
			jumlah++
		}

		if jumlahBarat != 0 {
			keamananBarat = append(keamananBarat, GroupingKeamananBarat{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananBarat,
				Jumlah:           jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keamananTimur = append(keamananTimur, GroupingKeamananTimur{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTimur,
				Jumlah:           jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keamananTengah = append(keamananTengah, GroupingKeamananTengah{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTengah,
				Jumlah:           jumlahTengah,
			})
		}

		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
			NamaKejadian:     jenisName,
			KejadianKeamanan: kejadianGroup,
			Jumlah:           jumlah,
			JumlahZonaBarat:  jumlahBarat,
			JumlahZonaTimur:  jumlahTimur,
			JumlahZonaTengah: jumlahTengah,
		})
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
	for _, kejadian := range data_keselamatan {
		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	sortedNames = make([]string, 0, len(groupedByJenisKeselamatan))
	for name := range groupedByJenisKeselamatan {
		sortedNames = append(sortedNames, name)
	}

	// Sort the slice of names
	sort.Strings(sortedNames)

	var groupKeselamatan []GroupingKeselamatan
	var keselamatanBarat []GroupingKeselamatanBarat
	var keselamatanTimur []GroupingKeselamatanTimur
	var keselamatanTengah []GroupingKeselamatanTengah
	var groupKeselamatanBarat []models.KejadianKeselamatanKorban
	var groupKeselamatanTimur []models.KejadianKeselamatanKorban
	var groupKeselamatanTengah []models.KejadianKeselamatanKorban
	// Print the grouped data
	for _, jenisName := range sortedNames {
		kejadianGroup := groupedByJenisKeselamatan[jenisName]
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		var list_korban []models.KejadianKeselamatanKorban

		for _, data := range kejadianGroup {
			var x models.ListKorban
			err := json.Unmarshal(data.Korban, &x)
			if err != nil {
				return
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			if data.Zona == "BARAT" {
				jumlahBarat++
				groupKeselamatanBarat = append(groupKeselamatanBarat, temp)
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
				groupKeselamatanTimur = append(groupKeselamatanTimur, temp)
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
				groupKeselamatanTengah = append(groupKeselamatanTengah, temp)
			}

			list_korban = append(list_korban, temp)
			jumlah++
		}

		if jumlahBarat != 0 {
			keselamatanBarat = append(keselamatanBarat, GroupingKeselamatanBarat{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanBarat,
				Jumlah:              jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keselamatanTimur = append(keselamatanTimur, GroupingKeselamatanTimur{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTimur,
				Jumlah:              jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keselamatanTengah = append(keselamatanTengah, GroupingKeselamatanTengah{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTengah,
				Jumlah:              jumlahTengah,
			})
		}

		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
			NamaKejadian:        jenisName,
			KejadianKeselamatan: list_korban,
			Jumlah:              jumlah,
			JumlahZonaBarat:     jumlahBarat,
			JumlahZonaTimur:     jumlahTimur,
			JumlahZonaTengah:    jumlahTengah,
		})
	}
	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	newDate := time.Now()
	date := fmt.Sprintf("%d %s %d", newDate.Day(), MonthNameIndonesia(newDate.Month()), newDate.Year())
	fmt.Printf("%d %s %d\n", newDate.Day(), MonthNameIndonesia(newDate.Month()), newDate.Year())
	// html template data
	templateData := struct {
		BaseURL                   string
		Tanggal                   string
		Jabatan                   string
		Nama                      string
		IsApproved                bool
		Ttd                       *string
		Nik                       string
		Bulan                     string
		BulanCapital              string
		Tahun                     string
		JumlahKejadianKeamanan    int
		JumlahKejadianKeselamatan int
		KejadianKeamanan          []GroupingKeamanan
		KejadianKeamananBarat     []GroupingKeamananBarat
		KejadianKeamananTimur     []GroupingKeamananTimur
		KejadianKeamananTengah    []GroupingKeamananTengah
		KejadianKeselamatan       []GroupingKeselamatan
		KejadianKeselamatanBarat  []GroupingKeselamatanBarat
		KejadianKeselamatanTimur  []GroupingKeselamatanTimur
		KejadianKeselamatanTengah []GroupingKeselamatanTengah
		KejadianKeamananWeek      map[string]map[string]int
		KejadianKeselamatanWeek   map[string]map[string]int
		JenisKejadianKeamanan     []models.JenisKejadian
		JenisKejadianKeselamatan  []models.JenisKejadian
		WeekName                  []string
		PeriodeTanggal            []string
	}{
		BaseURL:                   baseURL,
		Tanggal:                   date,
		Jabatan:                   deputi.Jabatan.Name,
		Nama:                      deputi.Name,
		IsApproved:                false == true,
		Ttd:                       deputi.Ttd,
		Nik:                       deputi.EmpNo,
		Bulan:                     bulan,
		BulanCapital:              strings.ToUpper(bulan),
		Tahun:                     year,
		JumlahKejadianKeamanan:    len(data_keamanan),
		JumlahKejadianKeselamatan: len(data_keselamatan),
		KejadianKeamanan:          groupKeamanan,
		KejadianKeamananBarat:     keamananBarat,
		KejadianKeamananTimur:     keamananTimur,
		KejadianKeamananTengah:    keamananTengah,
		KejadianKeselamatan:       groupKeselamatan,
		KejadianKeselamatanBarat:  keselamatanBarat,
		KejadianKeselamatanTimur:  keselamatanTimur,
		KejadianKeselamatanTengah: keselamatanTengah,
		KejadianKeamananWeek:      kejadianKeamananWeek,
		KejadianKeselamatanWeek:   kejadianKeselamatanWeek,
		JenisKejadianKeamanan:     jenisKejadianKeamanan,
		JenisKejadianKeselamatan:  jenisKejadianKeselamatan,
		WeekName:                  weekName,
		PeriodeTanggal:            periodeTanggal,
	}

	//html template path
	templatePath := "templates/laporan-bulanan.html"
	newTemplatePath := "laporan-bulanan.html"

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Bulanan/"
	fullPath := rootPath + dirPath

	_ = os.MkdirAll(fullPath, 0755)

	nameFile := strconv.Itoa(int(now.Month())) + ". Laporan bulan " + bulan + " Keamanan & Keselamatan di Wilayah Perairan Indonesia"
	path := fullPath + nameFile + ".pdf"

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(path, 0, "Bulanan_update/")
	} else {
		fmt.Println(err)
	}

	jenisFile := "Laporan Bulanan"
	document := dirPath + nameFile + ".pdf"

	// save to Laporan
	var laporan models.Laporan
	fmt.Println(id_laporan)
	facades.Orm().Query().Where("id_laporan=?", id_laporan).First(&laporan)
	laporan.NamaLaporan = nameFile
	laporan.JenisLaporan = jenisFile
	laporan.BulanKe = int(now.Month())
	laporan.TahunKe = now.Year()
	laporan.Dokumen = document
	// laporan.JenisLaporan = jenisFile
	// laporan.TahunKe = now.Year()
	// laporan.BulanKe = int(now.Month())
	// laporan.MingguKe = mingguKe
	// laporan.Dokumen = document

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", data_keamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	facades.Orm().Query().Save(&laporan)

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Delete(&models.Approval{}, "laporan_id = ?", laporan.IDLaporan)
	facades.Orm().Query().Create(&approval)

	fmt.Println("PDF created successfully!")
}

// LAPORAN TRIWULAN
func (r *Pdf) LaporanTriwulan() {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)
	const (
		templatePath    = "templates/laporan-triwulan.html"
		newTemplatePath = "laporan-triwulan.html"
	)

	now := time.Now()
	year := strconv.Itoa(now.Year())

	quarters := map[time.Month]struct {
		quarter      string
		periodFormat string
		months       []string
	}{
		time.January:   {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.February:  {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.March:     {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.April:     {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.May:       {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.June:      {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.July:      {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.August:    {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.September: {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.October:   {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.November:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.December:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
	}

	quarterInfo := quarters[now.Month()]
	triwulanKe := quarterInfo.quarter
	periodeBulan := fmt.Sprintf(quarterInfo.periodFormat, year)
	months := quarterInfo.months

	var dataKeamanan []models.KejadianKeamanan
	var default1, default2 []models.JenisKejadian
	var dataKeselamatan []models.KejadianKeselamatan

	monthNumbers := make([]int, len(months))
	for i, month := range months {
		monthNumbers[i] = MonthNameMap[month]
	}

	fmt.Println(monthNumbers)

	query1 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeamanan)
	query2 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeselamatan)
	query3 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
		Order("nama_kejadian asc").Find(&default1)
	query4 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
		Order("nama_kejadian asc").Find(&default2)

	if query1 != nil && query2 != nil && query3 != nil && query4 != nil {
		fmt.Println("failed to execute query")
		return
	}

	kejadianCountKeamanan := make(map[string]map[string]int)
	kejadianCountKeselamatan := make(map[string]map[string]int)
	locationCountKeamanan := make(map[string]map[string]int)

	for _, kejadian := range default1 {
		for _, month := range months {
			if kejadianCountKeamanan[kejadian.NamaKejadian] == nil {
				kejadianCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
			}
			if _, exists := kejadianCountKeamanan[kejadian.NamaKejadian][month]; !exists {
				kejadianCountKeamanan[kejadian.NamaKejadian][month] = 0
			}
		}

		if locationCountKeamanan[kejadian.NamaKejadian] == nil {
			locationCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
		}
	}
	for _, kejadian := range default2 {
		for _, month := range months {
			if kejadianCountKeselamatan[kejadian.NamaKejadian] == nil {
				kejadianCountKeselamatan[kejadian.NamaKejadian] = make(map[string]int)
			}
			if _, exists := kejadianCountKeselamatan[kejadian.NamaKejadian][month]; !exists {
				kejadianCountKeselamatan[kejadian.NamaKejadian][month] = 0
			}
		}
	}

	var id_keamanan_arr []int64
	for _, kejadian := range dataKeamanan {
		id_keamanan_arr = append(id_keamanan_arr, kejadian.IdKejadianKeamanan)
		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
		kejadianCountKeamanan[kejadian.JenisKejadian.NamaKejadian][month]++

		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "dermaga") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "pelabuhan") ||
			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Asal), "pelabuhan"))) ||
			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "pelabuhan"))) {
			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["dermaga/pelabuhan"]++
		}
		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "laut") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "perairan") ||
			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "laut") || strings.Contains(strings.ToLower(*kejadian.Asal), "perairan"))) ||
			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "laut") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "perairan"))) {
			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["laut/perairan"]++
		}
	}
	var id_keselamatan_arr []int64
	for _, kejadian := range dataKeselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, kejadian.IdKejadianKeselamatan)
		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
		kejadianCountKeselamatan[kejadian.JenisKejadian.NamaKejadian][month]++
	}

	var kejadianCountsKeamanan, kejadianCountsKeselamatan []MonthlyCount

	jenisKejadianKeysKeamanan := sortedKeys(kejadianCountKeamanan)
	jenisKejadianKeysKeselamatan := sortedKeys(kejadianCountKeselamatan)
	locationsKeysKeamanan := sortedKeys(locationCountKeamanan)

	for _, jenisKejadian := range jenisKejadianKeysKeamanan {
		monthCounts := kejadianCountKeamanan[jenisKejadian]
		entry := MonthlyCount{
			NamaKejadian: jenisKejadian,
			Bulan1:       months[0],
			Count1:       monthCounts[months[0]],
			Bulan2:       months[1],
			Count2:       monthCounts[months[1]],
			Bulan3:       months[2],
			Count3:       monthCounts[months[2]],
			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
		}
		kejadianCountsKeamanan = append(kejadianCountsKeamanan, entry)
	}
	for _, jenisKejadian := range jenisKejadianKeysKeselamatan {
		monthCounts := kejadianCountKeselamatan[jenisKejadian]
		entry := MonthlyCount{
			NamaKejadian: jenisKejadian,
			Bulan1:       months[0],
			Count1:       monthCounts[months[0]],
			Bulan2:       months[1],
			Count2:       monthCounts[months[1]],
			Bulan3:       months[2],
			Count3:       monthCounts[months[2]],
			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
		}
		kejadianCountsKeselamatan = append(kejadianCountsKeselamatan, entry)
	}

	var locationOutput []LocationOutput
	for _, jenisKejadian := range locationsKeysKeamanan {
		counts := locationCountKeamanan[jenisKejadian]
		locationOutput = append(locationOutput, LocationOutput{
			NamaKejadian:   jenisKejadian,
			JumlahDermaga:  counts["dermaga/pelabuhan"],
			JumlahPerairan: counts["laut/perairan"],
		})
	}

	var bulan []string
	for _, x := range monthNumbers {
		bulan = append(bulan, MonthNameIndonesia(time.Month(x)))
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range dataKeamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeamanan []GroupingKeamanan
	var keamananBarat []GroupingKeamananBarat
	var keamananTimur []GroupingKeamananTimur
	var keamananTengah []GroupingKeamananTengah
	var groupKeamananBarat []models.KejadianKeamanan
	var groupKeamananTimur []models.KejadianKeamanan
	var groupKeamananTengah []models.KejadianKeamanan
	// Print the grouped data
	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		for _, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				jumlahTengah++
				groupKeamananTengah = append(groupKeamananTengah, index)
			}
			jumlah++
		}

		if jumlahBarat != 0 {
			keamananBarat = append(keamananBarat, GroupingKeamananBarat{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananBarat,
				Jumlah:           jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keamananTimur = append(keamananTimur, GroupingKeamananTimur{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTimur,
				Jumlah:           jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keamananTengah = append(keamananTengah, GroupingKeamananTengah{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTengah,
				Jumlah:           jumlahTengah,
			})
		}

		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
			NamaKejadian:     jenisName,
			KejadianKeamanan: kejadianGroup,
			Jumlah:           jumlah,
			JumlahZonaBarat:  jumlahBarat,
			JumlahZonaTimur:  jumlahTimur,
			JumlahZonaTengah: jumlahTengah,
		})
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
	for _, kejadian := range dataKeselamatan {
		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeselamatan []GroupingKeselamatan
	var keselamatanBarat []GroupingKeselamatanBarat
	var keselamatanTimur []GroupingKeselamatanTimur
	var keselamatanTengah []GroupingKeselamatanTengah
	var groupKeselamatanBarat []models.KejadianKeselamatanKorban
	var groupKeselamatanTimur []models.KejadianKeselamatanKorban
	var groupKeselamatanTengah []models.KejadianKeselamatanKorban
	// Print the grouped data
	for jenisName, kejadianGroup := range groupedByJenisKeselamatan {
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		var list_korban []models.KejadianKeselamatanKorban

		for _, data := range kejadianGroup {
			var x models.ListKorban
			err := json.Unmarshal(data.Korban, &x)
			if err != nil {
				return
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			if data.Zona == "BARAT" {
				jumlahBarat++
				groupKeselamatanBarat = append(groupKeselamatanBarat, temp)
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
				groupKeselamatanTimur = append(groupKeselamatanTimur, temp)
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
				groupKeselamatanTengah = append(groupKeselamatanTengah, temp)
			}

			list_korban = append(list_korban, temp)
			jumlah++
		}

		if jumlahBarat != 0 {
			keselamatanBarat = append(keselamatanBarat, GroupingKeselamatanBarat{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanBarat,
				Jumlah:              jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keselamatanTimur = append(keselamatanTimur, GroupingKeselamatanTimur{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTimur,
				Jumlah:              jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keselamatanTengah = append(keselamatanTengah, GroupingKeselamatanTengah{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTengah,
				Jumlah:              jumlahTengah,
			})
		}

		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
			NamaKejadian:        jenisName,
			KejadianKeselamatan: list_korban,
			Jumlah:              jumlah,
			JumlahZonaBarat:     jumlahBarat,
			JumlahZonaTimur:     jumlahTimur,
			JumlahZonaTengah:    jumlahTengah,
		})
	}

	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	firstDayNextMonth := now.AddDate(0, 1, -now.Day()+1)
	date := fmt.Sprintf("%d %s %d", firstDayNextMonth.Day(), MonthNameIndonesia(firstDayNextMonth.Month()), firstDayNextMonth.Year())
	templateData := struct {
		BaseURL                                  string
		Tanggal                                  string
		Jabatan                                  string
		Nama                                     string
		IsApproved                               bool
		Ttd                                      *string
		Nik                                      string
		PeriodeTriwulan                          string
		BulanCapital                             string
		BulanSingkatan                           []string
		Bulan                                    []string
		TableKejadianKeamanan                    []MonthlyCount
		KejadianKeamanan                         []GroupingKeamanan
		KejadianKeamananBarat                    []GroupingKeamananBarat
		KejadianKeamananTimur                    []GroupingKeamananTimur
		KejadianKeamananTengah                   []GroupingKeamananTengah
		TableKejadianKeselamatan                 []MonthlyCount
		KejadianKeselamatan                      []GroupingKeselamatan
		KejadianKeselamatanBarat                 []GroupingKeselamatanBarat
		KejadianKeselamatanTimur                 []GroupingKeselamatanTimur
		KejadianKeselamatanTengah                []GroupingKeselamatanTengah
		TablePengelompokanLokasiKejadianKeamanan []LocationOutput
		Tahun                                    string
	}{
		BaseURL:                                  baseURL,
		Tanggal:                                  date,
		Jabatan:                                  deputi.Jabatan.Name,
		Nama:                                     deputi.Name,
		IsApproved:                               false == true,
		Ttd:                                      deputi.Ttd,
		Nik:                                      deputi.EmpNo,
		PeriodeTriwulan:                          triwulanKe,
		BulanCapital:                             periodeBulan,
		BulanSingkatan:                           months,
		Bulan:                                    bulan,
		TableKejadianKeamanan:                    kejadianCountsKeamanan,
		KejadianKeamanan:                         groupKeamanan,
		KejadianKeamananBarat:                    keamananBarat,
		KejadianKeamananTimur:                    keamananTimur,
		KejadianKeamananTengah:                   keamananTengah,
		TableKejadianKeselamatan:                 kejadianCountsKeselamatan,
		KejadianKeselamatan:                      groupKeselamatan,
		KejadianKeselamatanBarat:                 keselamatanBarat,
		KejadianKeselamatanTimur:                 keselamatanTimur,
		KejadianKeselamatanTengah:                keselamatanTengah,
		TablePengelompokanLokasiKejadianKeamanan: locationOutput,
		Tahun:                                    year,
	}

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Triwulan/"
	fullPath := rootPath + dirPath

	_ = os.MkdirAll(fullPath, 0755)

	nameFile := "Laporan Triwulan " + triwulanKe
	path := fullPath + nameFile + ".pdf"

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(path, 0, "Triwulan_create/")
	} else {
		fmt.Println(err)
	}

	jenisFile := "Laporan Triwulan"
	document := dirPath + nameFile + ".pdf"

	// save to Laporan
	laporan := models.Laporan{
		NamaLaporan:  nameFile,
		JenisLaporan: jenisFile,
		TahunKe:      now.Year(),
		BulanKe:      int(now.Month()),
		Dokumen:      document,
	}

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", dataKeamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	if err := facades.Orm().Query().Create(&laporan); err != nil {
		fmt.Println("Data Gagal Ditambahkan:", err)
		return
	}
	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Create(&approval)
}
func (r *Pdf) LaporanTriwulanUpdate(id_laporan int64, month int, years int) {
	appHost := facades.Config().Env("APP_HOST", "127.0.0.1")
	appPort := facades.Config().Env("APP_PORT", "3000")
	baseURL := fmt.Sprintf("http://%s:%s", appHost, appPort)
	const (
		templatePath    = "templates/laporan-triwulan.html"
		newTemplatePath = "laporan-triwulan.html"
	)

	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	year := strconv.Itoa(now.Year())

	quarters := map[time.Month]struct {
		quarter      string
		periodFormat string
		months       []string
	}{
		time.January:   {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.February:  {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.March:     {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.April:     {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.May:       {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.June:      {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.July:      {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.August:    {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.September: {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.October:   {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.November:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.December:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
	}

	quarterInfo := quarters[now.Month()]
	triwulanKe := quarterInfo.quarter
	periodeBulan := fmt.Sprintf(quarterInfo.periodFormat, year)
	months := quarterInfo.months

	var dataKeamanan []models.KejadianKeamanan
	var default1, default2 []models.JenisKejadian
	var dataKeselamatan []models.KejadianKeselamatan

	monthNumbers := make([]int, len(months))
	for i, month := range months {
		monthNumbers[i] = MonthNameMap[month]
	}

	query1 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeamanan)
	query2 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeselamatan)
	query3 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
		Order("nama_kejadian asc").Find(&default1)
	query4 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
		Order("nama_kejadian asc").Find(&default2)

	if query1 != nil && query2 != nil && query3 != nil && query4 != nil {
		fmt.Println("failed to execute query")
		return
	}

	kejadianCountKeamanan := make(map[string]map[string]int)
	kejadianCountKeselamatan := make(map[string]map[string]int)
	locationCountKeamanan := make(map[string]map[string]int)

	for _, kejadian := range default1 {
		for _, month := range months {
			if kejadianCountKeamanan[kejadian.NamaKejadian] == nil {
				kejadianCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
			}
			if _, exists := kejadianCountKeamanan[kejadian.NamaKejadian][month]; !exists {
				kejadianCountKeamanan[kejadian.NamaKejadian][month] = 0
			}
		}

		if locationCountKeamanan[kejadian.NamaKejadian] == nil {
			locationCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
		}
	}
	for _, kejadian := range default2 {
		for _, month := range months {
			if kejadianCountKeselamatan[kejadian.NamaKejadian] == nil {
				kejadianCountKeselamatan[kejadian.NamaKejadian] = make(map[string]int)
			}
			if _, exists := kejadianCountKeselamatan[kejadian.NamaKejadian][month]; !exists {
				kejadianCountKeselamatan[kejadian.NamaKejadian][month] = 0
			}
		}
	}

	var id_keamanan_arr []int64
	for _, kejadian := range dataKeamanan {
		id_keamanan_arr = append(id_keamanan_arr, kejadian.IdKejadianKeamanan)
		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
		kejadianCountKeamanan[kejadian.JenisKejadian.NamaKejadian][month]++

		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "dermaga") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "pelabuhan") ||
			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Asal), "pelabuhan"))) ||
			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "pelabuhan"))) {
			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["dermaga/pelabuhan"]++
		}
		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "laut") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "perairan") ||
			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "laut") || strings.Contains(strings.ToLower(*kejadian.Asal), "perairan"))) ||
			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "laut") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "perairan"))) {
			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["laut/perairan"]++
		}
	}
	var id_keselamatan_arr []int64
	for _, kejadian := range dataKeselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, kejadian.IdKejadianKeselamatan)
		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
		kejadianCountKeselamatan[kejadian.JenisKejadian.NamaKejadian][month]++
	}

	var kejadianCountsKeamanan, kejadianCountsKeselamatan []MonthlyCount

	jenisKejadianKeysKeamanan := sortedKeys(kejadianCountKeamanan)
	jenisKejadianKeysKeselamatan := sortedKeys(kejadianCountKeselamatan)
	locationsKeysKeamanan := sortedKeys(locationCountKeamanan)

	for _, jenisKejadian := range jenisKejadianKeysKeamanan {
		monthCounts := kejadianCountKeamanan[jenisKejadian]
		entry := MonthlyCount{
			NamaKejadian: jenisKejadian,
			Bulan1:       months[0],
			Count1:       monthCounts[months[0]],
			Bulan2:       months[1],
			Count2:       monthCounts[months[1]],
			Bulan3:       months[2],
			Count3:       monthCounts[months[2]],
			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
		}
		kejadianCountsKeamanan = append(kejadianCountsKeamanan, entry)
	}
	for _, jenisKejadian := range jenisKejadianKeysKeselamatan {
		monthCounts := kejadianCountKeselamatan[jenisKejadian]
		entry := MonthlyCount{
			NamaKejadian: jenisKejadian,
			Bulan1:       months[0],
			Count1:       monthCounts[months[0]],
			Bulan2:       months[1],
			Count2:       monthCounts[months[1]],
			Bulan3:       months[2],
			Count3:       monthCounts[months[2]],
			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
		}
		kejadianCountsKeselamatan = append(kejadianCountsKeselamatan, entry)
	}

	var locationOutput []LocationOutput
	for _, jenisKejadian := range locationsKeysKeamanan {
		counts := locationCountKeamanan[jenisKejadian]
		locationOutput = append(locationOutput, LocationOutput{
			NamaKejadian:   jenisKejadian,
			JumlahDermaga:  counts["dermaga/pelabuhan"],
			JumlahPerairan: counts["laut/perairan"],
		})
	}

	var bulan []string
	for _, x := range monthNumbers {
		bulan = append(bulan, MonthNameIndonesia(time.Month(x)))
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range dataKeamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeamanan []GroupingKeamanan
	var keamananBarat []GroupingKeamananBarat
	var keamananTimur []GroupingKeamananTimur
	var keamananTengah []GroupingKeamananTengah
	var groupKeamananBarat []models.KejadianKeamanan
	var groupKeamananTimur []models.KejadianKeamanan
	var groupKeamananTengah []models.KejadianKeamanan
	// Print the grouped data
	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		for _, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				jumlahTengah++
				groupKeamananTengah = append(groupKeamananTengah, index)
			}
			jumlah++
		}

		if jumlahBarat != 0 {
			keamananBarat = append(keamananBarat, GroupingKeamananBarat{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananBarat,
				Jumlah:           jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keamananTimur = append(keamananTimur, GroupingKeamananTimur{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTimur,
				Jumlah:           jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keamananTengah = append(keamananTengah, GroupingKeamananTengah{
				NamaKejadian:     jenisName,
				KejadianKeamanan: groupKeamananTengah,
				Jumlah:           jumlahTengah,
			})
		}

		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
			NamaKejadian:     jenisName,
			KejadianKeamanan: kejadianGroup,
			Jumlah:           jumlah,
			JumlahZonaBarat:  jumlahBarat,
			JumlahZonaTimur:  jumlahTimur,
			JumlahZonaTengah: jumlahTengah,
		})
	}

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
	for _, kejadian := range dataKeselamatan {
		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeselamatan []GroupingKeselamatan
	var keselamatanBarat []GroupingKeselamatanBarat
	var keselamatanTimur []GroupingKeselamatanTimur
	var keselamatanTengah []GroupingKeselamatanTengah
	var groupKeselamatanBarat []models.KejadianKeselamatanKorban
	var groupKeselamatanTimur []models.KejadianKeselamatanKorban
	var groupKeselamatanTengah []models.KejadianKeselamatanKorban
	// Print the grouped data
	for jenisName, kejadianGroup := range groupedByJenisKeselamatan {
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		var list_korban []models.KejadianKeselamatanKorban

		for _, data := range kejadianGroup {
			var x models.ListKorban
			err := json.Unmarshal(data.Korban, &x)
			if err != nil {
				return
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			if data.Zona == "BARAT" {
				jumlahBarat++
				groupKeselamatanBarat = append(groupKeselamatanBarat, temp)
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
				groupKeselamatanTimur = append(groupKeselamatanTimur, temp)
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
				groupKeselamatanTengah = append(groupKeselamatanTengah, temp)
			}

			list_korban = append(list_korban, temp)
			jumlah++
		}

		if jumlahBarat != 0 {
			keselamatanBarat = append(keselamatanBarat, GroupingKeselamatanBarat{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanBarat,
				Jumlah:              jumlahBarat,
			})
		}
		if jumlahTimur != 0 {
			keselamatanTimur = append(keselamatanTimur, GroupingKeselamatanTimur{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTimur,
				Jumlah:              jumlahTimur,
			})
		}
		if jumlahTengah != 0 {
			keselamatanTengah = append(keselamatanTengah, GroupingKeselamatanTengah{
				NamaKejadian:        jenisName,
				KejadianKeselamatan: groupKeselamatanTengah,
				Jumlah:              jumlahTengah,
			})
		}

		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
			NamaKejadian:        jenisName,
			KejadianKeselamatan: list_korban,
			Jumlah:              jumlah,
			JumlahZonaBarat:     jumlahBarat,
			JumlahZonaTimur:     jumlahTimur,
			JumlahZonaTengah:    jumlahTengah,
		})
	}

	var deputi models.Karyawan
	facades.Orm().Query().
		With("Jabatan").
		With("User.Role").
		Join("JOIN jabatan ON jabatan.id_jabatan = karyawan.jabatan_id").
		Where("jabatan.name = ?", "Deputi Informasi, Hukum dan Kerja Sama").
		First(&deputi)

	newDate := time.Now()
	date := fmt.Sprintf("%d %s %d", newDate.Day(), MonthNameIndonesia(newDate.Month()), newDate.Year())
	templateData := struct {
		BaseURL                                  string
		Tanggal                                  string
		Jabatan                                  string
		Nama                                     string
		IsApproved                               bool
		Ttd                                      *string
		Nik                                      string
		PeriodeTriwulan                          string
		BulanCapital                             string
		BulanSingkatan                           []string
		Bulan                                    []string
		TableKejadianKeamanan                    []MonthlyCount
		KejadianKeamanan                         []GroupingKeamanan
		KejadianKeamananBarat                    []GroupingKeamananBarat
		KejadianKeamananTimur                    []GroupingKeamananTimur
		KejadianKeamananTengah                   []GroupingKeamananTengah
		TableKejadianKeselamatan                 []MonthlyCount
		KejadianKeselamatan                      []GroupingKeselamatan
		KejadianKeselamatanBarat                 []GroupingKeselamatanBarat
		KejadianKeselamatanTimur                 []GroupingKeselamatanTimur
		KejadianKeselamatanTengah                []GroupingKeselamatanTengah
		TablePengelompokanLokasiKejadianKeamanan []LocationOutput
		Tahun                                    string
	}{
		BaseURL:                                  baseURL,
		Tanggal:                                  date,
		Jabatan:                                  deputi.Jabatan.Name,
		Nama:                                     deputi.Name,
		IsApproved:                               false == true,
		Ttd:                                      deputi.Ttd,
		Nik:                                      deputi.EmpNo,
		PeriodeTriwulan:                          triwulanKe,
		BulanCapital:                             periodeBulan,
		BulanSingkatan:                           months,
		Bulan:                                    bulan,
		TableKejadianKeamanan:                    kejadianCountsKeamanan,
		KejadianKeamanan:                         groupKeamanan,
		KejadianKeamananBarat:                    keamananBarat,
		KejadianKeamananTimur:                    keamananTimur,
		KejadianKeamananTengah:                   keamananTengah,
		TableKejadianKeselamatan:                 kejadianCountsKeselamatan,
		KejadianKeselamatan:                      groupKeselamatan,
		KejadianKeselamatanBarat:                 keselamatanBarat,
		KejadianKeselamatanTimur:                 keselamatanTimur,
		KejadianKeselamatanTengah:                keselamatanTengah,
		TablePengelompokanLokasiKejadianKeamanan: locationOutput,
		Tahun:                                    year,
	}

	// Save the PDF to a file
	rootPath := "storage/app/"
	dirPath := year + "/Laporan Triwulan/"
	fullPath := rootPath + dirPath

	_ = os.MkdirAll(fullPath, 0755)

	nameFile := "Laporan Triwulan " + triwulanKe
	path := fullPath + nameFile + ".pdf"

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(path, 0, "Triwulan_update/")
	} else {
		fmt.Println(err)
	}

	jenisFile := "Laporan Triwulan"
	document := dirPath + nameFile + ".pdf"

	var laporan models.Laporan
	facades.Orm().Query().Where("id_laporan=?", id_laporan).First(&laporan)
	// save to Laporan

	laporan.NamaLaporan = nameFile
	laporan.JenisLaporan = jenisFile
	laporan.TahunKe = now.Year()
	laporan.BulanKe = int(now.Month())
	laporan.Dokumen = document

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", dataKeamanan[0].CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan")
		return
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		fmt.Println("Data Tidak Ditemukan:")
		return
	}

	facades.Orm().Query().Save(&laporan)

	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": true,
	})

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Delete(&models.Approval{}, "laporan_id = ?", laporan.IDLaporan)
	facades.Orm().Query().Create(&approval)
}
