package generator

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"goravel/app/models"
// 	"html/template"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// 	"regexp"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/goravel/framework/contracts/http"
// 	"github.com/goravel/framework/facades"
// 	"github.com/jung-kurt/gofpdf"
// )

// type Pdf struct {
// 	body string
// 	//Dependent services
// }

// func NewPdf(body string) *Pdf {
// 	return &Pdf{
// 		body: body,
// 		//Inject services
// 	}
// }

// func (r *Pdf) ParseTemplate(templateFileName string, newTemplateFileName string, data interface{}) error {
// 	t, err := template.New(newTemplateFileName).Funcs(template.FuncMap{
// 		"add": Add,
// 		"seq": Seq,
// 		"sub": Sub,
// 		"inc": Increment,
// 		"last": func(index int, length int) bool {
// 			return index == length-1
// 		},
// 	}).ParseFiles(templateFileName)

// 	if err != nil {
// 		return err
// 	}

// 	t = template.Must(t, err)

// 	buf := new(bytes.Buffer)
// 	if err = t.Execute(buf, data); err != nil {
// 		return err
// 	}
// 	r.body = buf.String()
// 	return nil
// }

// // generate slide function
// func (r *Pdf) GenerateSlide(slidePath string) (bool, error) {
// 	t := time.Now().Unix()
// 	// write whole the body

// 	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
// 		errDir := os.Mkdir("cloneTemplate/", 0777)
// 		if errDir != nil {
// 			fmt.Printf("Error: %v\n", errDir)
// 			return false, errDir
// 		}
// 	}
// 	err1 := ioutil.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
// 	if err1 != nil {
// 		fmt.Printf("Error: %v\n", err1)
// 		return false, err1
// 	}

// 	// Define the input HTML file and output image file
// 	inputFile := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"

// 	// Check if the input file exists
// 	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
// 		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
// 		return false, err
// 	}

// 	// Construct the command
// 	cmd := exec.Command("wkhtmltoimage", "--enable-local-file-access", "--javascript-delay", "500", inputFile, slidePath)

// 	// Set the command's standard output and error to the current process's standard output and error
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	// Run the command
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return false, err
// 	}

// 	dir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return false, err
// 	}

// 	defer os.RemoveAll(dir + "/cloneTemplate")

// 	return true, nil
// }

// // generate laporan function
// func (r *Pdf) GenerateLaporan(laporanPath string) (bool, error) {
// 	t := time.Now().Unix()
// 	// write whole the body

// 	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
// 		errDir := os.Mkdir("cloneTemplate/", 0777)
// 		if errDir != nil {
// 			fmt.Printf("Error: %v\n", errDir)
// 			return false, errDir
// 		}
// 	}
// 	err1 := ioutil.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
// 	if err1 != nil {
// 		fmt.Printf("Error: %v\n", err1)
// 		return false, err1
// 	}

// 	// Define the input HTML file and output image file
// 	inputFile := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"

// 	// Check if the input file exists
// 	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
// 		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
// 		return false, err
// 	}

// 	// Construct the command
// 	cmd := exec.Command("wkhtmltopdf", "--disable-smart-shrinking", "--page-size", "A4", "--javascript-delay", "1000", inputFile, laporanPath)

// 	// Set the command's standard output and error to the current process's standard output and error
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	// Run the command
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return false, err
// 	}

// 	dir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return false, err
// 	}

// 	defer os.RemoveAll(dir + "/cloneTemplate")

// 	return true, nil
// }

// func (r *Pdf) Index(ctx http.Context) http.Response {
// 	//html template path
// 	templatePath := "templates/sample.html"
// 	newtemplatePath := "sample.html"

// 	//path for download pdf
// 	outputPath := "storage/example.png"

// 	//html template data
// 	templateData := struct {
// 		Title       string
// 		Description string
// 		Company     string
// 		Contact     string
// 		Country     string
// 	}{
// 		Title:       "HTML to PDF generator",
// 		Description: "This is the simple HTML to PDF file.",
// 		Company:     "Jhon Lewis",
// 		Contact:     "Maria Anders",
// 		Country:     "Germany",
// 	}

// 	if err := r.ParseTemplate(templatePath, newtemplatePath, templateData); err == nil {
// 		// Generate PDF
// 		ok, _ := r.GenerateSlide(outputPath)
// 		fmt.Println(ok, "pdf generated successfully")
// 	} else {
// 		fmt.Println(err)
// 	}
// 	return nil
// }

// func (r *Pdf) GenerateMingguan(ctx http.Context) http.Response {
// 	//html template path
// 	templateKeamananPath := "templates/keamanan.html"
// 	newTemplateKeamananPath := "keamanan.html"
// 	templateKeselamatanPath := "templates/keselamatan.html"
// 	newTemplateKeselamatanPath := "keselamatan.html"

// 	var data_keamanan []models.KejadianKeamanan
// 	query1 := facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query1.Order("k.nama_kejadian asc, tanggal asc").Find(&data_keamanan)

