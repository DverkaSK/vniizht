package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/errs"
	"vniizht/internal/middleware"
	"vniizht/internal/model"
	"vniizht/internal/service"
)

const maxAttachmentSize = 50 << 20

type AttachmentsHandler struct {
	svc *service.AttachmentService
}

func NewAttachmentsHandler(svc *service.AttachmentService) *AttachmentsHandler {
	return &AttachmentsHandler{svc: svc}
}

func (h *AttachmentsHandler) List(w http.ResponseWriter, r *http.Request) {
	targetType := model.AttachmentTarget(r.URL.Query().Get("target_type"))
	switch targetType {
	case model.AttachmentTargetQuestion, model.AttachmentTargetAnswer, model.AttachmentTargetComment:
	default:
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidTarget)
		return
	}

	targetID, err := strconv.ParseInt(r.URL.Query().Get("target_id"), 10, 64)
	if err != nil || targetID <= 0 {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidTarget)
		return
	}

	attachments, err := h.svc.List(r.Context(), targetType, targetID)
	if err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.InternalError)
		return
	}

	resp := make([]attachmentResponse, len(attachments))
	for i, attachment := range attachments {
		resp[i] = attachmentToResponse(attachment)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AttachmentsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	r.Body = http.MaxBytesReader(w, r.Body, maxAttachmentSize+4096)
	if err := r.ParseMultipartForm(maxAttachmentSize); err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentTooLarge)
		return
	}

	targetID, err := strconv.ParseInt(r.FormValue("target_id"), 10, 64)
	if err != nil || targetID <= 0 {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidTarget)
		return
	}

	targetType := model.AttachmentTarget(r.FormValue("target_type"))
	switch targetType {
	case model.AttachmentTargetQuestion, model.AttachmentTargetAnswer, model.AttachmentTargetComment:
	default:
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidTarget)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentNoFile)
		return
	}
	defer file.Close()

	if header.Size > maxAttachmentSize {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentTooLarge)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	attachment := &model.Attachment{
		UploaderID: user.ID,
		TargetType: targetType,
		TargetID:   targetID,
		Filename:   header.Filename,
		MimeType:   mimeType,
		SizeBytes:  header.Size,
	}
	if err := h.svc.Create(r.Context(), attachment, file); err != nil {
		errs.Write(w, http.StatusInternalServerError, errs.AttachmentUploadError)
		return
	}

	writeJSON(w, http.StatusCreated, attachmentToResponse(attachment))
}

func (h *AttachmentsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "attachmentID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidID)
		return
	}

	attachment, obj, err := h.svc.Download(r.Context(), id)
	if err != nil {
		mapErr(w, err, errs.AttachmentNotFound, "", "", errs.InternalError)
		return
	}
	defer obj.Close()

	w.Header().Set("Content-Type", attachment.MimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+attachment.Filename+`"`)
	w.Header().Set("Content-Length", strconv.FormatInt(attachment.SizeBytes, 10))
	_, _ = io.Copy(w, obj)
}

func (h *AttachmentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "attachmentID"), 10, 64)
	if err != nil {
		errs.Write(w, http.StatusBadRequest, errs.AttachmentInvalidID)
		return
	}

	if err := h.svc.DeleteWithObject(r.Context(), id, user); err != nil {
		mapErr(w, err, errs.AttachmentNotFound, errs.AttachmentForbidden, "", errs.AttachmentDeleteError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type attachmentResponse struct {
	ID         int64  `json:"id"`
	UploaderID int64  `json:"uploader_id"`
	TargetType string `json:"target_type"`
	TargetID   int64  `json:"target_id"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mime_type"`
	SizeBytes  int64  `json:"size_bytes"`
	URL        string `json:"url"`
}

func attachmentToResponse(a *model.Attachment) attachmentResponse {
	return attachmentResponse{
		ID:         a.ID,
		UploaderID: a.UploaderID,
		TargetType: string(a.TargetType),
		TargetID:   a.TargetID,
		Filename:   a.Filename,
		MimeType:   a.MimeType,
		SizeBytes:  a.SizeBytes,
		URL:        "/attachments/" + strconv.FormatInt(a.ID, 10),
	}
}
