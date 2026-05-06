package service

import (
	"context"
	"fmt"

	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type TrainingEntry struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Questions []string `json:"questions"`
	Positive  string   `json:"positive"`
	Original  string   `json:"original"`
}

type ImportService struct {
	questions *repository.QuestionRepo
	answers   *repository.AnswerRepo
}

func NewImportService(questions *repository.QuestionRepo, answers *repository.AnswerRepo) *ImportService {
	return &ImportService{questions: questions, answers: answers}
}

func (s *ImportService) Import(ctx context.Context, entries []TrainingEntry, authorID int64) (int, error) {
	existing, err := s.questions.ListTitles(ctx)
	if err != nil {
		return 0, fmt.Errorf("load existing titles: %w", err)
	}

	count := 0
	for _, entry := range entries {
		if entry.Title == "" || entry.Positive == "" {
			continue
		}
		if _, dup := existing[entry.Title]; dup {
			continue
		}

		body := entry.Title
		if len(entry.Questions) > 0 && entry.Questions[0] != "" {
			body = entry.Questions[0]
		}

		q := &model.Question{
			AuthorID: authorID,
			Title:    entry.Title,
			Body:     body,
		}
		if err := s.questions.Create(ctx, q, nil); err != nil {
			return count, fmt.Errorf("create question (entry id=%d): %w", entry.ID, err)
		}
		existing[entry.Title] = struct{}{}

		a := &model.Answer{
			QuestionID: q.ID,
			AuthorID:   authorID,
			Body:       entry.Positive,
		}
		if err := s.answers.Create(ctx, a); err != nil {
			return count, fmt.Errorf("create answer (entry id=%d): %w", entry.ID, err)
		}

		if err := s.answers.SetVerified(ctx, q.ID, a.ID); err != nil {
			return count, fmt.Errorf("verify answer (entry id=%d): %w", entry.ID, err)
		}

		if err := s.questions.Close(ctx, q.ID); err != nil {
			return count, fmt.Errorf("close question (entry id=%d): %w", entry.ID, err)
		}

		count++
	}
	return count, nil
}

func (s *ImportService) Export(ctx context.Context) ([]TrainingEntry, error) {
	rows, err := s.questions.ListForExport(ctx)
	if err != nil {
		return nil, fmt.Errorf("export questions: %w", err)
	}

	entries := make([]TrainingEntry, len(rows))
	for i, row := range rows {
		answerBody := ""
		if row.AnswerBody != nil {
			answerBody = *row.AnswerBody
		}
		entries[i] = TrainingEntry{
			ID:        i,
			Title:     row.Title,
			Questions: []string{row.Body},
			Positive:  answerBody,
			Original:  answerBody,
		}
	}
	return entries, nil
}
