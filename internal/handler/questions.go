package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/errs"
	"vniizht/internal/middleware"
	"vniizht/internal/model"
	"vniizht/internal/repository"
	"vniizht/internal/service"
)

type QuestionsHandler struct {
	svc *service.QuestionService
}

func NewQuestionsHandler(svc *service.QuestionService) *QuestionsHandler {
	return &QuestionsHandler{svc: svc}
}

func (h *QuestionsHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.QuestionFilter{
		Limit:  20,
		Offset: 0,
	}

	if s := r.URL.Query().Get("assigned_specialist_id"); s != "" {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
			return
		}
		f.AssignedSpecialistID = &id
	}

	if s := r.URL.Query().Get("status"); s != "" {
		st := model.QuestionStatus(s)
		switch st {
		case model.QuestionOpen, model.QuestionClosed, model.QuestionDuplicate:
			f.Status = &st
		default:
			errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidStatus)
			return
		}
	}

	if c := r.URL.Query().Get("category_id"); c != "" {
		id, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
			return
		}
		f.CategoryID = &id
	}

	if t := r.URL.Query().Get("tag_id"); t != "" {
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			errs.Write(w, http.StatusBadRequest, errs.TagInvalidID)
			return
		}
		f.TagID = &id
	}

	if p := r.URL.Query().Get("page"); p != "" {
		page, err := strconv.Atoi(p)
		if err != nil || page < 1 {
			errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidPage)
			return
		}
		f.Offset = (page - 1) * f.Limit
	}

	qs, total, err := h.svc.List(r.Context(), f)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.QuestionListError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  questionsToResponse(qs),
		"total":  total,
		"limit":  f.Limit,
		"offset": f.Offset,
	})
}

func (h *QuestionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	q, err := h.svc.Get(r.Context(), id)
	if err != nil {
		mapErr(w, err, errs.QuestionNotFound, "", "", errs.QuestionListError)
		return
	}

	writeJSON(w, http.StatusOK, questionToResponse(q))
}

func (h *QuestionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	var req struct {
		Title      string  `json:"title"`
		Body       string  `json:"body"`
		CategoryID *int64  `json:"category_id"`
		TagIDs     []int64 `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Title == "" || req.Body == "" {
		errs.Write(w, http.StatusBadRequest, errs.QuestionFieldsRequired)
		return
	}

	q := &model.Question{
		AuthorID:   user.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Body:       req.Body,
	}
	if err := h.svc.Create(r.Context(), q, req.TagIDs); err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.QuestionCreateError)
		return
	}

	q.AuthorUsername = user.Username
	writeJSON(w, http.StatusCreated, questionToResponse(q))
}

func (h *QuestionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	var req struct {
		Title  string  `json:"title"`
		Body   string  `json:"body"`
		TagIDs []int64 `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Title == "" || req.Body == "" {
		errs.Write(w, http.StatusBadRequest, errs.QuestionFieldsRequired)
		return
	}

	if err := h.svc.Update(r.Context(), id, user, req.Title, req.Body, req.TagIDs); err != nil {
		mapErr(w, err, errs.QuestionNotFound, errs.QuestionForbidden, errs.QuestionAlreadyClosed, errs.QuestionUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	if err := h.svc.Delete(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.QuestionNotFound, errs.QuestionForbidden, "", errs.QuestionUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionsHandler) Assign(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	var req struct {
		SpecialistID *int64 `json:"specialist_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}

	if err := h.svc.AssignSpecialist(r.Context(), id, req.SpecialistID, user); err != nil {
		mapErr(w, err, errs.QuestionNotFound, errs.QuestionForbidden, "", errs.QuestionAssignError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionsHandler) MarkDuplicate(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	var req struct {
		DuplicateOf int64 `json:"duplicate_of"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.DuplicateOf == 0 {
		errs.Write(w, http.StatusBadRequest, errs.QuestionDuplicateInvalid)
		return
	}

	if err := h.svc.MarkDuplicate(r.Context(), id, req.DuplicateOf, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalid):
			errs.Write(w, http.StatusBadRequest, errs.QuestionSelfDuplicate)
		default:
			mapErr(w, err, errs.QuestionNotFound, errs.QuestionForbidden, errs.QuestionAlreadyClosed, errs.QuestionDuplicateError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionsHandler) Close(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	if err := h.svc.Close(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.QuestionNotFound, errs.QuestionForbidden, "", errs.QuestionCloseError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionsHandler) History(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}
	history, err := h.svc.GetHistory(r.Context(), id)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.QuestionListError)
		return
	}
	type historyEntry struct {
		ID             int64     `json:"id"`
		EditorID       int64     `json:"editor_id"`
		EditorUsername string    `json:"editor_username"`
		Title          string    `json:"title"`
		Body           string    `json:"body"`
		EditedAt       time.Time `json:"edited_at"`
	}
	res := make([]historyEntry, len(history))
	for i, h := range history {
		res[i] = historyEntry{
			ID: h.ID, EditorID: h.EditorID, EditorUsername: h.EditorUsername,
			Title: h.Title, Body: h.Body, EditedAt: h.EditedAt,
		}
	}
	writeJSON(w, http.StatusOK, res)
}

type questionResponse struct {
	ID                 int64                `json:"id"`
	AuthorID           int64                `json:"author_id"`
	AuthorUsername     string               `json:"author_username"`
	CategoryID         *int64               `json:"category_id,omitempty"`
	SpecialistID       *int64               `json:"specialist_id,omitempty"`
	SpecialistUsername *string              `json:"specialist_username,omitempty"`
	DuplicateOf        *int64               `json:"duplicate_of,omitempty"`
	Title              string               `json:"title"`
	Body               string               `json:"body"`
	Status             model.QuestionStatus `json:"status"`
	ViewCount          int                  `json:"view_count"`
	AnswerCount        int                  `json:"answer_count"`
	HasVerified        bool                 `json:"has_verified"`
	Tags               []tagResponse        `json:"tags"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

func questionToResponse(q *model.Question) questionResponse {
	tags := make([]tagResponse, len(q.Tags))
	for i, t := range q.Tags {
		tags[i] = tagResponse{ID: t.ID, Name: t.Name}
	}
	return questionResponse{
		ID:                 q.ID,
		AuthorID:           q.AuthorID,
		AuthorUsername:     q.AuthorUsername,
		CategoryID:         q.CategoryID,
		SpecialistID:       q.SpecialistID,
		SpecialistUsername: q.SpecialistUsername,
		DuplicateOf:        q.DuplicateOf,
		Title:              q.Title,
		Body:               q.Body,
		Status:             q.Status,
		ViewCount:          q.ViewCount,
		AnswerCount:        q.AnswerCount,
		HasVerified:        q.HasVerified,
		Tags:               tags,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
	}
}

func questionsToResponse(qs []*model.Question) []questionResponse {
	res := make([]questionResponse, len(qs))
	for i, q := range qs {
		res[i] = questionToResponse(q)
	}
	return res
}
