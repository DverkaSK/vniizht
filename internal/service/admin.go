package service

import (
	"context"
	"fmt"
	"strings"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type AdminService struct {
	users      *repository.UserRepo
	categories *repository.CategoryRepo
	tags       *repository.TagRepo
}

func NewAdminService(users *repository.UserRepo, categories *repository.CategoryRepo, tags *repository.TagRepo) *AdminService {
	return &AdminService{users: users, categories: categories, tags: tags}
}

func (s *AdminService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.users.List(ctx)
}

func (s *AdminService) CreateUser(ctx context.Context, username, email, password string, role model.Role) (*model.User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (s *AdminService) ChangeRole(ctx context.Context, id int64, role model.Role) error {
	return s.users.UpdateRole(ctx, id, role)
}

func (s *AdminService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return s.categories.List(ctx)
}

func (s *AdminService) CreateCategory(ctx context.Context, name, description string) (*model.Category, error) {
	category := &model.Category{Name: name, Description: description}
	if err := s.categories.Create(ctx, category); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (s *AdminService) UpdateCategory(ctx context.Context, id int64, name, description string) error {
	category := &model.Category{ID: id, Name: name, Description: description}
	if err := s.categories.Update(ctx, category); err != nil {
		if isDuplicateAdminErr(err) {
			return errs.ErrConflict
		}
		return err
	}
	return nil
}

func (s *AdminService) DeleteCategory(ctx context.Context, id int64) error {
	return s.categories.Delete(ctx, id)
}

func (s *AdminService) ListTags(ctx context.Context) ([]*model.Tag, error) {
	return s.tags.List(ctx)
}

func (s *AdminService) CreateTag(ctx context.Context, name string) (*model.Tag, error) {
	tag := &model.Tag{Name: name}
	if err := s.tags.Create(ctx, tag); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create tag: %w", err)
	}
	return tag, nil
}

func (s *AdminService) UpdateTag(ctx context.Context, id int64, name string) error {
	tag := &model.Tag{ID: id, Name: name}
	if err := s.tags.Update(ctx, tag); err != nil {
		if isDuplicateAdminErr(err) {
			return errs.ErrConflict
		}
		return err
	}
	return nil
}

func (s *AdminService) DeleteTag(ctx context.Context, id int64) error {
	return s.tags.Delete(ctx, id)
}

func isDuplicateAdminErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