// 	var result_keamanan []models.KejadianKeamananImage
// 	for _, data := range data_keamanan {
// 		var data_keamanan_image []models.FileImage
// 		facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
// 			Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

// 		result_keamanan = append(result_keamanan, models.KejadianKeamananImage{
// 			KejadianKeamanan: data,
// 			FileImage:        data_keamanan_image,
// 		})
// 	}

// 	var data_keselamatan []models.KejadianKeselamatan
// 	query2 := facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query2.Order("k.nama_kejadian asc, tanggal asc").Find(&data_keselamatan)

// 	var result_keselamatan []models.KejadianKeselamatanImage
// 	for _, data := range data_keselamatan {
// 		var data_keselamatan_image []models.FileImage
// 		facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
// 			Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

// 		result_keselamatan = append(result_keselamatan, models.KejadianKeselamatanImage{
// 			KejadianKeselamatan: data,
// 			FileImage:           data_keselamatan_image,
// 		})
// 	}

// 	var images []string
// 	for _, data := range result_keamanan {
// 		// path for download pdf
// 		outputPath := fmt.Sprintf("storage/temp/pelanggaran%d.png", data.IdKejadianKeamanan)

// 		var abk string
// 		if strings.Contains(data.Muatan, "ABK") {
// 			re := regexp.MustCompile(`\b\d+\s+orang\b`)
// 			matches := re.FindAllString(data.Muatan, -1)
// 			if len(matches) > 0 {
// 				abk = matches[0]
// 			} else {
// 				abk = " - "
// 			}
// 		} else {
// 			abk = " - "
// 		}

// 		// html template data
// 		templateData := struct {
// 			Title            string
// 			NamaKapal        string
// 			Kejadian         string
// 			Penyebab         string
// 			Lokasi           string
// 			ABK              string
// 			Muatan           string
// 			InstansiPenindak string
// 			Keterangan       string
// 			Waktu            string
// 			SumberBerita     string
// 			Latitude         float64
// 			Longitude        float64
// 			Images           []models.FileImage
// 		}{
// 			Title:            data.JenisKejadian.NamaKejadian,
// 			NamaKapal:        data.NamaKapal,
// 			Kejadian:         data.JenisKejadian.NamaKejadian,
// 			Penyebab:         "-",
// 			Lokasi:           data.LokasiKejadian,
// 			ABK:              abk,
// 			Muatan:           data.Muatan,
// 			InstansiPenindak: data.SumberBerita,
// 			Keterangan:       data.TindakLanjut,
// 			Waktu:            data.Tanggal.ToDateString(),
// 			SumberBerita:     data.LinkBerita,
// 			Latitude:         data.Latitude,
// 			Longitude:        data.Longitude,
// 			Images:           data.FileImage,
// 		}

// 		if err := r.ParseTemplate(templateKeamananPath, newTemplateKeamananPath, templateData); err == nil {
// 			// Generate Image
// 			success, _ := r.GenerateSlide(outputPath)
// 			if success {
// 				images = append(images, outputPath)
// 			}
// 		} else {
// 			fmt.Printf("Error: %v\n", err)
// 			return nil
// 		}
// 	}

// 	for _, data := range result_keselamatan {
// 		// path for download pdf
// 		outputPath := fmt.Sprintf("storage/temp/kecelakaan%d.png", data.IdKejadianKeselamatan)

// 		var perpindahanAwal string
// 		if data.PelabuhanAsal != "-" && data.PelabuhanAsal != "" {
// 			perpindahanAwal = data.PelabuhanAsal
// 		}
// 		var perpindahanAkhir string
// 		if data.PelabuhanTujuan != "-" && data.PelabuhanTujuan != "" {
// 			perpindahanAwal = data.PelabuhanTujuan
// 		}

// 		var perpindahan string
// 		if perpindahanAwal != "" && perpindahanAkhir != "" {
// 			perpindahan = perpindahanAwal + " - " + perpindahanAkhir
// 		} else if perpindahanAwal != "" && perpindahanAkhir == "" {
// 			perpindahan = perpindahanAwal + " - "
// 		} else if perpindahanAwal == "" && perpindahanAkhir != "" {
// 			perpindahan = " - " + perpindahanAkhir
// 		} else {
// 			perpindahan = " - "
// 		}

// 		var korbanData models.ListKorban

// 		var korban string
// 		if err := json.Unmarshal(data.Korban, &korbanData); err != nil {
// 			fmt.Println("ERROR", err)
// 			return nil
// 		}

