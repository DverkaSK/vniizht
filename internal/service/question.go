package service

import (
	"context"
	"fmt"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type QuestionService struct {
	repo       *repository.QuestionRepo
	categories *repository.CategoryRepo
	notifs     *NotificationService
}

func NewQuestionService(repo *repository.QuestionRepo, categories *repository.CategoryRepo, notifs *NotificationService) *QuestionService {
	return &QuestionService{repo: repo, categories: categories, notifs: notifs}
}

func (s *QuestionService) List(ctx context.Context, f repository.QuestionFilter) ([]*model.Question, int, error) {
	return s.repo.List(ctx, f)
}

func (s *QuestionService) GetHistory(ctx context.Context, id int64) ([]model.QuestionHistory, error) {
	return s.repo.GetHistory(ctx, id)
}

func (s *QuestionService) Get(ctx context.Context, id int64) (*model.Question, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *QuestionService) Create(ctx context.Context, q *model.Question, tagIDs []int64) error {
	if q.CategoryID != nil {
		cat, err := s.categories.GetByID(ctx, *q.CategoryID)
		if err == nil && cat.SpecialistID != nil {
			q.SpecialistID = cat.SpecialistID
		}
	}
	if err := s.repo.Create(ctx, q, tagIDs); err != nil {
		return err
	}
	if q.SpecialistID != nil && *q.SpecialistID != q.AuthorID {
		s.notifs.Notify(ctx, *q.SpecialistID, model.NotifQuestionAssigned, model.NotificationPayload{
			QuestionID:    q.ID,
			QuestionTitle: q.Title,
			ActorUsername: q.AuthorUsername,
		})
	}
	return nil
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

func (s *QuestionService) AssignSpecialist(ctx context.Context, id int64, specialistID *int64, user *model.User) error {
	if user.Role != model.RoleAdmin {
		return errs.ErrForbidden
	}
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.AssignSpecialist(ctx, id, specialistID)
}

func (s *QuestionService) MarkDuplicate(ctx context.Context, id, duplicateOf int64, user *model.User) error {
	if id == duplicateOf {
		return errs.ErrInvalid
	}
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
	if _, err := s.repo.GetByID(ctx, duplicateOf); err != nil {
		return errs.ErrInvalid
	}
	return s.repo.MarkDuplicate(ctx, id, duplicateOf)
}

func (s *QuestionService) Delete(ctx context.Context, id int64, user *model.User) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

func (s *QuestionService) Close(ctx context.Context, id int64, user *model.User) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	if err := s.repo.Close(ctx, id); err != nil {
		return err
	}
	if existing.AuthorID != user.ID {
		s.notifs.Notify(ctx, existing.AuthorID, model.NotifQuestionClosed, model.NotificationPayload{
			QuestionID:    existing.ID,
			QuestionTitle: existing.Title,
			ActorUsername: user.Username,
		})
	}
	return nil
}
