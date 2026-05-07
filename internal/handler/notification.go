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
	"vniizht/internal/repository"
	"vniizht/internal/service"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	items, err := h.svc.List(r.Context(), user.ID, page)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationListError)
		return
	}
	writeJSON(w, http.StatusOK, notificationsToResponse(items))
}

func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	count, err := h.svc.CountUnread(r.Context(), user.ID)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationListError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "notifID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.NotificationInvalidID)
		return
	}
	if err := h.svc.MarkRead(r.Context(), id, user.ID); err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationUpdateError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if err := h.svc.MarkAllRead(r.Context(), user.ID); err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationUpdateError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	prefs, err := h.svc.GetPreferences(r.Context(), user.ID)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationListError)
		return
	}
	writeJSON(w, http.StatusOK, preferencesToResponse(prefs))
}

func (h *NotificationHandler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	var req []struct {
		Type         model.NotificationType `json:"type"`
		EmailEnabled bool                   `json:"email_enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}

	prefs := make([]repository.NotificationPreference, len(req))
	for i, p := range req {
		prefs[i] = repository.NotificationPreference{Type: p.Type, EmailEnabled: p.EmailEnabled}
	}

	if err := h.svc.SavePreferences(r.Context(), user.ID, prefs); err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.NotificationUpdateError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type preferenceResponse struct {
	Type         model.NotificationType `json:"type"`
	EmailEnabled bool                   `json:"email_enabled"`
}

func preferencesToResponse(prefs []repository.NotificationPreference) []preferenceResponse {
	res := make([]preferenceResponse, len(prefs))
	for i, p := range prefs {
		res[i] = preferenceResponse{Type: p.Type, EmailEnabled: p.EmailEnabled}
	}
	return res
}

type notificationResponse struct {
	ID        int64                    `json:"id"`
	Type      model.NotificationType  `json:"type"`
	Payload   model.NotificationPayload `json:"payload"`
	IsRead    bool                     `json:"is_read"`
	CreatedAt time.Time               `json:"created_at"`
}

func notificationToResponse(n *model.Notification) notificationResponse {
	return notificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		Payload:   n.Payload,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

func notificationsToResponse(ns []*model.Notification) []notificationResponse {
	res := make([]notificationResponse, len(ns))
	for i, n := range ns {
		res[i] = notificationToResponse(n)
	}
	return res
}
