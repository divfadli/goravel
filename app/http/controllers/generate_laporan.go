package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/unidoc/unioffice/color"

	"fmt"
	"log"

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