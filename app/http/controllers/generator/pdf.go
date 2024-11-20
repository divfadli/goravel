package generator

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/lib/pq"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type UintArray []uint8

func (a *UintArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return a.scanBytes(v)
	case string:
		return a.scanBytes([]byte(v))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("cannot scan %T into UintArray", src)
}

func (a *UintArray) scanBytes(src []byte) error {
	str := string(src)
	str = strings.Trim(str, "{}")
	parts := strings.Split(str, ",")
	*a = make([]uint8, len(parts))
	for i, s := range parts {
		num, err := strconv.ParseUint(strings.TrimSpace(s), 10, 8)
		if err != nil {
			return err
		}
		(*a)[i] = uint8(num)
	}
	return nil
}

func (a UintArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return strings.Join(strings.Fields(fmt.Sprint(a)), ","), nil
}

type Pdf struct {
	body string
	//Dependent services
}

func NewPdf(body string) *Pdf {
	return &Pdf{
		body: body,
		//Inject services
	}
}

func (r *Pdf) ParseTemplate(templateFileName string, newTemplateFileName string, data interface{}) error {
	t, err := template.New(newTemplateFileName).Funcs(template.FuncMap{
		"contains": strings.Contains,
		"split":    strings.Split,
		"trim":     strings.TrimSpace,
		"add":      Add,
		"seq":      Seq,
		"sub":      Sub,
		"inc":      Increment,
		"last": func(index int, length int) bool {
			return index == length-1
		},
	}).ParseFiles(templateFileName)

	if err != nil {
		return err
	}

	t = template.Must(t, err)

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

// generate slide function
func (r *Pdf) GenerateSlide(slidePath string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("cloneTemplate/", 0777)
		if errDir != nil {
			fmt.Printf("Error: %v\n", errDir)
			return false, errDir
		}
	}
	err1 := ioutil.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
		return false, err1
	}

	// Define the input HTML file and output image file
	inputFile := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"

	// Check if the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
		return false, err
	}

	// Construct the command
	cmd := exec.Command("wkhtmltopdf", "--enable-local-file-access", "--enable-smart-shrinking",
		"--enable-plugins", "--javascript-delay", "2000", "--no-stop-slow-scripts", "--debug-javascript",
		"--page-size", "A4", "--orientation", "Landscape",
		"--margin-top", "0", "--margin-bottom", "0", "--margin-left", "0", "--margin-right", "0",
		"--zoom", "0.8", inputFile, slidePath,
	)
	// cmd := exec.Command("wkhtmltopdf", "--enable-local-file-access", "--enable-plugins",
	// 	"--javascript-delay", "500", "--debug-javascript", "--page-size", "A4", "--orientation", "Landscape",
	// 	"--margin-top", "0", "--margin-bottom", "0", "--margin-left", "0", "--margin-right", "0", "--zoom", "0.8",
	// 	"--enable-smart-shrinking", inputFile, slidePath)
	// Set the command's standard output and error to the current process's standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	defer os.RemoveAll(dir + "/cloneTemplate")

	return true, nil
}

// generate laporan function
func (r *Pdf) GenerateLaporan(laporanPath string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("cloneTemplate/", 0777)
		if errDir != nil {
			fmt.Printf("Error: %v\n", errDir)
			return false, errDir
		}
	}
	err1 := ioutil.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
		return false, err1
	}

	// Define the input HTML file and output image file
	inputFile := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"

	// Check if the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
		return false, err
	}

	// Construct the command
	cmd := exec.Command("wkhtmltopdf", "--enable-local-file-access", "--zoom", "0.8",
		"--enable-smart-shrinking", "--page-size", "A4", "--javascript-delay", "5000",
		"--footer-center", "[page]", "--footer-font-size", "10",
		inputFile, laporanPath)

	// Set the command's standard output and error to the current process's standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	defer os.RemoveAll(dir + "/cloneTemplate")

	return true, nil
}