// 		if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas != 0 {
// 			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, selamat " +
// 				strconv.Itoa(korbanData.KorbanSelamat) + " orang, dan tewas " +
// 				strconv.Itoa(korbanData.KorbanTewas) + " orang"
// 		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas == 0 {
// 			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, dan selamat " +
// 				strconv.Itoa(korbanData.KorbanSelamat) + " orang"
// 		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas != 0 {
// 			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang, dan tewas " +
// 				strconv.Itoa(korbanData.KorbanTewas) + " orang"
// 		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas != 0 {
// 			korban = "Korban selamat " + strconv.Itoa(korbanData.KorbanSelamat) + " orang, dan tewas " +
// 				strconv.Itoa(korbanData.KorbanTewas) + " orang"
// 		} else if korbanData.KorbanHilang != 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas == 0 {
// 			korban = "Korban hilang " + strconv.Itoa(korbanData.KorbanHilang) + " orang"
// 		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat != 0 && korbanData.KorbanTewas == 0 {
// 			korban = "Korban selamat " + strconv.Itoa(korbanData.KorbanSelamat) + " orang"
// 		} else if korbanData.KorbanHilang == 0 && korbanData.KorbanSelamat == 0 && korbanData.KorbanTewas != 0 {
// 			korban = "Korban tewas " + strconv.Itoa(korbanData.KorbanTewas) + " orang"
// 		} else {
// 			korban = "tidak ada korban jiwa"
// 		}

// 		// html template data
// 		templateData := struct {
// 			Title            string
// 			NamaKapal        string
// 			Kejadian         string
// 			Penyebab         string
// 			Lokasi           string
// 			Korban           string
// 			Perpindahan      string
// 			Keterangan       string
// 			Waktu            string
// 			InstansiPenindak string
// 			SumberBerita     string
// 			Latitude         float64
// 			Longitude        float64
// 			Images           []models.FileImage
// 		}{
// 			Title:            data.JenisKejadian.NamaKejadian,
// 			NamaKapal:        data.NamaKapal,
// 			Kejadian:         data.JenisKejadian.NamaKejadian,
// 			Penyebab:         data.Penyebab,
// 			Lokasi:           data.LokasiKejadian,
// 			Korban:           korban,
// 			Perpindahan:      perpindahan,
// 			Keterangan:       data.TindakLanjut,
// 			Waktu:            data.Tanggal.ToDateString(),
// 			InstansiPenindak: data.SumberBerita,
// 			SumberBerita:     data.LinkBerita,
// 			Latitude:         data.Latitude,
// 			Longitude:        data.Longitude,
// 			Images:           data.FileImage,
// 		}

// 		if err := r.ParseTemplate(templateKeselamatanPath, newTemplateKeselamatanPath, templateData); err == nil {
// 			// Generate Image
// 			success, _ := r.GenerateSlide(outputPath)
// 			if success {
// 				images = append(images, outputPath)
// 			}
// 		} else {
// 			fmt.Printf("Error: %v\n", err)
// 			return nil
// 		}
// 	}

// 	// Create a new PDF document
// 	pdf := gofpdf.New("L", "mm", "A4", "")

// 	for _, image := range images {
// 		// Add a new page to the PDF
// 		pdf.AddPage()

// 		// Get the image dimensions
// 		options := gofpdf.ImageOptions{
// 			ReadDpi: true,
// 		}
// 		info := pdf.RegisterImageOptions(image, options)
// 		width, height := info.Extent()

// 		// Calculate the position to center the image on the page
// 		pageWidth, pageHeight := pdf.GetPageSize()
// 		x := (pageWidth - width) / 2
// 		y := (pageHeight - height) / 2

// 		// Add the image to the PDF
// 		pdf.ImageOptions(image, x, y, width, height, false, options, 0, "")
// 	}

// 	// Save the PDF to a file
// 	// path := strconv.Itoa(req.TahunKe) + "/" + req.JenisLaporan + "/Bulan " + monthName(req.BulanKe)
// 	err := pdf.OutputFileAndClose("storage/laporan-keamanan-mingguan.pdf")
// 	if err != nil {
// 		fmt.Printf("Error saving PDF: %s", err)
// 	}

// 	fmt.Println("PDF created successfully!")

// 	return nil
// }

// type GroupingKeamanan struct {
// 	NamaKejadian           string `json:"nama_kejadian"`
// 	KejadianKeamanan       []models.KejadianKeamanan
// 	KejadianKeamananBarat  []models.KejadianKeamanan
// 	KejadianKeamananTimur  []models.KejadianKeamanan
// 	KejadianKeamananTengah []models.KejadianKeamanan
// 	Jumlah                 int `json:"jumlah"`
// 	JumlahZonaBarat        int `json:"jumlah_zona_barat"`
// 	JumlahZonaTimur        int `json:"jumlah_zona_timur"`
// 	JumlahZonaTengah       int `json:"jumlah_zona_tengah"`
// }

// type GroupingKeselamatan struct {
// 	NamaKejadian              string `json:"nama_kejadian"`
// 	KejadianKeselamatan       []models.KejadianKeselamatanKorban
// 	KejadianKeselamatanBarat  []models.KejadianKeselamatanKorban
// 	KejadianKeselamatanTimur  []models.KejadianKeselamatanKorban
// 	KejadianKeselamatanTengah []models.KejadianKeselamatanKorban
// 	Jumlah                    int `json:"jumlah"`
// 	JumlahZonaBarat           int `json:"jumlah_zona_barat"`
// 	JumlahZonaTimur           int `json:"jumlah_zona_timur"`
// 	JumlahZonaTengah          int `json:"jumlah_zona_tengah"`
// }

