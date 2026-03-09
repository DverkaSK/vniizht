package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/errs"
	"vniizht/internal/middleware"
	"vniizht/internal/model"
	"vniizht/internal/service"
)

type CommentsHandler struct {
	svc *service.CommentService
}

func NewCommentsHandler(svc *service.CommentService) *CommentsHandler {
	return &CommentsHandler{svc: svc}
}

func (h *CommentsHandler) List(w http.ResponseWriter, r *http.Request) {
	answerID, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AnswerInvalidID)
		return
	}

	comments, err := h.svc.List(r.Context(), answerID)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.CommentCreateError)
		return
	}

	resp := make([]commentResponse, len(comments))
	for i, c := range comments {
		resp[i] = commentToResponse(c)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *CommentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	answerID, err := strconv.ParseInt(chi.URLParam(r, "answerID"), 10, 64)
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
		errs.Write(w, http.StatusBadRequest, errs.CommentBodyRequired)
		return
	}

	c, err := h.svc.Create(r.Context(), answerID, user, req.Body)
	if err != nil {
		mapErr(w, err, errs.AnswerNotFound, "", "", errs.CommentCreateError)
		return
	}

	writeJSON(w, http.StatusCreated, commentToResponse(c))
}

func (h *CommentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "commentID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.CommentInvalidID)
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
		errs.Write(w, http.StatusBadRequest, errs.CommentBodyRequired)
		return
	}

	if err := h.svc.Update(r.Context(), id, user, req.Body); err != nil {
		mapErr(w, err, errs.CommentNotFound, errs.CommentForbidden, "", errs.CommentUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "commentID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.CommentInvalidID)
		return
	}

	if err := h.svc.Delete(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.CommentNotFound, errs.CommentForbidden, "", errs.CommentDeleteError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type commentResponse struct {
	ID             int64     `json:"id"`
	AnswerID       int64     `json:"answer_id"`
	AuthorID       int64     `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func commentToResponse(c *model.Comment) commentResponse {
	return commentResponse{
		ID:             c.ID,
		AnswerID:       c.AnswerID,
		AuthorID:       c.AuthorID,
		AuthorUsername: c.AuthorUsername,
		Body:           c.Body,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
