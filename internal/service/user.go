package service

import (
	"context"

	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type UserService struct {
	users     *repository.UserRepo
	questions *repository.QuestionRepo
	answers   *repository.AnswerRepo
}

func NewUserService(users *repository.UserRepo, questions *repository.QuestionRepo, answers *repository.AnswerRepo) *UserService {
	return &UserService{users: users, questions: questions, answers: answers}
}

type UserProfile struct {
	User            *model.User
	QuestionCount   int
	AnswerCount     int
	RecentQuestions []*model.Question
	RecentAnswers   []*model.Answer
}

func (s *UserService) GetMe(ctx context.Context, user *model.User) (*UserProfile, error) {
	qCount, err := s.questions.CountByAuthor(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	aCount, err := s.answers.CountByAuthor(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		User:          user,
		QuestionCount: qCount,
		AnswerCount:   aCount,
	}, nil
}

func (s *UserService) GetPublic(ctx context.Context, id int64) (*UserProfile, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	qCount, err := s.questions.CountByAuthor(ctx, id)
	if err != nil {
		return nil, err
	}

	aCount, err := s.answers.CountByAuthor(ctx, id)
	if err != nil {
		return nil, err
	}

	recentQ, err := s.questions.ListRecentByAuthor(ctx, id, 5)
	if err != nil {
		return nil, err
	}

	recentA, err := s.answers.ListRecentByAuthor(ctx, id, 5)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		User:            user,
		QuestionCount:   qCount,
		AnswerCount:     aCount,
		RecentQuestions: recentQ,
		RecentAnswers:   recentA,
	}, nil
}
