package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type AnswerRepo struct {
	db *pgxpool.Pool
}

func NewAnswerRepo(db *pgxpool.Pool) *AnswerRepo {
	return &AnswerRepo{db: db}
}

func (r *AnswerRepo) Create(ctx context.Context, a *model.Answer) error {
	return r.db.QueryRow(ctx, queries.AnswerCreate,
		a.QuestionID, a.AuthorID, a.Body,
	).Scan(&a.ID, &a.IsVerified, &a.VoteScore, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AnswerRepo) GetByID(ctx context.Context, id int64) (*model.Answer, error) {
	a := &model.Answer{}
	err := r.db.QueryRow(ctx, queries.AnswerGetByID, id).
		Scan(&a.ID, &a.QuestionID, &a.AuthorID, &a.Body,
			&a.IsVerified, &a.VoteScore, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get answer: %w", err)
	}
	return a, nil
}

func (r *AnswerRepo) ListByQuestion(ctx context.Context, questionID int64, currentUserID *int64) ([]*model.Answer, error) {
	rows, err := r.db.Query(ctx, queries.AnswerListByQuestion, questionID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("list answers: %w", err)
	}
	defer rows.Close()

	var answers []*model.Answer
	for rows.Next() {
		a := &model.Answer{}
		var currentVote sql.NullString
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.AuthorID, &a.AuthorUsername, &a.Body,
			&a.IsVerified, &a.VoteScore, &currentVote, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		if currentVote.Valid {
			vote := model.VoteValue(currentVote.String)
			a.CurrentUserVote = &vote
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}

func (r *AnswerRepo) Update(ctx context.Context, a *model.Answer, editorID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.AnswerHistoryInsert, a.ID, editorID, a.Body); err != nil {
		return fmt.Errorf("save history: %w", err)
	}

	tag, err := tx.Exec(ctx, queries.AnswerUpdate, a.Body, a.ID)
	if err != nil {
		return fmt.Errorf("update answer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *AnswerRepo) CountByAuthor(ctx context.Context, userID int64) (int, error) {
	var count int
	if err := r.db.QueryRow(ctx, queries.AnswerCountByAuthor, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count answers by author: %w", err)
	}
	return count, nil
}

func (r *AnswerRepo) ListRecentByAuthor(ctx context.Context, userID int64, limit int) ([]*model.Answer, error) {
	rows, err := r.db.Query(ctx, queries.AnswerListRecentByAuthor, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent answers by author: %w", err)
	}
	defer rows.Close()

	var answers []*model.Answer
	for rows.Next() {
		a := &model.Answer{}
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.AuthorID, &a.AuthorUsername, &a.Body,
			&a.IsVerified, &a.VoteScore, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}

func (r *AnswerRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.AnswerDelete, id)
	if err != nil {
		return fmt.Errorf("delete answer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AnswerRepo) SetVerified(ctx context.Context, questionID, answerID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, queries.AnswerSetVerified, questionID, answerID)
	if err != nil {
		return fmt.Errorf("set verified answer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *AnswerRepo) ClearVerified(ctx context.Context, answerID int64) error {
	tag, err := r.db.Exec(ctx, queries.AnswerClearVerified, answerID)
	if err != nil {
		return fmt.Errorf("clear verified answer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AnswerRepo) CreateVote(ctx context.Context, answerID, userID int64, value model.VoteValue) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.AnswerVoteInsert, answerID, userID, value); err != nil {
		return fmt.Errorf("insert answer vote: %w", err)
	}
	if _, err := tx.Exec(ctx, queries.AnswerVoteRecount, answerID); err != nil {
		return fmt.Errorf("recount answer votes: %w", err)
	}
	if _, err := tx.Exec(ctx, queries.UserReputationRecountByAnswer, answerID); err != nil {
		return fmt.Errorf("recount reputation: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *AnswerRepo) UpdateVote(ctx context.Context, answerID, userID int64, value model.VoteValue) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, queries.AnswerVoteUpdate, answerID, userID, value)
	if err != nil {
		return fmt.Errorf("update answer vote: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, queries.AnswerVoteRecount, answerID); err != nil {
		return fmt.Errorf("recount answer votes: %w", err)
	}
	if _, err := tx.Exec(ctx, queries.UserReputationRecountByAnswer, answerID); err != nil {
		return fmt.Errorf("recount reputation: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *AnswerRepo) DeleteVote(ctx context.Context, answerID, userID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, queries.AnswerVoteDelete, answerID, userID)
	if err != nil {
		return fmt.Errorf("delete answer vote: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, queries.AnswerVoteRecount, answerID); err != nil {
		return fmt.Errorf("recount answer votes: %w", err)
	}
	if _, err := tx.Exec(ctx, queries.UserReputationRecountByAnswer, answerID); err != nil {
		return fmt.Errorf("recount reputation: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *AnswerRepo) GetHistory(ctx context.Context, answerID int64) ([]model.AnswerHistory, error) {
	rows, err := r.db.Query(ctx, queries.AnswerHistoryGet, answerID)
	if err != nil {
		return nil, fmt.Errorf("get answer history: %w", err)
	}
	defer rows.Close()

	var result []model.AnswerHistory
	for rows.Next() {
		var h model.AnswerHistory
		if err := rows.Scan(&h.ID, &h.AnswerID, &h.EditorID, &h.EditorUsername, &h.Body, &h.EditedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, rows.Err()
}

func (r *AnswerRepo) GetVote(ctx context.Context, answerID, userID int64) (*model.VoteValue, error) {
	var value string
	err := r.db.QueryRow(ctx, queries.AnswerVoteGet, answerID, userID).Scan(&value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get answer vote: %w", err)
	}
	vote := model.VoteValue(value)
	return &vote, nil
}
