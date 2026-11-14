package objectstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	internal *minio.Client
	public   *minio.Client
	bucket   string
}

func New(cfg Config) (*Store, error) {
	internal, err := client(cfg, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("object store endpoint: %w", err)
	}
	publicEndpoint := cfg.PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.Endpoint
	}
	public, err := client(cfg, publicEndpoint)
	if err != nil {
		return nil, fmt.Errorf("object store public endpoint: %w", err)
	}
	return &Store{internal: internal, public: public, bucket: cfg.Bucket}, nil
}

func (s *Store) EnsureBucket(ctx context.Context) error {
	exists, err := s.internal.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if exists {
		return nil
	}
	if err := s.internal.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	exists, err := s.internal.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		return ErrMissing
	}
	return nil
}

func (s *Store) PresignPut(ctx context.Context, key, contentType string, size int64, ttl time.Duration) (Upload, error) {
	headers := http.Header{}
	headers.Set("Content-Type", contentType)
	headers.Set("Content-Length", strconv.FormatInt(size, 10))
	u, err := s.public.PresignHeader(ctx, http.MethodPut, s.bucket, key, ttl, nil, headers)
	if err != nil {
		return Upload{}, fmt.Errorf("presign upload: %w", err)
	}
	return Upload{
		URL:       u.String(),
		Method:    http.MethodPut,
		Headers:   map[string]string{"Content-Type": contentType},
		ExpiresAt: time.Now().Add(ttl).UTC(),
	}, nil
}

func (s *Store) PresignGet(ctx context.Context, key, name, contentType string, ttl time.Duration) (string, error) {
	params := url.Values{}
	params.Set("response-content-disposition", "attachment; filename*=UTF-8''"+url.PathEscape(name))
	params.Set("response-content-type", contentType)
	u, err := s.public.PresignedGetObject(ctx, s.bucket, key, ttl, params)
	if err != nil {
		return "", fmt.Errorf("presign download: %w", err)
	}
	return u.String(), nil
}

func (s *Store) Size(ctx context.Context, key string) (int64, error) {
	info, err := s.internal.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == http.StatusNotFound {
			return 0, ErrMissing
		}
		return 0, fmt.Errorf("stat object: %w", err)
	}
	return info.Size, nil
}

func (s *Store) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.internal.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("open object: %w", err)
	}
	return obj, nil
}

func (s *Store) Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error {
	if _, err := s.internal.PutObject(ctx, s.bucket, key, body, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	if err := s.internal.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func client(cfg Config, endpoint string) (*minio.Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("%q must look like https://host:port", endpoint)
	}
	return minio.New(u.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:       u.Scheme == "https",
		Region:       cfg.Region,
		BucketLookup: minio.BucketLookupPath,
	})
}
