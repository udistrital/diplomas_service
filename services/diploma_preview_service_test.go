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

func TestFirmaPreviewConservaDosRolesHistoricosDeLaMismaPersona(t *testing.T) {
	const documento = int64(123456789)
	decano := construirFirmaPreviewHistorica(2, documento, "Persona Compartida", "rubrica")
	rector := construirFirmaPreviewHistorica(4, documento, "Persona Compartida", "rubrica")

	if decano.Orden != 2 || decano.Cargo != "DECANO" || rector.Orden != 4 || rector.Cargo != "RECTOR" {
		t.Fatalf("se perdieron posiciones historicas: decano=%+v rector=%+v", decano, rector)
	}
	if decano.Imagen != rector.Imagen || decano.Imagen == "" {
		t.Fatalf("la misma rubrica debe poder renderizarse en ambos cargos")
	}

	html, err := construirHTMLDiplomaPreview(DiplomaPreviewMetadatos{}, []diplomaSignatureRender{decano, rector}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "DECANO") || !strings.Contains(html, "RECTOR") || strings.Count(html, "data:image/png;base64,rubrica") != 2 {
		t.Fatalf("el preview que alimenta el PDF definitivo no renderizo ambos cargos: %s", html)
	}
}

func TestFirmaPreviewConservaOrdenDeCuatroFirmantes(t *testing.T) {
	firmas := []diplomaSignatureRender{
		construirFirmaPreviewHistorica(1, 1, "Secretario", "firma-1"),
		construirFirmaPreviewHistorica(2, 2, "Decano", "firma-2"),
		construirFirmaPreviewHistorica(3, 3, "Secretario General", "firma-3"),
		construirFirmaPreviewHistorica(4, 4, "Rector", "firma-4"),
	}
	for index, firma := range firmas {
		if firma.Orden != index+1 || firma.Cargo != cargoPreviewPorOrden(index+1) {
			t.Fatalf("firma %d fuera de su posicion historica: %+v", index, firma)
		}
	}
}
