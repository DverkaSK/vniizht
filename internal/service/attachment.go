package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type ObjectStorage interface {
	PutObject(ctx context.Context, objectKey string, reader io.Reader, size int64, mimeType string) error
	GetObject(ctx context.Context, objectKey string) (io.ReadCloser, error)
	RemoveObject(ctx context.Context, objectKey string) error
}

type AttachmentService struct {
	repo    *repository.AttachmentRepo
	storage ObjectStorage
}

func NewAttachmentService(repo *repository.AttachmentRepo, storage ObjectStorage) *AttachmentService {
	return &AttachmentService{repo: repo, storage: storage}
}

func (s *AttachmentService) List(ctx context.Context, targetType model.AttachmentTarget, targetID int64) ([]*model.Attachment, error) {
	return s.repo.ListByTarget(ctx, targetType, targetID)
}

func (s *AttachmentService) GetByID(ctx context.Context, id int64) (*model.Attachment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AttachmentService) Create(ctx context.Context, att *model.Attachment, file io.Reader) error {
	att.ObjectKey = fmt.Sprintf("%s/%d/%s-%s", att.TargetType, att.TargetID, randomHex(8), att.Filename)

	if err := s.storage.PutObject(ctx, att.ObjectKey, file, att.SizeBytes, att.MimeType); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, att); err != nil {
		_ = s.storage.RemoveObject(ctx, att.ObjectKey)
		return err
	}

	return nil
}

func (s *AttachmentService) CheckDeleteAccess(ctx context.Context, id int64, user *model.User) (*model.Attachment, error) {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if att.UploaderID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return nil, errs.ErrForbidden
	}
	return att, nil
}

func (s *AttachmentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *AttachmentService) Download(ctx context.Context, id int64) (*model.Attachment, io.ReadCloser, error) {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	obj, err := s.storage.GetObject(ctx, att.ObjectKey)
	if err != nil {
		return nil, nil, err
	}

	return att, obj, nil
}

func (s *AttachmentService) DeleteWithObject(ctx context.Context, id int64, user *model.User) error {
	att, err := s.CheckDeleteAccess(ctx, id, user)
	if err != nil {
		return err
	}

	_ = s.storage.RemoveObject(ctx, att.ObjectKey)
	return s.repo.Delete(ctx, id)
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
