package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"vniizht/internal/errs"
	"vniizht/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, errs.UserFieldsRequired)
		return
	}

	token, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, errs.InvalidCredentials)
		return
	}
	if errors.Is(err, service.ErrUserInactive) {
		writeError(w, http.StatusForbidden, errs.AccountDisabled)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, errs.InternalError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session"); err == nil {
		_ = h.auth.Logout(r.Context(), c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
