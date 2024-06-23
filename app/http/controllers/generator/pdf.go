package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/jung-kurt/gofpdf"
)

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
		"add": Add,
		"seq": Seq,
		"sub": Sub,
		"inc": Increment,
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
	cmd := exec.Command("wkhtmltoimage", "--enable-local-file-access", "--javascript-delay", "500", inputFile, slidePath)

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
	cmd := exec.Command("wkhtmltopdf", "--disable-smart-shrinking", "--page-size", "A4", "--javascript-delay", "1000", inputFile, laporanPath)

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

func (r *Pdf) Index(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/sample.html"
	newtemplatePath := "sample.html"

	//path for download pdf
	outputPath := "storage/example.png"

	//html template data
	templateData := struct {
		Title       string
		Description string
		Company     string
		Contact     string
		Country     string
	}{
		Title:       "HTML to PDF generator",
		Description: "This is the simple HTML to PDF file.",
		Company:     "Jhon Lewis",
		Contact:     "Maria Anders",
		Country:     "Germany",
	}

	if err := r.ParseTemplate(templatePath, newtemplatePath, templateData); err == nil {
		// Generate PDF
		ok, _ := r.GenerateSlide(outputPath)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}
	return nil
}

func (r *Pdf) GenerateMingguan(ctx http.Context) http.Response {
	//html template path
	templateKeamananPath := "templates/keamanan.html"
	newTemplateKeamananPath := "keamanan.html"
	templateKeselamatanPath := "templates/keselamatan.html"
	newTemplateKeselamatanPath := "keselamatan.html"

	var data_keamanan []models.KejadianKeamanan
	query1 := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query1.Order("tanggal asc").Find(&data_keamanan)

	var result_keamanan []models.KejadianKeamananImage
	for _, data := range data_keamanan {
		var data_keamanan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

		result_keamanan = append(result_keamanan, models.KejadianKeamananImage{
			KejadianKeamanan: data,
			FileImage:        data_keamanan_image,
		})
	}

	var data_keselamatan []models.KejadianKeselamatan
	query2 := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query2.Order("tanggal asc").Find(&data_keselamatan)

	var result_keselamatan []models.KejadianKeselamatanImage
	for _, data := range data_keselamatan {
		var data_keselamatan_image []models.FileImage
		facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
			Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

		result_keselamatan = append(result_keselamatan, models.KejadianKeselamatanImage{
			KejadianKeselamatan: data,
			FileImage:           data_keselamatan_image,
		})
	}

	var images []string
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
			return nil
		}
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
			return nil
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

	return nil
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

func (r *Pdf) GenerateBulanan(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/laporan-bulanan.html"
	newTemplatePath := "laporan-bulanan.html"

	var data_keamanan []models.KejadianKeamanan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query.Order("tanggal asc").Find(&data_keamanan)

	var data_keselamatan []models.KejadianKeselamatan
	query = facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query.Order("tanggal asc").Find(&data_keselamatan)

	outputPath := "storage/output-laporan-bulanan.pdf"

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range data_keamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeamanan []GroupingKeamanan
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
			} else if index.Zona == "TIMUR" {
				fmt.Println("MASUK TIMUR-", i)
				jumlahTimur++
			} else if index.Zona == "TENGAH" {
				fmt.Println("MASUK TENGAH-", i)
				jumlahTengah++
			}
			jumlah++
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

	var groupKeselamatan []GroupingKeselamatan
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

			if data.Zona == "BARAT" {
				jumlahBarat++
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			list_korban = append(list_korban, temp)
			jumlah++
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
		KejadianKeselamatan       []GroupingKeselamatan
	}{
		Bulan:                     "Mei",
		BulanCapital:              "MEI",
		Tahun:                     "2024",
		JumlahKejadianKeamanan:    len(data_keamanan),
		JumlahKejadianKeselamatan: len(data_keselamatan),
		KejadianKeamanan:          groupKeamanan,
		KejadianKeselamatan:       groupKeselamatan,
	}

	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
		r.GenerateLaporan(outputPath)
	} else {
		fmt.Println(err)
	}

	// fmt.Println("PDF created successfully!")
	return ctx.Response().Success().Json(map[string]interface{}{
		"Status": "success",
		"data-1": groupKeamanan,
		"data-2": groupKeselamatan,
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

func (r *Pdf) GenerateCoba(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/coba.html"
	newtemplatePath := "coba.html"

	var data_keamanan []models.KejadianKeamanan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query.Order("tanggal asc").Find(&data_keamanan)

	var data_keselamatan []models.KejadianKeselamatan
	query = facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
	query.Order("tanggal asc").Find(&data_keselamatan)

	outputPath := "storage/output-laporan-bulanan.pdf"

	// Group the incidents by 'jenis_kejadian_id'
	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
	for _, kejadian := range data_keamanan {
		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
	}

	var groupKeamanan []GroupingKeamanan
	// Print the grouped data
	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
		jumlah := 0
		jumlahBarat := 0
		jumlahTimur := 0
		jumlahTengah := 0
		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
		for i, index := range kejadianGroup {
			fmt.Println(index)
			if index.Zona == "BARAT" {
				fmt.Println("MASUK BARAT-", i)
				jumlahBarat++
			} else if index.Zona == "TIMUR" {
				fmt.Println("MASUK TIMUR-", i)
				jumlahTimur++
			} else if index.Zona == "TENGAH" {
				fmt.Println("MASUK TENGAH-", i)
				jumlahTengah++
			}
			jumlah++
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

	var groupKeselamatan []GroupingKeselamatan
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

			if data.Zona == "BARAT" {
				jumlahBarat++
			} else if data.Zona == "TIMUR" {
				jumlahTimur++
			} else if data.Zona == "TENGAH" {
				jumlahTengah++
			}

			temp := models.KejadianKeselamatanKorban{
				KejadianKeselamatan: data,
				ListKorban:          x,
			}

			list_korban = append(list_korban, temp)
			jumlah++
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
		KejadianKeselamatan       []GroupingKeselamatan
	}{
		Bulan:                     "Mei",
		BulanCapital:              "MEI",
		Tahun:                     "2024",
		JumlahKejadianKeamanan:    len(data_keamanan),
		JumlahKejadianKeselamatan: len(data_keselamatan),
		KejadianKeamanan:          groupKeamanan,
		KejadianKeselamatan:       groupKeselamatan,
	}

	if err := r.ParseTemplate(templatePath, newtemplatePath, templateData); err == nil {
		r.GenerateLaporan(outputPath)
	} else {
		fmt.Println(err)
	}

	// fmt.Println("PDF created successfully!")
	return ctx.Response().Success().Json(map[string]interface{}{
		"Status": "success",
		"data-1": groupKeamanan,
		"data-2": groupKeselamatan,
	})
}
