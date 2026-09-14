package services

import (
	"os"
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

func pazYSalvosCrudURL() string {
	return configString("PazYSalvosCrudURL", "PAZ_Y_SALVOS_CRUD_URL")
}

func academicaCoreServiceURL() string {
	return configString("AcademicaCoreServiceURL", "ACADEMICA_CORE_SERVICE_URL")
}

func awsRegion() string {
	value, _ := beego.AppConfig.String("AWSRegion")
	value = strings.TrimSpace(value)
	if value == "" {
		return "us-east-1"
	}
	return value
}

func awsEndpointURL() string {
	value, _ := beego.AppConfig.String("AWSEndpointURL")
	return strings.TrimSpace(value)
}

func awsAccessKeyID() string {
	value, _ := beego.AppConfig.String("AWSAccessKeyID")
	return strings.TrimSpace(value)
}

func awsSecretAccessKey() string {
	value, _ := beego.AppConfig.String("AWSSecretAccessKey")
	return strings.TrimSpace(value)
}

func diplomasS3Bucket() string {
	value, _ := beego.AppConfig.String("DiplomasS3Bucket")
	return strings.TrimSpace(value)
}
