package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

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

// generate slide function
func (r *Pdf) GenerateSlide(slidePath string) (bool, error) {
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

	// Define the input HTML file and output image file
	inputFile := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"

	// Check if the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
		panic(err)
	}

	// Construct the command
	cmd := exec.Command("wkhtmltoimage", "--javascript-delay", "500", inputFile, slidePath)

	// Set the command's standard output and error to the current process's standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		panic(err)
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
