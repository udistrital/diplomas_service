package services

import (
	"testing"
	"time"
)

func TestNombreArchivoDiploma(t *testing.T) {
	preview := &DiplomaPreviewResult{
		DocumentoDigital: documentoDigitalResponse{CodigoEstudiante: 20242781013},
		DiplomaDigital:   &diplomaDigitalResponse{FechaGrado: time.Date(2029, 11, 26, 0, 0, 0, 0, time.UTC)},
	}

	if got := nombreArchivoDiploma(preview); got != "20291126_Diploma_20242781013.pdf" {
		t.Fatalf("nombreArchivoDiploma() = %q", got)
	}
}
