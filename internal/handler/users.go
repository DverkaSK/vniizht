package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/errs"
	"vniizht/internal/middleware"
	"vniizht/internal/repository"
	"vniizht/internal/service"
)

type UsersHandler struct {
	svc *service.UserService
}

func NewUsersHandler(svc *service.UserService) *UsersHandler {
	return &UsersHandler{svc: svc}
}

func (h *UsersHandler) Me(w http.ResponseWriter, r *http.Request) {
	profile, err := h.svc.GetMe(r.Context(), middleware.CurrentUser(r))
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ProfileError)
		return
	}

	writeJSON(w, http.StatusOK, privateUserResponse{
		ID:            profile.User.ID,
		Username:      profile.User.Username,
		Email:         profile.User.Email,
		Role:          string(profile.User.Role),
		Reputation:    profile.User.Reputation,
		IsActive:      profile.User.IsActive,
		CreatedAt:     profile.User.CreatedAt,
		QuestionCount: profile.QuestionCount,
		AnswerCount:   profile.AnswerCount,
	})
}

func (h *UsersHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.UserInvalidID)
		return
	}

	profile, err := h.svc.GetPublic(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		errs.Write(w, http.StatusNotFound, errs.UserNotFound)
		return
	}
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ProfileError)
		return
	}

	resp := publicUserResponse{
		ID:              profile.User.ID,
		Username:        profile.User.Username,
		Role:            string(profile.User.Role),
		Reputation:      profile.User.Reputation,
		CreatedAt:       profile.User.CreatedAt,
		QuestionCount:   profile.QuestionCount,
		AnswerCount:     profile.AnswerCount,
		RecentQuestions: []recentQuestionItem{},
		RecentAnswers:   []recentAnswerItem{},
	}

	for _, q := range profile.RecentQuestions {
		resp.RecentQuestions = append(resp.RecentQuestions, recentQuestionItem{
			ID:        q.ID,
			Title:     q.Title,
			Status:    string(q.Status),
			CreatedAt: q.CreatedAt,
		})
	}

	for _, a := range profile.RecentAnswers {
		resp.RecentAnswers = append(resp.RecentAnswers, recentAnswerItem{
			ID:         a.ID,
			QuestionID: a.QuestionID,
			Snippet:    truncate(a.Body, 120),
			CreatedAt:  a.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

type privateUserResponse struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	Reputation    int       `json:"reputation"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	QuestionCount int       `json:"question_count"`
	AnswerCount   int       `json:"answer_count"`
}

type publicUserResponse struct {
	ID              int64                `json:"id"`
	Username        string               `json:"username"`
	Role            string               `json:"role"`
	Reputation      int                  `json:"reputation"`
	CreatedAt       time.Time            `json:"created_at"`
	QuestionCount   int                  `json:"question_count"`
	AnswerCount     int                  `json:"answer_count"`
	RecentQuestions []recentQuestionItem `json:"recent_questions"`
	RecentAnswers   []recentAnswerItem   `json:"recent_answers"`
}

type recentQuestionItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type recentAnswerItem struct {
	ID         int64     `json:"id"`
	QuestionID int64     `json:"question_id"`
	Snippet    string    `json:"snippet"`
	CreatedAt  time.Time `json:"created_at"`
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
