package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/goravel/framework/contracts/http"
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

// generate pdf function
func (r *Pdf) GeneratePDF(pdfPath string, args []string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("cloneTemplate/", 0777)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	err1 := ioutil.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}

	f, err := os.Open("cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Use arguments to customize PDF generation process
	for _, arg := range args {
		switch arg {
		case "low-quality":
			pdfg.LowQuality.Set(true)
		case "no-pdf-compression":
			pdfg.NoPdfCompression.Set(true)
		case "grayscale":
			pdfg.Grayscale.Set(true)
			// Add other arguments as needed
		}
	}

	page := wkhtmltopdf.NewPageReader(f)
	page.EnableLocalFileAccess.Set(true) // Enable local file access
	page.EnablePlugins.Set(true)
	page.EnableTocBackLinks.Set(true)
	page.EnableForms.Set(true)

	pdfg.AddPage(page)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	// Retry mechanism
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = pdfg.Create()
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to create PDF: %v", i+1, err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir + "/cloneTemplate")

	return true, nil
}

func (r *Pdf) Index(ctx http.Context) http.Response {
	//html template path
	templatePath := "templates/sample.html"

	//path for download pdf
	outputPath := "storage/example.pdf"

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

		// Generate PDF with custom arguments
		args := []string{"no-pdf-compression"}

		// Generate PDF
		ok, _ := r.GeneratePDF(outputPath, args)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}
	return nil
}
