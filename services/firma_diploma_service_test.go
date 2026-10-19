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

	if got := nombreArchivoDiploma(preview); got != "20291125_Diploma_20242781013.pdf" {
		t.Fatalf("nombreArchivoDiploma() = %q", got)
	}
}

func TestNombreArchivoDiplomaUsaMismaFechaBogotaDelDocumento(t *testing.T) {
	fechaUTC := time.Date(2026, 10, 9, 2, 0, 0, 0, time.UTC)
	preview := &DiplomaPreviewResult{
		DocumentoDigital: documentoDigitalResponse{CodigoEstudiante: 20181016035},
		DiplomaDigital:   &diplomaDigitalResponse{FechaGrado: fechaUTC},
	}

	if fechaTexto := fechaCeremoniaBogota(fechaUTC); fechaTexto != "En la ciudad de Bogota a los 8 dias del mes de octubre de 2026" {
		t.Fatalf("fechaCeremoniaBogota() = %q", fechaTexto)
	}
	if got := nombreArchivoDiploma(preview); got != "20261008_Diploma_20181016035.pdf" {
		t.Fatalf("nombreArchivoDiploma() = %q", got)
	}
}