// func (r *Pdf) GenerateBulanan(ctx http.Context) http.Response {
// 	//html template path
// 	now := time.Now()
// 	bulan := monthNameIndonesia(now.Month())
// 	year := strconv.Itoa(now.Year())

// 	// lastMonth1 := int(now.Month())
// 	templatePath := "templates/laporan-bulanan.html"
// 	newTemplatePath := "laporan-bulanan.html"

// 	var data_keamanan []models.KejadianKeamanan
// 	query := facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query.Order("tanggal asc").Find(&data_keamanan)

// 	var data_keselamatan []models.KejadianKeselamatan
// 	query = facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query.Order("tanggal asc").Find(&data_keselamatan)

// 	outputPath := "storage/output-laporan-bulanan.pdf"

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
// 	for _, kejadian := range data_keamanan {
// 		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeamanan []GroupingKeamanan
// 	var groupKeamananBarat []models.KejadianKeamanan
// 	// var jenisNamaKeamananBarat []string
// 	// var jenisNamaKeamananTimur []string
// 	// var jenisNamaKeamananTengah []string
// 	// if jumlahBarat != 0 {
// 	// 	jenisNamaKeamananBarat = append(jenisNamaKeamananBarat, jenisName)
// 	// } else if jumlahTimur != 0 {
// 	// 	jenisNamaKeamananTimur = append(jenisNamaKeamananTimur, jenisName)
// 	// } else if jumlahTengah != 0 {
// 	// 	jenisNamaKeamananTengah = append(jenisNamaKeamananTengah, jenisName)
// 	// }
// 	var groupKeamananTimur []models.KejadianKeamanan
// 	var groupKeamananTengah []models.KejadianKeamanan
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0

// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		for i, index := range kejadianGroup {
// 			if index.Zona == "BARAT" {
// 				fmt.Println("MASUK BARAT-", i)
// 				jumlahBarat++
// 				groupKeamananBarat = append(groupKeamananBarat, index)
// 			} else if index.Zona == "TIMUR" {
// 				fmt.Println("MASUK TIMUR-", i)
// 				jumlahTimur++
// 				groupKeamananTimur = append(groupKeamananTimur, index)
// 			} else if index.Zona == "TENGAH" {
// 				fmt.Println("MASUK TENGAH-", i)
// 				jumlahTengah++
// 				groupKeamananTengah = append(groupKeamananTengah, index)
// 			}
// 			jumlah++
// 		}

// 		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
// 			NamaKejadian:           jenisName,
// 			KejadianKeamanan:       kejadianGroup,
// 			KejadianKeamananBarat:  groupKeamananBarat,
// 			KejadianKeamananTimur:  groupKeamananTimur,
// 			KejadianKeamananTengah: groupKeamananTengah,
// 			Jumlah:                 jumlah,
// 			JumlahZonaBarat:        jumlahBarat,
// 			JumlahZonaTimur:        jumlahTimur,
// 			JumlahZonaTengah:       jumlahTengah,
// 		})
// 	}

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
// 	for _, kejadian := range data_keselamatan {
// 		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeselamatan []GroupingKeselamatan
// 	var groupKeselamatanBarat []models.KejadianKeselamatanKorban
// 	var groupKeselamatanTimur []models.KejadianKeselamatanKorban
// 	var groupKeselamatanTengah []models.KejadianKeselamatanKorban
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeselamatan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0

// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		var list_korban []models.KejadianKeselamatanKorban

// 		for _, data := range kejadianGroup {
// 			var x models.ListKorban
// 			err := json.Unmarshal(data.Korban, &x)
// 			if err != nil {
// 				return nil
// 			}

// 			temp := models.KejadianKeselamatanKorban{
// 				KejadianKeselamatan: data,
// 				ListKorban:          x,
// 			}

// 			if data.Zona == "BARAT" {
// 				jumlahBarat++
// 				groupKeselamatanBarat = append(groupKeselamatanBarat, temp)
// 			} else if data.Zona == "TIMUR" {
// 				jumlahTimur++
// 				groupKeselamatanTimur = append(groupKeselamatanTimur, temp)
// 			} else if data.Zona == "TENGAH" {
// 				jumlahTengah++
// 				groupKeselamatanTengah = append(groupKeselamatanTengah, temp)
// 			}

// 			list_korban = append(list_korban, temp)
// 			jumlah++
// 		}

