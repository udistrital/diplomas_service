package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestValidarDocumentoFirmado(t *testing.T) {
	pdf := []byte("%PDF-1.7\nfinal document")
	valid := firmaElectronicaResponse{
		DocumentoID: 7, RepositorioDocumental: "diplomas", FirmaID: "firma",
		File: base64.StdEncoding.EncodeToString(pdf), HashSHA256: fmt.Sprintf("%x", sha256.Sum256(pdf)),
	}
	if err := validarDocumentoFirmado(7, valid); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"documento", "repositorio", "firma", "base64", "hash", "bytes"} {
		t.Run(field, func(t *testing.T) {
			response := valid
			switch field {
			case "documento":
				response.DocumentoID++
			case "repositorio":
				response.RepositorioDocumental = "nuxeo"
			case "firma":
				response.FirmaID = ""
			case "base64":
				response.File = "!"
			case "hash":
				response.HashSHA256 = ""
			case "bytes":
				response.File = base64.StdEncoding.EncodeToString(append(pdf, 'x'))
			}
			if err := validarDocumentoFirmado(7, response); err == nil {
				t.Fatal("expected rejection before storage")
			}
		})
	}
}
