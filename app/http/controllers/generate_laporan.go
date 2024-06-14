package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/unidoc/unioffice/color"

	"fmt"
	"log"
	"strings"

	template "html/template"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/schema/soo/pml"
)

type GenerateLaporan struct {
	//Dependent services
}

func NewGenerateLaporan() *GenerateLaporan {
	return &GenerateLaporan{
		//Inject services
	}
}

func (r *GenerateLaporan) Index(ctx http.Context) http.Response {
	ppt, err := presentation.OpenTemplate("coba_template.pptx")
	if err != nil {
		log.Fatalf("unable to open template: %s", err)
	}
	for i, layout := range ppt.SlideLayouts() {
		fmt.Println(i, " LL ", layout.Name(), "/", layout.Type())
	}

	// remove any existing slides
	for _, s := range ppt.Slides() {
		ppt.RemoveSlide(s)
	}
	l, err := ppt.GetLayoutByName("Title and Content")
	if err != nil {
		log.Fatalf("error retrieving layout: %s", err)
	}
	sld, err := ppt.AddDefaultSlideWithLayout(l)
	if err != nil {
		log.Fatalf("error adding slide: %s", err)
	}

	ph, _ := sld.GetPlaceholder(pml.ST_PlaceholderTypeTitle)
	ph.SetText("Using unioffice")
	ph, _ = sld.GetPlaceholder(pml.ST_PlaceholderTypeBody)
	ph.SetText("Created with github.com/unidoc/unioffice/")

	tac, _ := ppt.GetLayoutByName("Title and Content")

	sld, err = ppt.AddDefaultSlideWithLayout(tac)
	if err != nil {
		log.Fatalf("error adding slide: %s", err)
	}

	ph, _ = sld.GetPlaceholder(pml.ST_PlaceholderTypeTitle)
	ph.SetText("Placeholders")
	ph, _ = sld.GetPlaceholderByIndex(1)
	ph.ClearAll()
	para := ph.AddParagraph()

	run := para.AddRun()
	run.SetText("Adding paragraphs can create bullets depending on the placeholder")
	para.AddBreak()
	run = para.AddRun()
	run.SetText("Line breaks work as expected within a paragraph")

	for i := 1; i < 5; i++ {
		para = ph.AddParagraph()
		para.Properties().SetLevel(int32(i))
		run = para.AddRun()
		run.SetText("Level controls indentation")
	}

	para = ph.AddParagraph()
	run = para.AddRun()
	run.SetText("One Last Paragraph in a different font")
	run.Properties().SetSize(20)
	run.Properties().SetFont("Courier")
	run.Properties().SetSolidFill(color.Red)

	if err != nil {
		log.Fatalf("error opening template: %s", err)
	}
	ppt.SaveToFile("mod.pptx")
	return nil
}

func (r *GenerateLaporan) TestingConvert(ctx http.Context) http.Response {
	data := struct {
		Title string
	}{
		Title: "Sample PDF from Template",
	}
	tmpl, err := template.ParseFiles("template/sample.html")
	if err != nil {
		log.Fatal("X", err)
	}

	fmt.Println("STEP-1")

	var htmlContent strings.Builder
	err = tmpl.Execute(&htmlContent, data)
	if err != nil {
		log.Fatal("Y", err)
	}
	fmt.Println("STEP-2")

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Fatal("A", err)
	}
	fmt.Println("STEP-3")

	// Add the generated HTML content to the PDF generator
	page := wkhtml.NewPageReader(strings.NewReader(htmlContent.String()))
	pdfg.AddPage(page)

	fmt.Println("STEP-4")
	// Generate the PDF
	err = pdfg.Create()
	if err != nil {
		log.Fatal("B", err)
	}
	fmt.Println("STEP-5")

	// Save the PDF to a file
	outputFileName := "output_template.pdf"
	err = pdfg.WriteFile(outputFileName)
	if err != nil {
		log.Fatal("C", err)
	}
	fmt.Println("STEP-6")

	return ctx.Response().Success().Json(http.Json{
		"message": "PDF generated successfully",
		"file":    outputFileName,
	})
}
