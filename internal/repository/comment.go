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

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, c *model.Comment) error {
	return r.db.QueryRow(ctx, queries.CommentCreate,
		c.AnswerID, c.AuthorID, c.Body,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CommentRepo) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	c := &model.Comment{}
	err := r.db.QueryRow(ctx, queries.CommentGetByID, id).
		Scan(&c.ID, &c.AnswerID, &c.AuthorID, &c.Body, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return c, nil
}

func (r *CommentRepo) ListByAnswer(ctx context.Context, answerID int64) ([]*model.Comment, error) {
	rows, err := r.db.Query(ctx, queries.CommentListByAnswer, answerID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.AnswerID, &c.AuthorID, &c.AuthorUsername, &c.Body,
			&c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *CommentRepo) Update(ctx context.Context, c *model.Comment) error {
	tag, err := r.db.Exec(ctx, queries.CommentUpdate, c.Body, c.ID)
	if err != nil {
		return fmt.Errorf("update comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.CommentDelete, id)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
