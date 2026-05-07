package queries

const (
	NotificationCreate = `
		INSERT INTO notifications (user_id, type, payload)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	NotificationListByUser = `
		SELECT id, user_id, type, payload, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	NotificationCountUnread = `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE`

	NotificationMarkRead = `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1 AND user_id = $2`

	NotificationMarkAllRead = `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND is_read = FALSE`

	// notification_preferences

	NotificationPreferencesGet = `
		SELECT t.type, COALESCE(p.email_enabled, TRUE) AS email_enabled
		FROM (VALUES
		    ('NEW_ANSWER'::notification_type),
		    ('ANSWER_VERIFIED'::notification_type),
		    ('QUESTION_ASSIGNED'::notification_type),
		    ('QUESTION_CLOSED'::notification_type),
		    ('NEW_COMMENT'::notification_type)
		) AS t(type)
		LEFT JOIN notification_preferences p
		    ON p.user_id = $1 AND p.type = t.type
		ORDER BY t.type`

	NotificationPreferenceUpsert = `
		INSERT INTO notification_preferences (user_id, type, email_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, type) DO UPDATE SET email_enabled = EXCLUDED.email_enabled`

	NotificationPreferenceEmailEnabled = `
		SELECT COALESCE(
		    (SELECT email_enabled FROM notification_preferences WHERE user_id = $1 AND type = $2),
		    TRUE
		)`
)
