package services

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestFormatoTituloDiploma(t *testing.T) {
	cases := map[string]string{
		"INGENIERO DE SISTEMAS": "Ingeniero de Sistemas",
		"ESPECIALISTA EN DISEÑO DE VIAS URBANAS TRANSITO Y TRANSPORTE": "Especialista en Diseño de Vias Urbanas Transito y Transporte",
	}
	for input, expected := range cases {
		if got := formatoTituloDiploma(input); got != expected {
			t.Fatalf("formatoTituloDiploma(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestFormatoNombreDiploma(t *testing.T) {
	cases := map[string]string{
		"APELLIDO NOMBRE":                             "Nombre Apellido",
		"APELLIDO1 APELLIDO2 NOMBRE":                  "Nombre Apellido1 Apellido2",
		"APELLIDO1 APELLIDO2 NOMBRE1 NOMBRE2":         "Nombre1 Nombre2 Apellido1 Apellido2",
		"APELLIDO1 APELLIDO2 NOMBRE1 NOMBRE2 NOMBRE3": "Nombre1 Nombre2 Nombre3 Apellido1 Apellido2",
		"APELLIDO1 APELLIDO2, NOMBRE1 NOMBRE2":        "Nombre1 Nombre2 Apellido1 Apellido2",
	}
	for input, expected := range cases {
		if got := formatoNombreDiploma(input); got != expected {
			t.Fatalf("formatoNombreDiploma(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestConstruirHTMLDiplomaPreviewUsaTemplate(t *testing.T) {
	html, err := construirHTMLDiplomaPreview(DiplomaPreviewMetadatos{
		Titulo:                  "INGENIERO DE SISTEMAS",
		NombreEstudiante:        "GUTIERREZ MEJIA JACKY ESTEFANNY",
		TipoDocumentoEstudiante: "CC",
		NumeroDocumento:         "1070925421",
		LugarExpedicion:         "BOGOTA D.C.",
		FechaCeremonia:          "En la ciudad de Bogota a los 16 dias del mes de diciembre de 2026",
		DiplomaID:               1,
		NumeroRegistro:          8723,
		Folio:                   10,
		Libro:                   1,
		Facultad:                "FACULTAD DE MEDIO AMBIENTE",
	}, []diplomaSignatureRender{{
		Orden:  4,
		Imagen: "iVBORw0KGgo=",
	}}, true)
	if err != nil {
		t.Fatalf("construirHTMLDiplomaPreview() error = %v", err)
	}

	for _, expected := range []string{
		`assets/templates/diploma.css`,
		`class="qr-demo"`,
		`data:image/png;base64,`,
		`Jacky Estefanny Gutierrez Mejia`,
		`Ingeniero de Sistemas`,
		`Registro No. FMA 8723 Folio No. 010 Libro No. 1`,
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("html no contiene %q", expected)
		}
	}
	if strings.Contains(html, "ZgotmplZ") {
		t.Fatalf("html contiene URL bloqueada por html/template")
	}
	if strings.Contains(html, "FACULTAD DE MEDIO AMBIENTE") {
		t.Fatalf("html contiene la facultad en renglon independiente")
	}
}

func TestSiglaFacultadRegistro(t *testing.T) {
	cases := map[string]string{
		"FACULTAD DE INGENIERIA":      "FI",
		"Facultad de Medio Ambiente":  "FMA",
		"FACULTAD 65":                 "F65",
		"Facultad de Artes ASAB":      "FAA",
		"Facultad del Medio Ambiente": "FMA",
	}
	for input, expected := range cases {
		if got := siglaFacultadRegistro(input); got != expected {
			t.Fatalf("siglaFacultadRegistro(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestFacultadTextoUsaNombreInstitucional(t *testing.T) {
	id := int64(14)
	dependenciaOikosCache.Lock()
	dependenciaOikosCache.items[id] = dependenciaOikosCacheEntry{nombre: "FACULTAD DE INGENIERIA", expira: time.Now().Add(time.Hour)}
	dependenciaOikosCache.Unlock()
	t.Cleanup(func() {
		dependenciaOikosCache.Lock()
		delete(dependenciaOikosCache.items, id)
		dependenciaOikosCache.Unlock()
	})

	if got := siglaFacultadRegistro(facultadTexto(context.Background(), &id)); got != "FI" {
		t.Fatalf("sigla facultad 14 = %q, want FI", got)
	}
}
