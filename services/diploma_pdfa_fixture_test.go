package services

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// Export an entirely fictional diploma using the current production template.
// The Python integration suite stamps its QR and validates the definitive PDF.
func TestExportDiplomaPDFAFixture(t *testing.T) {
	output := os.Getenv("PDFA_TEST_INPUT")
	if output == "" {
		t.Skip("set PDFA_TEST_INPUT to export the institutional fixture")
	}
	im := image.NewNRGBA(image.Rect(0, 0, 300, 80))
	for x := 15; x < 285; x++ {
		y := 25 + (x/15)%3*10
		for dy := 0; dy < 3; dy++ {
			im.Set(x, y+dy, color.NRGBA{A: 255})
		}
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, im); err != nil {
		t.Fatal(err)
	}
	var signatures []diplomaSignatureRender
	for order := 1; order <= 4; order++ {
		signatures = append(signatures, diplomaSignatureRender{
			Orden: order, Imagen: base64.StdEncoding.EncodeToString(encoded.Bytes()),
		})
	}
	pdf, err := construirPDFDiplomaPreview(DiplomaPreviewMetadatos{
		Titulo: "INGENIERO DE SISTEMAS", NombreEstudiante: "PRUEBA EJEMPLO PERSONA FICTICIA",
		TipoDocumentoEstudiante: "CC", NumeroDocumento: "0000000000",
		LugarExpedicion: "BOGOTA D.C.", FechaCeremonia: "En la ciudad de Bogota a los 15 dias del mes de octubre de 2026",
		DiplomaID: 999999, NumeroRegistro: 999999, Folio: 1, Libro: 1, Facultad: "FACULTAD DE INGENIERIA",
	}, signatures, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, pdf, 0600); err != nil {
		t.Fatal(err)
	}
}
