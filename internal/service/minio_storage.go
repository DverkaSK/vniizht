package service

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(client *minio.Client, bucket string) *MinioStorage {
	return &MinioStorage{client: client, bucket: bucket}
}

func (s *MinioStorage) PutObject(ctx context.Context, objectKey string, reader io.Reader, size int64, mimeType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	return err
}

func (s *MinioStorage) GetObject(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucket, objectKey, minio.GetObjectOptions{})
}

func (s *MinioStorage) RemoveObject(ctx context.Context, objectKey string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{})
}