func (r *Pdf) GenerateMingguanLastMonth(ctx http.Context, date time.Time) http.Response {
	baseURL := "http://" + ctx.Request().Host()

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
		bulan = monthNameIndonesia(12)
		bulanInt = 12
		year = strconv.Itoa(date.Year() - 1)
		yearInt = date.Year() - 1
	} else {
		bulan = monthNameIndonesia(date.Month() - 1)
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

	jumlahHari := daysInMonth(startOfMonth)

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

	var weeklyDataKeamanans, weeklyDataKeselamatans, allWeeklyDataKeamanans, allWeeklyDataKeselamatans weeklyData

	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
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
		Where(tanggal, startOfWeek, endOfWeek).
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
	outputPath := fmt.Sprintf("storage/temp/pelanggaran-%s.pdf", "head-last-month")

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
		return nil
	}

	for _, data := range result_keamanan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/pelanggaran-%d.pdf", data.IdKejadianKeamanan)

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
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return nil
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
	outputPath = fmt.Sprintf("storage/temp/kecelakaan-%s.pdf", "head-last-month")

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
		success, _ := r.GenerateSlide(outputPath)
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return nil
	}

	for _, data := range result_keselamatan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/kecelakaan-%d.pdf", data.IdKejadianKeselamatan)

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
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return nil
		}
	}

	// Save the PDF to a file
	// Create directory path
	dirPath := "storage/app/" + strconv.Itoa(yearInt) + "/Laporan Mingguan/Bulan " + bulan
	_ = os.MkdirAll(dirPath, 0755)

	yearStr := strconv.Itoa(yearInt)
	lastTwoDigits := yearStr[len(yearStr)-2:]
	path := dirPath + "/LAP MING KE-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishMap[bulanInt] + "'" + lastTwoDigits + ".pdf"
	err := api.MergeCreateFile(images, path, false, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return nil
	}

	fmt.Println("PDF created successfully!")

	return ctx.Response().Success().Json(map[string]interface{}{
		"Status": "success",
	})
}

