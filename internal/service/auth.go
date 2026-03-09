package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"vniizht/internal/model"
	"vniizht/internal/repository"
)

const (
	sessionTTL      = 24 * time.Hour
	sessionInactive = 2 * time.Hour
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is disabled")
	ErrSessionExpired     = errors.New("session expired")
)

type AuthService struct {
	users    *repository.UserRepo
	sessions *repository.SessionRepo
}

func NewAuthService(users *repository.UserRepo, sessions *repository.SessionRepo) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if errors.Is(err, repository.ErrNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}
	if !user.IsActive {
		return "", ErrUserInactive
	}
	if !CheckPassword(user.PasswordHash, password) {
		return "", ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	sess := &model.Session{
		ID:        token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	if err := s.sessions.Create(ctx, sess); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessions.Delete(ctx, token)
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (*model.User, error) {
	sess, err := s.sessions.Get(ctx, token)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrSessionExpired
	}
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	now := time.Now()
	if now.After(sess.ExpiresAt) {
		_ = s.sessions.Delete(ctx, token)
		return nil, ErrSessionExpired
	}
	if now.Sub(sess.LastActive) > sessionInactive {
		_ = s.sessions.Delete(ctx, token)
		return nil, ErrSessionExpired
	}

	_ = s.sessions.Touch(ctx, token, now)

	user, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("get session user: %w", err)
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return user, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
