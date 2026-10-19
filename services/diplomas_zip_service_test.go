package services

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestNormalizarFechaZIP(t *testing.T) {
	if got := normalizarFechaZIP("2026-10-09T00:00:00Z"); got != "20261008" {
		t.Fatalf("normalizarFechaZIP() = %q", got)
	}
	if got := normalizarFechaZIP("2026-10-08"); got != "20261008" {
		t.Fatalf("normalizarFechaZIP fecha civil = %q", got)
	}
	if got := normalizarFechaZIP(""); got != "sin_fecha" {
		t.Fatalf("normalizarFechaZIP vacia = %q", got)
	}
}

func TestEscribirDiplomaZIPConservaBytesYOrganizaCarpetas(t *testing.T) {
	original := []byte("%PDF-1.7\ncontenido-original-s3\n%%EOF")
	var output bytes.Buffer
	zw := zip.NewWriter(&output)
	err := escribirDiplomaZIP(zw, diplomaZIPItem{
		FacultadID: 14, CeremoniaID: 4, Codigo: 20181016035, FechaCeremonia: "2026-10-09T00:00:00Z",
	}, bytes.NewReader(original))
	if err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
	if err != nil || len(zr.File) != 1 {
		t.Fatalf("ZIP invalido: archivos=%d error=%v", len(zr.File), err)
	}
	if got := zr.File[0].Name; got != "Facultad_14/Ceremonia_20261008_4/20261008_Diploma_20181016035.pdf" {
		t.Fatalf("ruta ZIP = %q", got)
	}
	archivo, err := zr.File[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer archivo.Close()
	contenido, err := io.ReadAll(archivo)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(contenido, original) {
		t.Fatal("el PDF original fue alterado al incluirlo en el ZIP")
	}
}
