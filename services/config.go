package services

import (
	"os"
	"path/filepath"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

func configString(confKey, envKey string) string {
	value, _ := beego.AppConfig.String(confKey)
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "${") {
		value = strings.TrimSpace(os.Getenv(envKey))
	}
	return strings.TrimRight(value, "/")
}

func firmaElectronicaMidURL() string {
	return configString("FirmaElectronicaMidURL", "FIRMA_ELECTRONICA_MID_URL")
}

func diplomasCrudURL() string {
	return configString("DiplomasCrudURL", "DIPLOMAS_CRUD_URL")
}

func administrativaAmazonAPIURL() string {
	return configString("AdministrativaAmazonAPIURL", "ADMINISTRATIVA_AMAZON_API_URL")
}

func documentosCrudGestorDocumentalURL() string {
	return configString("DocumentosCrudGestorDocumentalURL", "DOCUMENTOS_CRUD_GESTOR_DOCUMENTAL_URL")
}

func gestorDocumentalDocumentURL() string {
	return configString("GestorDocumentalDocumentURL", "GESTOR_DOCUMENTAL_DOCUMENT_URL")
}

func pazYSalvosCrudURL() string {
	return configString("PazYSalvosCrudURL", "PAZ_Y_SALVOS_CRUD_URL")
}

func pazYSalvosAprobadosQuery() string {
	value := configString("PazYSalvosAprobadosQuery", "PAZ_Y_SALVOS_APROBADOS_QUERY")
	if value == "" {
		return "Orc:true"
	}
	return value
}

func academicaCoreServiceURL() string {
	return configString("AcademicaCoreServiceURL", "ACADEMICA_CORE_SERVICE_URL")
}

func parametrosCrudURL() string {
	return configString("ParametrosCrudURL", "PARAMETROS_CRUD_URL")
}

func diplomaPreviewQRDemoURL() string {
	return configString("DiplomaPreviewQRDemoURL", "DIPLOMA_PREVIEW_QR_DEMO_URL")
}

func diplomaInstitutionalLogoPath() string {
	value := configString("DiplomaInstitutionalLogoPath", "DIPLOMA_INSTITUTIONAL_LOGO_PATH")
	if value == "" {
		return resolveLocalPath("assets/images/logo-udistrital-escudo.svg")
	}
	return resolveLocalPath(value)
}

func diplomaCambriaFontPath() string {
	value := configString("DiplomaCambriaFontPath", "DIPLOMA_CAMBRIA_FONT_PATH")
	if value == "" {
		return resolveLocalPath("assets/fonts/Cambria.ttf")
	}
	return resolveLocalPath(value)
}

func diplomaCambriaBoldFontPath() string {
	value := configString("DiplomaCambriaBoldFontPath", "DIPLOMA_CAMBRIA_BOLD_FONT_PATH")
	if value == "" {
		return resolveLocalPath("assets/fonts/Cambriab.ttf")
	}
	return resolveLocalPath(value)
}

func diplomaEngraversFontPath() string {
	value := configString("DiplomaEngraversFontPath", "DIPLOMA_ENGRAVERS_FONT_PATH")
	if value == "" {
		return resolveLocalPath("assets/fonts/EngraversOldEnglishBold.otf")
	}
	return resolveLocalPath(value)
}

func resolveLocalPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	candidates := []string{
		path,
		filepath.Join("..", path),
		filepath.Join("..", "..", path),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return path
}

func awsRegion() string {
	value := configString("AWSRegion", "AWS_REGION")
	if value == "" {
		return "us-east-1"
	}
	return value
}

func awsEndpointURL() string {
	return configString("AWSEndpointURL", "AWS_ENDPOINT_URL")
}

func awsAccessKeyID() string {
	return configString("AWSAccessKeyID", "AWS_ACCESS_KEY_ID")
}

func awsSecretAccessKey() string {
	return configString("AWSSecretAccessKey", "AWS_SECRET_ACCESS_KEY")
}

func diplomasS3Bucket() string {
	return configString("DiplomasS3Bucket", "DIPLOMAS_S3_BUCKET")
}