// 		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
// 			NamaKejadian:              jenisName,
// 			KejadianKeselamatan:       list_korban,
// 			KejadianKeselamatanBarat:  groupKeselamatanBarat,
// 			KejadianKeselamatanTimur:  groupKeselamatanTimur,
// 			KejadianKeselamatanTengah: groupKeselamatanTengah,
// 			Jumlah:                    jumlah,
// 			JumlahZonaBarat:           jumlahBarat,
// 			JumlahZonaTimur:           jumlahTimur,
// 			JumlahZonaTengah:          jumlahTengah,
// 		})
// 	}
// 	// html template data
// 	templateData := struct {
// 		Bulan                     string
// 		BulanCapital              string
// 		Tahun                     string
// 		JumlahKejadianKeamanan    int
// 		JumlahKejadianKeselamatan int
// 		KejadianKeamanan          []GroupingKeamanan
// 		KejadianKeselamatan       []GroupingKeselamatan
// 	}{
// 		Bulan:                     bulan,
// 		BulanCapital:              strings.ToUpper(bulan),
// 		Tahun:                     year,
// 		JumlahKejadianKeamanan:    len(data_keamanan),
// 		JumlahKejadianKeselamatan: len(data_keselamatan),
// 		KejadianKeamanan:          groupKeamanan,
// 		KejadianKeselamatan:       groupKeselamatan,
// 	}

// 	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
// 		r.GenerateLaporan(outputPath)
// 	} else {
// 		fmt.Println(err)
// 	}

// 	// fmt.Println("PDF created successfully!")
// 	return ctx.Response().Success().Json(map[string]interface{}{
// 		"Status": "success",
// 		"data-1": groupKeamanan,
// 		"data-2": groupKeselamatan,
// 	})
// }

// func Add(x, y int) int {
// 	return x + y
// }

// func Seq(start, end int) []int {
// 	s := make([]int, end-start+1)
// 	for i := range s {
// 		s[i] = start + i
// 	}
// 	return s
// }

// // sub subtracts b from a.
// func Sub(a, b int) int {
// 	return a - b
// }

// func Increment(val int) int {
// 	return val + 1
// }

// func (r *Pdf) GenerateCoba(ctx http.Context) http.Response {
// 	//html template path
// 	templatePath := "templates/coba.html"
// 	newtemplatePath := "coba.html"

// 	var data_keamanan []models.KejadianKeamanan
// 	query := facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query.Order("tanggal asc").Find(&data_keamanan)

// 	var data_keselamatan []models.KejadianKeselamatan
// 	query = facades.Orm().Query().
// 		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
// 		With("JenisKejadian")
// 	query.Order("tanggal asc").Find(&data_keselamatan)

// 	outputPath := "storage/output-laporan-bulanan.pdf"

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
// 	for _, kejadian := range data_keamanan {
// 		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeamanan []GroupingKeamanan
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0
// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		for i, index := range kejadianGroup {
// 			fmt.Println(index)
// 			if index.Zona == "BARAT" {
// 				fmt.Println("MASUK BARAT-", i)
// 				jumlahBarat++
// 			} else if index.Zona == "TIMUR" {
// 				fmt.Println("MASUK TIMUR-", i)
// 				jumlahTimur++
// 			} else if index.Zona == "TENGAH" {
// 				fmt.Println("MASUK TENGAH-", i)
// 				jumlahTengah++
// 			}
// 			jumlah++
// 		}
// 		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
// 			NamaKejadian:     jenisName,
// 			KejadianKeamanan: kejadianGroup,
// 			Jumlah:           jumlah,
// 			JumlahZonaBarat:  jumlahBarat,
// 			JumlahZonaTimur:  jumlahTimur,
// 			JumlahZonaTengah: jumlahTengah,
// 		})
// 	}

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
// 	for _, kejadian := range data_keselamatan {
// 		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeselamatan []GroupingKeselamatan
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeselamatan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0

// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		var list_korban []models.KejadianKeselamatanKorban

// 		for _, data := range kejadianGroup {
// 			var x models.ListKorban
// 			err := json.Unmarshal(data.Korban, &x)
// 			if err != nil {
// 				return nil
// 			}

// 			if data.Zona == "BARAT" {
// 				jumlahBarat++
// 			} else if data.Zona == "TIMUR" {
// 				jumlahTimur++
// 			} else if data.Zona == "TENGAH" {
// 				jumlahTengah++
// 			}

// 			temp := models.KejadianKeselamatanKorban{
// 				KejadianKeselamatan: data,
// 				ListKorban:          x,
// 			}

// 			list_korban = append(list_korban, temp)
// 			jumlah++
// 		}

// 		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
// 			NamaKejadian:        jenisName,
// 			KejadianKeselamatan: list_korban,
// 			Jumlah:              jumlah,
// 			JumlahZonaBarat:     jumlahBarat,
// 			JumlahZonaTimur:     jumlahTimur,
// 			JumlahZonaTengah:    jumlahTengah,
// 		})
// 	}
// 	// html template data
// 	templateData := struct {
// 		Bulan                     string
// 		BulanCapital              string
// 		Tahun                     string
// 		JumlahKejadianKeamanan    int
// 		JumlahKejadianKeselamatan int
// 		KejadianKeamanan          []GroupingKeamanan
// 		KejadianKeselamatan       []GroupingKeselamatan
// 	}{
// 		Bulan:                     "Mei",
// 		BulanCapital:              "MEI",
// 		Tahun:                     "2024",
// 		JumlahKejadianKeamanan:    len(data_keamanan),
// 		JumlahKejadianKeselamatan: len(data_keselamatan),
// 		KejadianKeamanan:          groupKeamanan,
// 		KejadianKeselamatan:       groupKeselamatan,
// 	}

