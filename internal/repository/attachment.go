package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type AttachmentRepo struct {
	db *pgxpool.Pool
}

func NewAttachmentRepo(db *pgxpool.Pool) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

func (r *AttachmentRepo) Create(ctx context.Context, a *model.Attachment) error {
	return r.db.QueryRow(ctx, queries.AttachmentCreate,
		a.UploaderID, string(a.TargetType), a.TargetID,
		a.Filename, a.ObjectKey, a.MimeType, a.SizeBytes,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *AttachmentRepo) GetByID(ctx context.Context, id int64) (*model.Attachment, error) {
	a := &model.Attachment{}
	var targetType string
	err := r.db.QueryRow(ctx, queries.AttachmentGetByID, id).
		Scan(&a.ID, &a.UploaderID, &targetType, &a.TargetID,
			&a.Filename, &a.ObjectKey, &a.MimeType, &a.SizeBytes, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get attachment: %w", err)
	}
	a.TargetType = model.AttachmentTarget(targetType)
	return a, nil
}

func (r *AttachmentRepo) ListByTarget(ctx context.Context, targetType model.AttachmentTarget, targetID int64) ([]*model.Attachment, error) {
	rows, err := r.db.Query(ctx, queries.AttachmentListByTarget, string(targetType), targetID)
	if err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	defer rows.Close()

	var attachments []*model.Attachment
	for rows.Next() {
		a := &model.Attachment{}
		var tt string
		if err := rows.Scan(&a.ID, &a.UploaderID, &tt, &a.TargetID,
			&a.Filename, &a.ObjectKey, &a.MimeType, &a.SizeBytes, &a.CreatedAt); err != nil {
			return nil, err
		}
		a.TargetType = model.AttachmentTarget(tt)
		attachments = append(attachments, a)
	}
	return attachments, rows.Err()
}

func (r *AttachmentRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.AttachmentDelete, id)
	if err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
