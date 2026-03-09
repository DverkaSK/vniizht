package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"vniizht/internal/errs"
	"vniizht/internal/repository"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	errs.Write(w, code, msg)
}

func mapErr(w http.ResponseWriter, err error, notFound, forbidden, conflict, internal string) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errs.Write(w, http.StatusNotFound, notFound)
	case errors.Is(err, errs.ErrForbidden):
		errs.Write(w, http.StatusForbidden, forbidden)
	case errors.Is(err, errs.ErrConflict):
		errs.Write(w, http.StatusConflict, conflict)
	default:
		errs.Write(w, http.StatusInternalServerError, internal)
	}
}