// 	if err := r.ParseTemplate(templatePath, newtemplatePath, templateData); err == nil {
// 		r.GenerateLaporan(outputPath)
// 	} else {
// 		fmt.Println(err)
// 	}

// 	// fmt.Println("PDF created successfully!")
// 	return ctx.Response().Success().Json(map[string]interface{}{
// 		"Status": "success",
// 		"data-1": groupKeamanan,
// 		"data-2": groupKeselamatan,
// 	})
// }

// type MonthlyCount struct {
// 	NamaKejadian  string `json:"nama_kejadian"`
// 	Bulan1        string `json:"bulan_1"`
// 	KorbanTewas   int    `json:"korban_tewas"`
// 	KorbanSelamat int    `json:"korban_selamat"`
// 	KorbanHilang  int    `json:"korban_hilang"`
// 	Count1        int    `json:"count_1"`
// 	Bulan2        string `json:"bulan_2"`
// 	Count2        int    `json:"count_2"`
// 	Bulan3        string `json:"bulan_3"`
// 	Count3        int    `json:"count_3"`
// 	Total         int    `json:"total"`
// }

// type LocationOutput struct {
// 	NamaKejadian   string `json:"nama_kejadian"`
// 	JumlahDermaga  int    `json:"jumlah_dermaga"`
// 	JumlahPerairan int    `json:"jumlah_perairan"`
// }

// var monthNameMap = map[string]int{
// 	"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MEI": 5, "JUN": 6,
// 	"JUL": 7, "AGT": 8, "SEP": 9, "OKT": 10, "NOV": 11, "DES": 12,
// }

// var monthNameEnglishMap = map[time.Month]string{
// 	time.January: "JAN", time.February: "FEB", time.March: "MAR",
// 	time.April: "APR", time.May: "MEI", time.June: "JUN",
// 	time.July: "JUL", time.August: "AGT", time.September: "SEP",
// 	time.October: "OKT", time.November: "NOV", time.December: "DES",
// }

// var monthNameIndonesiaMap = map[time.Month]string{
// 	time.January: "Januari", time.February: "Februari", time.March: "Maret",
// 	time.April: "April", time.May: "Mei", time.June: "Juni",
// 	time.July: "Juli", time.August: "Agustus", time.September: "September",
// 	time.October: "Oktober", time.November: "November", time.December: "Desember",
// }

// func monthNameEnglish(month time.Month) string {
// 	return monthNameEnglishMap[month]
// }

// func monthNameIndonesia(month time.Month) string {
// 	return monthNameIndonesiaMap[month]
// }

// func (r *Pdf) GenerateTriwulan(ctx http.Context) http.Response {
// 	const (
// 		templatePath    = "templates/laporan-triwulan.html"
// 		newTemplatePath = "laporan-triwulan.html"
// 		outputPath      = "storage/output-laporan-triwulan.pdf"
// 	)

// 	now := time.Now()
// 	year := strconv.Itoa(now.Year())

// 	quarters := map[time.Month]struct {
// 		quarter      string
// 		periodFormat string
// 		months       []string
// 	}{
// 		time.January:   {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
// 		time.February:  {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
// 		time.March:     {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
// 		time.April:     {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
// 		time.May:       {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
// 		time.June:      {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
// 		time.July:      {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
// 		time.August:    {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
// 		time.September: {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
// 		time.October:   {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
// 		time.November:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
// 		time.December:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
// 	}

// 	quarterInfo := quarters[now.Month()]
// 	triwulanKe := quarterInfo.quarter
// 	periodeBulan := fmt.Sprintf(quarterInfo.periodFormat, year)
// 	months := quarterInfo.months

// 	var dataKeamanan []models.KejadianKeamanan
// 	var default1, default2 []models.JenisKejadian
// 	var dataKeselamatan []models.KejadianKeselamatan

// 	monthNumbers := make([]int, len(months))
// 	for i, month := range months {
// 		monthNumbers[i] = monthNameMap[month]
// 	}

// 	query1 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
// 		With("JenisKejadian").
// 		Where("DATE_PART('month', tanggal) IN (?)", monthNumbers).
// 		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeamanan)
// 	query2 := facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
// 		With("JenisKejadian").
// 		Where("DATE_PART('month', tanggal) IN (?)", monthNumbers).
// 		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeselamatan)
// 	query3 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keamanan Laut").
// 		Order("nama_kejadian asc").Find(&default1)
// 	query4 := facades.Orm().Query().Where("klasifikasi_name = ?", "Keselamatan Laut").
// 		Order("nama_kejadian asc").Find(&default2)

