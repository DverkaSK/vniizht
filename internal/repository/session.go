package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, s *model.Session) error {
	_, err := r.db.Exec(ctx, queries.SessionCreate, s.ID, s.UserID, s.ExpiresAt)
	return err
}

func (r *SessionRepo) Get(ctx context.Context, id string) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRow(ctx, queries.SessionGet, id).
		Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.ExpiresAt, &s.LastActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return s, nil
}

func (r *SessionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, queries.SessionDelete, id)
	return err
}

func (r *SessionRepo) Touch(ctx context.Context, id string, now time.Time) error {
	_, err := r.db.Exec(ctx, queries.SessionTouch, now, id)
	return err
}

func (r *SessionRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.db.Exec(ctx, queries.SessionDeleteExpired)
	return err
}
