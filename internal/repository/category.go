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

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.db.Query(ctx, queries.CategoryList)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var cats []*model.Category
	for rows.Next() {
		c := &model.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.SpecialistID, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	c := &model.Category{}
	err := r.db.QueryRow(ctx, queries.CategoryGetByID, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.SpecialistID, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	return c, nil
}

func (r *CategoryRepo) Create(ctx context.Context, c *model.Category) error {
	return r.db.QueryRow(ctx, queries.CategoryCreate, c.Name, c.Description, c.SpecialistID).
		Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepo) Update(ctx context.Context, c *model.Category) error {
	tag, err := r.db.Exec(ctx, queries.CategoryUpdate, c.Name, c.Description, c.SpecialistID, c.ID)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.CategoryDelete, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