// 	if query1 != nil && query2 != nil && query3 != nil && query4 != nil {
// 		fmt.Println("failed to execute query")
// 		return nil
// 	}

// 	kejadianCountKeamanan := make(map[string]map[string]int)
// 	kejadianCountKeselamatan := make(map[string]map[string]int)
// 	locationCountKeamanan := make(map[string]map[string]int)

// 	for _, kejadian := range default1 {
// 		for _, month := range months {
// 			if kejadianCountKeamanan[kejadian.NamaKejadian] == nil {
// 				kejadianCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
// 			}
// 			if _, exists := kejadianCountKeamanan[kejadian.NamaKejadian][month]; !exists {
// 				kejadianCountKeamanan[kejadian.NamaKejadian][month] = 0
// 			}
// 		}

// 		if locationCountKeamanan[kejadian.NamaKejadian] == nil {
// 			locationCountKeamanan[kejadian.NamaKejadian] = make(map[string]int)
// 		}
// 	}
// 	for _, kejadian := range default2 {
// 		for _, month := range months {
// 			if kejadianCountKeselamatan[kejadian.NamaKejadian] == nil {
// 				kejadianCountKeselamatan[kejadian.NamaKejadian] = make(map[string]int)
// 			}
// 			if _, exists := kejadianCountKeselamatan[kejadian.NamaKejadian][month]; !exists {
// 				kejadianCountKeselamatan[kejadian.NamaKejadian][month] = 0
// 			}
// 		}
// 	}

// 	for _, kejadian := range dataKeamanan {
// 		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
// 		kejadianCountKeamanan[kejadian.JenisKejadian.NamaKejadian][month]++

// 		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "dermaga") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "pelabuhan") ||
// 			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Asal), "pelabuhan"))) ||
// 			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "dermaga") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "pelabuhan"))) {
// 			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["dermaga/pelabuhan"]++
// 		}
// 		if strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "laut") || strings.Contains(strings.ToLower(kejadian.LokasiKejadian), "perairan") ||
// 			(kejadian.Asal != nil && (strings.Contains(strings.ToLower(*kejadian.Asal), "laut") || strings.Contains(strings.ToLower(*kejadian.Asal), "perairan"))) ||
// 			(kejadian.Tujuan != nil && (strings.Contains(strings.ToLower(*kejadian.Tujuan), "laut") || strings.Contains(strings.ToLower(*kejadian.Tujuan), "perairan"))) {
// 			locationCountKeamanan[kejadian.JenisKejadian.NamaKejadian]["laut/perairan"]++
// 		}
// 	}
// 	for _, kejadian := range dataKeselamatan {
// 		month := monthNameEnglish(time.Month(kejadian.Tanggal.Month()))
// 		kejadianCountKeselamatan[kejadian.JenisKejadian.NamaKejadian][month]++
// 	}

// 	var kejadianCountsKeamanan, kejadianCountsKeselamatan []MonthlyCount

// 	jenisKejadianKeysKeamanan := sortedKeys(kejadianCountKeamanan)
// 	jenisKejadianKeysKeselamatan := sortedKeys(kejadianCountKeselamatan)
// 	locationsKeysKeamanan := sortedKeys(locationCountKeamanan)

// 	for _, jenisKejadian := range jenisKejadianKeysKeamanan {
// 		monthCounts := kejadianCountKeamanan[jenisKejadian]
// 		entry := MonthlyCount{
// 			NamaKejadian: jenisKejadian,
// 			Bulan1:       months[0],
// 			Count1:       monthCounts[months[0]],
// 			Bulan2:       months[1],
// 			Count2:       monthCounts[months[1]],
// 			Bulan3:       months[2],
// 			Count3:       monthCounts[months[2]],
// 			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
// 		}
// 		kejadianCountsKeamanan = append(kejadianCountsKeamanan, entry)
// 	}
// 	for _, jenisKejadian := range jenisKejadianKeysKeselamatan {
// 		monthCounts := kejadianCountKeselamatan[jenisKejadian]
// 		entry := MonthlyCount{
// 			NamaKejadian: jenisKejadian,
// 			Bulan1:       months[0],
// 			Count1:       monthCounts[months[0]],
// 			Bulan2:       months[1],
// 			Count2:       monthCounts[months[1]],
// 			Bulan3:       months[2],
// 			Count3:       monthCounts[months[2]],
// 			Total:        monthCounts[months[0]] + monthCounts[months[1]] + monthCounts[months[2]],
// 		}
// 		kejadianCountsKeselamatan = append(kejadianCountsKeselamatan, entry)
// 	}

// 	var locationOutput []LocationOutput
// 	for _, jenisKejadian := range locationsKeysKeamanan {
// 		counts := locationCountKeamanan[jenisKejadian]
// 		locationOutput = append(locationOutput, LocationOutput{
// 			NamaKejadian:   jenisKejadian,
// 			JumlahDermaga:  counts["dermaga/pelabuhan"],
// 			JumlahPerairan: counts["laut/perairan"],
// 		})
// 	}

