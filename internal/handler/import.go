package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// ExportPeriod отдаёт вопросы за указанный период в формате CSV или JSON.
// Query-параметры: from, to — YYYY-MM-DD (необязательны); format — "csv" (по умолчанию) или "json".
func (h *ImportHandler) ExportPeriod(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time

	if s := r.URL.Query().Get("from"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			errs.Write(w, http.StatusBadRequest, "некорректный параметр from (ожидается YYYY-MM-DD)")
			return
		}
		from = &t
	}
	if s := r.URL.Query().Get("to"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			errs.Write(w, http.StatusBadRequest, "некорректный параметр to (ожидается YYYY-MM-DD)")
			return
		}
		t = t.AddDate(0, 0, 1) // включаем весь день «to»
		to = &t
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	rows, err := h.svc.ExportByPeriod(r.Context(), from, to)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ExportError)
		return
	}

	// Базовое имя файла из дат
	filename := "questions"
	if from != nil {
		filename += "_" + from.Format("2006-01-02")
	}
	if to != nil {
		filename += "_" + to.AddDate(0, 0, -1).Format("2006-01-02")
	}

	if format == "json" {
		entries := make([]service.TrainingEntry, 0, len(rows))
		for i, row := range rows {
			positive := ""
			if row.VerifiedAnswer != nil {
				positive = *row.VerifiedAnswer
			}
			entries = append(entries, service.TrainingEntry{
				ID:       i,
				Title:    row.Title,
				Questions: []string{row.Body},
				Positive: positive,
				Original: positive,
			})
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename+".json"))
		writeJSON(w, http.StatusOK, map[string]any{"training_data": entries})
		return
	}

	// CSV
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename+".csv"))

	// BOM — нужен Excel для корректного отображения UTF-8
	_, _ = w.Write([]byte("\xEF\xBB\xBF"))

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{
		"ID", "Дата создания", "Заголовок", "Описание",
		"Категория", "Теги", "Автор", "Статус", "Специалист",
		"Кол-во ответов", "Есть верифицированный ответ", "Верифицированный ответ",
	})

	statusLabel := map[string]string{
		"OPEN": "Открыт", "CLOSED": "Закрыт", "DUPLICATE": "Дубликат",
	}
	str := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	for _, row := range rows {
		verified := "Нет"
		if row.HasVerified {
			verified = "Да"
		}
		label := statusLabel[string(row.Status)]
		if label == "" {
			label = string(row.Status)
		}
		_ = cw.Write([]string{
			fmt.Sprintf("%d", row.ID),
			row.CreatedAt.Format("02.01.2006 15:04"),
			row.Title,
			row.Body,
			str(row.Category),
			str(row.Tags),
			row.Author,
			label,
			str(row.Specialist),
			fmt.Sprintf("%d", row.AnswerCount),
			verified,
			str(row.VerifiedAnswer),
		})
	}
	cw.Flush()
}

func (h *ImportHandler) Export(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.Export(r.Context())
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.ExportError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"training_data": entries})
}
