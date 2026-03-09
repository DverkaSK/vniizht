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

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(db *pgxpool.Pool) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	rows, err := r.db.Query(ctx, queries.TagList)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []*model.Tag
	for rows.Next() {
		t := &model.Tag{}
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	t := &model.Tag{}
	err := r.db.QueryRow(ctx, queries.TagGetByID, id).
		Scan(&t.ID, &t.Name, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tag: %w", err)
	}
	return t, nil
}

func (r *TagRepo) Create(ctx context.Context, t *model.Tag) error {
	return r.db.QueryRow(ctx, queries.TagCreate, t.Name).
		Scan(&t.ID, &t.CreatedAt)
}

func (r *TagRepo) Update(ctx context.Context, t *model.Tag) error {
	tag, err := r.db.Exec(ctx, queries.TagUpdate, t.Name, t.ID)
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TagRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.TagDelete, id)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
