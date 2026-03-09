package service

import (
	"context"
	"fmt"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type AnswerService struct {
	repo      *repository.AnswerRepo
	questions *repository.QuestionRepo
}

func NewAnswerService(repo *repository.AnswerRepo, questions *repository.QuestionRepo) *AnswerService {
	return &AnswerService{repo: repo, questions: questions}
}

func (s *AnswerService) List(ctx context.Context, questionID int64, user *model.User) ([]*model.Answer, error) {
	var currentUserID *int64
	if user != nil {
		currentUserID = &user.ID
	}
	return s.repo.ListByQuestion(ctx, questionID, currentUserID)
}

func (s *AnswerService) GetByID(ctx context.Context, id int64) (*model.Answer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AnswerService) Create(ctx context.Context, questionID int64, user *model.User, body string) (*model.Answer, error) {
	q, err := s.questions.GetByID(ctx, questionID)
	if err != nil {
		return nil, err
	}
	if q.Status != model.QuestionOpen {
		return nil, errs.ErrConflict
	}
	a := &model.Answer{
		QuestionID: questionID,
		AuthorID:   user.ID,
		Body:       body,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("create answer: %w", err)
	}
	a.AuthorUsername = user.Username
	return a, nil
}

func (s *AnswerService) Update(ctx context.Context, id int64, user *model.User, body string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	existing.Body = body
	if err := s.repo.Update(ctx, existing, user.ID); err != nil {
		return fmt.Errorf("update answer: %w", err)
	}
	return nil
}

func (s *AnswerService) Delete(ctx context.Context, id int64, user *model.User) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != user.ID && user.Role != model.RoleAdmin && user.Role != model.RoleSpecialist {
		return errs.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

func (s *AnswerService) Verify(ctx context.Context, id int64, user *model.User) error {
	if user.Role != model.RoleSpecialist && user.Role != model.RoleAdmin {
		return errs.ErrForbidden
	}

	answer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.SetVerified(ctx, answer.QuestionID, id)
}

func (s *AnswerService) Unverify(ctx context.Context, id int64, user *model.User) error {
	if user.Role != model.RoleSpecialist && user.Role != model.RoleAdmin {
		return errs.ErrForbidden
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	return s.repo.ClearVerified(ctx, id)
}

func (s *AnswerService) CreateVote(ctx context.Context, id int64, user *model.User, value model.VoteValue) error {
	if !isValidVoteValue(value) {
		return errs.ErrInvalid
	}

	answer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if answer.AuthorID == user.ID {
		return errs.ErrForbidden
	}

	if _, err := s.repo.GetVote(ctx, id, user.ID); err == nil {
		return errs.ErrConflict
	} else if err != repository.ErrNotFound {
		return err
	}

	return s.repo.CreateVote(ctx, id, user.ID, value)
}

func (s *AnswerService) UpdateVote(ctx context.Context, id int64, user *model.User, value model.VoteValue) error {
	if !isValidVoteValue(value) {
		return errs.ErrInvalid
	}

	answer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if answer.AuthorID == user.ID {
		return errs.ErrForbidden
	}

	return s.repo.UpdateVote(ctx, id, user.ID, value)
}

func (s *AnswerService) DeleteVote(ctx context.Context, id int64, user *model.User) error {
	answer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if answer.AuthorID == user.ID {
		return errs.ErrForbidden
	}

	return s.repo.DeleteVote(ctx, id, user.ID)
}

func isValidVoteValue(value model.VoteValue) bool {
	return value == model.VoteUp || value == model.VoteDown
}
