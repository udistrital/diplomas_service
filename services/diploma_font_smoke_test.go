package services

import (
	"bytes"
	"testing"

	"github.com/phpdave11/gofpdf"
)

func TestDiplomaFontsSmoke(t *testing.T) {
	t.Skip("gofpdf smoke test is kept only for TTF debugging; diploma rendering uses Chrome for OTF fonts")
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: 400, Ht: 280},
	})
	pdf.AddUTF8Font("Cambria", "", diplomaCambriaFontPath())
	pdf.AddUTF8Font("Cambria", "B", diplomaCambriaBoldFontPath())
	pdf.AddUTF8Font("Engravers", "B", diplomaEngraversFontPath())
	pdf.AddPage()
	pdf.SetFont("Cambria", "", 12)
	pdf.Text(10, 10, "Cambria")
	pdf.SetFont("Cambria", "B", 12)
	pdf.Text(10, 20, "Cambria Bold")
	pdf.SetFont("Engravers", "B", 12)
	pdf.Text(10, 30, "Engravers")
	var out bytes.Buffer
	if err := pdf.Output(&out); err != nil {
		t.Fatal(err)
	}
}
