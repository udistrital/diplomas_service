package services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// validarDocumentoFirmado checks the MID response before any S3/CRUD writes.
func validarDocumentoFirmado(documentoID int64, response firmaElectronicaResponse) error {
	if response.DocumentoID != documentoID || response.RepositorioDocumental != "diplomas" || response.FirmaID == "" {
		return fmt.Errorf("%w: signed document identity mismatch", ErrExternalService)
	}
	pdf, err := base64.StdEncoding.DecodeString(response.File)
	if err != nil || len(pdf) == 0 {
		return fmt.Errorf("%w: invalid signed document encoding", ErrExternalService)
	}
	digest := sha256.Sum256(pdf)
	if hex.EncodeToString(digest[:]) != response.HashSHA256 {
		return fmt.Errorf("%w: signed document SHA-256 mismatch", ErrExternalService)
	}
	return nil
}
