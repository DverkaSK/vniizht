package queries

const (
	AttachmentCreate = `
		INSERT INTO attachments (uploader_id, target_type, target_id, filename, object_key, mime_type, size_bytes)
		VALUES ($1, $2::attachment_target, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	AttachmentGetByID = `
		SELECT id, uploader_id, target_type, target_id, filename, object_key, mime_type, size_bytes, created_at
		FROM attachments
		WHERE id = $1`

	AttachmentListByTarget = `
		SELECT id, uploader_id, target_type, target_id, filename, object_key, mime_type, size_bytes, created_at
		FROM attachments
		WHERE target_type = $1::attachment_target AND target_id = $2
		ORDER BY created_at ASC`

	AttachmentDelete = `
		DELETE FROM attachments
		WHERE id = $1`
)
