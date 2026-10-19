package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestS3FinalBytesChecksumVersionAndNoOverwrite(t *testing.T) {
	finalPDF := []byte("%PDF-1.7\nfinal bytes after validation")
	if input := os.Getenv("PDFA_FINAL_INPUT"); input != "" {
		var err error
		finalPDF, err = os.ReadFile(input)
		if err != nil {
			t.Fatal(err)
		}
	}
	digest := sha256.Sum256(finalPDF)
	stored := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/diplomas-test/diplomas/new-uuid/Diploma.pdf" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != string(finalPDF) || r.Header.Get("Content-Type") != "application/pdf" {
			t.Error("final document bytes changed in transport")
		}
		if got := r.Header.Get("Content-Disposition"); got != `attachment; filename="Diploma.pdf"` {
			t.Errorf("unexpected download filename: %q", got)
		}
		if r.Header.Get("X-Amz-Checksum-Sha256") != base64.StdEncoding.EncodeToString(digest[:]) || r.Header.Get("If-None-Match") != "*" {
			t.Error("missing checksum or overwrite precondition")
		}
		for _, header := range []string{"X-Amz-Object-Lock-Mode", "X-Amz-Object-Lock-Retain-Until-Date", "X-Amz-Bypass-Governance-Retention"} {
			if r.Header.Get(header) != "" {
				t.Errorf("must preserve existing retention policy: %s", header)
			}
		}
		if stored {
			w.WriteHeader(http.StatusPreconditionFailed)
			_, _ = io.WriteString(w, "<Error><Code>PreconditionFailed</Code></Error>")
			return
		}
		stored = true
		w.Header().Set("X-Amz-Version-Id", "immutable-version-1")
		w.Header().Set("ETag", `"etag"`)
	}))
	defer server.Close()
	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("DIPLOMAS_S3_BUCKET", "diplomas-test")
	encoded := base64.StdEncoding.EncodeToString(finalPDF)
	ref, err := (S3Service{}).PutDiploma(context.Background(), "new-uuid", "Diploma.pdf", encoded)
	if err != nil {
		t.Fatal(err)
	}
	if ref.VersionID != "immutable-version-1" || ref.Key != "diplomas/new-uuid/Diploma.pdf" {
		t.Fatalf("missing exact version reference: %+v", ref)
	}
	if _, err := (S3Service{}).PutDiploma(context.Background(), "new-uuid", "Diploma.pdf", encoded); err == nil {
		t.Fatal("must reject overwrite; no deletion or retention bypass")
	}
}

func TestNombreDescargaSeguroEliminaRutaSinCambiarClaveS3(t *testing.T) {
	if got := nombreDescargaSeguro("diplomas/uuid/20261008_Diploma_20181016035.pdf"); got != "20261008_Diploma_20181016035.pdf" {
		t.Fatalf("nombreDescargaSeguro() = %q", got)
	}
}
