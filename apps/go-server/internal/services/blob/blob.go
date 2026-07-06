// Package blob ports services/s3.py + services/cloudinary_service.py to a single
// Azure Blob Storage-backed uploader, per the rewrite's storage decision
// (Azure Blob replaces S3/Cloudinary; aligns billing with Container Apps).
package blob

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// Store uploads/deletes course media. Falls back to an inline base64 data URI
// when Azure Blob isn't configured or the upload fails, mirroring the legacy
// try/except cascade in course.py's _upload_media so local dev without Azure
// credentials still works end-to-end.
type Store struct {
	client    *azblob.Client
	container string
}

// New builds a Store. If connectionString is empty, Upload always falls back
// to inline data URIs (useful for local dev without an Azure Storage account).
func New(connectionString, container string) (*Store, error) {
	if strings.TrimSpace(connectionString) == "" {
		return &Store{container: container}, nil
	}
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("azure blob client: %w", err)
	}
	return &Store{client: client, container: container}, nil
}

func dataURI(fileBytes []byte, contentType string) string {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(fileBytes))
}

// Upload stores fileBytes under folder/ and returns a public URL (or a data
// URI fallback). Blob name is randomized to avoid collisions, matching the
// legacy services' behaviour of namespacing uploads by folder.
func (s *Store) Upload(ctx context.Context, fileBytes []byte, folder, filename, contentType string) (string, error) {
	if s.client == nil {
		return dataURI(fileBytes, contentType), nil
	}

	blobName := fmt.Sprintf("%s/%d-%s-%s", folder, time.Now().UnixNano(), uuid.NewString()[:8], filename)
	if _, err := s.client.UploadBuffer(ctx, s.container, blobName, fileBytes, nil); err != nil {
		// Degrade gracefully instead of failing the request outright.
		return dataURI(fileBytes, contentType), nil
	}
	base := strings.TrimSuffix(s.client.URL(), "/")
	return fmt.Sprintf("%s/%s/%s", base, s.container, blobName), nil
}

// Delete removes a previously uploaded blob given its full public URL.
// No-ops for data URIs or when Azure Blob isn't configured (matches the
// legacy _delete_media_url short-circuit on "data:" URLs).
func (s *Store) Delete(ctx context.Context, fileURL string) error {
	if s.client == nil || fileURL == "" || strings.HasPrefix(fileURL, "data:") {
		return nil
	}
	prefix := strings.TrimSuffix(s.client.URL(), "/") + "/" + s.container + "/"
	if !strings.HasPrefix(fileURL, prefix) {
		return nil // not one of ours (e.g. a leftover Cloudinary/S3 URL from before migration)
	}
	blobName := strings.TrimPrefix(fileURL, prefix)
	_, err := s.client.DeleteBlob(ctx, s.container, blobName, nil)
	return err
}
