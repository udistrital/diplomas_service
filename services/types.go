package services

import "time"

type PersonaFirma struct {
	Nombre         string   `json:"nombre"`
	Cargo          string   `json:"cargo,omitempty"`
	Oficina        string   `json:"oficina,omitempty"`
	TipoID         string   `json:"tipoId,omitempty"`
	Identificacion string   `json:"identificacion,omitempty"`
	FirmaURL       string   `json:"firma_url,omitempty"`
	OrdenCampos    []string `json:"orden_campos,omitempty"`
}

type FirmaDiplomaInput struct {
	Nombre                     string                 `json:"nombre"`
	Descripcion                string                 `json:"descripcion"`
	DocumentoIdentidadFirmante int64                  `json:"documento_identidad_firmante"`
	Metadatos                  map[string]interface{} `json:"metadatos"`
	Firmantes                  []PersonaFirma         `json:"firmantes,omitempty"`
	Representantes             []PersonaFirma         `json:"representantes"`
	File                       string                 `json:"file"`
}

type firmaElectronicaRequest struct {
	RepositorioDocumental string                 `json:"repositorio_documental"`
	DocumentoID           int64                  `json:"documento_id"`
	Nombre                string                 `json:"nombre"`
	Descripcion           string                 `json:"descripcion"`
	Metadatos             map[string]interface{} `json:"metadatos"`
	Firmantes             []PersonaFirma         `json:"firmantes"`
	Representantes        []PersonaFirma         `json:"representantes"`
	File                  string                 `json:"file"`
}

type firmaElectronicaResponse struct {
	Status                string                 `json:"Status"`
	FirmaID               string                 `json:"firma_id"`
	RepositorioDocumental string                 `json:"repositorio_documental"`
	DocumentoID           int64                  `json:"documento_id"`
	UUIDDocumento         string                 `json:"uuid_documento"`
	HashSHA256            string                 `json:"hash_sha256"`
	CodigoAutenticidad    string                 `json:"codigo_autenticidad"`
	FirmaEncriptada       string                 `json:"firma_encriptada"`
	Llaves                map[string]interface{} `json:"llaves"`
	Firmantes             map[string]interface{} `json:"firmantes"`
	QRURLSegura           string                 `json:"qr_url_segura"`
	DynamoDB              map[string]interface{} `json:"dynamodb"`
	File                  string                 `json:"file"`
}

type actualizarUUIDDocumentoRequest struct {
	UUIDDocumento string `json:"uuid_documento"`
}

type crearDocumentoDigitalRequest struct {
	TipoDocumentoDigitalID    int64  `json:"tipo_documento_digital_id"`
	EstadoDocumentoID         int64  `json:"estado_documento_id"`
	CodigoEstudiante          int64  `json:"codigo_estudiante"`
	FacultadID                int64  `json:"facultad_id"`
	ProgramaAcademicoID       int64  `json:"programa_academico_id"`
	PeriodoID                 int64  `json:"periodo_id"`
	Vigencia                  int    `json:"vigencia"`
	NombreEstudiante          string `json:"nombre_estudiante,omitempty"`
	TipoDocumentoEstudiante   string `json:"tipo_documento_estudiante,omitempty"`
	NumeroDocumentoEstudiante string `json:"numero_documento_estudiante,omitempty"`
	MunicipioExpedicion       string `json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string `json:"titulo_otorgado,omitempty"`
	Activo                    bool   `json:"activo"`
}

type documentoDigitalResponse struct {
	ID                        int64   `json:"id"`
	TipoDocumentoDigitalID    int64   `json:"tipo_documento_digital_id"`
	EstadoDocumentoID         int64   `json:"estado_documento_id"`
	CodigoEstudiante          int64   `json:"codigo_estudiante"`
	FacultadID                *int64  `json:"facultad_id,omitempty"`
	ProgramaAcademicoID       *int64  `json:"programa_academico_id,omitempty"`
	PeriodoID                 *int64  `json:"periodo_id,omitempty"`
	Vigencia                  *int    `json:"vigencia,omitempty"`
	NombreEstudiante          string  `json:"nombre_estudiante,omitempty"`
	TipoDocumentoEstudiante   string  `json:"tipo_documento_estudiante,omitempty"`
	NumeroDocumentoEstudiante string  `json:"numero_documento_estudiante,omitempty"`
	MunicipioExpedicion       string  `json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string  `json:"titulo_otorgado,omitempty"`
	UUIDDocumento             *string `json:"uuid_documento,omitempty"`
	Activo                    bool    `json:"activo"`
}

