package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type secretariaFacultadItem struct {
	Documento   string `json:"EMP_NRO_IDEN"`
	Dependencia string `json:"SEC_DEP_COD"`
	Facultad    string `json:"DEP_NOMBRE"`
	Nombre      string `json:"SEC_NOMBRE"`
	Apellido    string `json:"SEC_APELLIDO"`
}

type facultadSecretariaResponse struct {
	Facultades struct {
		Secretaria []secretariaFacultadItem `json:"secretaria"`
	} `json:"facultades"`
}

type facultadOikosGedepResponse struct {
	Homologacion struct {
		IDOikos string `json:"id_oikos"`
		IDGedep string `json:"id_gedep"`
	} `json:"homologacion"`
}

type secretariaFacultadAutorizada struct {
	Documento   int64
	Nombre      string
	Dependencia string
	Facultad    string
	FacultadID  int64
}

func consultarSecretariaFacultadOikos(ctx context.Context, documento int64) (*secretariaFacultadAutorizada, error) {
	if documento <= 0 {
		return nil, fmt.Errorf("%w: documento_identidad_secretario is required", ErrInvalidInput)
	}
	secretaria, err := consultarFacultadSecretaria(ctx, documento)
	if err != nil {
		return nil, err
	}
	facultadID, err := consultarFacultadOikosGedep(ctx, secretaria.Dependencia)
	if err != nil {
		return nil, err
	}
	secretaria.FacultadID = facultadID
	return secretaria, nil
}

func consultarFacultadSecretaria(ctx context.Context, documento int64) (*secretariaFacultadAutorizada, error) {
	baseURL := academicaPruebasAPIURL()
	if baseURL == "" {
		return nil, fmt.Errorf("%w: ACADEMICA_PRUEBAS_API_URL no configurado", ErrInvalidState)
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/facultad_secretaria/" + url.PathEscape(strconv.FormatInt(documento, 10))
	var response facultadSecretariaResponse
	if err := getJSON(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("%w: facultad_secretaria: %v", ErrExternalService, err)
	}
	if len(response.Facultades.Secretaria) == 0 {
		return nil, fmt.Errorf("%w: el documento no pertenece a una secretaria de facultad", ErrNoActiveRole)
	}
	item := response.Facultades.Secretaria[0]
	dependencia := strings.TrimSpace(item.Dependencia)
	if dependencia == "" {
		return nil, fmt.Errorf("%w: facultad_secretaria sin SEC_DEP_COD", ErrInvalidState)
	}
	return &secretariaFacultadAutorizada{
		Documento:   documento,
		Nombre:      strings.TrimSpace(strings.TrimSpace(item.Nombre) + " " + strings.TrimSpace(item.Apellido)),
		Dependencia: dependencia,
		Facultad:    strings.TrimSpace(item.Facultad),
	}, nil
}

func consultarFacultadOikosGedep(ctx context.Context, depGedep string) (int64, error) {
	baseURL := homologacionDependenciasAPIURL()
	if baseURL == "" {
		return 0, fmt.Errorf("%w: HOMOLOGACION_DEPENDENCIAS_API_URL no configurado", ErrInvalidState)
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/facultad_oikos_gedep/" + url.PathEscape(strings.TrimSpace(depGedep))
	var response facultadOikosGedepResponse
	if err := getJSON(ctx, endpoint, &response); err != nil {
		return 0, fmt.Errorf("%w: facultad_oikos_gedep: %v", ErrExternalService, err)
	}
	id, err := strconv.ParseInt(strings.TrimSpace(response.Homologacion.IDOikos), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: homologacion sin id_oikos valido", ErrInvalidState)
	}
	return id, nil
}

func getJSON(ctx context.Context, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}

func firmanteDesdeSecretaria(secretaria *secretariaFacultadAutorizada) (*FirmanteActivo, error) {
	if secretaria == nil || secretaria.FacultadID <= 0 {
		return nil, errors.New("secretaria de facultad invalida")
	}
	return &FirmanteActivo{
		Rol:                "secretario_academico",
		Orden:              1,
		Cargo:              "SECRETARIO ACADEMICO",
		Nombre:             secretaria.Nombre,
		DocumentoIdentidad: secretaria.Documento,
		Dependencia:        secretaria.Dependencia,
		FacultadID:         secretaria.FacultadID,
	}, nil
}
