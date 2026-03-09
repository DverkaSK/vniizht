package service

import (
	"context"
	"fmt"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type QuestionService struct {
	repo *repository.QuestionRepo
}

func NewQuestionService(repo *repository.QuestionRepo) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) List(ctx context.Context, f repository.QuestionFilter) ([]*model.Question, int, error) {
	return s.repo.List(ctx, f)
}

func (s *QuestionService) Get(ctx context.Context, id int64) (*model.Question, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *QuestionService) Create(ctx context.Context, q *model.Question, tagIDs []int64) error {
	return s.repo.Create(ctx, q, tagIDs)
}

func (s *QuestionService) Update(ctx context.Context, id int64, user *model.User, title, body string, tagIDs []int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	if existing.Status != model.QuestionOpen {
		return errs.ErrConflict
	}
	existing.Title = title
	existing.Body = body
	if err := s.repo.Update(ctx, existing, tagIDs, user.ID); err != nil {
		return fmt.Errorf("update question: %w", err)
	}
	return nil
}

func (s *QuestionService) Close(ctx context.Context, id, userID int64) error {
	return s.repo.Close(ctx, id, userID)
}
