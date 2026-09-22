package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/udistrital/utils_oas/v2/request"
)

type FirmaDiplomaService struct {
	S3 S3Service
}

func (s FirmaDiplomaService) Firmar(ctx context.Context, documentoID int64, input *FirmaDiplomaInput) (*FirmaDiplomaResult, error) {
	if documentoID <= 0 {
		return nil, fmt.Errorf("%w: documento id is required", ErrInvalidInput)
	}
	if input == nil {
		return nil, fmt.Errorf("%w: request body is required", ErrInvalidInput)
	}
	if input.DocumentoIdentidadFirmante <= 0 {
		return nil, fmt.Errorf("%w: documento_identidad_firmante is required", ErrInvalidInput)
	}

	firmanteActivo, err := FirmanteService{}.ConsultarRolActivo(ctx, input.DocumentoIdentidadFirmante, time.Now())
	if err != nil {
		return nil, err
	}
	if err := validarFlujoFirmaDocumento(ctx, documentoID, firmanteActivo); err != nil {
		return nil, err
	}
	preview := firmanteActivo.Orden < 4
	if preview {
		firmaDocumento, err := registrarFirmaDocumentoDiploma(ctx, documentoID, firmanteActivo)
		if err != nil {
			return nil, err
		}
		previewDocumento, err := DiplomaPreviewService{}.Generar(ctx, documentoID)
		if err != nil {
			return nil, err
		}
		qrDemoURL := diplomaPreviewQRDemoURL()
		return &FirmaDiplomaResult{
			Preview:          true,
			FirmaID:          "QR-DEMO",
			DocumentoID:      documentoID,
			QRURLSegura:      qrDemoURL,
			Firmante:         *firmanteActivo,
			DocumentoDigital: previewDocumento.DocumentoDigital,
			FirmaDocumento:   firmaDocumento,
			File:             previewDocumento.File,
		}, nil
	}

	previewDocumento, err := generarDiplomaPreview(ctx, documentoID, firmanteActivo)
	if err != nil {
		return nil, err
	}
	metadatos := cloneMetadatos(input.Metadatos)
	metadatos["rol_firmante"] = firmanteActivo.Rol
	metadatos["orden_firmante"] = firmanteActivo.Orden
	metadatos["cargo_id_firmante"] = firmanteActivo.CargoID
	metadatos["documento_identidad_firmante"] = firmanteActivo.DocumentoIdentidad

	firmantes := []PersonaFirma{{
		Nombre:         firmanteActivo.Nombre,
		Cargo:          firmanteActivo.Cargo,
		TipoID:         "CC",
		Identificacion: strconv.FormatInt(firmanteActivo.DocumentoIdentidad, 10),
	}}

	nombre := input.Nombre
	if nombre == "" {
		nombre = fmt.Sprintf("Diploma documento %d", documentoID)
	}
	descripcion := input.Descripcion
	if descripcion == "" {
		descripcion = "Diploma digital final"
	}

	firmaPayload := []firmaElectronicaRequest{{
		RepositorioDocumental: "diplomas",
		DocumentoID:           documentoID,
		Nombre:                nombre,
		Descripcion:           descripcion,
		Metadatos:             metadatos,
		Firmantes:             firmantes,
		Representantes:        input.Representantes,
		File:                  previewDocumento.File,
	}}

	var firmaResponse firmaElectronicaResponse
	status, err := request.PostWithContext(ctx, firmaElectronicaMidURL()+"/v2/firma_electronica", firmaPayload, &firmaResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: firma_electronica_mid status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: firma_electronica_mid returned status %d", ErrExternalService, status)
	}
	if firmaResponse.File == "" {
		return nil, fmt.Errorf("%w: firma_electronica_mid response missing file", ErrExternalService)
	}
	if firmaResponse.UUIDDocumento == "" {
		return nil, fmt.Errorf("%w: firma_electronica_mid response missing uuid_documento", ErrExternalService)
	}

	s3Service := s.S3
	s3Ref, err := s3Service.PutDiploma(ctx, firmaResponse.UUIDDocumento, firmaResponse.File)
	if err != nil {
		return nil, err
	}

	var documentoResponse documentoDigitalResponse
	updatePayload := actualizarUUIDDocumentoRequest{UUIDDocumento: firmaResponse.UUIDDocumento}
	status, err = request.PutWithContext(
		ctx,
		fmt.Sprintf("%s/documento_digital/%d/uuid", diplomasCrudURL(), documentoID),
		updatePayload,
		&documentoResponse,
	)
	if err != nil {
		return nil, fmt.Errorf("diplomas_crud actualizar uuid status %d: %w", status, err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: diplomas_crud actualizar uuid returned status %d", ErrExternalService, status)
	}

	firmaDocumento, err := registrarFirmaDocumentoDiploma(ctx, documentoID, firmanteActivo)
	if err != nil {
		return nil, err
	}

	return &FirmaDiplomaResult{
		Preview:               false,
		FirmaID:               firmaResponse.FirmaID,
		RepositorioDocumental: firmaResponse.RepositorioDocumental,
		DocumentoID:           documentoID,
		UUIDDocumento:         firmaResponse.UUIDDocumento,
		HashSHA256:            firmaResponse.HashSHA256,
		CodigoAutenticidad:    firmaResponse.CodigoAutenticidad,
		QRURLSegura:           firmaResponse.QRURLSegura,
		Firmante:              *firmanteActivo,
		S3:                    &s3Ref,
		DynamoDB:              firmaResponse.DynamoDB,
		DocumentoDigital:      documentoResponse,
		FirmaDocumento:        firmaDocumento,
	}, nil
}

func cloneMetadatos(input map[string]interface{}) map[string]interface{} {
	metadatos := make(map[string]interface{}, len(input)+4)
	for key, value := range input {
		metadatos[key] = value
	}
	return metadatos
}

func registrarFirmaDocumentoDiploma(ctx context.Context, documentoID int64, firmante *FirmanteActivo) (*firmaDocumentoResponse, error) {
	estadoCodigo := "FIRMA_PAR"
	if firmante.Orden >= 4 {
		estadoCodigo = "FIRMADO"
	}

	estadoID, err := resolverParametroID(ctx, 0, estadoCodigo)
	if err != nil {
		return nil, err
	}

	firmaFirmante, err := buscarFirmaFirmanteActiva(ctx, firmante.DocumentoIdentidad)
	if err != nil {
		return nil, err
	}

	var firmaFirmanteID *int64
	if firmaFirmante != nil {
		firmaFirmanteID = &firmaFirmante.ID
	}
	payload := registrarFirmaDocumentoRequest{
		FirmaFirmanteID:    firmaFirmanteID,
		RolFirmanteID:      int64(firmante.CargoID),
		DocumentoIdentidad: firmante.DocumentoIdentidad,
		EstadoFirmadoID:    &estadoID,
		Observacion:        fmt.Sprintf("Firma %s registrada en flujo de diploma", firmante.Rol),
	}

	var firmaDocumento firmaDocumentoResponse
	status, err := request.PostWithContext(
		ctx,
		fmt.Sprintf("%s/documento_digital/%d/firma", diplomasCrudURL(), documentoID),
		payload,
		&firmaDocumento,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: diplomas_crud registrar firma_documento status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: diplomas_crud registrar firma_documento returned status %d", ErrExternalService, status)
	}

	return &firmaDocumento, nil
}

func validarFlujoFirmaDocumento(ctx context.Context, documentoID int64, firmante *FirmanteActivo) error {
	if firmante == nil {
		return fmt.Errorf("%w: firmante is required", ErrInvalidInput)
	}

	estadoFirmadoID, err := resolverParametroID(ctx, 0, "FIRMADO")
	if err != nil {
		return err
	}

	var documento documentoDigitalResponse
	status, err := request.GetWithContext(
		ctx,
		fmt.Sprintf("%s/documento_digital/%d", diplomasCrudURL(), documentoID),
		&documento,
	)
	if err != nil {
		return fmt.Errorf("%w: diplomas_crud documento_digital status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("%w: diplomas_crud documento_digital returned status %d", ErrExternalService, status)
	}
	if documento.EstadoDocumentoID == estadoFirmadoID {
		return fmt.Errorf("%w: el documento ya esta firmado", ErrInvalidState)
	}

	firmas, err := listarFirmasDocumento(ctx, documentoID)
	if err != nil {
		return err
	}
	ordenesFirmadas := make(map[int]bool, len(firmas))
	for _, firma := range firmas {
		orden, ok := ordenRolFirmante(firma.RolFirmanteID)
		if !ok {
			continue
		}
		if orden == firmante.Orden {
			return fmt.Errorf("%w: ya existe firma registrada para el orden %d", ErrInvalidState, orden)
		}
		ordenesFirmadas[orden] = true
	}

	ordenEsperado := 1
	for ordenesFirmadas[ordenEsperado] {
		ordenEsperado++
	}
	if firmante.Orden != ordenEsperado {
		return fmt.Errorf("%w: el siguiente firmante debe ser orden %d", ErrInvalidState, ordenEsperado)
	}

	return nil
}

func listarFirmasDocumento(ctx context.Context, documentoID int64) ([]firmaDocumentoResponse, error) {
	rawQuery := fmt.Sprintf("documento_digital_id:%d,activo:true", documentoID)
	endpoint := fmt.Sprintf("%s/firma_documento/?query=%s", diplomasCrudURL(), url.QueryEscape(rawQuery))

	var firmas []firmaDocumentoResponse
	status, err := request.GetWithContext(ctx, endpoint, &firmas)
	if err != nil {
		return nil, fmt.Errorf("%w: diplomas_crud firma_documento status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: diplomas_crud firma_documento returned status %d", ErrExternalService, status)
	}
	return firmas, nil
}

func ordenRolFirmante(rolFirmanteID int64) (int, bool) {
	rol, ok := rolesFirmaPermitidos[int(rolFirmanteID)]
	if !ok {
		return 0, false
	}
	return rol.Orden, true
}

type gestorDocumentalDocumentoResponse struct {
	Thumbnail   gestorDocumentalArchivo `json:"thumb:thumbnail"`
	FileContent gestorDocumentalArchivo `json:"file:content"`
	File        string                  `json:"file"`
}

type gestorDocumentalArchivo struct {
	Data     string `json:"data"`
	MimeType string `json:"mime-type"`
}

func consultarFirmaGraficaFirmante(ctx context.Context, documentoIdentidad int64) (string, error) {
	firmaFirmante, err := buscarFirmaFirmanteActiva(ctx, documentoIdentidad)
	if err != nil {
		return "", err
	}
	if firmaFirmante == nil || firmaFirmante.EnlaceFirma == "" {
		return "", fmt.Errorf("%w: el firmante no tiene firma registrada", ErrInvalidState)
	}

	baseURL := gestorDocumentalDocumentURL()
	if baseURL == "" {
		return "", fmt.Errorf("%w: GestorDocumentalDocumentURL is required", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(firmaFirmante.EnlaceFirma))

	var documento gestorDocumentalDocumentoResponse
	status, err := request.GetWithContext(ctx, endpoint, &documento)
	if err != nil {
		return "", fmt.Errorf("%w: gestor_documental document firma status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("%w: gestor_documental document firma returned status %d", ErrExternalService, status)
	}

	if documento.File != "" {
		return normalizarFirmaGraficaBase64(documento.File), nil
	}
	if documento.Thumbnail.Data != "" {
		return documento.Thumbnail.Data, nil
	}
	if documento.FileContent.Data != "" {
		return documento.FileContent.Data, nil
	}
	return "", fmt.Errorf("%w: gestor_documental document firma missing file url", ErrExternalService)
}

func normalizarFirmaGraficaBase64(file string) string {
	file = strings.TrimSpace(file)
	if file == "" || strings.HasPrefix(file, "data:") {
		return file
	}
	mimeType := "image/png"
	if strings.HasPrefix(file, "/9j/") {
		mimeType = "image/jpeg"
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, file)
}

func buscarFirmaFirmanteActiva(ctx context.Context, documentoIdentidad int64) (*FirmaFirmanteCRUDResponse, error) {
	rawQuery := fmt.Sprintf("documento_identidad:%d,activo:true", documentoIdentidad)
	endpoint := fmt.Sprintf("%s/firma_firmante/?query=%s", diplomasCrudURL(), url.QueryEscape(rawQuery))

	var firmas []FirmaFirmanteCRUDResponse
	status, err := request.GetWithContext(ctx, endpoint, &firmas)
	if err != nil {
		return nil, fmt.Errorf("%w: diplomas_crud firma_firmante status %d: %v", ErrExternalService, status, err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: diplomas_crud firma_firmante returned status %d", ErrExternalService, status)
	}
	if len(firmas) == 0 {
		return nil, nil
	}
	return &firmas[0], nil
}
