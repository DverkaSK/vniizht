package handler

import (
	"encoding/json"
	"net/http"

	"vniizht/internal/errs"
	"vniizht/internal/middleware"
	"vniizht/internal/service"
)

type ImportHandler struct {
	svc *service.ImportService
}

func NewImportHandler(svc *service.ImportService) *ImportHandler {
	return &ImportHandler{svc: svc}
}

func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TrainingData []service.TrainingEntry `json:"training_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.ImportInvalidBody)
		return
	}
	if len(req.TrainingData) == 0 {
		errs.Write(w, http.StatusBadRequest, errs.ImportInvalidBody)
		return
	}

	user := middleware.CurrentUser(r)
	count, err := h.svc.Import(r.Context(), req.TrainingData, user.ID)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ImportError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"imported": count})
}

func (h *ImportHandler) Export(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.Export(r.Context())
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ExportError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"training_data": entries})
}
