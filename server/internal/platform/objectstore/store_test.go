package objectstore_test

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

func newStore(t *testing.T) *objectstore.Store {
	t.Helper()
	s, err := objectstore.New(objectstore.Config{
		Endpoint:       "http://minio:9000",
		PublicEndpoint: "https://files.example.edu",
		Region:         "us-east-1",
		Bucket:         "lodestar",
		AccessKey:      "key",
		SecretKey:      "secret",
	})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	return s
}

func TestPresignedUploadUsesPublicEndpointAndSignsSize(t *testing.T) {
	upload, err := newStore(t).PresignPut(context.Background(), "t/f", "application/pdf", 42, time.Minute)
	if err != nil {
		t.Fatalf("presign: %v", err)
	}
	u, err := url.Parse(upload.URL)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	signed := u.Query().Get("X-Amz-SignedHeaders")
	if u.Host != "files.example.edu" || u.Path != "/lodestar/t/f" || !strings.Contains(signed, "content-length") || !strings.Contains(signed, "content-type") {
		t.Fatalf("url = %s", upload.URL)
	}
	if upload.Method != "PUT" || upload.Headers["Content-Type"] != "application/pdf" {
		t.Fatalf("upload = %+v", upload)
	}
}

func TestPresignedDownloadNamesFile(t *testing.T) {
	link, err := newStore(t).PresignGet(context.Background(), "t/f", "Курсовая работа.pdf", "application/pdf", time.Minute)
	if err != nil {
		t.Fatalf("presign: %v", err)
	}
	u, _ := url.Parse(link)
	if got := u.Query().Get("response-content-disposition"); !strings.Contains(got, "%D0%9A") {
		t.Fatalf("disposition = %q", got)
	}
}

func TestRejectsMalformedEndpoint(t *testing.T) {
	if _, err := objectstore.New(objectstore.Config{Endpoint: "minio:9000", Bucket: "b"}); err == nil {
		t.Fatal("endpoint without scheme accepted")
	}
}
