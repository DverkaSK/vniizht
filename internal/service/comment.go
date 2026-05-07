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
	notifs  *NotificationService
}

func NewCommentService(repo *repository.CommentRepo, answers *repository.AnswerRepo, notifs *NotificationService) *CommentService {
	return &CommentService{repo: repo, answers: answers, notifs: notifs}
}

func (s *CommentService) List(ctx context.Context, answerID int64) ([]*model.Comment, error) {
	return s.repo.ListByAnswer(ctx, answerID)
}

func (s *CommentService) Create(ctx context.Context, answerID int64, user *model.User, body string) (*model.Comment, error) {
	answer, err := s.answers.GetByID(ctx, answerID)
	if err != nil {
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

	if answer.AuthorID != user.ID {
		s.notifs.Notify(ctx, answer.AuthorID, model.NotifNewComment, model.NotificationPayload{
			QuestionID:    answer.QuestionID,
			AnswerID:      answer.ID,
			CommentID:     c.ID,
			ActorUsername: user.Username,
		})
	}
	return c, nil
}

func (s *CommentService) Update(ctx context.Context, id int64, user *model.User, body string) (*model.Comment, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return nil, errs.ErrForbidden
	}
	existing.Body = body
	updatedAt, err := s.repo.Update(ctx, existing, user.ID)
	if err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}
	existing.UpdatedAt = updatedAt
	return existing, nil
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