type crearDiplomaDigitalRequest struct {
	FacultadID                int64     `json:"facultad_id"`
	Vigencia                  int       `json:"vigencia"`
	FechaGrado                time.Time `json:"fecha_grado"`
	NombreEstudiante          string    `json:"nombre_estudiante"`
	TipoDocumentoEstudiante   string    `json:"tipo_documento_estudiante"`
	NumeroDocumentoEstudiante string    `json:"numero_documento_estudiante"`
	MunicipioExpedicion       string    `json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string    `json:"titulo_otorgado"`
	EstadoDocumentoCreadoID   int64     `json:"estado_documento_creado_id"`
}

type diplomaDigitalResponse struct {
	ID                        int64                     `json:"id"`
	DocumentoDigital          *documentoDigitalResponse `json:"documento_digital,omitempty"`
	FacultadID                int64                     `json:"facultad_id"`
	Vigencia                  int                       `json:"vigencia"`
	FechaGrado                time.Time                 `json:"fecha_grado"`
	NombreEstudiante          string                    `json:"nombre_estudiante"`
	TipoDocumentoEstudiante   string                    `json:"tipo_documento_estudiante"`
	NumeroDocumentoEstudiante string                    `json:"numero_documento_estudiante"`
	MunicipioExpedicion       string                    `json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string                    `json:"titulo_otorgado"`
	ConsecutivoDiploma        int64                     `json:"consecutivo_diploma"`
	ConsecutivoFacultad       int                       `json:"consecutivo_facultad"`
	Folio                     int                       `json:"folio"`
	Acta                      int                       `json:"acta"`
	Libro                     int                       `json:"libro"`
	Activo                    bool                      `json:"activo"`
}

type S3ObjectRef struct {
	Bucket    string `json:"bucket"`
	Key       string `json:"key"`
	VersionID string `json:"version_id,omitempty"`
	ETag      string `json:"etag,omitempty"`
	URI       string `json:"uri"`
}

type FirmanteActivo struct {
	Rol                string `json:"rol"`
	Orden              int    `json:"orden"`
	CargoID            int    `json:"cargo_id"`
	Cargo              string `json:"cargo"`
	Nombre             string `json:"nombre"`
	DocumentoIdentidad int64  `json:"documento_identidad"`
}

type SubirFirmaFirmanteInput struct {
	DocumentoIdentidadFirmante int64                  `json:"documento_identidad_firmante"`
	IdTipoDocumento            int64                  `json:"IdTipoDocumento,omitempty"`
	Nombre                     string                 `json:"nombre,omitempty"`
	Descripcion                string                 `json:"descripcion,omitempty"`
	Metadatos                  map[string]interface{} `json:"metadatos,omitempty"`
	File                       string                 `json:"file"`
}

type subirDocumentoFirmaRequest struct {
	IdTipoDocumento int64                  `json:"IdTipoDocumento"`
	Nombre          string                 `json:"nombre"`
	Metadatos       map[string]interface{} `json:"metadatos"`
	Descripcion     string                 `json:"descripcion"`
	File            string                 `json:"file"`
}

type registrarFirmaFirmanteRequest struct {
	DocumentoIdentidad int64  `json:"documento_identidad"`
	EnlaceFirma        string `json:"enlace_firma"`
	Activo             bool   `json:"activo"`
}

type FirmaFirmanteCRUDResponse struct {
	ID                 int64  `json:"id"`
	DocumentoIdentidad int64  `json:"documento_identidad"`
	EnlaceFirma        string `json:"enlace_firma"`
	Activo             bool   `json:"activo"`
}

type SubirFirmaFirmanteResult struct {
	Firmante      FirmanteActivo            `json:"firmante"`
	EnlaceFirma   string                    `json:"enlace_firma"`
	Documento     map[string]interface{}    `json:"documento"`
	FirmaFirmante FirmaFirmanteCRUDResponse `json:"firma_firmante"`
}

type FirmaDiplomaResult struct {
	Preview               bool                     `json:"preview"`
	FirmaID               string                   `json:"firma_id"`
	RepositorioDocumental string                   `json:"repositorio_documental"`
	DocumentoID           int64                    `json:"documento_id"`
	UUIDDocumento         string                   `json:"uuid_documento"`
	HashSHA256            string                   `json:"hash_sha256"`
	CodigoAutenticidad    string                   `json:"codigo_autenticidad"`
	QRURLSegura           string                   `json:"qr_url_segura"`
	Firmante              FirmanteActivo           `json:"firmante"`
	S3                    *S3ObjectRef             `json:"s3,omitempty"`
	DynamoDB              map[string]interface{}   `json:"dynamodb"`
	DocumentoDigital      documentoDigitalResponse `json:"documento_digital"`
	FirmaDocumento        *firmaDocumentoResponse  `json:"firma_documento,omitempty"`
	File                  string                   `json:"file,omitempty"`
}

type registrarFirmaDocumentoRequest struct {
	FirmaFirmanteID    *int64 `json:"firma_firmante_id,omitempty"`
	RolFirmanteID      int64  `json:"rol_firmante_id"`
	DocumentoIdentidad int64  `json:"documento_identidad"`
	EstadoFirmadoID    *int64 `json:"estado_firmado_id,omitempty"`
	Observacion        string `json:"observacion,omitempty"`
}

type firmaDocumentoResponse struct {
	ID                 int64                      `json:"id"`
	DocumentoDigital   *documentoDigitalResponse  `json:"documento_digital,omitempty"`
	FirmaFirmante      *FirmaFirmanteCRUDResponse `json:"firma_firmante,omitempty"`
	RolFirmanteID      int64                      `json:"rol_firmante_id"`
	DocumentoIdentidad int64                      `json:"documento_identidad"`
	Activo             bool                       `json:"activo"`
}