func (r *Pdf) GenerateMingguan(ctx http.Context) http.Response {
	baseURL := "http://" + ctx.Request().Host()
	fmt.Println("baseURL:", baseURL)

	//html template path
	templateKeamananPath := "templates/keamanan.html"
	newTemplateKeamananPath := "keamanan.html"
	templateKeamananHeadPath := "templates/keamanan-head.html"
	newTemplateKeamananHeadPath := "keamanan-head.html"
	templateKeselamatanHeadPath := "templates/keselamatan-head.html"
	newTemplateKeselamatanHeadPath := "keselamatan-head.html"
	templateKeselamatanPath := "templates/keselamatan.html"
	newTemplateKeselamatanPath := "keselamatan.html"

	now := time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC)
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

	fmt.Println(minggu)

	var startOfWeek time.Time
	var endOfWeek time.Time
	var mingguKe int
	for i, weekEnd := range minggu {
		if weekEnd.After(now) {
			fmt.Println(i)
			if i == 0 {
				return r.GenerateMingguanLastMonth(ctx, now)
			} else if i == 1 {
				if minggu[i-1].Day() < 7 {
					r.GenerateMingguanLastMonth(ctx, now)
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

	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
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
		Where(tanggal, startOfWeek, endOfWeek).
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
	outputPath := fmt.Sprintf("storage/temp/pelanggaran-%s.pdf", "head")

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
		return nil
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
	///x

	for _, data := range result_keamanan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/pelanggaran-%d.pdf", data.IdKejadianKeamanan)

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
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return nil
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
	outputPath = fmt.Sprintf("storage/temp/kecelakaan-%s.pdf", "head")

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
		success, _ := r.GenerateSlide(outputPath)
		if success {
			images = append(images, outputPath)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
		return nil
	}

	for _, data := range result_keselamatan {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/temp/kecelakaan-%d.pdf", data.IdKejadianKeselamatan)

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
			success, _ := r.GenerateSlide(outputPath)
			if success {
				images = append(images, outputPath)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
			return nil
		}
	}

	// Save the PDF to a file
	dirPath := "storage/app/" + year + "/Laporan Mingguan/Bulan " + bulan
	_ = os.MkdirAll(dirPath, 0755)

	lastTwoDigits := year[len(year)-2:]
	path := dirPath + "/LAP MING KE-" + strconv.Itoa(mingguKe) + " " + monthNameEnglishMap[now.Month()] + "'" + lastTwoDigits + ".pdf"
	err := api.MergeCreateFile(images, path, false, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return nil
	}

	fmt.Println("PDF created successfully!")

	return ctx.Response().Success().Json(map[string]interface{}{
		"Status": "success",
		"week":   mingguKe,
	})
}

type GroupingKeamananBarat struct {
	NamaKejadian     string `json:"nama_kejadian"`
	KejadianKeamanan []models.KejadianKeamanan
	Jumlah           int `json:"jumlah"`
}
type GroupingKeamananTimur struct {
	NamaKejadian     string `json:"nama_kejadian"`
	KejadianKeamanan []models.KejadianKeamanan
	Jumlah           int `json:"jumlah"`
}
type GroupingKeamananTengah struct {
	NamaKejadian     string `json:"nama_kejadian"`
	KejadianKeamanan []models.KejadianKeamanan
	Jumlah           int `json:"jumlah"`
}

type GroupingKeamanan struct {
	NamaKejadian     string `json:"nama_kejadian"`
	KejadianKeamanan []models.KejadianKeamanan
	Jumlah           int `json:"jumlah"`
	JumlahZonaBarat  int `json:"jumlah_zona_barat"`
	JumlahZonaTimur  int `json:"jumlah_zona_timur"`
	JumlahZonaTengah int `json:"jumlah_zona_tengah"`
}

type GroupingKeselamatan struct {
	NamaKejadian        string `json:"nama_kejadian"`
	KejadianKeselamatan []models.KejadianKeselamatanKorban
	Jumlah              int `json:"jumlah"`
	JumlahZonaBarat     int `json:"jumlah_zona_barat"`
	JumlahZonaTimur     int `json:"jumlah_zona_timur"`
	JumlahZonaTengah    int `json:"jumlah_zona_tengah"`
}
type GroupingKeselamatanBarat struct {
	NamaKejadian        string `json:"nama_kejadian"`
	KejadianKeselamatan []models.KejadianKeselamatanKorban
	Jumlah              int `json:"jumlah"`
}
type GroupingKeselamatanTimur struct {
	NamaKejadian        string `json:"nama_kejadian"`
	KejadianKeselamatan []models.KejadianKeselamatanKorban
	Jumlah              int `json:"jumlah"`
}
type GroupingKeselamatanTengah struct {
	NamaKejadian        string `json:"nama_kejadian"`
	KejadianKeselamatan []models.KejadianKeselamatanKorban
	Jumlah              int `json:"jumlah"`
}

type weeklyData []struct {
	WeekStart   time.Time     `gorm:"column:week_start"`
	KejadianIDs pq.Int64Array `gorm:"column:kejadian_ids"`
}

func daysInMonth(tanggal time.Time) int {
	// Create a time object for the first day of the next month
	firstDayNextMonth := time.Date(tanggal.Year(), tanggal.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day to get the last day of the current month
	lastDayCurrentMonth := firstDayNextMonth.AddDate(0, 0, -1)
	// Return the day of the month, which is the number of days in the month
	return lastDayCurrentMonth.Day()
}

func (r *Pdf) GenerateBulanan(ctx http.Context) http.Response {
	//html template path
	// now := time.Now()
	now := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	bulan := monthNameIndonesia(now.Month())
	// intBulan := int(now.Month())
	year := strconv.Itoa(now.Year())
	dayperweek := 7

	jumlahHari := daysInMonth(now)
	// fmt.Printf("Jumlah hari dalam bulan ini: %d\n", jumlahHari)

	var minggu []time.Time

	// // Calculate the first day of the month
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

	// minggu = append(minggu, firstDay)
	nextWeek := firstDay.AddDate(0, 0, daysToAdd)
	minggu = append(minggu, nextWeek)

	jumlahHari -= daysToAdd + 1
	// fmt.Println(jumlahHari)
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// // // Iterate through each week and calculate the start date
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

	var weeklyDataKeamanans, weeklyDataKeselamatans weeklyData

	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Scan(&weeklyDataKeselamatans)

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

	// for _, key := range weeklyDataKeamananSorted {
	// 	value := weeklyDataKeamanan[key]
	// 	fmt.Printf("Key: %s, Kejadian: %d\n", key, len(value))
	// 	for _, kejadian := range value {
	// 		fmt.Println("ID:", kejadian.JenisKejadian.IDJenisKejadian, kejadian.JenisKejadian.NamaKejadian)
	// 	}
	// }

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
		// Where("DATE_PART('month', tanggal) IN (?)", 6)
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

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		for i, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				fmt.Println("MASUK BARAT-", i)
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				fmt.Println("MASUK TIMUR-", i)
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				fmt.Println("MASUK TENGAH-", i)
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
			fmt.Println("TIMURRR")
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

		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		var list_korban []models.KejadianKeselamatanKorban

		for _, data := range kejadianGroup {
			var x models.ListKorban
			err := json.Unmarshal(data.Korban, &x)
			if err != nil {
				return nil
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
	// html template data
	templateData := struct {
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

	templatePath := "templates/laporan-bulanan.html"
	newTemplatePath := "laporan-bulanan.html"
	outputPath := "storage/output-laporan-bulanan.pdf"

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(outputPath)
	} else {
		fmt.Println(err)
	}

	// fmt.Println("PDF created successfully!")
	return ctx.Response().Success().Json(map[string]interface{}{
		"Status":                    "success",
		"data-1":                    kejadianKeamananWeek["Human Trafficking"]["10-16 Juni"],
		"kejadian_keamanan_week":    kejadianKeamananWeek,
		"jenis_kejadian_keamanan":   jenisKejadianKeamanan,
		"kejadian_keselamatan_week": kejadianKeselamatanWeek,
		"week":                      weekName,
		"kejadian_keamanan":         groupKeamanan,
		"kejadian_keselamatan":      groupKeselamatan,
	})
}

func Add(x, y int) int {
	return x + y
}

func Seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

// sub subtracts b from a.
func Sub(a, b int) int {
	return a - b
}

func Increment(val int) int {
	return val + 1
}

type MonthlyCount struct {
	NamaKejadian  string `json:"nama_kejadian"`
	Bulan1        string `json:"bulan_1"`
	KorbanTewas   int    `json:"korban_tewas"`
	KorbanSelamat int    `json:"korban_selamat"`
	KorbanHilang  int    `json:"korban_hilang"`
	Count1        int    `json:"count_1"`
	Bulan2        string `json:"bulan_2"`
	Count2        int    `json:"count_2"`
	Bulan3        string `json:"bulan_3"`
	Count3        int    `json:"count_3"`
	Total         int    `json:"total"`
}

type LocationOutput struct {
	NamaKejadian   string `json:"nama_kejadian"`
	JumlahDermaga  int    `json:"jumlah_dermaga"`
	JumlahPerairan int    `json:"jumlah_perairan"`
}

var monthNameMap = map[string]int{
	"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MEI": 5, "JUN": 6,
	"JUL": 7, "AGT": 8, "SEP": 9, "OKT": 10, "NOV": 11, "DES": 12,
}

func monthNameEnglishTitleCase(month time.Month) string {
	monthUpper := monthNameEnglishMap[month]
	return strings.Title(strings.ToLower(monthUpper))
}

var monthNameEnglishMap = map[time.Month]string{
	time.January: "JAN", time.February: "FEB", time.March: "MAR",
	time.April: "APR", time.May: "MEI", time.June: "JUN",
	time.July: "JUL", time.August: "AGT", time.September: "SEP",
	time.October: "OKT", time.November: "NOV", time.December: "DES",
}

var monthNameIndonesiaMap = map[time.Month]string{
	time.January: "Januari", time.February: "Februari", time.March: "Maret",
	time.April: "April", time.May: "Mei", time.June: "Juni",
	time.July: "Juli", time.August: "Agustus", time.September: "September",
	time.October: "Oktober", time.November: "November", time.December: "Desember",
}

func monthNameEnglish(month time.Month) string {
	return monthNameEnglishMap[month]
}

func monthNameIndonesia(month time.Month) string {
	return monthNameIndonesiaMap[month]
}

func (r *Pdf) GenerateTriwulan(ctx http.Context) http.Response {
	const (
		templatePath    = "templates/laporan-triwulan.html"
		newTemplatePath = "laporan-triwulan.html"
		outputPath      = "storage/output-laporan-triwulan.pdf"
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
		monthNumbers[i] = monthNameMap[month]
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
		return nil
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

	for _, kejadian := range dataKeamanan {
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
	for _, kejadian := range dataKeselamatan {
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
		bulan = append(bulan, monthNameIndonesia(time.Month(x)))
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
		for i, index := range kejadianGroup {
			if index.Zona == "BARAT" {
				fmt.Println("MASUK BARAT-", i)
				jumlahBarat++
				groupKeamananBarat = append(groupKeamananBarat, index)
			} else if index.Zona == "TIMUR" {
				fmt.Println("MASUK TIMUR-", i)
				jumlahTimur++
				groupKeamananTimur = append(groupKeamananTimur, index)
			} else if index.Zona == "TENGAH" {
				fmt.Println("MASUK TENGAH-", i)
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
			fmt.Println("TIMURRR")
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
				return nil
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

	templateData := struct {
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

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(outputPath)
	} else {
		fmt.Println(err)
	}

	return ctx.Response().Success().Json(map[string]interface{}{
		"Status":   "success",
		"triwulan": templateData,
	})
}

func sortedKeys(m map[string]map[string]int) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
