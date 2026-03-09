package handler

import (
	"net/http"
	"strconv"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/service"
)

type SearchHandler struct {
	svc *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		errs.Write(w, http.StatusBadRequest, errs.SearchQueryRequired)
		return
	}

	filter := model.SearchFilter{Query: q}

	if c := r.URL.Query().Get("category_id"); c != "" {
		id, err := strconv.ParseInt(c, 10, 64)
		if err == nil && id > 0 {
			filter.CategoryID = &id
		}
	}

	if t := r.URL.Query().Get("tag_id"); t != "" {
		id, err := strconv.ParseInt(t, 10, 64)
		if err == nil && id > 0 {
			filter.TagID = &id
		}
	}

	if s := r.URL.Query().Get("status"); s != "" {
		status := model.QuestionStatus(s)
		switch status {
		case model.QuestionOpen, model.QuestionClosed, model.QuestionDuplicate:
			filter.Status = &status
		}
	}

	results, err := h.svc.Search(r.Context(), filter)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.SearchError)
		return
	}

	if results == nil {
		results = []*model.SearchResult{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": results,
		"total": len(results),
	})
}
