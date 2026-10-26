package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func ceremoniaInputValido() *CrearCeremoniaInput {
	return &CrearCeremoniaInput{
		FechaCeremonia:                "2026-11-20",
		Lugar:                         "Auditorio Mayor",
		CantidadEstudiantesProgramada: 65,
		Observacion:                   "Ceremonia ordinaria",
		DocumentoIdentidadSecretario:  1023026828,
	}
}

func secretariaFacultadJSON() string {
	return `{"facultades":{"secretaria":[{"EMP_NRO_IDEN":"1023026828","SEC_DEP_COD":"33","DEP_NOMBRE":"FACULTAD DE INGENIERIA","SEC_NOMBRE":"ORLANDO","SEC_APELLIDO":"RIOS LEON"}]}}`
}

func homologacionGedepJSON() string {
	return `{"homologacion":{"id_oikos":"14","id_gedep":"33"}}`
}

func TestCrearCeremoniaExitosa(t *testing.T) {
	var ceremoniaPosts, relacionPosts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_secretaria/1023026828":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(secretariaFacultadJSON()))
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_oikos_gedep/33":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(homologacionGedepJSON()))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/ceremonia_grado/":
			ceremoniaPosts++
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":99,"lugar":"Auditorio Mayor","cantidad_estudiantes_programada":65,"activo":true}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/ceremonia_facultad_programa/":
			relacionPosts++
			var payload ceremoniaFacultadProgramaCRUDRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("invalid relation payload: %v", err)
			}
			if payload.FacultadID != 14 || payload.ProgramaAcademicoID != 0 || payload.CantidadEstudiantesProgramada == nil || *payload.CantidadEstudiantesProgramada != 65 {
				t.Fatalf("unexpected relation payload: %+v", payload)
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": int64(100), "facultad_id": int64(14), "programa_academico_id": int64(0), "activo": true})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("DIPLOMAS_CRUD_URL", server.URL)
	t.Setenv("ACADEMICA_PRUEBAS_API_URL", server.URL)
	t.Setenv("HOMOLOGACION_DEPENDENCIAS_API_URL", server.URL)

	result, err := CeremoniaService{}.Crear(t.Context(), ceremoniaInputValido())
	if err != nil {
		t.Fatalf("Crear returned error: %v", err)
	}
	if result.Ceremonia.ID != 99 || result.Facultad.FacultadID != 14 || result.SecretarioAcademico.Rol != "secretario_academico" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if ceremoniaPosts != 1 || relacionPosts != 1 {
		t.Fatalf("unexpected posts ceremonia=%d relaciones=%d", ceremoniaPosts, relacionPosts)
	}
}

func TestCrearCeremoniaPayloadInvalido(t *testing.T) {
	_, err := CeremoniaService{}.Crear(t.Context(), nil)
	if err == nil || !strings.Contains(err.Error(), "request body") {
		t.Fatalf("expected invalid payload error, got %v", err)
	}
}

func TestCrearCeremoniaSinDocumentoSecretario(t *testing.T) {
	input := ceremoniaInputValido()
	input.DocumentoIdentidadSecretario = 0
	_, err := CeremoniaService{}.Crear(t.Context(), input)
	if err == nil || !strings.Contains(err.Error(), "documento_identidad_secretario") {
		t.Fatalf("expected documento error, got %v", err)
	}
}

func TestCrearCeremoniaRechazaNoSecretario(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/facultad_secretaria/1023026828" {
			_, _ = w.Write([]byte(`{"facultades":{"secretaria":[]}}`))
			return
		}
		t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()
	t.Setenv("ACADEMICA_PRUEBAS_API_URL", server.URL)
	t.Setenv("HOMOLOGACION_DEPENDENCIAS_API_URL", server.URL)

	_, err := CeremoniaService{}.Crear(t.Context(), ceremoniaInputValido())
	if err == nil || !strings.Contains(err.Error(), "secretaria de facultad") {
		t.Fatalf("expected secretary error, got %v", err)
	}
}

func TestCrearCeremoniaFalloCrearCeremoniaEnCRUD(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_secretaria/1023026828":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(secretariaFacultadJSON()))
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_oikos_gedep/33":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(homologacionGedepJSON()))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/ceremonia_grado/":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("DIPLOMAS_CRUD_URL", server.URL)
	t.Setenv("ACADEMICA_PRUEBAS_API_URL", server.URL)
	t.Setenv("HOMOLOGACION_DEPENDENCIAS_API_URL", server.URL)

	_, err := CeremoniaService{}.Crear(t.Context(), ceremoniaInputValido())
	if err == nil || !strings.Contains(err.Error(), "ceremonia_grado") {
		t.Fatalf("expected CRUD error, got %v", err)
	}
}

func TestCrearCeremoniaFalloRelacionEjecutaRollback(t *testing.T) {
	var deleteCeremonia bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_secretaria/1023026828":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(secretariaFacultadJSON()))
		case r.Method == http.MethodGet && r.URL.Path == "/facultad_oikos_gedep/33":
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept application/json")
			}
			_, _ = w.Write([]byte(homologacionGedepJSON()))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/ceremonia_grado/":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":99,"activo":true}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/ceremonia_facultad_programa/":
			w.WriteHeader(http.StatusInternalServerError)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/ceremonia_grado/99":
			deleteCeremonia = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("DIPLOMAS_CRUD_URL", server.URL)
	t.Setenv("ACADEMICA_PRUEBAS_API_URL", server.URL)
	t.Setenv("HOMOLOGACION_DEPENDENCIAS_API_URL", server.URL)

	_, err := CeremoniaService{}.Crear(t.Context(), ceremoniaInputValido())
	if err == nil || !strings.Contains(err.Error(), "ceremonia_facultad_programa") {
		t.Fatalf("expected relation error, got %v", err)
	}
	if !deleteCeremonia {
		t.Fatal("expected rollback delete ceremonia")
	}
}
