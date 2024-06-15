package generator

import (
	"bytes"
	"fmt"
	"goravel/app/models"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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

func (r *Pdf) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
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

func (r *Pdf) Index(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/sample.html"

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

	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		// Generate PDF
		ok, _ := r.GenerateSlide(outputPath)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}
	return nil
}

func (r *Pdf) GenerateKeamanan(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/keamanan.html"

	var data_keamanan []models.KejadianKeamanan
	query := facades.Orm().Query().
		Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
		With("JenisKejadian")
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

	var images []string
	for _, data := range results {
		// path for download pdf
		outputPath := fmt.Sprintf("storage/%d.png", data.IdKejadianKeamanan)
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
			ABK:              "-",
			Muatan:           data.Muatan,
			InstansiPenindak: "-",
			Keterangan:       "-",
			Waktu:            data.Tanggal.ToDateString(),
			SumberBerita:     data.SumberBerita,
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Images:           data.FileImage,
		}

		if err := r.ParseTemplate(templatePath, templateData); err == nil {
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
	err := pdf.OutputFileAndClose("storage/laporan-keamanan-mingguan.pdf")
	if err != nil {
		fmt.Printf("Error saving PDF: %s", err)
	}

	fmt.Println("PDF created successfully!")

	return nil
}
