package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type NotificationRepo struct {
	db *pgxpool.Pool
}

func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, n *model.Notification) error {
	payload, err := json.Marshal(n.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return r.db.QueryRow(ctx, queries.NotificationCreate, n.UserID, n.Type, payload).
		Scan(&n.ID, &n.CreatedAt)
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.Notification, error) {
	rows, err := r.db.Query(ctx, queries.NotificationListByUser, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var result []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		var rawPayload []byte
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &rawPayload, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(rawPayload, &n.Payload)
		result = append(result, n)
	}
	return result, rows.Err()
}

func (r *NotificationRepo) CountUnread(ctx context.Context, userID int64) (int, error) {
	var count int
	if err := r.db.QueryRow(ctx, queries.NotificationCountUnread, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count unread: %w", err)
	}
	return count, nil
}

func (r *NotificationRepo) MarkRead(ctx context.Context, id, userID int64) error {
	_, err := r.db.Exec(ctx, queries.NotificationMarkRead, id, userID)
	return err
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, queries.NotificationMarkAllRead, userID)
	return err
}

type NotificationPreference struct {
	Type         model.NotificationType
	EmailEnabled bool
}

func (r *NotificationRepo) GetPreferences(ctx context.Context, userID int64) ([]NotificationPreference, error) {
	rows, err := r.db.Query(ctx, queries.NotificationPreferencesGet, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	defer rows.Close()

	var result []NotificationPreference
	for rows.Next() {
		var p NotificationPreference
		if err := rows.Scan(&p.Type, &p.EmailEnabled); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *NotificationRepo) UpsertPreference(ctx context.Context, userID int64, typ model.NotificationType, emailEnabled bool) error {
	_, err := r.db.Exec(ctx, queries.NotificationPreferenceUpsert, userID, typ, emailEnabled)
	return err
}

func (r *NotificationRepo) IsEmailEnabled(ctx context.Context, userID int64, typ model.NotificationType) (bool, error) {
	var enabled bool
	err := r.db.QueryRow(ctx, queries.NotificationPreferenceEmailEnabled, userID, typ).Scan(&enabled)
	return enabled, err
}
