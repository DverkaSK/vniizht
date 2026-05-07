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
	"vniizht/internal/service"
)

type AnswersHandler struct {
	svc *service.AnswerService
}

func NewAnswersHandler(svc *service.AnswerService) *AnswersHandler {
	return &AnswersHandler{svc: svc}
}

func (h *AnswersHandler) List(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	answers, err := h.svc.List(r.Context(), questionID, middleware.CurrentUser(r))
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.AnswerListError)
		return
	}

	resp := make([]answerResponse, len(answers))
	for i, a := range answers {
		resp[i] = answerToResponse(a)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AnswersHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	questionID, err := strconv.ParseInt(chi.URLParam(r, "questionID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.QuestionInvalidID)
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Body == "" {
		errs.Write(w, http.StatusBadRequest, errs.AnswerBodyRequired)
		return
	}

	a, err := h.svc.Create(r.Context(), questionID, user, req.Body)
	if err != nil {
		mapErr(w, err, errs.QuestionNotFound, "", errs.QuestionAlreadyClosed, errs.AnswerCreateError)
		return
	}

	writeJSON(w, http.StatusCreated, answerToResponse(a))
}

func (h *AnswersHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Body == "" {
		errs.Write(w, http.StatusBadRequest, errs.AnswerBodyRequired)
		return
	}

	if err := h.svc.Update(r.Context(), id, user, req.Body); err != nil {
		mapErr(w, err, errs.AnswerNotFound, errs.AnswerForbidden, "", errs.AnswerUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnswersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	if err := h.svc.Delete(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.AnswerNotFound, errs.AnswerForbidden, "", errs.AnswerDeleteError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnswersHandler) Verify(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	if err := h.svc.Verify(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.AnswerNotFound, errs.AnswerForbidden, "", errs.AnswerVerifyError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnswersHandler) Unverify(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	if err := h.svc.Unverify(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.AnswerNotFound, errs.AnswerForbidden, "", errs.AnswerUnverifyError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnswersHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	h.writeVote(w, r, voteActionCreate)
}

func (h *AnswersHandler) UpdateVote(w http.ResponseWriter, r *http.Request) {
	h.writeVote(w, r, voteActionUpdate)
}

func (h *AnswersHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	if err := h.svc.DeleteVote(r.Context(), id, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalid):
			errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidVote)
		default:
			mapErr(w, err, errs.AnswerNotFound, errs.AnswerSelfVote, errs.AnswerVoteExists, errs.AnswerVoteError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type voteAction int

const (
	voteActionCreate voteAction = iota
	voteActionUpdate
)

func (h *AnswersHandler) writeVote(w http.ResponseWriter, r *http.Request, action voteAction) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	var req struct {
		Value model.VoteValue `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}

	switch action {
	case voteActionCreate:
		err = h.svc.CreateVote(r.Context(), id, user, req.Value)
	case voteActionUpdate:
		err = h.svc.UpdateVote(r.Context(), id, user, req.Value)
	}
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalid):
			errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidVote)
		default:
			mapErr(w, err, errs.AnswerNotFound, errs.AnswerSelfVote, errs.AnswerVoteExists, errs.AnswerVoteError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnswersHandler) History(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}
	history, err := h.svc.GetHistory(r.Context(), id)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.AnswerListError)
		return
	}
	type historyEntry struct {
		ID             int64     `json:"id"`
		EditorID       int64     `json:"editor_id"`
		EditorUsername string    `json:"editor_username"`
		Body           string    `json:"body"`
		EditedAt       time.Time `json:"edited_at"`
	}
	res := make([]historyEntry, len(history))
	for i, h := range history {
		res[i] = historyEntry{
			ID: h.ID, EditorID: h.EditorID, EditorUsername: h.EditorUsername,
			Body: h.Body, EditedAt: h.EditedAt,
		}
	}
	writeJSON(w, http.StatusOK, res)
}

type answerResponse struct {
	ID             int64     `json:"id"`
	QuestionID     int64     `json:"question_id"`
	AuthorID       int64     `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Body           string    `json:"body"`
	IsVerified     bool      `json:"is_verified"`
	VoteScore      int       `json:"vote_score"`
	CurrentUserVote *model.VoteValue `json:"current_user_vote,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func answerToResponse(a *model.Answer) answerResponse {
	return answerResponse{
		ID:             a.ID,
		QuestionID:     a.QuestionID,
		AuthorID:       a.AuthorID,
		AuthorUsername: a.AuthorUsername,
		Body:           a.Body,
		IsVerified:     a.IsVerified,
		VoteScore:      a.VoteScore,
		CurrentUserVote: a.CurrentUserVote,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}