// 	var bulan []string
// 	for _, x := range monthNumbers {
// 		bulan = append(bulan, monthNameIndonesia(time.Month(x)))
// 	}

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeamanan := make(map[string][]models.KejadianKeamanan)
// 	for _, kejadian := range dataKeamanan {
// 		groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeamanan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeamanan []GroupingKeamanan
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeamanan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0

// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		for i, index := range kejadianGroup {
// 			if index.Zona == "BARAT" {
// 				fmt.Println("MASUK BARAT-", i)
// 				jumlahBarat++
// 			} else if index.Zona == "TIMUR" {
// 				fmt.Println("MASUK TIMUR-", i)
// 				jumlahTimur++
// 			} else if index.Zona == "TENGAH" {
// 				fmt.Println("MASUK TENGAH-", i)
// 				jumlahTengah++
// 			}
// 			jumlah++
// 		}
// 		groupKeamanan = append(groupKeamanan, GroupingKeamanan{
// 			NamaKejadian:     jenisName,
// 			KejadianKeamanan: kejadianGroup,
// 			Jumlah:           jumlah,
// 			JumlahZonaBarat:  jumlahBarat,
// 			JumlahZonaTimur:  jumlahTimur,
// 			JumlahZonaTengah: jumlahTengah,
// 		})
// 	}

// 	// Group the incidents by 'jenis_kejadian_id'
// 	groupedByJenisKeselamatan := make(map[string][]models.KejadianKeselamatan)
// 	for _, kejadian := range dataKeselamatan {
// 		groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian] = append(groupedByJenisKeselamatan[kejadian.JenisKejadian.NamaKejadian], kejadian)
// 	}

// 	var groupKeselamatan []GroupingKeselamatan
// 	// Print the grouped data
// 	for jenisName, kejadianGroup := range groupedByJenisKeselamatan {
// 		jumlah := 0
// 		jumlahBarat := 0
// 		jumlahTimur := 0
// 		jumlahTengah := 0

// 		fmt.Printf("Jenis Kejadian ID: %s\n", jenisName)
// 		var list_korban []models.KejadianKeselamatanKorban

// 		for _, data := range kejadianGroup {
// 			var x models.ListKorban
// 			err := json.Unmarshal(data.Korban, &x)
// 			if err != nil {
// 				return nil
// 			}

// 			if data.Zona == "BARAT" {
// 				jumlahBarat++
// 			} else if data.Zona == "TIMUR" {
// 				jumlahTimur++
// 			} else if data.Zona == "TENGAH" {
// 				jumlahTengah++
// 			}

// 			temp := models.KejadianKeselamatanKorban{
// 				KejadianKeselamatan: data,
// 				ListKorban:          x,
// 			}

// 			list_korban = append(list_korban, temp)
// 			jumlah++
// 		}

// 		groupKeselamatan = append(groupKeselamatan, GroupingKeselamatan{
// 			NamaKejadian:        jenisName,
// 			KejadianKeselamatan: list_korban,
// 			Jumlah:              jumlah,
// 			JumlahZonaBarat:     jumlahBarat,
// 			JumlahZonaTimur:     jumlahTimur,
// 			JumlahZonaTengah:    jumlahTengah,
// 		})
// 	}

// 	templateData := struct {
// 		PeriodeTriwulan                          string
// 		BulanCapital                             string
// 		BulanSingkatan                           []string
// 		Bulan                                    []string
// 		TableKejadianKeamanan                    []MonthlyCount
// 		KejadianKeamanan                         []GroupingKeamanan
// 		TableKejadianKeselamatan                 []MonthlyCount
// 		KejadianKeselamatan                      []GroupingKeselamatan
// 		TablePengelompokanLokasiKejadianKeamanan []LocationOutput
// 		Tahun                                    string
// 	}{
// 		PeriodeTriwulan:                          triwulanKe,
// 		BulanCapital:                             periodeBulan,
// 		BulanSingkatan:                           months,
// 		Bulan:                                    bulan,
// 		TableKejadianKeamanan:                    kejadianCountsKeamanan,
// 		KejadianKeamanan:                         groupKeamanan,
// 		TableKejadianKeselamatan:                 kejadianCountsKeselamatan,
// 		KejadianKeselamatan:                      groupKeselamatan,
// 		TablePengelompokanLokasiKejadianKeamanan: locationOutput,
// 		Tahun:                                    year,
// 	}

// 	if err := r.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
// 		r.GenerateLaporan(outputPath)
// 	} else {
// 		fmt.Println(err)
// 	}

// 	return ctx.Response().Success().Json(map[string]interface{}{
// 		"Status":   "success",
// 		"triwulan": templateData,
// 	})
// }

// func sortedKeys(m map[string]map[string]int) []string {
// 	keys := make([]string, 0, len(m))
// 	for key := range m {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)
// 	return keys
// }
