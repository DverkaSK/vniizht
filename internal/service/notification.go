package service

import (
	"context"
	"log/slog"

	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type NotificationService struct {
	repo  *repository.NotificationRepo
	users *repository.UserRepo
	email *EmailService
}

func NewNotificationService(repo *repository.NotificationRepo, users *repository.UserRepo, email *EmailService) *NotificationService {
	return &NotificationService{repo: repo, users: users, email: email}
}

// Notify creates a notification and sends email if the user has it enabled.
func (s *NotificationService) Notify(ctx context.Context, userID int64, typ model.NotificationType, payload model.NotificationPayload) {
	n := &model.Notification{UserID: userID, Type: typ, Payload: payload}
	if err := s.repo.Create(ctx, n); err != nil {
		slog.Warn("failed to create notification", "err", err, "type", typ, "user_id", userID)
		return
	}

	if !s.email.enabled {
		return
	}

	emailEnabled, err := s.repo.IsEmailEnabled(ctx, userID, typ)
	if err != nil || !emailEnabled {
		return
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return
	}

	s.email.SendNotification(user.Email, typ, payload)
}

func (s *NotificationService) List(ctx context.Context, userID int64, page int) ([]*model.Notification, error) {
	const limit = 20
	offset := (page - 1) * limit
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *NotificationService) CountUnread(ctx context.Context, userID int64) (int, error) {
	return s.repo.CountUnread(ctx, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, id, userID int64) error {
	return s.repo.MarkRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *NotificationService) GetPreferences(ctx context.Context, userID int64) ([]repository.NotificationPreference, error) {
	return s.repo.GetPreferences(ctx, userID)
}

func (s *NotificationService) SavePreferences(ctx context.Context, userID int64, prefs []repository.NotificationPreference) error {
	for _, p := range prefs {
		if err := s.repo.UpsertPreference(ctx, userID, p.Type, p.EmailEnabled); err != nil {
			return err
		}
	}
	return nil
}
