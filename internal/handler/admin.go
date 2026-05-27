package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
	"vniizht/internal/service"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.UserListError)
		return
	}
	writeJSON(w, http.StatusOK, usersToResponse(users))
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string     `json:"username"`
		Email    string     `json:"email"`
		Password string     `json:"password"`
		Role     model.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		errs.Write(w, http.StatusBadRequest, errs.UserFieldsRequired)
		return
	}
	if req.Role == "" {
		req.Role = model.RoleUser
	}

	user, err := h.svc.CreateUser(r.Context(), req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		mapErr(w, err, "", "", errs.UserDuplicate, errs.UserCreateError)
		return
	}

	writeJSON(w, http.StatusCreated, userToResponse(user))
}

func (h *AdminHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.UserInvalidID)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}

	if err := h.svc.SetUserActive(r.Context(), id, req.IsActive); errors.Is(err, repository.ErrNotFound) {
		errs.Write(w, http.StatusNotFound, errs.UserNotFound)
		return
	} else if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.UserUpdateRoleError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.UserInvalidID)
		return
	}

	var req struct {
		Role model.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if !isValidRole(req.Role) {
		errs.Write(w, http.StatusBadRequest, errs.UserInvalidRole)
		return
	}

	if err := h.svc.ChangeRole(r.Context(), id, req.Role); errors.Is(err, repository.ErrNotFound) {
		errs.Write(w, http.StatusNotFound, errs.UserNotFound)
		return
	} else if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.UserUpdateRoleError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategories(r.Context())
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.CategoryListError)
		return
	}
	writeJSON(w, http.StatusOK, categoriesToResponse(categories))
}

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		SpecialistID *int64 `json:"specialist_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Name == "" {
		errs.Write(w, http.StatusBadRequest, errs.CategoryNameRequired)
		return
	}

	category, err := h.svc.CreateCategory(r.Context(), req.Name, req.Description, req.SpecialistID)
	if err != nil {
		mapAdminConflict(w, err, errs.CategoryDuplicate, errs.CategoryCreateError)
		return
	}

	writeJSON(w, http.StatusCreated, categoryToResponse(category))
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "categoryID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.CategoryInvalidID)
		return
	}

	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		SpecialistID *int64 `json:"specialist_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Name == "" {
		errs.Write(w, http.StatusBadRequest, errs.CategoryNameRequired)
		return
	}

	if err := h.svc.UpdateCategory(r.Context(), id, req.Name, req.Description, req.SpecialistID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errs.Write(w, http.StatusNotFound, errs.CategoryNotFound)
			return
		}
		mapAdminConflict(w, err, errs.CategoryDuplicate, errs.CategoryUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "categoryID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.CategoryInvalidID)
		return
	}

	if err := h.svc.DeleteCategory(r.Context(), id); errors.Is(err, repository.ErrNotFound) {
		errs.Write(w, http.StatusNotFound, errs.CategoryNotFound)
		return
	} else if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.CategoryDeleteError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListTags(r.Context())
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.TagListError)
		return
	}
	writeJSON(w, http.StatusOK, tagsToResponse(tags))
}

func (h *AdminHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Name == "" {
		errs.Write(w, http.StatusBadRequest, errs.TagNameRequired)
		return
	}

	tag, err := h.svc.CreateTag(r.Context(), req.Name)
	if err != nil {
		mapAdminConflict(w, err, errs.TagDuplicate, errs.TagCreateError)
		return
	}

	writeJSON(w, http.StatusCreated, tagToResponse(tag))
}

func (h *AdminHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "tagID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.TagInvalidID)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.InvalidBody)
		return
	}
	if req.Name == "" {
		errs.Write(w, http.StatusBadRequest, errs.TagNameRequired)
		return
	}

	if err := h.svc.UpdateTag(r.Context(), id, req.Name); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errs.Write(w, http.StatusNotFound, errs.TagNotFound)
			return
		}
		mapAdminConflict(w, err, errs.TagDuplicate, errs.TagUpdateError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "tagID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.TagInvalidID)
		return
	}

	if err := h.svc.DeleteTag(r.Context(), id); errors.Is(err, repository.ErrNotFound) {
		errs.Write(w, http.StatusNotFound, errs.TagNotFound)
		return
	} else if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.TagDeleteError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type userResponse struct {
	ID         int64      `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Role       model.Role `json:"role"`
	IsActive   bool       `json:"is_active"`
	Reputation int        `json:"reputation"`
}

func userToResponse(u *model.User) userResponse {
	return userResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Role:       u.Role,
		IsActive:   u.IsActive,
		Reputation: u.Reputation,
	}
}

func usersToResponse(users []*model.User) []userResponse {
	res := make([]userResponse, len(users))
	for i, u := range users {
		res[i] = userToResponse(u)
	}
	return res
}

func isValidRole(r model.Role) bool {
	switch r {
	case model.RoleUser, model.RoleSpecialist, model.RoleAdmin:
		return true
	}
	return false
}

type categoryResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SpecialistID *int64 `json:"specialist_id,omitempty"`
}

func categoryToResponse(c *model.Category) categoryResponse {
	return categoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Description:  c.Description,
		SpecialistID: c.SpecialistID,
	}
}

func categoriesToResponse(cats []*model.Category) []categoryResponse {
	res := make([]categoryResponse, len(cats))
	for i, c := range cats {
		res[i] = categoryToResponse(c)
	}
	return res
}

type tagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func tagToResponse(t *model.Tag) tagResponse {
	return tagResponse{ID: t.ID, Name: t.Name}
}

func tagsToResponse(tags []*model.Tag) []tagResponse {
	res := make([]tagResponse, len(tags))
	for i, t := range tags {
		res[i] = tagToResponse(t)
	}
	return res
}

func (h *AdminHandler) Suggest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := h.svc.Suggest(r.Context(), q)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.SearchError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func mapAdminConflict(w http.ResponseWriter, err error, conflictMsg, internalMsg string) {
	if errors.Is(err, errs.ErrConflict) {
		errs.Write(w, http.StatusConflict, conflictMsg)
		return
	}
	errs.Write(w, http.StatusInternalServerError, internalMsg)
}
