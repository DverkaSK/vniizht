package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vniizht/internal/model"
	"vniizht/internal/repository/queries"
)

type QuestionFilter struct {
	Status              *model.QuestionStatus
	CategoryID          *int64
	TagID               *int64
	AssignedSpecialistID *int64
	Limit               int
	Offset              int
}

type QuestionRepo struct {
	db *pgxpool.Pool
}

func NewQuestionRepo(db *pgxpool.Pool) *QuestionRepo {
	return &QuestionRepo{db: db}
}

func (r *QuestionRepo) Create(ctx context.Context, q *model.Question, tagIDs []int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, queries.QuestionCreate,
		q.AuthorID, q.CategoryID, q.SpecialistID, q.Title, q.Body,
	).Scan(&q.ID, &q.Status, &q.ViewCount, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create question: %w", err)
	}

	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx, queries.QuestionTagInsert, q.ID, tagID); err != nil {
			return fmt.Errorf("attach tag %d: %w", tagID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	q.Tags, err = r.getTags(ctx, q.ID)
	return err
}

func (r *QuestionRepo) GetByID(ctx context.Context, id int64) (*model.Question, error) {
	q := &model.Question{}
	err := r.db.QueryRow(ctx, queries.QuestionGetByID, id).
		Scan(&q.ID, &q.AuthorID, &q.AuthorUsername, &q.CategoryID, &q.SpecialistID, &q.SpecialistUsername,
			&q.Title, &q.Body, &q.Status, &q.DuplicateOf, &q.ViewCount,
			&q.CreatedAt, &q.UpdatedAt, &q.AnswerCount, &q.HasVerified)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get question: %w", err)
	}

	tags, err := r.getTags(ctx, id)
	if err != nil {
		return nil, err
	}
	q.Tags = tags

	_, _ = r.db.Exec(ctx, queries.QuestionIncrementView, id)

	return q, nil
}

func (r *QuestionRepo) List(ctx context.Context, f QuestionFilter) ([]*model.Question, int, error) {
	rows, err := r.db.Query(ctx, queries.QuestionList,
		f.Status, f.CategoryID, f.TagID, f.AssignedSpecialistID, f.Limit, f.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list questions: %w", err)
	}
	defer rows.Close()

	var qs []*model.Question
	for rows.Next() {
		q := &model.Question{}
		if err := rows.Scan(&q.ID, &q.AuthorID, &q.AuthorUsername, &q.CategoryID, &q.SpecialistID, &q.SpecialistUsername,
			&q.Title, &q.Body, &q.Status, &q.DuplicateOf, &q.ViewCount,
			&q.CreatedAt, &q.UpdatedAt, &q.AnswerCount, &q.HasVerified); err != nil {
			return nil, 0, err
		}
		qs = append(qs, q)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.db.QueryRow(ctx, queries.QuestionCount, f.Status, f.CategoryID, f.TagID, f.AssignedSpecialistID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count questions: %w", err)
	}

	return qs, total, nil
}

func (r *QuestionRepo) Update(ctx context.Context, q *model.Question, tagIDs []int64, editorID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.QuestionHistoryInsert,
		q.ID, editorID, q.Title, q.Body); err != nil {
		return fmt.Errorf("save history: %w", err)
	}

	tag, err := tx.Exec(ctx, queries.QuestionUpdate, q.Title, q.Body, q.ID)
	if err != nil {
		return fmt.Errorf("update question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	if _, err := tx.Exec(ctx, queries.QuestionTagsDelete, q.ID); err != nil {
		return fmt.Errorf("delete tags: %w", err)
	}
	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx, queries.QuestionTagInsert, q.ID, tagID); err != nil {
			return fmt.Errorf("attach tag %d: %w", tagID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *QuestionRepo) Close(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.QuestionClose, id)
	if err != nil {
		return fmt.Errorf("close question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *QuestionRepo) AssignSpecialist(ctx context.Context, id int64, specialistID *int64) error {
	tag, err := r.db.Exec(ctx, queries.QuestionAssign, id, specialistID)
	if err != nil {
		return fmt.Errorf("assign specialist: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *QuestionRepo) MarkDuplicate(ctx context.Context, id, duplicateOf int64) error {
	tag, err := r.db.Exec(ctx, queries.QuestionMarkDuplicate, id, duplicateOf)
	if err != nil {
		return fmt.Errorf("mark duplicate: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *QuestionRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, queries.QuestionDelete, id)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *QuestionRepo) CountByAuthor(ctx context.Context, userID int64) (int, error) {
	var count int
	if err := r.db.QueryRow(ctx, queries.QuestionCountByAuthor, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count questions by author: %w", err)
	}
	return count, nil
}

func (r *QuestionRepo) ListRecentByAuthor(ctx context.Context, userID int64, limit int) ([]*model.Question, error) {
	rows, err := r.db.Query(ctx, queries.QuestionListRecentByAuthor, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent questions by author: %w", err)
	}
	defer rows.Close()

	var qs []*model.Question
	for rows.Next() {
		q := &model.Question{}
		if err := rows.Scan(&q.ID, &q.AuthorID, &q.AuthorUsername, &q.CategoryID, &q.SpecialistID, &q.SpecialistUsername,
			&q.Title, &q.Body, &q.Status, &q.DuplicateOf, &q.ViewCount,
			&q.CreatedAt, &q.UpdatedAt, &q.AnswerCount, &q.HasVerified); err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}
	return qs, rows.Err()
}

type ExportQuestion struct {
	ID         int64
	Title      string
	Body       string
	AnswerBody *string
}

func (r *QuestionRepo) ListTitles(ctx context.Context) (map[string]struct{}, error) {
	rows, err := r.db.Query(ctx, queries.QuestionListTitles)
	if err != nil {
		return nil, fmt.Errorf("list question titles: %w", err)
	}
	defer rows.Close()

	titles := make(map[string]struct{})
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		titles[title] = struct{}{}
	}
	return titles, rows.Err()
}

func (r *QuestionRepo) ListForExport(ctx context.Context) ([]ExportQuestion, error) {
	rows, err := r.db.Query(ctx, queries.QuestionListForExport)
	if err != nil {
		return nil, fmt.Errorf("list questions for export: %w", err)
	}
	defer rows.Close()

	var result []ExportQuestion
	for rows.Next() {
		var q ExportQuestion
		if err := rows.Scan(&q.ID, &q.Title, &q.Body, &q.AnswerBody); err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, rows.Err()
}

func (r *QuestionRepo) GetHistory(ctx context.Context, questionID int64) ([]model.QuestionHistory, error) {
	rows, err := r.db.Query(ctx, queries.QuestionHistoryGet, questionID)
	if err != nil {
		return nil, fmt.Errorf("get question history: %w", err)
	}
	defer rows.Close()

	var result []model.QuestionHistory
	for rows.Next() {
		var h model.QuestionHistory
		if err := rows.Scan(&h.ID, &h.QuestionID, &h.EditorID, &h.EditorUsername, &h.Title, &h.Body, &h.EditedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, rows.Err()
}

func (r *QuestionRepo) getTags(ctx context.Context, questionID int64) ([]model.Tag, error) {
	rows, err := r.db.Query(ctx, queries.QuestionTagsGet, questionID)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}
