package service

import (
	"context"
	"fmt"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type CommentService struct {
	repo    *repository.CommentRepo
	answers *repository.AnswerRepo
}

func NewCommentService(repo *repository.CommentRepo, answers *repository.AnswerRepo) *CommentService {
	return &CommentService{repo: repo, answers: answers}
}

func (s *CommentService) List(ctx context.Context, answerID int64) ([]*model.Comment, error) {
	return s.repo.ListByAnswer(ctx, answerID)
}

func (s *CommentService) Create(ctx context.Context, answerID int64, user *model.User, body string) (*model.Comment, error) {
	if _, err := s.answers.GetByID(ctx, answerID); err != nil {
		return nil, err
	}
	c := &model.Comment{
		AnswerID: answerID,
		AuthorID: user.ID,
		Body:     body,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	c.AuthorUsername = user.Username
	return c, nil
}

func (s *CommentService) Update(ctx context.Context, id int64, user *model.User, body string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	existing.Body = body
	if err := s.repo.Update(ctx, existing); err != nil {
		return fmt.Errorf("update comment: %w", err)
	}
	return nil
}

func (s *CommentService) Delete(ctx context.Context, id int64, user *model.User) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}
