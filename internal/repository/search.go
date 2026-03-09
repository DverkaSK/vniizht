package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type SearchRepo struct {
	db *pgxpool.Pool
}

func NewSearchRepo(db *pgxpool.Pool) *SearchRepo {
	return &SearchRepo{db: db}
}

func (r *SearchRepo) Search(ctx context.Context, f model.SearchFilter) ([]*model.SearchResult, error) {
	var statusStr *string
	if f.Status != nil {
		s := string(*f.Status)
		statusStr = &s
	}

	var results []*model.SearchResult

	qRows, err := r.db.Query(ctx, queries.SearchQuestions, f.Query, f.CategoryID, statusStr, f.TagID)
	if err != nil {
		return nil, fmt.Errorf("search questions: %w", err)
	}
	defer qRows.Close()

	for qRows.Next() {
		sr := &model.SearchResult{ResultType: "question"}
		var status string
		if err := qRows.Scan(&sr.QuestionID, &sr.QuestionTitle, &sr.Snippet, &status,
			&sr.AuthorID, &sr.CategoryID, &sr.CreatedAt, &sr.Rank); err != nil {
			return nil, err
		}
		sr.Status = model.QuestionStatus(status)
		results = append(results, sr)
	}
	if err := qRows.Err(); err != nil {
		return nil, err
	}
	qRows.Close()

	aRows, err := r.db.Query(ctx, queries.SearchAnswers, f.Query, f.CategoryID, statusStr, f.TagID)
	if err != nil {
		return nil, fmt.Errorf("search answers: %w", err)
	}
	defer aRows.Close()

	for aRows.Next() {
		sr := &model.SearchResult{ResultType: "answer"}
		var answerID int64
		var status string
		if err := aRows.Scan(&answerID, &sr.QuestionID, &sr.QuestionTitle, &sr.Snippet, &status,
			&sr.AuthorID, &sr.CategoryID, &sr.CreatedAt, &sr.Rank); err != nil {
			return nil, err
		}
		sr.AnswerID = &answerID
		sr.Status = model.QuestionStatus(status)
		results = append(results, sr)
	}
	if err := aRows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Rank > results[j].Rank
	})

	return results, nil
}
